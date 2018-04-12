package gmushroom

import (
	"sync/atomic"
	"fmt"
)

type Latch struct {
	users int64 // should operatio atomic
	id    PageId
	sLock SharedLock
}

func newLatch() *Latch {
	return &Latch{
		id: 0x7FFFFFFF,
	}
}

func (l *Latch) Id() PageId {
	return l.id
}

func (l *Latch) SetId(id PageId) {
	l.id = id
}

func (l *Latch) Pin() {
	atomic.AddInt64(&l.users, 1)
}

func (l *Latch) UnPin() {
	atomic.AddInt64(&l.users, -1)
}

func (l *Latch) Free() bool {
	return atomic.LoadInt64(&l.users) == 0
}

func (l *Latch) LockShared() {
	Assert(atomic.LoadInt64(&l.users) > 0, "latch.users <= 0")
	l.sLock.LockShared()
}

func (l *Latch) UnlockShared() {
	Assert(atomic.LoadInt64(&l.users) > 0, "latch.users <= 0")
	l.sLock.UnlockShared()
	l.UnPin()
}

func (l *Latch) Lock() {
	Assert(atomic.LoadInt64(&l.users) > 0, "latch.users <= 0")
	l.sLock.Lock()
}

func (l *Latch) Unlock() {
	Assert(atomic.LoadInt64(&l.users) > 0, "latch.users <= 0")
	l.sLock.Unlock()
	l.UnPin()
}

func (l *Latch) Upgrade() {
	Assert(atomic.LoadInt64(&l.users) > 0, "latch.users <= 0")
	l.sLock.Upgrade()
}

func (l *Latch) Degrade() {
	Assert(atomic.LoadInt64(&l.users) > 0, "latch.users <= 0")
	l.sLock.Degrade()
}

func (l *Latch) String() string {
	if l.id == 0x7FFFFFFF {
		return ""
	}

	return fmt.Sprintf("%d:%d", l.id, l.users)
}
