package backup

import (
	"archive/zip"
	"github.com/TeaWeb/code/teatesting"
	"os"
	"testing"
	"time"
)

func TestBackupZip(t *testing.T) {
	if teatesting.IsGlobal() {
		return
	}

	tmp := "/tmp/backup.test.zip"
	fp, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer fp.Close()

	z := zip.NewWriter(fp)
	defer z.Close()

	{
		h := &zip.FileHeader{
			Name:     "test.txt",
			Modified: time.Now(),
		}
		w, err := z.CreateHeader(h)
		if err != nil {
			t.Fatal(err)
		}
		w.Write([]byte("Hello, World"))
	}

	{
		h := &zip.FileHeader{
			Name:     "1/2/3/test.txt",
			Modified: time.Now(),
		}
		w, err := z.CreateHeader(h)
		if err != nil {
			t.Fatal(err)
		}
		w.Write([]byte("Hello, Hello"))
	}
}

func TestBackup(t *testing.T) {
	if teatesting.IsGlobal() {
		return
	}

	err := backupTask()
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("OK")
	}
}
