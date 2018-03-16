package util

import (
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
	"container/list"
	"fmt"
)

// TODO startChan的作用??? 检查goroutine启动成功
func StartRoutine(f func()) {
	startChan := make(chan bool)
	go func() {
		startChan <- true
		f()
	}()

	<- startChan
	fmt.Println("start Routine done ")
}


type TimerThread struct {
	// save timer id
	nowTimerId uint32

	// for check timer id existing
	existTimerIdMap map[uint32]bool

	// mutex protect timer exist map
	mapMutex sync.Mutex

	// mutex protect timer lists
	mutex sync.Mutex

	// new added timers saved in newTimerList
	newTimerList *list.List

	// being added into head saved in currentTimerList
	currentTimerList *list.List

	// thread end flag
	end bool

	// now time (in ms)
	now uint64
}

type TimerObj interface {
	OnTimeout(timer *Timer)
}

type Timer struct {
	Id        uint32
	Obj       TimerObj
	AbsTime   uint64
	TimerType int
}

func NewTimerThread() *TimerThread {
	timerThread := &TimerThread{
		nowTimerId: 1,
		existTimerIdMap:make(map[uint32]bool,0),
		newTimerList: list.New(),
		currentTimerList: list.New(),
		end: false,
		now: NowTimeMs(),
	}

	StartRoutine(timerThread.main)
	return timerThread
}

func (timerThread *TimerThread) Stop() {
	timerThread.end = true
}

func (timerThread *TimerThread) main() {
	for !timerThread.end {
		// fire every 1 ms
		timerChan := time.NewTimer(1 * time.Millisecond).C
		<- timerChan
		timerThread.now = NowTimeMs()

	again:
	// iterator current timer list
		len := timerThread.currentTimerList.Len()
		for i:=0; i < len; i++ {
			obj := timerThread.currentTimerList.Front()
			timer := obj.Value.(*Timer)
			if timer.AbsTime > timerThread.now { // 时间还未到
				break
			}

			timerThread.currentTimerList.Remove(obj)
			timerThread.fireTimeout(timer)
		}

		// TODO 这个逻辑，每个加入的timer的timeOut都应该相同吧，可以参考java DelayQueue
		// if current timer list is empty, then exchange two list
		if timerThread.currentTimerList.Len() == 0 {
			timerThread.mutex.Lock()
			tmp := timerThread.currentTimerList
			timerThread.currentTimerList = timerThread.newTimerList
			timerThread.newTimerList = tmp
			timerThread.mutex.Unlock()
			// check timeout agant
			goto again
		}
	}
}

func (timerThread *TimerThread) fireTimeout(timer *Timer) {
	id := timer.Id
	timerThread.mapMutex.Lock()
	_, ok := timerThread.existTimerIdMap[id]
	timerThread.mapMutex.Unlock()

	if ok {
		log.Debug("fire timeout:%v, %d", timer.Obj, timer.TimerType)
		timer.Obj.OnTimeout(timer)
	}
}

func (timerThread *TimerThread) AddTimer(timeoutMs uint32, timeType int, obj TimerObj) uint32 {
	timerThread.mutex.Lock()
	absTime := timerThread.now + uint64(timeoutMs) // 过期时间
	timer := newTimer(timerThread.nowTimerId, absTime, timeType, obj)
	timerId := timerThread.nowTimerId
	timerThread.nowTimerId += 1

	timerThread.newTimerList.PushBack(timer)

	// add into exist timer map
	timerThread.mapMutex.Lock()
	timerThread.existTimerIdMap[timerId] = true
	timerThread.mapMutex.Unlock()

	timerThread.mutex.Unlock()

	return timerId
}

func (timerThread *TimerThread) DelTimer(timerId uint32) {
	timerThread.mapMutex.Lock()
	delete(timerThread.existTimerIdMap, timerId)
	timerThread.mapMutex.Unlock()
}

func newTimer(timerId uint32, absTime uint64, timeType int, obj TimerObj) *Timer {
	return &Timer{
		Id:        timerId,
		AbsTime:   absTime,
		Obj:       obj,
		TimerType: timeType,
	}
}