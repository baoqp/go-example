package gmushroom

import "sync"

type SharedLock struct {
	mutex sync.RWMutex
}

func (l *SharedLock) LockShared() {
	l.mutex.RLock()
}

func (l *SharedLock) UnlockShared() {
	l.mutex.RUnlock()
}


func (l *SharedLock) Lock() {
	l.mutex.Lock()
}

func (l *SharedLock) Unlock() {
	l.mutex.Unlock()
}

// TODO 这里是否有并发安全问题
func (l *SharedLock) Upgrade() {
	l.UnlockShared()
	l.Lock()
}

// TODO 这里是否有并发安全问题
func (l *SharedLock) Degrade() {
	l.Unlock()
	l.LockShared()
}