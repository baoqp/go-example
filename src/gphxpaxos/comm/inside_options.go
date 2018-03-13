package comm

import (
	"sync"
	"gphxpaxos/util"
)

// TODO 各种Get方法部分参数的设置理由
type InsideOptions struct {
	isLargeBufferMode bool
	isIMFollower      bool
	groupCount        int
}

var once sync.Once
var iinsideOptions *InsideOptions

func GetInstance() *InsideOptions {
	once.Do(func() {
		iinsideOptions = &InsideOptions{
			isLargeBufferMode: false,
			isIMFollower:      false,
			groupCount:        1}

	})
	return iinsideOptions
}

func (insideOptions *InsideOptions) SetAsLargeBufferMode() {
	insideOptions.isLargeBufferMode = true
}

func (insideOptions *InsideOptions) SetAsFollower() {
	insideOptions.isIMFollower = true
}

func (insideOptions *InsideOptions) SetGroupCount(iGroupCount int) {
	insideOptions.groupCount = iGroupCount
}

func (insideOptions *InsideOptions) GetMaxBufferSize() int {
	if insideOptions.isLargeBufferMode {
		return 52428800
	}
	return 10485760
}

func (insideOptions *InsideOptions) GetStartPrepareTimeoutMs() int {
	if insideOptions.isLargeBufferMode {
		return 15000
	}
	return 2000
}

func (insideOptions *InsideOptions) GetStartAcceptTimeoutMs() int {
	if insideOptions.isLargeBufferMode {
		return 15000
	}
	return 2000
}

func (insideOptions *InsideOptions) GetMaxPrepareTimeoutMs() int {
	if insideOptions.isLargeBufferMode {
		return 90000
	}
	return 8000
}

func (insideOptions *InsideOptions) GetMaxAcceptTimeoutMs() int {
	if insideOptions.isLargeBufferMode {
		return 90000
	}
	return 8000
}

func (insideOptions *InsideOptions) GetMaxIOLoopQueueLen() int {
	if insideOptions.isLargeBufferMode {
		return 1024/insideOptions.groupCount + 100
	}
	return 10240/insideOptions.groupCount + 1000

}

func (insideOptions *InsideOptions) GetMaxQueueLen() int {
	if insideOptions.isLargeBufferMode {
		return 1024
	}
	return 10240
}

func (insideOptions *InsideOptions) GetAskforLearnerval() int {
	if !insideOptions.isIMFollower {
		if insideOptions.isLargeBufferMode {
			return 50000 + util.Rand(10000)
		} else {
			return 2500 + util.Rand(500)
		}
	} else {
		if insideOptions.isLargeBufferMode {
			return 30000 + util.Rand(15000)
		} else {
			return 2000 + util.Rand(1000)
		}
	}
}

func (insideOptions *InsideOptions) GetLearnerReceiverAckLead() int {
	if insideOptions.isLargeBufferMode {
		return 2
	}
	return 4
}

func (insideOptions *InsideOptions) GetLearnerSenderPrepareTimeoutMs() int {
	if iinsideOptions.isLargeBufferMode {
		return 6000
	}

	return 5000
}

func (insideOptions *InsideOptions) GetLearnerSenderAckTimeoutMs() int {
	if iinsideOptions.isLargeBufferMode {
		return 6000
	}

	return 5000
}

func (insideOptions *InsideOptions) GetLearnerSenderAckLead() int {
	if iinsideOptions.isLargeBufferMode {
		return 5
	}

	return 21
}

func (insideOptions *InsideOptions) GetTcpOutQueueDropTimeMs() int {
	if iinsideOptions.isLargeBufferMode {
		return 20000
	}

	return 5000
}

func (insideOptions *InsideOptions) GetLogFileMaxSize() int {
	if iinsideOptions.isLargeBufferMode {
		return 524288000
	}

	return 104857600
}

func (insideOptions *InsideOptions) GetTcpConnectionNonActiveTimeout() int {
	if iinsideOptions.isLargeBufferMode {
		return 600000
	}

	return 60000
}

func (insideOptions *InsideOptions) GetLearnerSenderSendQps() int {
	if iinsideOptions.isLargeBufferMode {
		return 10000 / insideOptions.groupCount
	}

	return 100000 / insideOptions.groupCount
}

func (insideOptions *InsideOptions) GetCleanerDeleteQps() int {
	if iinsideOptions.isLargeBufferMode {
		return 30000 / insideOptions.groupCount
	}

	return 300000 / insideOptions.groupCount
}
