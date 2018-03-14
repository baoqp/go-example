package logstorage

import (
	"strconv"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	log "github.com/sirupsen/logrus"
	"fmt"
	"gphxpaxos/comm"
	"math/rand"
	"gphxpaxos/util"
	"os"
	"sync"
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

func (database *Database) GetDBPath() string {
	return database.dbPath
}

func (database *Database) ClearAllLog() error {
	var systemVariablesBuffer string
	err := database.GetSystemVariables(&systemVariablesBuffer)
	if err != nil && err != comm.ErrKeyNotFound {
		log.Error("GetSystemVariables fail, ret %v", err)
		return err
	}

	var masterVariablesBuffer string
	err = database.GetMasterVariables(&masterVariablesBuffer)
	if err != nil && err != comm.ErrKeyNotFound {
		log.Error("GetMasterVariables fail, ret %v", err)
		return err
	}

	database.hasInit = false
	database.leveldb = nil
	database.valueStore = nil

	bakPath := database.dbPath + ".bak"
	err = util.DeleteDir(bakPath)
	if err != nil {
		log.Error("delete bak dir %s fail:%v", bakPath, err)
		return err
	}

	os.Rename(database.dbPath, bakPath)

	err = database.Init(database.dbPath, database.myGroupIdx)
	if err != nil {
		log.Error("init again fail:%v", err)
		return err
	}

	options := WriteOptions{
		Sync: true,
	}
	if len(systemVariablesBuffer) > 0 {
		err = database.SetSystemVariables(&options, systemVariablesBuffer)
		if err != nil {
			log.Error("SetSystemVariables fail:%v", err)
			return err
		}
	}
	if len(masterVariablesBuffer) > 0 {
		err = database.SetMasterVariables(&options, masterVariablesBuffer)
		if err != nil {
			log.Error("SetMasterVariables fail:%v", err)
			return err
		}
	}

	return nil
}

func (database *Database) Get(instanceId uint64) ([]byte, error) {
	var err error

	if !database.hasInit {
		err = fmt.Errorf("not init yet")
		return nil, err
	}

	var fileId string
	err = database.getFromLevelDb(instanceId, &fileId) // 从LevelDB中获取fileid
	if err != nil {
		return nil, err
	}

	var fileinstanceId uint64
	value, err := database.fileIdToValue(fileId, &fileinstanceId) // 从vfile中获取value
	if err != nil {
		return nil, err
	}

	if fileinstanceId != instanceId {
		log.Error("file instance id %d not equal to instance id %d", fileinstanceId, instanceId)
		return nil, comm.ErrInvalidInstanceId
	}

	return value, nil
}

func (database *Database) Put(options *WriteOptions, instanceId uint64, value []byte) error {
	var err error

	if !database.hasInit {
		err = fmt.Errorf("not init yet")
		return err
	}

	var fileId string
	err = database.valueToFileId(options, instanceId, value, &fileId)
	if err != nil {
		return err
	}

	return database.putToLevelDB(false, instanceId, []byte(fileId))
}

func (database *Database) Del(options *WriteOptions, instanceId uint64) error {
	if !database.hasInit {
		log.Error("no init yet")
		return comm.ErrDbNotInit
	}

	key := database.genKey(instanceId)

	// vfile并不用每次都删除，只要把LevelDB中的删除就访问不到vfile里面的value了
	if rand.Intn(100) < 10 {
		fileId, err := database.leveldb.Get([]byte(key), &opt.ReadOptions{})
		if err != nil {
			if err == leveldb.ErrNotFound {
				log.Error("leveldb.get not found, instance:%d", instanceId)
				return nil
			}
			log.Error("leveldb.get fail:%v", err)
			return comm.ErrGetFail
		}

		err = database.valueStore.Del(string(fileId), instanceId)
		if err != nil {
			return err
		}

	}
	writeOptions := opt.WriteOptions{
		Sync: options.Sync,
	}
	err := database.leveldb.Delete([]byte(key), &writeOptions)
	if err != nil {
		log.Error("leveldb.delete fail, instanceId %d, err:%v", instanceId, err)
		return err
	}
	return nil
}

func (database *Database) ForceDel(options WriteOptions, instanceId uint64) error {

	if !database.hasInit {
		log.Error("no init yet")
		return comm.ErrDbNotInit
	}

	key := database.genKey(instanceId)
	fileId, err := database.leveldb.Get([]byte(key), &opt.ReadOptions{})
	if err != nil {
		if err == leveldb.ErrNotFound {
			log.Error("leveldb.get not found, instance:%d", instanceId)
			return nil
		}
		log.Error("leveldb.get fail:%v", err)
		return comm.ErrGetFail
	}

	err = database.valueStore.ForceDel(string(fileId), instanceId)
	if err != nil {
		return err
	}

	writeOptions := opt.WriteOptions{
		Sync: options.Sync,
	}
	err = database.leveldb.Delete([]byte(key), &writeOptions)
	if err != nil {
		log.Error("leveldb.delete fail, instanceId %d, err:%v", instanceId, err)
		return err
	}
	return nil
}

// 获取最大的instanceId，其实就是LevelDB最大的key
func (database *Database) GetMaxinstanceId() (uint64, error) {
	var instanceId uint64 = MINCHOSEN_KEY
	iter := database.leveldb.NewIterator(nil, &opt.ReadOptions{})

	iter.Last()

	for {
		if !iter.Valid() {
			break
		}

		instanceId = database.getinstanceIdFromKey(string(iter.Key()))
		if instanceId == MINCHOSEN_KEY || instanceId == SYSTEMVARIABLES_KEY || instanceId == MASTERVARIABLES_KEY {
			iter.Prev()
		} else {
			return instanceId, nil
		}
	}

	return comm.INVALID_INSTANCEID, comm.ErrKeyNotFound
}

func (database *Database) GetMaxinstanceIdFileID() (string, uint64, error) {
	maxinstanceId, err := database.GetMaxinstanceId()
	if err != nil {
		return "", 0, nil
	}

	key := database.genKey(maxinstanceId)
	value, err := database.leveldb.Get([]byte(key), &opt.ReadOptions{})
	if err != nil {
		if err == leveldb.ErrNotFound {
			return "", 0, comm.ErrKeyNotFound
		}

		log.Error("leveldb.get fail:%v", err)
		return "", 0, comm.ErrGetFail
	}

	return string(value), maxinstanceId, nil
}

// 替换一条记录，LevelDB中的替换使用追加就行，LevelDB默认会读取相同key的最新的值
func (database *Database) rebuildOneIndex(instanceId uint64, fileIdstr string) error {

	key := database.genKey(instanceId)

	opt := opt.WriteOptions{
		Sync: false,
	}

	err := database.leveldb.Put([]byte(key), []byte(fileIdstr), &opt)
	if err != nil {
		log.Error("leveldb.Put fail, instanceId %d valuelen %d", instanceId, len(fileIdstr))
		return err
	}
	return nil

}

// TODO for what ????
func (database *Database) SetMinChoseninstanceId(writeOptions *WriteOptions, mininstanceId uint64) error {
	if !database.hasInit {
		log.Error("no init yet")
		return comm.ErrDbNotInit
	}

	var value = make([]byte, comm.UINT64SIZE)
	util.EncodeUint64(value, 0, mininstanceId)

	err := database.putToLevelDB(true, MINCHOSEN_KEY, value)
	if err != nil {
		return err
	}

	log.Info("ok, min chosen instanceId %d", mininstanceId)
	return nil
}

func (database *Database) GetMinChoseninstanceId() (uint64, error) {
	if !database.hasInit {
		log.Error("db not init yet")
		return comm.INVALID_INSTANCEID, comm.ErrDbNotInit
	}


	value, err := database.getFromLevelDb(MINCHOSEN_KEY)
	if err != nil && err != comm.ErrKeyNotFound {
		return comm.INVALID_INSTANCEID, err
	}

	if err == comm.ErrKeyNotFound {
		log.Error("no min chosen instanceId")
		return 0, nil
	}

	if len(value) != comm.UINT64SIZE {
		log.Error("fail, mininstanceId size wrong")
		return comm.INVALID_INSTANCEID, comm.ErrInvalidInstanceId
	}

	var mininstanceId uint64
	util.DecodeUint64(value, 0, &mininstanceId)
	log.Info("ok, min chosen instanceId:%d", mininstanceId)
	return mininstanceId, nil
}

func (database *Database) SetSystemVariables(writeOptions *WriteOptions, value string) error {
	return database.putToLevelDB(true, SYSTEMVARIABLES_KEY, []byte(value))
}

func (database *Database) GetSystemVariables() ([]byte, error) {
	return database.getFromLevelDb(SYSTEMVARIABLES_KEY)
}

func (database *Database) SetMasterVariables(writeOptions *WriteOptions, value string) error {
	return database.putToLevelDB(true, MASTERVARIABLES_KEY, []byte(value))
}

func (database *Database) GetMasterVariables() ([]byte, error) {
	return database.getFromLevelDb(MASTERVARIABLES_KEY)
}

// 从LevelDB中获取fileidStr
func (database *Database) getFromLevelDb(instanceId uint64) ([]byte, error) {
	key := database.genKey(instanceId)
	ret, err := database.leveldb.Get([]byte(key), nil)

	if err != nil {
		if err == leveldb.ErrNotFound {
			log.Debug("leveldb.get not found, instanceId %d", instanceId)
			return nil, comm.ErrKeyNotFound
		}

		log.Error("leveldb.get fail, instanceId %d", instanceId)
		return nil, err
	}


	return ret, nil
}

// 写入LevelDB
func (database *Database) putToLevelDB(sync bool, instanceId uint64, value []byte) error {
	key := database.genKey(instanceId)

	options := opt.WriteOptions{
		Sync: sync,
	}

	err := database.leveldb.Put([]byte(key), value, &options)
	if err != nil {
		log.Error("leveldb put fail, instanceId %d value len %d", instanceId, len(value))
		return err
	}

	return nil
}

// 从vfile读取
func (database *Database) fileIdToValue(fileId string, instanceId *uint64) ([]byte, error) {
	value, err := database.valueStore.Read(fileId, instanceId)
	if err != nil {
		log.Error("fieldIdToValue fail, ret %v", err)
		return nil, err
	}

	return value, nil
}

// 写入vfile
func (database *Database) valueToFileId(options *WriteOptions, instanceId uint64, value []byte, fileId *string) error {
	err := database.valueStore.Append(options, instanceId, value, fileId)
	if err != nil {
		log.Error("valueStore append fail:%v", err)
	}
	return err
}

func (database *Database) genKey(instanceId uint64) string {
	return fmt.Sprintf("%d", instanceId)
}

func (database *Database) getinstanceIdFromKey(key string) uint64 {
	instanceId, _ := strconv.ParseUint(key, 10, 64)
	return instanceId
}

//----------------------------------MultiDatabase 多个Database的封装，实现LogStorage接口----------------------------------//

type MultiDatabase struct {
	dbList []*Database
}

func (multiDatabase *MultiDatabase) Init(dbPath string, groupCount int) error {
	exists, err := util.Exists(dbPath)

	if err != nil {
		return fmt.Errorf("access dbpath error")
	}

	if !exists {
		err := os.MkdirAll(dbPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("create dir %s error: %v", dbPath, err)
		}
	}

	if groupCount < 1 || groupCount > 10000 {
		return fmt.Errorf("groupCount wrong %d", groupCount)
	}

	newDbPath := dbPath

	if dbPath[len(dbPath)-1] != os.PathSeparator {
		newDbPath = newDbPath + string(os.PathSeparator)
	}

	var waitGroup sync.WaitGroup

	for i := 0; i < groupCount; i++ {
		waitGroup.Add(1)
		go func(idx int) {
			defer waitGroup.Done()
			dbPath := fmt.Sprintf("%sg%d", newDbPath, i)
			db := &Database{}
			err = db.Init(dbPath, idx)
			if err == nil {
				multiDatabase.dbList = append(multiDatabase.dbList, db)
			}
		}(i)
	}

	waitGroup.Wait()

	if err != nil {
		return err
	}

	return nil
}

func (multiDatabase *MultiDatabase) GetLogStorageDirPath(groupIdx int) (string, error) {
	if groupIdx > len(multiDatabase.dbList) {
		return "", fmt.Errorf("groupIdx out of bround")
	}
	return multiDatabase.dbList[groupIdx].GetDBPath(), nil
}

func (multiDatabase *MultiDatabase) Get(groupIdx int, instanceId uint64) ([]byte, error) {
	if groupIdx > len(multiDatabase.dbList) {
		return nil, fmt.Errorf("groupIdx out of bround")
	}

	return multiDatabase.dbList[groupIdx].Get(instanceId)
}

func (multiDatabase *MultiDatabase) Put(writeOptions *WriteOptions, groupIdx int, instanceId uint64, value []byte) error {
	if groupIdx > len(multiDatabase.dbList) {
		return fmt.Errorf("groupIdx out of bround")
	}

	return multiDatabase.dbList[groupIdx].Put(writeOptions, instanceId, value)
}

func (multiDatabase *MultiDatabase) Del(writeOptions *WriteOptions, groupIdx int, instanceId uint64) error {
	if groupIdx > len(multiDatabase.dbList) {
		return fmt.Errorf("groupIdx out of bround")
	}

	return multiDatabase.dbList[groupIdx].Del(writeOptions, instanceId)
}

func (multiDatabase *MultiDatabase) GetMaxinstanceId(groupIdx int) (uint64, error) {
	if groupIdx > len(multiDatabase.dbList) {
		return -1, fmt.Errorf("groupIdx out of bround")
	}

	return multiDatabase.dbList[groupIdx].GetMaxinstanceId()
}

func (multiDatabase *MultiDatabase) SetMinChoseninstanceId(writeOptions *WriteOptions, groupIdx int, mininstanceId uint64) error {
	if groupIdx > len(multiDatabase.dbList) {
		return fmt.Errorf("groupIdx out of bround")
	}

	return multiDatabase.dbList[groupIdx].SetMinChoseninstanceId(writeOptions, mininstanceId)
}

func (multiDatabase *MultiDatabase) GetMinChoseninstanceId(groupIdx int)  (uint64, error) {
	if groupIdx > len(multiDatabase.dbList) {
		return -1, fmt.Errorf("groupIdx out of bround")
	}

	return multiDatabase.dbList[groupIdx].GetMinChoseninstanceId()
}

func (multiDatabase *MultiDatabase) ClearAllLog(groupIdx int) error {
	if groupIdx > len(multiDatabase.dbList) {
		return fmt.Errorf("groupIdx out of bround")
	}

	return multiDatabase.dbList[groupIdx].ClearAllLog()
}

func (multiDatabase *MultiDatabase) SetSystemVariables(writeOptions *WriteOptions, groupIdx int, value string) error {
	if groupIdx > len(multiDatabase.dbList) {
		return fmt.Errorf("groupIdx out of bround")
	}

	return multiDatabase.dbList[groupIdx].SetSystemVariables(writeOptions, value)

}

func (multiDatabase *MultiDatabase) GetSystemVariables(groupIdx int) ([]byte, error) {

	if groupIdx > len(multiDatabase.dbList) {
		return nil, fmt.Errorf("groupIdx out of bround")
	}

	return multiDatabase.dbList[groupIdx].GetSystemVariables()

}

func (multiDatabase *MultiDatabase) SetMasterVariables(writeOptions *WriteOptions, groupIdx int, value string) error {
	if groupIdx > len(multiDatabase.dbList) {
		return   fmt.Errorf("groupIdx out of bround")
	}

	return multiDatabase.dbList[groupIdx].SetMasterVariables(writeOptions, value)
}

func (multiDatabase *MultiDatabase) GetMasterVariables(groupIdx int) ([]byte, error) {
	if groupIdx > len(multiDatabase.dbList) {
		return   nil, fmt.Errorf("groupIdx out of bround")
	}

	return multiDatabase.dbList[groupIdx].GetMasterVariables()
}
