// Copyright 2016 DeepFabric, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package raftstore

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
	"multi-raft-kv/pb/metapb"
	"multi-raft-kv/util"
	"raft"
	"multi-raft-kv/pb/raftpb"
	"github.com/deepfabric/elasticell/pkg/pb/mraft"
	"errors"
)

const (
	// When we create a region peer, we should initialize its log term/index > 0,
	// so that we can force the follower peer to sync the snapshot first.
	raftInitLogTerm  = 5
	raftInitLogIndex = 5

	maxSnapTryCnt = 5
)

type snapshotState int

var (
	relax        = snapshotState(1)
	generating   = snapshotState(2)
	applying     = snapshotState(3)
	applyAborted = snapshotState(4)
)

const (
	pending    = iota
	running    
	cancelling 
	cancelled  
	finished   
	failed     
)

type peerStorage struct {
	store *Store
	cell  metapb.Cell

	lastTerm         uint64
	appliedIndexTerm uint64
	lastReadyIndex   uint64 // 上一次处理Ready结构体时的Indexz
	lastCompactIndex uint64
	raftState        raftpb.RaftLocalState
	applyState       raftpb.RaftApplyState

	snapTriedCnt     int
	genSnapJob       *util.Job
	applySnapJob     *util.Job
	applySnapJobLock sync.RWMutex

	//pendingReads *readIndexQueue
}

func newPeerStorage(store *Store, cell metapb.Cell) (*peerStorage, error) {
	s := new(peerStorage)
	s.store = store
	s.cell = cell
	s.appliedIndexTerm = raftInitLogTerm

	err := s.initRaftState() // 保存RaftState在handleRaftReadyAppend()方法中
	if err != nil {
		return nil, err
	}
	log.Infof("raftstore[cell-%d]: init raft state, state=<%+v>",
		cell.Id,
		s.raftState)

	err = s.initApplyState()
	if err != nil {
		return nil, err
	}
	log.Infof("raftstore[cell-%d]: init apply state, state=<%+v>",
		cell.Id,
		s.applyState)

	err = s.initLastTerm()
	if err != nil {
		return nil, err
	}
	log.Infof("raftstore[cell-%d]: init last term, last term=<%d>",
		cell.Id,
		s.lastTerm)

	s.lastReadyIndex = s.getAppliedIndex()
	//s.pendingReads = new(readIndexQueue)

	return s, nil
}

// 如果底层存储已经保存了之前的状态则用其初始化，否则初始化为最初状态
// TODO 为什么不同快照和wal恢复 ???
func (ps *peerStorage) initRaftState() error {
	v, err := ps.store.getMetaEngine().Get(getRaftStateKey(ps.getCell().Id))
	if err != nil {
		return err
	}

	if len(v) > 0 {
		s := &mraft.RaftLocalState{}
		err = s.Unmarshal(v)
		if err != nil {
			return err
		}

		ps.raftState = *s
		return nil
	}

	s := &mraft.RaftLocalState{}
	if len(ps.getCell().Peers) > 0 {
		s.LastIndex = raftInitLogIndex
	}

	ps.raftState = *s
	return nil
}

func (ps *peerStorage) initApplyState() error {
	v, err := ps.store.getMetaEngine().Get(getApplyStateKey(ps.getCell().Id))
	if err != nil {
		return errors.Wrap(err, "")
	}

	if len(v) > 0 && len(ps.getCell().Peers) > 0 {
		s := &mraft.RaftApplyState{}
		err = s.Unmarshal(v)
		if err != nil {
			return errors.Wrap(err, "")
		}

		ps.applyState = *s
		return nil
	}

	if len(ps.getCell().Peers) > 0 {
		ps.applyState.AppliedIndex = raftInitLogIndex
		ps.applyState.TruncatedState.Index = raftInitLogIndex
		ps.applyState.TruncatedState.Term = raftInitLogTerm
	}

	return nil
}

