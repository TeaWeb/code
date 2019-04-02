package backup

import (
	"archive/zip"
	"os"
	"testing"
	"time"
)

func TestBackupZip(t *testing.T) {
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
	err := backupTask()
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("OK")
	}
}
