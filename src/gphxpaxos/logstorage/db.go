package logstorage

import (
	"strconv"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	log "github.com/sirupsen/logrus"
)

//---------------------PaxosComparator 实现goleveldb的comparer接口--------------------//
type PaxosComparator struct {
}

func (comparator *PaxosComparator) Compare(a, b []byte) int {
	ua, _ := strconv.ParseUint(string(a), 10, 64)
	ub, _ := strconv.ParseUint(string(b), 10, 64)

	if ua == ub {
		return 0
	}

	if ua < ub {
		return -1
	}
	return 1
}

func (comparator *PaxosComparator) Name() string {
	return "PaxosComparator"
}

func (comparator *PaxosComparator) Separator(dst, a, b []byte) []byte {
	return nil
}

func (comparator *PaxosComparator) Successor(dst, b []byte) []byte {
	return nil
}

//-------------------------------------Database---------------------------------------//
// 一个Database就是对一个LevelDB实例的封装
type Database struct {
	leveldb    *leveldb.DB
	comparator PaxosComparator
	hasInit    bool
	valueStore *LogStore
	dbPath     string
	myGroupIdx int
}

func (database *Database) Init(dbPath string, myGroupIdx int) error {

	if database.hasInit {
		return nil
	}

	database.myGroupIdx = myGroupIdx

	options := opt.Options{
		ErrorIfMissing: false,
		Comparer:       &database.comparator,
		// every group have different buffer size to avoid all group compact at the same time.
		WriteBuffer: int(1024*1024 + myGroupIdx + 10 + 1024),
	}
	var err error
	database.leveldb, err = leveldb.OpenFile(dbPath, &options) // 打开LevelDB
	if err != nil {
		log.Error("open leveldb fail, db path:%s", dbPath)
		return err
	}

	database.valueStore = NewLogStore()
	err = database.valueStore.Init(dbPath, database)
	if err != nil {
		log.Error("value store init fail:%v", err)
		return err
	}
	database.hasInit = true

	log.Info("db init OK, db path:%s", dbPath)

	return nil
}

func (database *Database) GetLogStorageDirPath(groupIdx int) {

}

func (database *Database) Get(iGroupIdx int, llInstanceID uint64, sValue string) {

}

func (database *Database) Put(oWriteOptions *WriteOptions, iGroupIdx int, llInstanceID uint64, sValue *string) {

}

func (database *Database) Del(oWriteOptions *WriteOptions, iGroupIdx WriteOptions, llInstanceID uint64) {

}

func (database *Database) ForceDel(oWriteOptions *WriteOptions, iGroupIdx int, llInstanceID *uint64) {
}

func (database *Database) GetMaxInstanceID(groupIdx int, instanceID uint64) {

}

func (database *Database) GetMaxInstanceIDFileID() (string, uint64, error) {
	return "", 0, nil
}

func (database *Database) rebuildOneIndex(instanceId uint64, fileIdstr string) error {
	return nil
}

func (database *Database) SetMinChosenInstanceID(writeOptions *WriteOptions, groupIdx int, minInstanceID uint64) {

}

func (database *Database) GetMinChosenInstanceID(groupIdx int, minInstanceID uint64) {

}

func (database *Database) ClearAllLog(groupIdx int) {

}

func (database *Database) SetSystemVariables(writeOptions *WriteOptions, groupIdx int, sBuffer *string) {

}

func (database *Database) GetSystemVariables(groupIdx int, sBuffer *string) {

}

func (database *Database) SetMasterVariables(writeOptions *WriteOptions, groupIdx int, sBuffer *string) {

}

func (database *Database) GetMasterVariables(groupIdx int, sBuffer *string) {

}