func (ps *peerStorage) initLastTerm() error {
	lastIndex := ps.raftState.LastIndex

	if lastIndex == 0 {
		ps.lastTerm = lastIndex
		return nil
	} else if lastIndex == raftInitLogIndex {
		ps.lastTerm = raftInitLogTerm
		return nil
	} else if lastIndex == ps.applyState.TruncatedState.Index {
		ps.lastTerm = ps.applyState.TruncatedState.Term
		return nil
	} else if lastIndex < raftInitLogIndex {
		log.Fatalf("raftstore[cell-%d]: error raft last index, index=<%d>",
			ps.getCell().Id,
			lastIndex)
		return nil
	}

	v, err := ps.store.getMetaEngine().Get(getRaftLogKey(ps.getCell().Id, lastIndex))
	if err != nil {
		return errors.Wrap(err, "")
	}

	if nil == v {
		return fmt.Errorf("raftstore[cell-%d]: entry at index<%d> doesn't exist, may lose data",
			ps.getCell().Id,
			lastIndex)
	}

	s := &raftpb.Entry{}
	err = s.Unmarshal(v)
	if err != nil {
		return errors.Wrap(err, "")
	}

	ps.lastTerm = s.Term
	return nil
}

func (ps *peerStorage) isApplyComplete() bool {
	return ps.getCommittedIndex() == ps.getAppliedIndex()
}

func (ps *peerStorage) setApplyState(applyState *mraft.RaftApplyState) {
	ps.applyState = *applyState
}

func (ps *peerStorage) getApplyState() *mraft.RaftApplyState {
	return &ps.applyState
}

func (ps *peerStorage) getAppliedIndex() uint64 {
	return ps.getApplyState().AppliedIndex
}

func (ps *peerStorage) getCommittedIndex() uint64 {
	return ps.raftState.HardState.Commit
}

// TruncatedIndex之前的数据在快照的时候被截断，已经不存在
func (ps *peerStorage) getTruncatedIndex() uint64 {
	return ps.getApplyState().TruncatedState.Index
}

func (ps *peerStorage) getTruncatedTerm() uint64 {
	return ps.getApplyState().TruncatedState.Term
}

func (ps *peerStorage) getAppliedIndexTerm() uint64 {
	return atomic.LoadUint64(&ps.appliedIndexTerm)
}

func (ps *peerStorage) setAppliedIndexTerm(appliedIndexTerm uint64) {
	atomic.StoreUint64(&ps.appliedIndexTerm, appliedIndexTerm)
}

func (ps *peerStorage) valIdateSnap(snap *raftpb.Snapshot) bool {
	Idx := snap.Metadata.Index

	if Idx < ps.getTruncatedIndex() {
		// stale snapshot, should generate again.
		log.Infof("raftstore[cell-%d]: snapshot is stale, generate again, snapIndex=<%d> currIndex=<%d>",
			ps.getCell().Id,
			Idx,
			ps.getTruncatedIndex())
		return false
	}

	snapData := &mraft.SnapshotMessage{}
	err := snapData.Unmarshal(snap.Data)
	if err != nil {
		log.Errorf("raftstore[cell-%d]: decode snapshot fail, errors:\n %+v",
			ps.getCell().Id,
			err)
		return false
	}

	snapEpoch := snapData.Header.Cell.Epoch
	lastEpoch := ps.getCell().Epoch

	if snapEpoch.ConfVer < lastEpoch.ConfVer {
		log.Infof("raftstore[cell-%d]: snapshot epoch stale, generate again. snap=<%s> curr=<%s>",
			ps.getCell().Id,
			snapEpoch.String(),
			lastEpoch.String())
		return false
	}

	return true
}

func (ps *peerStorage) isInitialized() bool {
	return len(ps.getCell().Peers) != 0
}

func (ps *peerStorage) isApplyingSnapshot() bool {
	return ps.applySnapJob != nil && ps.applySnapJob.IsNotComplete()
}

func (ps *peerStorage) getCell() metapb.Cell {
	return ps.cell
}

