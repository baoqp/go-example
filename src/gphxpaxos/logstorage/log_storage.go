package logstorage

import "math"

const MINCHOSEN_KEY = math.MaxUint64 - 1
const SYSTEMVARIABLES_KEY = MINCHOSEN_KEY - 1
const MASTERVARIABLES_KEY = MINCHOSEN_KEY - 2

type WriteOptions struct {
	Sync bool
}

/*
  LogStorage接口，可以根据需要有不同的实现。

  在phxpaxos的实现中, 实际数据value保存在文件（vfile)中（log_store封装了相关操作），在LevelDB中保存了value的索引

  LevelDB中数据格式：
    key - instance
    value format - fileid(int32)+file offset(uint64)+cksum of file value(uint32)

  元文件保存了当前正在使用的vfile的fileid
  meta file format(data path/vpath/meta):
    current file id(int32)
    file id cksum(uint32)

  value文件格式
  data file(data path/vpath/fileid.f) data format:
    data len(int32)
    value(data len) format:
      instance id(uint64)
      acceptor state data(data len - sizeof(uint64))
 */
type LogStorage interface {
	GetMaxInstanceIDFileID() (string, uint64, error)
	rebuildOneIndex(instanceId uint64, fileIdstr string) error
}
