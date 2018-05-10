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

package storage

import "github.com/deepfabric/elasticell/pkg/pb/raftcmdpb"

// WriteBatch batch operation
type WriteBatch interface {
	Delete(key []byte) error
	Set(key []byte, value []byte) error
}

// Driver is def storage interface
type Driver interface {
	GetEngine() Engine
	GetDataEngine() DataEngine
	GetKVEngine() KVEngine
	NewWriteBatch() WriteBatch
	Write(wb WriteBatch, sync bool) error
}

// KVEngine is the storage of KV
type KVEngine interface {
	Set(key, value []byte) error
	Get(key []byte) ([]byte, error)
	NewWriteBatch() WriteBatch
	Write(wb WriteBatch) error
}



// Seekable support seek
type Seekable interface {
	Seek(key []byte) ([]byte, []byte, error)
}

// Scanable support scan
type Scanable interface {
	// Scan scans the range and execute the handler fun.
	// returns false means end the scan.
	Scan(start, end []byte, handler func(key, value []byte) (bool, error), pooledKey bool) error
	// Free free the pooled bytes
	Free(pooled []byte)
}

// RangeDeleteable support range delete
type RangeDeleteable interface {
	RangeDelete(start, end []byte) error
}

// DataEngine is the storage of redis data
type DataEngine interface {
	RangeDeleteable
	// GetTargetSizeKey Find a key in the range [startKey, endKey) that sum size over target
	// if found returns the key
	GetTargetSizeKey(startKey []byte, endKey []byte, size uint64) (uint64, []byte, error)
	// CreateSnapshot create a snapshot file under the giving path
	CreateSnapshot(path string, start, end []byte) error
	// ApplySnapshot apply a snapshort file from giving path
	ApplySnapshot(path string) error

	// ScanIndexInfo scans the range and execute the handler fun. Returens a tuple (error count, first error)
	ScanIndexInfo(startKey []byte, endKey []byte, skipEmpty bool, handler func(key, idxInfo []byte) error) (int, error)
	SetIndexInfo(key, idxInfo []byte) error
	GetIndexInfo(key []byte) (idxInfo []byte, err error)
}

// Engine is the storage of meta data
type Engine interface {
	Scanable
	Seekable
	RangeDeleteable

	Set(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	Delete(key []byte) error
}
