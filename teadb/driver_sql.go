package teadb

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var jsonArrayIndexReg = regexp.MustCompile(`\.(\d+)`)

type SQLDriver struct {
	BaseDriver

	driver string

	db          *sql.DB
	conn        *sql.Conn
	queryLocker sync.Mutex
}

// 查找单条记录
func (this *SQLDriver) FindOne(query *Query, modelPtr interface{}) (interface{}, error) {
	ones, err := this.FindOnes(query.Limit(1), modelPtr)
	if err != nil {
		return nil, err
	}
	if len(ones) == 0 {
		return nil, nil
	}
	return ones[0], nil
}

// 查找多条记录
func (this *SQLDriver) FindOnes(query *Query, modelPtr interface{}) ([]interface{}, error) {
	conn, err := this.connect()
	if err != nil {
		return nil, err
	}

	this.queryLocker.Lock()
	defer this.queryLocker.Unlock()

	holder := NewSQLParamsHolder(this.driver)
	sqlString, err := this.asSQL(SQLSelect, query, holder, "", nil)
	if err != nil {
		return nil, err
	}

	if query.debug {
		logs.Println("sql:", sqlString)
	}

	stmt, err := conn.PrepareContext(context.Background(), sqlString)
	if err != nil {
		return nil, this.processError(err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	rows, err := stmt.Query(holder.Args...)
	if err != nil {
		return nil, this.processError(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	modelType := reflect.TypeOf(modelPtr)
	modelElem := modelType.Elem()
	method, methodExists := modelType.MethodByName("SetDBColumns")
	result := []interface{}{}
	for rows.Next() {
		cols, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		colPtrList := []interface{}{}
		for range cols {
			p := interface{}(nil)
			colPtrList = append(colPtrList, &p)
		}
		err = rows.Scan(colPtrList...)
		if err != nil {
			return nil, err
		}
		values := maps.Map{}
		for index, col := range cols {
			v := reflect.Indirect(reflect.ValueOf(colPtrList[index])).Interface()
			if v != nil {
				if _, ok := v.([]byte); ok {
					v = string(v.([]byte))
				}
			}
			values[col] = v
		}
		one := reflect.New(modelElem)
		if methodExists {
			method.Func.Call([]reflect.Value{one, reflect.ValueOf(values)})
		}
		result = append(result, one.Interface())
	}

	return result, nil
}

// 插入一条记录
func (this *SQLDriver) InsertOne(table string, modelPtr interface{}) error {
	conn, err := this.connect()
	if err != nil {
		return err
	}

	if modelPtr == nil {
		return errors.New("modelPtr should not be nil")
	}

	this.queryLocker.Lock()
	defer this.queryLocker.Unlock()

	modelType := reflect.TypeOf(modelPtr)
	method, methodExists := modelType.MethodByName("DBColumns")
	if !methodExists {
		return errors.New("'DBColumns() maps.Map' method not exist in '" + modelType.String() + "'")
	}
	result := method.Func.Call([]reflect.Value{reflect.ValueOf(modelPtr)})
	if len(result) != 1 {
		return errors.New("'DBColumns() maps.Map' method not exist in '" + modelType.String() + "'")
	}
	v := result[0].Interface()
	m, ok := v.(maps.Map)
	if !ok {
		return errors.New("'DBColumns() maps.Map' method not exist in '" + modelType.String() + "'")
	}

	b := strings.Builder{}
	b.WriteString("INSERT INTO " + this.quoteKeyword(table) + " (")
	index := 0
	args := []interface{}{}
	for k, v := range m {
		if index > 0 {
			b.WriteString(", ")
		}
		b.WriteString(this.quoteKeyword(k))
		args = append(args, v)
		index++
	}
	b.WriteString(") ")
	b.WriteString("VALUES (")
	for index := range args {
		if index > 0 {
			switch this.driver {
			case "mysql":
				b.WriteString(", ?")
			case "postgres":
				b.WriteString(", $" + strconv.Itoa(index+1))
			default:
				b.WriteString(", ?")
			}
		} else {
			switch this.driver {
			case "mysql":
				b.WriteString("?")
			case "postgres":
				b.WriteString("$" + strconv.Itoa(index+1))
			default:
				b.WriteString("?")
			}
		}
	}
	b.WriteString(")")
	stmt, err := conn.PrepareContext(context.Background(), b.String())
	if err != nil {
		return this.processError(err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.Exec(args...)

	return this.processError(err)
}

// 插入多条记录
func (this *SQLDriver) InsertOnes(table string, modelPtrSlice interface{}) error {
	conn, err := this.connect()
	if err != nil {
		return err
	}
	if modelPtrSlice == nil {
		return nil
	}

	this.queryLocker.Lock()
	defer this.queryLocker.Unlock()

	sliceType := reflect.TypeOf(modelPtrSlice)
	if sliceType.Kind() != reflect.Slice {
		return errors.New("only slice can be accepted in 'InsertOnes' method")
	}

	modelValues := reflect.ValueOf(modelPtrSlice)
	countValues := modelValues.Len()
	if modelValues.Len() == 0 {
		return nil
	}

	modelPtr := modelValues.Index(0).Interface()
	modelType := reflect.TypeOf(modelPtr)
	method, methodExists := modelType.MethodByName("DBColumns")
	if !methodExists {
		return errors.New("'DBColumns() maps.Map' method not exist in '" + modelType.String() + "'")
	}

	b := strings.Builder{}
	result := method.Func.Call([]reflect.Value{reflect.ValueOf(modelPtr)})
	if len(result) != 1 {
		return errors.New("'DBColumns() maps.Map' method not exist in '" + modelType.String() + "'")
	}
	v := result[0].Interface()
	m, ok := v.(maps.Map)
	if !ok {
		return errors.New("'DBColumns() maps.Map' method not exist in '" + modelType.String() + "'")
	}
	b.WriteString("INSERT INTO " + this.quoteKeyword(table) + " (")
	keys := []string{}
	index := 0
	for k := range m {
		if index > 0 {
			b.WriteString(", ")
		}
		b.WriteString(this.quoteKeyword(k))
		keys = append(keys, k)
		index++
	}
	b.WriteString(") ")
	b.WriteString("VALUES ")

	args := []interface{}{}
	paramIndex := 0
	for i := 0; i < countValues; i++ {
		if i > 0 {
			b.WriteString(", ")
		}
		modelValue := modelValues.Index(i)
		result := method.Func.Call([]reflect.Value{reflect.ValueOf(modelValue.Interface())})
		if len(result) != 1 {
			return errors.New("'DBColumns() maps.Map' method not exist in '" + modelType.String() + "'")
		}
		v := result[0].Interface()
		m, ok := v.(maps.Map)
		if !ok {
			return errors.New("'DBColumns() maps.Map' method not exist in '" + modelType.String() + "'")
		}

		b.WriteString("(")
		for index, key := range keys {
			if index > 0 {
				switch this.driver {
				case "mysql":
					b.WriteString(", ?")
				case "postgres":
					b.WriteString(", $" + strconv.Itoa(paramIndex+1))
					paramIndex++
				default:
					b.WriteString(", ?")
				}
			} else {
				switch this.driver {
				case "mysql":
					b.WriteString("?")
				case "postgres":
					b.WriteString("$" + strconv.Itoa(paramIndex+1))
					paramIndex++
				default:
					b.WriteString("?")
				}
			}
			args = append(args, m.Get(key))
		}
		b.WriteString(")")
	}

	stmt, err := conn.PrepareContext(context.Background(), b.String())
	if err != nil {
		return this.processError(err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.Exec(args...)

	return this.processError(err)
}

// 删除多条记录
func (this *SQLDriver) DeleteOnes(query *Query) error {
	conn, err := this.connect()
	if err != nil {
		return err
	}

	this.queryLocker.Lock()
	defer this.queryLocker.Unlock()
	holder := NewSQLParamsHolder(this.driver)
	sqlString, err := this.asSQL(SQLDelete, query, holder, "", nil)
	if err != nil {
		return err
	}

	stmt, err := conn.PrepareContext(context.Background(), sqlString)
	if err != nil {
		return this.processError(err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.Exec(holder.Args...)
	return this.processError(err)
}

// 修改多条记录
func (this *SQLDriver) UpdateOnes(query *Query, values map[string]interface{}) error {
	conn, err := this.connect()
	if err != nil {
		return err
	}

	this.queryLocker.Lock()
	defer this.queryLocker.Unlock()

	holder := NewSQLParamsHolder(this.driver)
	sqlString, err := this.asSQL(SQLUpdate, query, holder, "", values)
	if err != nil {
		return err
	}

	stmt, err := conn.PrepareContext(context.Background(), sqlString)
	if err != nil {
		return this.processError(err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.Exec(holder.Args...)

	return this.processError(err)
}

// 计算总数量
func (this *SQLDriver) Count(query *Query) (int64, error) {
	conn, err := this.connect()
	if err != nil {
		return 0, err
	}

	this.queryLocker.Lock()
	defer this.queryLocker.Unlock()

	holder := NewSQLParamsHolder(this.driver)
	query.Result("COUNT(*)")
	sqlString, err := this.asSQL(SQLSelect, query, holder, "", nil)

	if query.debug {
		logs.Println("sql:", sqlString)
	}

	if err != nil {
		return 0, err
	}
	stmt, err := conn.PrepareContext(context.Background(), sqlString)
	if err != nil {
		return 0, this.processError(err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	row := stmt.QueryRow(holder.Args...)
	if row == nil {
		return 0, nil
	}

	result := interface{}(nil)
	err = row.Scan(&result)
	if err != nil {
		return 0, this.processError(err)
	}

	return types.Int64(result), nil
}

// 计算总和
func (this *SQLDriver) Sum(query *Query, field string) (float64, error) {
	conn, err := this.connect()
	if err != nil {
		return 0, err
	}

	this.queryLocker.Lock()
	defer this.queryLocker.Unlock()

	holder := NewSQLParamsHolder(this.driver)
	query.Result("SUM(" + this.quoteKeyword(field) + ")")
	sqlString, err := this.asSQL(SQLSelect, query, holder, "", nil)

	if query.debug {
		logs.Println("sql:", sqlString)
	}

	if err != nil {
		return 0, err
	}
	stmt, err := conn.PrepareContext(context.Background(), sqlString)
	if err != nil {
		return 0, this.processError(err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	row := stmt.QueryRow(holder.Args...)
	if row == nil {
		return 0, nil
	}

	result := interface{}(nil)
	err = row.Scan(&result)
	if err != nil {
		return 0, this.processError(err)
	}

	return types.Float64(result), nil
}

// 计算平均值
func (this *SQLDriver) Avg(query *Query, field string) (float64, error) {
	conn, err := this.connect()
	if err != nil {
		return 0, err
	}

	this.queryLocker.Lock()
	defer this.queryLocker.Unlock()

	holder := NewSQLParamsHolder(this.driver)
	query.Result("AVG(" + this.quoteKeyword(field) + ")")
	sqlString, err := this.asSQL(SQLSelect, query, holder, "", nil)

	if query.debug {
		logs.Println("sql:", sqlString)
	}

	if err != nil {
		return 0, err
	}
	stmt, err := conn.PrepareContext(context.Background(), sqlString)
	if err != nil {
		return 0, this.processError(err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	row := stmt.QueryRow(holder.Args...)
	if row == nil {
		return 0, nil
	}

	result := interface{}(nil)
	err = row.Scan(&result)
	if err != nil {
		return 0, this.processError(err)
	}

	return types.Float64(result), nil
}

// 计算最小值
func (this *SQLDriver) Min(query *Query, field string) (float64, error) {
	conn, err := this.connect()
	if err != nil {
		return 0, err
	}

	this.queryLocker.Lock()
	defer this.queryLocker.Unlock()

	holder := NewSQLParamsHolder(this.driver)
	query.Result("MIN(" + this.quoteKeyword(field) + ")")
	sqlString, err := this.asSQL(SQLSelect, query, holder, "", nil)

	if query.debug {
		logs.Println("sql:", sqlString)
	}

	if err != nil {
		return 0, err
	}
	stmt, err := conn.PrepareContext(context.Background(), sqlString)
	if err != nil {
		return 0, this.processError(err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	row := stmt.QueryRow(holder.Args...)
	if row == nil {
		return 0, nil
	}

	result := interface{}(nil)
	err = row.Scan(&result)
	if err != nil {
		return 0, this.processError(err)
	}

	return types.Float64(result), nil
}

// 计算最大值
func (this *SQLDriver) Max(query *Query, field string) (float64, error) {
	conn, err := this.connect()
	if err != nil {
		return 0, err
	}

	this.queryLocker.Lock()
	defer this.queryLocker.Unlock()

	holder := NewSQLParamsHolder(this.driver)
	query.Result("MAX(" + this.quoteKeyword(field) + ")")
	sqlString, err := this.asSQL(SQLSelect, query, holder, "", nil)

	if query.debug {
		logs.Println("sql:", sqlString)
	}

	if err != nil {
		return 0, err
	}
	stmt, err := conn.PrepareContext(context.Background(), sqlString)
	if err != nil {
		return 0, this.processError(err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	row := stmt.QueryRow(holder.Args...)
	if row == nil {
		return 0, nil
	}

	result := interface{}(nil)
	err = row.Scan(&result)
	if err != nil {
		return 0, this.processError(err)
	}

	return types.Float64(result), nil
}

// 对数据进行分组统计
func (this *SQLDriver) Group(query *Query, groupField string, result map[string]Expr) ([]maps.Map, error) {
	conn, err := this.connect()
	if err != nil {
		return nil, err
	}

	this.queryLocker.Lock()
	defer this.queryLocker.Unlock()

	switch this.driver {
	case "mysql":
		_, err = conn.ExecContext(context.Background(), "SET SESSION sql_mode=(SELECT REPLACE(@@sql_mode,'ONLY_FULL_GROUP_BY',''));")
		if err != nil {
			return nil, this.processError(err)
		}
	case "postgres":
		// do nothing
	default:
		return nil, errors.New("unknown database type")
	}

	for field, expr := range result {
		switch e := expr.(type) {
		case *SumExpr:
			index := strings.Index(e.Field, ".")
			if index > -1 {
				e.Field = this.JSONExtractNumeric(e.Field[:index], e.Field[index+1:])
			}
			query.Result("SUM(" + this.quoteKeyword(e.Field) + ") AS " + this.quoteKeyword(field))
		case *AvgExpr:
			index := strings.Index(e.Field, ".")
			if index > -1 {
				e.Field = this.JSONExtractNumeric(e.Field[:index], e.Field[index+1:])
			}
			query.Result("AVG(" + this.quoteKeyword(e.Field) + ") AS " + this.quoteKeyword(field))
		case *MaxExpr:
			index := strings.Index(e.Field, ".")
			if index > -1 {
				e.Field = this.JSONExtractNumeric(e.Field[:index], e.Field[index+1:])
			}
			query.Result("MAX(" + this.quoteKeyword(e.Field) + ") AS " + this.quoteKeyword(field))
		case *MinExpr:
			index := strings.Index(e.Field, ".")
			if index > -1 {
				e.Field = this.JSONExtractNumeric(e.Field[:index], e.Field[index+1:])
			}
			query.Result("MIN(" + this.quoteKeyword(e.Field) + ") AS " + this.quoteKeyword(field))
		case string:
			index := strings.Index(e, ".")
			isAgg := false
			if index > -1 {
				if this.driver == "postgres" {
					isAgg = true
					e = this.JSONExtract("(array_agg("+this.quoteKeyword(e[:index])+"))[1]", e[index+1:])
				} else {
					e = this.JSONExtract(e[:index], e[index+1:])
				}
			}
			if this.driver == "postgres" && !isAgg {
				query.Result("(array_agg(" + this.quoteKeyword(e) + "))[1] AS " + this.quoteKeyword(field))
			} else {
				query.Result(this.quoteKeyword(e) + " AS " + this.quoteKeyword(field))
			}
		}
	}

	holder := NewSQLParamsHolder(this.driver)
	sqlString, err := this.asSQL(SQLSelect, query, holder, groupField, nil)
	if err != nil {
		return nil, err
	}

	if query.debug {
		logs.Println("sql:", sqlString)
	}

	stmt, err := conn.PrepareContext(context.Background(), sqlString)
	if err != nil {
		return nil, this.processError(err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	rows, err := stmt.Query(holder.Args...)
	if err != nil {
		return nil, this.processError(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	ones := []maps.Map{}
	for rows.Next() {
		columns, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		result := []interface{}{}
		for range columns {
			v := interface{}(nil)
			result = append(result, &v)
		}
		err = rows.Scan(result...)
		if err != nil {
			return nil, this.processError(err)
		}
		m := maps.Map{}
		for index, column := range columns {
			v := reflect.Indirect(reflect.ValueOf(result[index])).Interface()
			if v != nil {
				switch v1 := v.(type) {
				case []byte:
					v = string(v1)
				}
			}

			keys := strings.Split(column, ".")
			if len(keys) > 1 {
				this.setMapValue(m, keys, v)
			} else {
				m[column] = v
			}
		}

		ones = append(ones, m)
	}

	return ones, nil
}

// 测试数据库连接
func (this *SQLDriver) Test() error {
	if this.db == nil {
		return errors.New("no db available")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	done := make(chan bool, 1)
	isClosed := false

	var err error

	go func() {
		select {
		case <-ctx.Done():
			if !isClosed {
				done <- true
			}
			cancel()
			err = errors.New("connection timeout")
		}
	}()

	go func() {
		conn, err1 := this.db.Conn(ctx)
		if err1 != nil {
			err = err1
			_ = this.processError(err)
			return
		}
		if !isClosed {
			done <- true
		}
		err = conn.Close()
		if err != nil {
			logs.Error(err)
		}
	}()

	<-done
	close(done)
	isClosed = true

	return err
}

// 删除表
func (this *SQLDriver) DropTable(table string) error {
	conn, err := this.connect()
	if err != nil {
		return err
	}
	_, err = conn.ExecContext(context.Background(), "DROP TABLE "+this.quoteKeyword(table))
	return this.processError(err)
}

// 读取JSON字段
func (this *SQLDriver) JSONExtract(field string, path string) string {
	switch this.driver {
	case "mysql":
		return "JSON_EXTRACT(" + this.quoteKeyword(field) + ", \"$." + jsonArrayIndexReg.ReplaceAllString(path, "[$1]") + "\")"
	case "postgres":
		return "JSON_EXTRACT_PATH_TEXT(" + this.quoteKeyword(field) + ", '" + strings.ReplaceAll(path, ".", "', '") + "')"
	}
	return ""
}

// 读取JSON字段
func (this *SQLDriver) JSONExtractNumeric(field string, path string) string {
	switch this.driver {
	case "mysql":
		return "JSON_EXTRACT(" + this.quoteKeyword(field) + ", \"$." + jsonArrayIndexReg.ReplaceAllString(path, "[$1]") + "\")"
	case "postgres":
		return "JSON_EXTRACT_PATH_TEXT(" + this.quoteKeyword(field) + ", '" + strings.ReplaceAll(path, ".", "', '") + "')::\"float8\""
	}
	return ""
}

// 连接
func (this *SQLDriver) connect() (*sql.Conn, error) {
	this.queryLocker.Lock()
	defer this.queryLocker.Unlock()

	if this.db == nil {
		return nil, errors.New("no db available")
	}

	if this.conn != nil {
		return this.conn, nil
	}

	if this.conn != nil {
		return this.conn, nil
	}

	logs.Println("[db]connecting to database")

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := this.db.Conn(ctx)
	if err != nil {
		return nil, errors.New("[db]connect: " + err.Error())
	}
	this.conn = conn

	return this.conn, nil
}

// 处理错误
func (this *SQLDriver) processError(err error) error {
	if err == nil {
		return nil
	}
	if err == sql.ErrConnDone || err == driver.ErrBadConn {
		this.conn = nil
	}
	return err
}

func (this *SQLDriver) setMapValue(m maps.Map, keys []string, v interface{}) {
	l := len(keys)
	if l == 0 {
		return
	}
	if l == 1 {
		m[keys[0]] = v
		return
	}
	subM, ok := m[keys[0]]
	if ok {
		subV, ok := subM.(maps.Map)
		if ok {
			this.setMapValue(subV, keys[1:], v)
		} else {
			m[keys[0]] = maps.Map{}
			this.setMapValue(m[keys[0]].(maps.Map), keys[1:], v)
		}
	} else {
		m[keys[0]] = maps.Map{}
		this.setMapValue(m[keys[0]].(maps.Map), keys[1:], v)
	}
}

func (this *SQLDriver) asSQL(action SQLAction, query *Query, paramsHolder *SQLParamsHolder, groupField string, updateValues map[string]interface{}) (string, error) {
	b := strings.Builder{}

	switch action {
	case SQLSelect:
		b.WriteString("SELECT ")

		// result
		if len(query.resultFields) == 0 {
			b.WriteString("* ")
		} else {
			for index, field := range query.resultFields {
				if index > 0 {
					b.WriteString(", ")
				}
				b.WriteString(this.quoteKeyword(field))
			}
			b.WriteString(" ")
		}
	case SQLDelete:
		b.WriteString("DELETE ")
	case SQLUpdate:
		b.WriteString("UPDATE ")
	}

	// table
	if action == SQLSelect || action == SQLDelete {
		b.WriteString("FROM ")
	}
	b.WriteString(this.quoteKeyword(query.table))
	b.WriteString(" ")

	// set
	if action == SQLUpdate {
		b.WriteString("SET ")
		index := 0
		for k, v := range updateValues {
			if index > 0 {
				b.WriteString(", ")
			}
			b.WriteString(this.quoteKeyword(k))
			b.WriteString("=")
			b.WriteString(paramsHolder.Add(v))
			index++
		}
		b.WriteString(" ")
	}

	// where
	if len(query.operandMap) > 0 {
		where, err := this.buildWhere(query.operandMap, query.fieldMapping, paramsHolder)
		if err != nil {
			return "", err
		}
		if len(where) > 0 {
			b.WriteString("WHERE ")
			b.WriteString(where)
			b.WriteString(" ")
		}
	}

	// group
	hasGroups := false
	if action == SQLSelect && len(groupField) > 0 {
		if query.fieldMapping != nil {
			groupField = query.fieldMapping(groupField)
		}
		b.WriteString("GROUP BY " + this.quoteKeyword(groupField))
		b.WriteString(" ")
		hasGroups = true
	}

	// order
	if action == SQLSelect && len(query.sortFields) > 0 {
		b.WriteString("ORDER BY ")
		for index, field := range query.sortFields {
			if index > 0 {
				b.WriteString(", ")
			}
			if query.fieldMapping != nil {
				field.Name = query.fieldMapping(field.Name)
			}

			// 支持点符号
			if strings.IndexAny(field.Name, "( ") == -1 {
				dotIndex := strings.Index(field.Name, ".")
				if dotIndex > -1 {
					field.Name = this.JSONExtract(field.Name[:dotIndex], field.Name[:dotIndex+1])
				}
			}
			if hasGroups && this.driver == "postgres" && !strings.ContainsAny(field.Name, "({:") {
				b.WriteString("(array_agg(" + this.quoteKeyword(field.Name) + "))[1]")
			} else {
				b.WriteString(this.quoteKeyword(field.Name))
			}
			if field.IsAsc() {
				b.WriteString(" ASC ")
			} else {
				b.WriteString(" DESC ")
			}
		}
	}

	// limit
	if query.size > 0 {
		b.WriteString("LIMIT " + strconv.Itoa(query.size) + " ")
	}
	if query.offset > 0 {
		b.WriteString("OFFSET " + strconv.Itoa(query.offset) + " ")
	}

	if len(paramsHolder.Params) > 0 {
		return paramsHolder.Parse(b.String()), nil
	}

	return b.String(), nil
}

func (this *SQLDriver) quoteKeyword(s string) string {
	if strings.IndexAny(s, "( :{") > -1 {
		return s
	}
	switch this.driver {
	case "mysql":
		return "`" + s + "`"
	case "postgres":
		return "\"" + s + "\""
	}
	return "\"" + s + "\""
}

// 构造where
func (this *SQLDriver) buildWhere(operandMap OperandMap, fieldMapping func(field string) string, paramsHolder *SQLParamsHolder) (string, error) {
	b := strings.Builder{}
	hasPrefix := false
	for field, operands := range operandMap {
		if fieldMapping != nil {
			field = fieldMapping(field)
		}
		for _, op := range operands {
			if !hasPrefix {
				hasPrefix = true
			} else {
				b.WriteString(" AND ")
			}
			switch op.Code {
			case OperandEq:
				b.WriteString(this.quoteKeyword(field) + "=" + paramsHolder.Add(op.Value))
			case OperandLt:
				b.WriteString(this.quoteKeyword(field) + "<" + paramsHolder.Add(op.Value))
			case OperandLte:
				b.WriteString(this.quoteKeyword(field) + "<=" + paramsHolder.Add(op.Value))
			case OperandGt:
				b.WriteString(this.quoteKeyword(field) + ">" + paramsHolder.Add(op.Value))
			case OperandGte:
				b.WriteString(this.quoteKeyword(field) + ">=" + paramsHolder.Add(op.Value))
			case OperandIn:
				b.WriteString(this.quoteKeyword(field) + " IN " + paramsHolder.AddSlice(op.Value))
			case OperandNotIn:
				b.WriteString(this.quoteKeyword(field) + " NOT IN " + paramsHolder.AddSlice(op.Value))
			case OperandNeq:
				b.WriteString(this.quoteKeyword(field) + "!=" + paramsHolder.AddSlice(op.Value))
			case operandSQLCond:
				if op.Value != nil {
					cond, ok := op.Value.(*SQLCond)
					if ok {
						b.WriteString(cond.Expr)
						for k, v := range cond.Params {
							paramsHolder.AddHolder(k, v)
						}
					} else {
						return "", errors.New("operand 'operandSQLCond' value must be '*SQLCond'")
					}
				}
			case OperandOr:
				if op.Value != nil {
					operandMaps, ok := op.Value.([]OperandMap)
					if ok {
						if len(operandMaps) > 1 {
							b.WriteString("(")
						}
						for index, operandMap := range operandMaps {
							f, err := this.buildWhere(operandMap, fieldMapping, paramsHolder)
							if err != nil {
								return "", err
							}
							if index > 0 {
								b.WriteString("OR ")
							}
							b.WriteString("(")
							b.WriteString(f)
							b.WriteString(") ")
						}
						if len(operandMaps) > 1 {
							b.WriteString(") ")
						}
					} else {
						return "", errors.New("or: should be a valid []OperandMap")
					}
				} else {
					return "", errors.New("or: should be a valid []OperandMap")
				}
			default:
				return "", errors.New("invalid operand '" + op.Code + "'")
			}
		}
	}

	return b.String(), nil
}
