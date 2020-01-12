package schedule

import (
	"io"
	"os/exec"
	"testing"
	"time"
)

var (
	logcache LogCacher
	readcount int
	initcache bool
)

func TestLogCache(t *testing.T) {
	logcache = NewLogCache()
	go t.Run("read", TestLogCache_ReadOnly)
	t.Run("write", TestLogCache_Write)
	
}

func TestLogCache_ReadOnly(t *testing.T) {
	if !initcache {
		return
	}
	offset := 0
	var out = make([]byte, 25)
	for {
		n, err := logcache.ReadOnly(out, offset)
		if err == nil {
			if n > 0 {
				offset += n
				readcount += n
			}
			continue
		}
		if err == io.EOF {
			t.Log("read done")
			break
		} else if err == ErrNoReadData {
			t.Log("no data pls wait")
			time.Sleep(time.Millisecond * 500)
		} else {
			t.Fatal("read failed", err)
		}
	}
	buftotal,err := logcache.ReadAll()
	if err != nil {
		t.Fatal("logcache ReadAll failed",err)
	}
	if len(buftotal) != readcount  {
		t.Errorf("writecount byte %d not equal readcount byte %d",buftotal,readcount)
	}
}


func TestLogCache_Write(t *testing.T) {
	if !initcache {
		return
	}
	cmd := exec.Command("ping", "127.0.0.1", "-c", "3")
	pw, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	cmd.Start()
	var out = make([]byte, 1024)
	for {
		n, err := pw.Read(out)
		if err != nil {
			if err == io.EOF {
				logcache.Close()
				break
		
			}
			t.Error("read failed", err)
		}
		// fmt.Printf("%s",out[:n])
		_, err = logcache.Write(out[:n])
		if err != nil {
			t.Fatal("logcache write failed", err)
		}
	}
	cmd.Wait()
}

func TestLogCache_Close(t *testing.T) {
	logcachetmp := NewLogCache()
	logcachetmp.Write([]byte("testetstettttetttstet"))
	err := logcachetmp.Close()
	if err != nil {
		t.Fatalf("close logcache failed %v",err)
	}
}

func TestLogCache_ReadAll(t *testing.T) {
	readdata := []byte("etst1 test2t1tt1t11")
	logcachetmp := NewLogCache()
	logcachetmp.Write(readdata)
	res, err := logcachetmp.ReadAll()
	if err != nil {
		t.Fatalf("logcachetmp.ReadAll  failed %v",err)
	}
	if len(readdata)!=len(res) {
		t.Error("read data not equal write data")
	}
}