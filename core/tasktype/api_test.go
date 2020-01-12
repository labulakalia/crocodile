package tasktype

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestDataAPI_Run(t *testing.T) {
	returnstr := RandStringRunes(3000)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, returnstr)
	}))
	defer ts.Close()


	var dataapi = DataAPI{
		URL:    ts.URL,
		Method: http.MethodGet,
		Header: make(map[string]string),
	}
	reader := dataapi.Run(context.Background())
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
	if code != 200 {
		t.Errorf("status code is %d, not 200", code)
	}
}
