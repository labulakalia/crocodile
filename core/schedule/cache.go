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
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// recv task run log cache
// 需要实时的从rpc返回的stream中接收返回的数据 Write
// 还需要提供一个ReadOnly方法，以便其他地方读取
// Read读取的数据为copy，并且当开始读时，必须先从buf里将已经接收的数据拷贝一份返回后，再实时的返回新接收到的数据
// LogCacher interface
const (
	defaultCacheCap = 1e3
)

var (
	// cachepool
	cachepool = sync.Pool{
		New: func() interface{} {
			return NewLogCache()
		},
	}

	bytepool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 0, defaultCacheCap)
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
	Clean()
	// Close will stop task
	Close()
	// ReadAll data from buf
	ReadAll() (p string)
	// GetCode return task return code
	GetCode() int
	// SetRunHost save task run host
	Save(interface{})
	// GetRunHost task run host addr
	Get() interface{}
	// SetTaskStatus set task status
	// wait running finish fail cancel
	SetTaskStatus(define.TaskStatus)
	GetTaskStatus() define.TaskStatus
}

var _ LogCacher = &LogCache{}

// LogCache otehr could read latest data from buf and do not clean it
type LogCache struct {
	buf *bytes.Buffer
	// task log cache
	buffer []byte
	// count buffer
	count int
	// save task
	close bool
	// save task returncode,tasktype,if task could find run host,run host will be save hear
	// actual pb.TaskResp
	resdata interface{}
	// task is running
	status define.TaskStatus
}

// NewLogCache return impl LogCache struct
func NewLogCache() LogCacher {
	buf := make([]byte, 0, defaultCacheCap)
	var logcache = LogCache{
		buf:    bytes.NewBuffer(buf),
		buffer: bytepool.Get().([]byte),
		close:  false,
		status: define.TsWait,
	}
	return &logcache
}

// ReadOnly Get data from buf
func (l *LogCache) ReadOnly(p []byte, off int) (n int, err error) {
	bufcount := l.count
	countp := len(p)
	status := l.GetTaskStatus()
	log.Debug("start read log", zap.Int("off", off), zap.Int("count", bufcount), zap.String("status", status.String()))
	defer func() {
		time.Sleep(time.Second)
	}()
	if off >= bufcount {
		// bug
		// l.GetTaskStatus()
		// 需要结合任务的状态来获取任务日志
		if status == define.TsFinish {
			// 主要是任务完成后还可以获取日志，
			log.Debug("task is finished and read log finished")
			return 0, io.EOF
		}
		log.Info("offset ranther than total buf byte")
		return 0, ErrNoReadData
	}

	// countp + off > bufcount && countp > bufcount
	// 总数小于off + countp
	// 大于off buf count
	if off+countp > bufcount {
		// copy from off to l.buf.Len()
		copy(p, l.buffer[off:bufcount])
		// p = append(p, l.buffer[off:bufcount]...)
		return bufcount - off, nil
	}

	// 总数大于off + countp
	// if bufcount >= off+countp {
	copy(p, l.buffer[off:off+countp])
	// p = append(p, l.buffer[off:off+countp]...)
	return countp, nil
	// }
	// bufcount >
	// copy(p, l.buffer[off:off+len(p)])

	// bufio.ScanBytes(data []byte, atEOF bool)

}

// Write recv byte and write buf
func (l *LogCache) Write(p []byte) (n int, err error) {
	// n, err = l.buf.Write(p)
	l.count += len(p)
	l.buffer = append(l.buffer, p...)
	return
}

// WriteString Write String to buf
func (l *LogCache) WriteString(p string) (n int, err error) {
	now := time.Now().Local().Format("2006-01-02 15:04:05: ")
	w := now + p + "\n"
	l.Write([]byte(w))
	// l.buffer = append(l.buffer, []byte(w)...)
	// n, err = l.buf.WriteString(w)
	return
}

// WriteStringf Write Tmpl string format to buf
func (l *LogCache) WriteStringf(tmpl string, args ...interface{}) (n int, err error) {
	now := time.Now().Local().Format("2006-01-02 15:04:05: ") + fmt.Sprintf(tmpl, args...) +"\n"
	n, err = l.buf.WriteString(now)
	l.Write([]byte(now))
	return
}

// Close will stop read data
func (l *LogCache) Close() {
	l.close = true
	bytepool.Put(l.buffer[0:0])
}

// Clean clean all data
func (l *LogCache) Clean() {
	// l.buf.Reset()
	l.count = 0
	l.buffer = bytepool.Get().([]byte)
	l.resdata = nil
	l.close = false
}

// ReadAll will Get All recv data
func (l *LogCache) ReadAll() (p string) {
	p = string(l.buffer)
	return
}

// GetCode return task return code
func (l *LogCache) GetCode() int {
	if l.count >= 3 {
		codebyte := l.buffer[l.count-3:]
		code, err := strconv.Atoi(strings.TrimLeft(string(codebyte), " "))
		if err != nil {
			// if err != nil ,it is bug
			log.Error("Change str to int failed", zap.Error(err))
			return tasktype.DefaultExitCode
		}
		return code
	}
	// if code run there,this is bug
	log.Error("thia is bug,recv buf is nether than 3, get code failed")
	return tasktype.DefaultExitCode

}

// Save save task run data
func (l *LogCache) Save(data interface{}) {
	l.resdata = data
}

// Get task run host data
func (l *LogCache) Get() interface{} {
	return l.resdata
}

// // TaskStatus task run status
// type TaskStatus uint

// const (
// 	// waiting task is waiting pre task is running
// 	wait TaskStatus = iota + 1
// 	// running tassk is running
// 	run
// 	// finish task is run finish
// 	finish
// 	// fail task run fail
// 	fail
// 	// cancel task is cancel ,because pre task is run fail
// 	cancel
// 	// nodata parenttasks or childtasks no task
// 	nodata
// )

// func (t TaskStatus) String() string {
// 	switch t {
// 	case wait:
// 		return "wait"
// 	case run:
// 		return "run"
// 	case finish:
// 		return "finish"
// 	case fail:
// 		return "fail"
// 	case cancel:
// 		return "cancel"
// 	case nodata:
// 		return "nodata"
// 	default:
// 		return "unknown"
// 	}
// }

// // GetTasksTreeStatus return a slice
// func GetTasksTreeStatus() []*define.TaskStatusTree {
// 	retTasksStatus := make([]*define.TaskStatusTree, 0, 3)
// 	parentTasksStatus := &define.TaskStatusTree{
// 		Name:     "ParentTasks",
// 		Status:   nodata.String(),
// 		Children: make([]*define.TaskStatusTree, 0),
// 	}

// 	mainTaskStatus := &define.TaskStatusTree{
// 		// Name:   task.name,
// 		// ID:     taskid,
// 		Status: nodata.String(),
// 	}

// 	childTasksStatus := &define.TaskStatusTree{
// 		Name:     "ChildTasks",
// 		Status:   nodata.String(),
// 		Children: make([]*define.TaskStatusTree, 0),
// 	}

// 	retTasksStatus = append(retTasksStatus,
// 		parentTasksStatus,
// 		mainTaskStatus,
// 		childTasksStatus)
// 	return retTasksStatus
// }

// SetTaskStatus will be set task status is running or stop
func (l *LogCache) SetTaskStatus(status define.TaskStatus) {
	l.status = status
}

// GetTaskStatus will be set task status is running or stop
func (l *LogCache) GetTaskStatus() define.TaskStatus {
	return l.status
}