func (ps *peerStorage) setCell(cell metapb.Cell) {
	ps.cell = cell
}

func (ps *peerStorage) checkRange(low, high uint64) error {
	if low > high {
		return fmt.Errorf("raftstore[cell-%d]: low is greater that high, low=<%d> high=<%d>",
			ps.getCell().Id,
			low,
			high)
	} else if low <= ps.getTruncatedIndex() {
		return raft.ErrCompacted
	} else {
		i, err := ps.LastIndex()
		if err != nil {
			return err
		}

		if high > i+1 {
			return fmt.Errorf("raftstore[cell-%d]: entries' high is out of bound lastindex, hight=<%d> lastindex=<%d>",
				ps.getCell().Id,
				high,
				i)
		}
	}

	return nil
}

func (ps *peerStorage) loadLogEntry(index uint64) (*raftpb.Entry, error) {
	key := getRaftLogKey(ps.getCell().Id, index)
	v, err := ps.store.getMetaEngine().Get(key)
	if err != nil {
		log.Errorf("raftstore[cell-%d]: load entry failure, index=<%d> errors:\n %+v",
			ps.getCell().Id,
			index,
			err)
		return nil, err
	} else if len(v) == 0 {
		log.Errorf("raftstore[cell-%d]: entry not found, index=<%d>",
			ps.getCell().Id,
			index)
		return nil, fmt.Errorf("log entry at <%d> not found", index)
	}

	return ps.unmarshal(v, index)
}

func (ps *peerStorage) loadCellLocalState(job *util.Job) (*mraft.CellLocalState, error) {
	if nil != job &&
		job.IsCancelling() {
		return nil, util.ErrJobCancelled
	}

	return loadCellLocalState(ps.getCell().Id, ps.store.engine, false)
}

func (ps *peerStorage) applySnapshot(job *util.Job) error {
	if nil != job &&
		job.IsCancelling() {
		return util.ErrJobCancelled
	}

	snap := &mraft.SnapshotMessage{}
	snap.Header = mraft.SnapshotMessageHeader{
		Cell:  ps.getCell(),
		Term:  ps.applyState.TruncatedState.Term,
		Index: ps.applyState.TruncatedState.Index,
	}

	return ps.store.snapshotManager.Apply(snap)
}

func (ps *peerStorage) loadApplyState() (*mraft.RaftApplyState, error) {
	key := getApplyStateKey(ps.getCell().Id)
	v, err := ps.store.getMetaEngine().Get(key)
	if err != nil {
		log.Errorf("raftstore[cell-%d]: load apply state failed, errors:\n %+v",
			ps.getCell().Id,
			err)
		return nil, err
	}

	if nil == v {
		return nil, errors.New("cell apply state not found")
	}

	applyState := &mraft.RaftApplyState{}
	err = applyState.Unmarshal(v)
	return applyState, err
}

func (ps *peerStorage) unmarshal(v []byte, expectIndex uint64) (*raftpb.Entry, error) {
	e := acquireEntry()
	if err := e.Unmarshal(v); err != nil {
		log.Errorf("raftstore[cell-%d]: unmarshal entry failure, index=<%d>, v=<%+v> errors:\n %+v",
			ps.getCell().Id,
			expectIndex,
			v,
			err)
		releaseEntry(e)
		return nil, err
	}

	if e.Index != expectIndex {
		log.Fatalf("raftstore[cell-%d]: raft log index not match, logIndex=<%d> expect=<%d>",
			ps.getCell().Id,
			e.Index,
			expectIndex)
	}

	return e, nil
}

/// Delete all data belong to the region.
/// If return Err, data may get partial deleted.
func (ps *peerStorage) clearData() error {
	cell := ps.getCell()

	cellId := cell.Id
	startKey := encStartKey(&cell)
	endKey := encEndKey(&cell)

	err := ps.store.addSnapJob(func() error {
		log.Infof("raftstore-destroy[cell-%d]: deleting data, start=<%v> end=<%v>",
			cellId, startKey, endKey)
		err := ps.deleteAllInRange(startKey, endKey, nil)
		if err != nil {
			log.Errorf("raftstore-destroy[cell-%d]: failed to delete data, start=<%v> end=<%v> errors:\n %+v",
				cellId, startKey, endKey, err)
		}

		return err
	}, nil)

	return err
}

