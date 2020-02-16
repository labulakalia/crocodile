package schedule

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/utils/define"
	"go.uber.org/zap"
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
	logcachetmp.Close()

}

func TestLogCache_ReadAll(t *testing.T) {
	// readdata := []byte("etst1 test2t1tt1t11")
	// logcachetmp := NewLogCache()
	// logcachetmp.Write(readdata)
	// res := logcachetmp.ReadAll()

	// if len(readdata) != len(res) {
	// 	t.Error("read data not equal write data")
	// }
	id := "233903600084979712"
	a := define.MasterTask

	fmt.Printf("%s_%d\n", id, a)

}

func TestLogCache_GetCode(t *testing.T) {
	data := `2020-02-15 15:42:00: Start Prepare Task testrun[233903600084979712]
	2020-02-15 15:42:00: Start Conn Worker Host For Task testrun[233903600084979712]
	2020-02-15 15:42:00: Success Conn Worker Host[127.0.0.1:8081]
	2020-02-15 15:42:00: Start Get Task testrun[233903600084979712] Run Data
	2020-02-15 15:42:00: Success Get Task testrun[233903600084979712] Run Data
	2020-02-15 15:42:00: Start Run Task testrun[233903600084979712] On Host[127.0.0.1:8081]
	2020-02-15 15:42:00: Task testrun[233903600084979712] Start Output----------------
	0
	1
	2
	3
	4
	5
	6
	7
	8
	9
	10
	11
	12
	13
	14
	15
	16
	17
	18
	19
	20
	  2`
	var buf bytes.Buffer
	buf.WriteString(data)
	if buf.Len() >= 3 {
		codebyte := buf.Bytes()[buf.Len()-3:]
		code, err := strconv.Atoi(strings.TrimLeft(string(codebyte), " "))
		if err != nil {
			// if err != nil ,it is bug
			log.Error("Change str to int failed", zap.Error(err))
			return
		}
		buf.Truncate(buf.Len() - 3)
		t.Log(code)
	}
	t.Logf("%s", buf.Bytes())
}

func Benchmark_main(b *testing.B) {
	b.Run("Benchmark_append", Benchmark_append)
	b.Run("Benchmark_bytes", Benchmark_bytes)
}

func Benchmark_bytes(b *testing.B) {
	a := bytes.NewBuffer(make([]byte,0, 1e2))
	for i := 0; i < b.N; i++ {
		a.WriteString("111111111111111111111111111")
		// b.Log(a.Len())
		c := a.Len()
		c = c
	}
}

func Benchmark_append(b *testing.B) {
	a := make([]byte,0, 1e2)
	for i := 0; i < b.N; i++ {
		a = append(a, []byte("111111111111111111111111111")...)
		c := len(a)
		c = c
	}
}

func TestLogCach11e(t *testing.T) {
	// a := make([]int,0,1000)
	// a = append(a,1,2,3)
	// t.Log(a[1:3])
	// for i:=0;i<100;i++{
	// 	a = append(a,[]byte("111111")...)
	// }
	// t.Log(len(a),cap(a))
	// a = a[0:0:200]
	// t.Log(len(a),cap(a))
	ticker := time.NewTimer(time.Second)
	
	for {
		select {
		case <- ticker.C:
			fmt.Println("111")
		}
	}
}
