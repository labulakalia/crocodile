package utils

import (
	"fmt"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"sync"
	"time"
)

const (
	numberBits  uint8 = 12                      // 每个集群的节点生成的ID数最大位数
	workerBits  uint8 = 10                      // 工作机器的ID位数
	numberMax   int64 = -1 ^ (-1 << numberBits) // ID序号的最大值  4096
	workerIdMax int64 = -1 ^ (-1 << workerBits) // 工作机器的ID最大值 1024
	timeShift   uint8 = workerBits + numberBits // 时间戳向左偏移量
	workerShift uint8 = numberBits              // 节点ID向左偏移数
	sub         int64 = 1525705533000           // 减去现在的时间戳

	defaultWorkerId = 1 // 默认worker
)

var (
	_worker *Worker
	_once   sync.Once
)

type Worker struct {
	mu        sync.RWMutex
	timestamp int64 // 上一次生成ID的时间戳
	workerId  int64 // 节点ID
	number    int64 // 已经生成的ID数
}

func newWorker(workerId int64) (*Worker, error) {
	if workerId < 0 || workerId > workerIdMax {
		return nil, errors.New(fmt.Sprintf("unvalid workid 0~%d", workerIdMax))
	}
	return &Worker{
		timestamp: 0,
		workerId:  workerId,
		number:    0,
	}, nil
}

// 生成ID
func (w *Worker) generateId() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now().UnixNano() / 1e6
	if w.timestamp == now {
		w.number++
		// 判断是否已经超出最大的限制的ID
		if w.number > numberMax {
			// 等待下一毫秒
			for now <= w.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		// 新的一毫秒将number 修改为0 timestamp修改为now
		w.number = 0
		w.timestamp = now
	}
	id := (now-sub)<<timeShift | w.workerId<<workerShift | w.number
	return id
}

func GetId() int64 {
	_once.Do(func() {
		w, err := newWorker(defaultWorkerId)
		if err != nil {
			log.Fatal("NewWorker failed", zap.Error(err))
		}
		_worker = w
	})

	return _worker.generateId()
}
