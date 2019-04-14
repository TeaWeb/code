package teacache

import (
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/assert"
	"github.com/iwind/TeaGo/logs"
	"sync"
	"testing"
	"time"
)

func TestFileManager(t *testing.T) {
	a := assert.NewAssertion(t)

	m := NewFileManager()
	err := m.Write("123456", []byte("abc"))
	t.Log(err)
	a.IsNotNil(err)

	m.dir = Tea.TmpDir() + "/cache"

	err = m.Write("123456", []byte("abcd"))
	a.IsNil(err)
	if err != nil {
		t.Log(err)
	}

	data, err := m.Read("123456")
	a.IsNil(err, func() string {
		return err.Error()
	})
	a.IsTrue(string(data) == "abcd", func() string {
		return "data:" + string(data)
	})
	t.Log(string(data))

	time.Sleep(3 * time.Second)
}

func TestFileManagerConcurrent(t *testing.T) {
	m := NewFileManager()
	m.dir = Tea.TmpDir() + "/cache"

	wg := sync.WaitGroup{}
	wg.Add(1000)
	for i := 0; i < 1000; i ++ {
		go func(i int) {
			_, err := m.Read("123456")
			if err != nil {
				//logs.Println(err)
			} else {
				logs.Println("read success")
			}

			err = m.Write("123456", []byte("abc"))
			if err != nil {
				//logs.Println(i, err)
			} else {
				logs.Println("write success")
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestFileManager_Stat(t *testing.T) {
	m := NewFileManager()
	m.dir = Tea.Root + "/./cache"
	t.Log(m.Stat())
}

func TestFileManager_Clean(t *testing.T) {
	m := NewFileManager()
	m.dir = Tea.Root + "/./cache"
	t.Log(m.Clean())
}