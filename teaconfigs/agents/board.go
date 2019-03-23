package agents

import (
	"errors"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
)

// 看板图表定义
type BoardChart struct {
	AppId   string `yaml:"appId" json:"appId"`
	ItemId  string `yaml:"itemId" json:"itemId"`
	ChartId string `yaml:"chartId" json:"chartId"`
}

// 看板
type Board struct {
	Filename string        `yaml:"filename" json:"filename"`
	Charts   []*BoardChart `yaml:"charts" json:"charts"`
}

// 取得Agent看板
func NewAgentBoard(agentId string) *Board {
	filename := "board." + agentId + ".conf"
	file := files.NewFile(Tea.ConfigFile("agents/" + filename))
	if file.Exists() {
		reader, err := file.Reader()
		if err != nil {
			logs.Error(err)
			return nil
		}
		defer reader.Close()
		board := &Board{}
		err = reader.ReadYAML(board)
		if err != nil {
			return nil
		}
		board.Filename = filename
		return board
	} else {
		return &Board{
			Filename: filename,
			Charts:   []*BoardChart{},
		}
	}
}

// 添加图表
func (this *Board) AddChart(appId, itemId, chartId string) {
	if this.FindChart(chartId) != nil {
		return
	}
	this.Charts = append(this.Charts, &BoardChart{
		AppId:   appId,
		ItemId:  itemId,
		ChartId: chartId,
	})
}

// 查找图表
func (this *Board) FindChart(chartId string) *BoardChart {
	for _, c := range this.Charts {
		if c.ChartId == chartId {
			return c
		}
	}
	return nil
}

// 删除图表
func (this *Board) RemoveChart(chartId string) {
	result := []*BoardChart{}
	for _, c := range this.Charts {
		if c.ChartId == chartId {
			continue
		}
		result = append(result, c)
	}
	this.Charts = result
}

// 删除App相关的所有图表
func (this *Board) RemoveApp(appId string) {
	result := []*BoardChart{}
	for _, c := range this.Charts {
		if c.AppId == appId {
			continue
		}
		result = append(result, c)
	}
	this.Charts = result
}

// 保存
func (this *Board) Save() error {
	if len(this.Filename) == 0 {
		return errors.New("filename should be specified")
	}
	writer, err := files.NewWriter(Tea.ConfigFile("agents/" + this.Filename))
	if err != nil {
		return err
	}
	defer writer.Close()
	_, err = writer.WriteYAML(this)
	return err
}
