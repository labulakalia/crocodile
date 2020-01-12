package tasktype

import (
	"bytes"

	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"testing"
)

func TestDataShell_Run(t *testing.T) {
	var data TaskRuner = &DataShell{
		Command: "ping 127.0.0.1 -c 3",
	}
	reader := data.Run(context.Background())
	defer reader.Close()
	var (
		lastrecv []byte
		out      = make([]byte, 1024)
		buf      bytes.Buffer
	)
	for {
		n, err := reader.Read(out)
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Error(err)
			break
		}
		if n > 0 {
			fmt.Printf("%s", out[:n])
			lastrecv = out[:n]
			buf.Write(out[:n])
		}
	}
	code, err := strconv.Atoi(strings.TrimLeft(string(lastrecv), " "))
	if err != nil {
		t.Errorf("change str %s to int failed", strings.TrimLeft(string(lastrecv), " "))
		return
	}
	buf.Truncate(len(buf.Bytes()) - 4) // remove return code in run write
	if code != 0 {
		t.Errorf("status code is %d, not 0", code)
	}
}