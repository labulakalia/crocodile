package schedule

import (
	"io"
	"fmt"
	"os/exec"

	"testing"
	"time"
)

var (
	logcache  LogCacher
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
	buftotal := logcache.ReadAll()
	if len(buftotal) != readcount {
		t.Errorf("writecount byte %s not equal readcount byte %d", buftotal, readcount)
	}
}

func TestLogCache_Write(t *testing.T) {
	tmplogcache := NewLogCache()
	// tmplogcache.WriteStringf("Start Prepare Task %s[%s]", "taskdata.Name", "id")
	// tmplogcache.WriteStringf("Start Conn Worker Host For Task %s[%s]", "tashhshshs", "1111")
	// tmplogcache.WriteStringf("Start Conn Worker Host For Task %s[%s]", "tashhshshs", "1111")
	// tmplogcache.WriteString("Start Conn Worker Host For Task ")
	// tmplogcache.WriteString("Start Conn Worker Host For Task")
	// fmt.Println(tmplogcache.ReadAll())

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
				tmplogcache.Close()
				break

			}
			t.Error("read failed", err)
		}
		
		_, err = tmplogcache.Write(out[:n])
		if err != nil {
			t.Fatal("logcache write failed", err)
		}
		fmt.Printf("%s", out[:n])
	}
	cmd.Wait()
}

func TestLogCache_Close(t *testing.T) {
	logcachetmp := NewLogCache()
	logcachetmp.Write([]byte("testetstettttetttstet"))
	err := logcachetmp.Close()
	if err != nil {
		t.Fatalf("close logcache failed %v", err)
	}
}

func TestLogCache_ReadAll(t *testing.T) {
	readdata := []byte("etst1 test2t1tt1t11")
	logcachetmp := NewLogCache()
	logcachetmp.Write(readdata)
	res := logcachetmp.ReadAll()

	if len(readdata) != len(res) {
		t.Error("read data not equal write data")
	}
}
