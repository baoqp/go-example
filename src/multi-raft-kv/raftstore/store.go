package raftstore

import (

	"sync"

	"regexp"

	"multi-raft-kv/pb/metapb"
	"multi-raft-kv/util"

	"multi-raft-kv/storage"
)

type Store struct {

	id                 uint64
	clusterID          uint64
	startAt            uint32
	meta               metapb.Store


	/*snapshotManager    SnapshotManager
	pdClient           *pd.Client
	keyConvertFun      func([]byte, func([]byte) metapb.Cell) metapb.Cell
	replicatesMap      *cellPeersMap // cellid -> peer replicate  一个store上有多个raft group的副本
	keyRanges          *util.CellTree
	peerCache          *peerCacheMap
	delegates          *applyDelegateMap // TODO ???
	pendingLock        sync.RWMutex
	pendingSnapshots   map[uint64]mraft.SnapshotMessageHeader*/


	trans              *transport
	engine             storage.Driver // 存储引擎
	runner             *util.Runner   // 执行具体的任务

	/*sendingSnapCount   uint32
	reveivingSnapCount uint32
	rwlock             sync.RWMutex
	indices            map[string]*pdpb.IndexDef   // index name -> IndexDef TODO ???
	reExps             map[string]*regexp.Regexp   // index name -> Regexp   TODO ???
	docProts           map[string]*cql.Document    // index name -> cql.Document  TODO ???
	indexers           map[uint64]*indexer.Indexer // cell id -> Indexer TODO ???
	cellIDToStores     map[uint64][]uint64         // cell id -> non-leader store ids, fetched from PD
	storeIDToCells     map[uint64][]uint64         // store id -> leader cell ids, fetched from PD TODO ???
	syncEpoch          uint64
	queryStates        map[string]*QueryState // query UUID -> query state
	queryReqChan       chan *QueryRequestCb
	queryRspChan       chan *querypb.QueryRsp

	droppedLock     sync.Mutex
	droppedVoteMsgs map[uint64]raftpb.Message*/
}


