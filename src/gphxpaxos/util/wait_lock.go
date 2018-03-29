package util

import (
	"time"
	"sync"
	"errors"
)

var Waitlock_Timeout = errors.New("waitlock timeout")

type Waitlock struct {
	mutex sync.Mutex
	inUse bool
	waitChan chan bool
	WaitCount int32
}

func NewWaitlock() *Waitlock {
	return &Waitlock{
		inUse:false,
		waitChan: make(chan bool, 0),
		WaitCount:0,
	}
}

func (waitLock *Waitlock) Lock(waitMs int) (int, error) {
	timeOut := false

	now := NowTimeMs()
	getLock := false

	waitLock.mutex.Lock()
	if !waitLock.inUse {
		waitLock.inUse = true
		getLock = true
	} else {
		waitLock.WaitCount += 1
	}
	waitLock.mutex.Unlock()

	if getLock {
		// assume there is no time cost
		return 0, nil
	}

	timer := time.NewTimer(time.Duration(waitMs) * time.Millisecond)
	select {
	case <- timer.C:
		timeOut = true
		break
	case <- waitLock.waitChan:
		break
	}

	waitLock.mutex.Lock()
	waitLock.WaitCount -= 1
	waitLock.mutex.Unlock()
	if timeOut {
		return -1, Waitlock_Timeout
	}
	timer.Stop()

	return int(NowTimeMs() - now), nil
}

func (waitLock *Waitlock) Unlock() {
	waitLock.mutex.Lock()
	waitLock.inUse = false
	if waitLock.WaitCount == 0 {
		waitLock.mutex.Unlock()
		return
	}
	waitLock.mutex.Unlock()

	timeOut := false
	timer := time.NewTimer(time.Duration(1) * time.Millisecond)
	select {
	case <- timer.C:
		timeOut = true
		break
	case waitLock.waitChan <- true:
		break
	}

	if !timeOut {
		timer.Stop()
	}
}