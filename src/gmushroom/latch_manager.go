package gmushroom

const (
	LatchMax = 8
	Hash     = 16
	Mask     = Hash - 1
)

type LatchSet struct {
	latches [LatchMax]Latch
	lock    SharedLock
}

func (ls *LatchSet) String() string {
	var res string
	for _, latch := range ls.latches {
		res += latch.String()
	}
	return res
}

// TODO 为什么是这种加锁方式
func (ls *LatchSet) GetLatch(id PageId) *Latch {
	var latch *Latch = nil

	ls.lock.LockShared()
	for i := 0; i < LatchMax; i++ {
		if ls.latches[i].Id() == id {
			ls.latches[i].Pin()
			latch = &ls.latches[i]
			break
		}
	}
	ls.lock.UnlockShared()

	if latch != nil {
		return latch
	}

	ls.lock.Lock()
	for i := 0; i < LatchMax; i++ {
		if ls.latches[i].Id() == id {
			latch = &ls.latches[i]
			break
		}

		if ls.latches[i].Free() && latch == nil {
			latch = &ls.latches[i]
		}
	}
	latch.SetId(id)
	latch.Pin()
	ls.lock.Unlock()
	return latch

}

// 锁管理器
type LathcManager struct {
	latchSets [Hash]LatchSet
}

// 分段管理
func (lm *LathcManager) GetLatch(id PageId) *Latch {
	return lm.latchSets[id & Mask].GetLatch(id)
}

func (lm *LathcManager) String() string {
	var res string
	for _, ls := range lm.latchSets {
		res += ls.String() + "\n"
	}
	return res
}