package msgqueue

import (
	"log"
	"ngame/msgqueue/queue"
	"runtime"
	"sync/atomic"
	"time"
)

const (
	idle int32 = iota
	running
)

type Dispatcher interface {
	Schedule(f func())
	Throughput() int
}

type goRoutineDispatcher int

func (g goRoutineDispatcher) Schedule(fn func()) {

}

func (g goRoutineDispatcher) Throughput() int {
	return int(g)
}

type synchronizedDispatcher int

func NewDefaultDispatcher(val int) Dispatcher {
	return goRoutineDispatcher(val)
}

func (s synchronizedDispatcher) Schedule(f func()) {

}

func (s synchronizedDispatcher) Throughput() int {
	return int(s)
}

type Statistics interface {
	MessageReceived(msg interface{})
}

type MsgQueue struct {
	mailBox         *queue.Queue
	msgNumber       int32
	dispatcher      Dispatcher
	mailBoxStats    []Statistics
	schedulerStatus int32
	suspended       bool
	bStop           bool
}

func (m *MsgQueue) PostMsg(msg interface{}) {
	if m.isStop() {
		return
	}

	for _, ms := range m.mailBoxStats {
		ms.MessageReceived(msg)
	}
	m.mailBox.Push(msg)
	atomic.AddInt32(&m.msgNumber, 1)
	if atomic.CompareAndSwapInt32(&m.schedulerStatus, idle, running) {
		m.dispatcher.Schedule(m.processMessages)
	}
}

func (m *MsgQueue) processMessages() {
process:
	m.run()

	atomic.StoreInt32(&m.schedulerStatus, idle)
	val := atomic.LoadInt32(&m.msgNumber)
	if val > 0 {
		if atomic.CompareAndSwapInt32(&m.schedulerStatus, idle, running) {
			goto process
		}
	}
}

func (m *MsgQueue) run() {
	var msg interface{}

	defer func() {
		if r := recover(); r != nil {
			log.Println("(m *MsgQueue) run() err : ", r)
		}
	}()

	i, t := 0, m.dispatcher.Throughput()
	if i > t {
		i = 0
		runtime.Gosched()
	}

	i++

	if msg = m.mailBox.Pop(); msg != nil {
		atomic.AddInt32(&m.msgNumber, -1)
		if m.mailBoxStats[0] != nil {
			m.mailBoxStats[0].MessageReceived(msg)
		}
	} else {
		return
	}
}

func (m *MsgQueue) isStop() bool {
	return m.bStop
}

func (m *MsgQueue) stop() {
	if !m.bStop {
		m.bStop = true
	}
}

func UnboundedMsgQueue(dispatcher Dispatcher, mailBoxStats ...Statistics) *MsgQueue {
	return &MsgQueue{
		mailBox:      queue.New(),
		mailBoxStats: mailBoxStats,
		dispatcher:   dispatcher,
	}
}

type SyncMsgChan struct {
	MsgData   interface{}
	MsgChan   chan interface{}
	ChanStats *int32
}

func GenerateSyncMsg(data interface{}) *SyncMsgChan {
	chanState := new(int32)
	*chanState = 1

	return &SyncMsgChan{data, make(chan interface{}, 1), chanState}
}

func (v *SyncMsgChan) WaitSyncReply(d time.Duration) (interface{}, bool) {
	if d < time.Microsecond {
		d = time.Microsecond
	}

	timer := time.NewTimer(d)
	for {
		select {
		case <-timer.C:
			atomic.StoreInt32(v.ChanStats, 0)
			return nil, false
		case msg, ok := <-v.MsgChan:
			atomic.StoreInt32(v.ChanStats, 0)
			timer.Reset(0)
			if ok {
				return msg, true
			}
		}
	}
}
