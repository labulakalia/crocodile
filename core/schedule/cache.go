package schedule

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/tasktype"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// recv task run log cache
// 需要实时的从rpc返回的stream中接收返回的数据 Write
// 还需要提供一个ReadOnly方法，以便其他地方读取
// Read读取的数据为copy，并且当开始读时，必须先从buf里将已经接收的数据拷贝一份返回后，再实时的返回新接收到的数据
// LogCacher interface
const (
	defaultCacheCap = 10240
)

var (
	// cachepool
	cachepool = sync.Pool{
		New: func() interface{} {
			return NewLogCache()
		},
	}
	// ErrNoReadData define when read offset rather than len(bug), return this error
	ErrNoReadData = errors.New("no read data from cache")
)

// LogCacher Write buf,ReadOnly from buf
// a stream write and read
type LogCacher interface {
	// Read byte, but it not clean reader data
	// ReadOnly read data from buf
	// 从off后开始读，最多可以读取len(p)个字节，
	// - close为true,返回0,EOF
	// - off > len(l.buf),返回 noreaddata 这时读取方重新读取可以等待一会
	// - off+len(p) > len(l.buf),返回 p.buf.Len-off,nil
	// - off+len(p) < len(l.buf),返回 len(p), nil
	ReadOnly(p []byte, off int) (n int, err error)
	// Write byte buf
	Write(p []byte) (n int, err error)
	// Write byte buf,do not use it write recv task
	WriteString(p string) (n int, err error)
	// WriteStringf ,do not use it write recv task
	WriteStringf(tmpl string, args ...interface{}) (n int, err error)
	// Close will can not Read, and clean buf
	Close() (err error)
	// ReadAll data from buf
	ReadAll() (p string)
	// GetCode return task return code
	GetCode() int
	// SetRunHost save task run host
	Save(interface{})
	// GetRunHost task run host addr
	Get() interface{}
}

var _ LogCacher = &LogCache{}

// LogCache otehr could read latest data from buf and do not clean it
type LogCache struct {
	buf     *bytes.Buffer
	close   bool
	runhost interface{}
}

// NewLogCache return impl LogCache struct
func NewLogCache() LogCacher {
	buf := make([]byte, 0, defaultCacheCap)
	var logcache = LogCache{
		buf:   bytes.NewBuffer(buf),
		close: false,
	}
	return &logcache
}

// ReadOnly Get data from buf
func (l *LogCache) ReadOnly(p []byte, off int) (n int, err error) {
	if l.close {
		log.Info("cache close")
		return 0, io.EOF
	}
	bufcount := l.buf.Len()
	if off > bufcount {
		log.Info("offset ranther than total buf byte")
		return 0, ErrNoReadData
	}

	if off+len(p) > bufcount {
		// copy from off to l.buf.Len()
		copy(p, l.buf.Bytes()[off:bufcount])
		return bufcount - off, nil
	}
	copy(p, l.buf.Bytes()[off:off+len(p)])
	return len(p), nil
}

// Write recv byte and write buf
func (l *LogCache) Write(p []byte) (n int, err error) {
	n, err = l.buf.Write(p)
	return
}

// WriteString Write String to buf
func (l *LogCache) WriteString(p string) (n int, err error) {
	now := time.Now().Local().Format("2006-01-02 15:04:05")
	w := fmt.Sprintf("%s: %s\n", now, p)
	n, err = l.buf.WriteString(w)
	return
}

// WriteStringf Write Tmpl string format to buf
func (l *LogCache) WriteStringf(tmpl string, args ...interface{}) (n int, err error) {
	n, err = l.buf.WriteString(fmt.Sprintf(tmpl, args...))
	return
}

// Close will stop read data and clean all data
func (l *LogCache) Close() (err error) {
	l.close = true
	l.buf.Reset()
	// logcache.P
	return
}

// ReadAll will Get All recv data
func (l *LogCache) ReadAll() (p string) {
	p = l.buf.String()
	return
}

// GetCode return task return code
func (l *LogCache) GetCode() int {
	if l.buf.Len() >= 3 {
		codebyte := l.buf.Bytes()[l.buf.Len()-3:]
		code, err := strconv.Atoi(strings.TrimLeft(string(codebyte), " "))
		if err != nil {
			// if err != nil ,it is bug
			log.Error("Change str to int failed", zap.Error(err))
			return tasktype.DefaultExitCode
		}
		l.buf.Truncate(l.buf.Len() - 3)
		return code
	}
	// if code run there,this is bug
	log.Error("thia is bug,recv buf is nether than 3, get code failed")
	return tasktype.DefaultExitCode

}

// Save save task run data 
func (l *LogCache) Save(data interface{}) {
	l.runhost = data
}

// Get task run host data
func (l *LogCache) Get() interface{} {
	return l.runhost
}