// Delete all data that is not covered by `newCell`.
func (ps *peerStorage) clearExtraData(newCell metapb.Cell) error {
	cell := ps.getCell()

	oldStartKey := encStartKey(&cell)
	oldEndKey := encEndKey(&cell)

	newStartKey := encStartKey(&newCell)
	newEndKey := encEndKey(&newCell)

	if bytes.Compare(oldStartKey, newStartKey) < 0 {
		err := ps.startDestroyDataJob(newCell.Id, oldStartKey, newStartKey)

		if err != nil {
			return err
		}
	}

	if bytes.Compare(newEndKey, oldEndKey) < 0 {
		err := ps.startDestroyDataJob(newCell.Id, newEndKey, oldEndKey)

		if err != nil {
			return err
		}
	}

	return nil
}

func (ps *peerStorage) updatePeerState(cell metapb.Cell, state mraft.PeerState, wb storage.WriteBatch) error {
	cellState := &mraft.CellLocalState{}
	cellState.State = state
	cellState.Cell = cell

	data, _ := cellState.Marshal()

	if wb != nil {
		return wb.Set(getCellStateKey(cell.Id), data)
	}

	return ps.store.getMetaEngine().Set(getCellStateKey(cell.Id), data)
}

func (ps *peerStorage) writeInitialState(cellId uint64, wb storage.WriteBatch) error {
	raftState := new(mraft.RaftLocalState)
	raftState.LastIndex = raftInitLogIndex
	raftState.HardState.Term = raftInitLogTerm
	raftState.HardState.Commit = raftInitLogIndex

	applyState := new(mraft.RaftApplyState)
	applyState.AppliedIndex = raftInitLogIndex
	applyState.TruncatedState.Index = raftInitLogIndex
	applyState.TruncatedState.Term = raftInitLogTerm

	err := wb.Set(getRaftStateKey(cellId), util.MustMarshal(raftState))
	if err != nil {
		return err
	}

	return wb.Set(getApplyStateKey(cellId), util.MustMarshal(applyState))
}

func (ps *peerStorage) deleteAllInRange(start, end []byte, job *util.Job) error {
	if job != nil &&
		job.IsCancelling() {
		return util.ErrJobCancelled
	}

	return ps.store.getDataEngine().RangeDelete(start, end)
}

func compactRaftLog(cellId uint64, state *mraft.RaftApplyState, compactIndex, compactTerm uint64) error {
	log.Debugf("raftstore-compact[cell-%d]: compact log entries to index, index=<%d>",
		cellId,
		compactIndex)
	if compactIndex <= state.TruncatedState.Index {
		return errors.New("try to truncate compacted entries")
	} else if compactIndex > state.AppliedIndex {
		return fmt.Errorf("compact index %d > applied index %d", compactIndex, state.AppliedIndex)
	}

	// we don't actually delete the logs now, we add an async task to do it.
	state.TruncatedState.Index = compactIndex
	state.TruncatedState.Term = compactTerm

	return nil
}

func loadCellLocalState(cellId uint64, driver storage.Driver, allowNotFound bool) (*mraft.CellLocalState, error) {
	key := getCellStateKey(cellId)
	v, err := driver.GetEngine().Get(key)
	if err != nil {
		log.Errorf("raftstore[cell-%d]: load raft state failed, errors:\n %+v",
			cellId,
			err)
		return nil, err
	} else if len(v) == 0 {
		if allowNotFound {
			return nil, nil
		}

		return nil, errors.New("cell state not found")
	}

	stat := &mraft.CellLocalState{}
	err = stat.Unmarshal(v)

	return stat, err
}
