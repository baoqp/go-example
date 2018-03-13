package graft

import (
	"os"
	"io/ioutil"
	"github.com/golang/protobuf/proto"
)

type LocalRaftMetaStorage struct {
	isInited bool
	path     string
	term     int64
	votedFor PeerId
}

func (metaStorage *LocalRaftMetaStorage) setTerm(term int64) {
	if metaStorage.isInited {
		metaStorage.term = term
		metaStorage.save()
	} else {
		panic("setTerm error, metaStore not inited ")
	}

}

func (metaStorage *LocalRaftMetaStorage) getTerm() int64 {
	if metaStorage.isInited {
		return metaStorage.term
	} else {
		panic("setTerm error, metaStore not inited ")
	}

}


func (metaStorage *LocalRaftMetaStorage) setVotedFor(peerId PeerId) {
	if metaStorage.isInited {
		metaStorage.votedFor = peerId
		metaStorage.save()
	} else {
		panic("setTerm error, metaStore not inited ")
	}

}


func (metaStorage *LocalRaftMetaStorage) getVotedFor() PeerId {
	if metaStorage.isInited {
		return metaStorage.votedFor
	} else {
		panic("setTerm error, metaStore not inited ")
	}

}


func (metaStorage *LocalRaftMetaStorage) setTermAndVotedFor(term int64, peerId PeerId) {
	if metaStorage.isInited {
		metaStorage.term = term
		metaStorage.votedFor = peerId
		metaStorage.save()
	} else {
		panic("setTerm error, metaStore not inited ")
	}

}

func (metaStorage *LocalRaftMetaStorage) init() bool {

	if metaStorage.isInited {
		return false
	}

	pathExists, _ := PathExists(metaStorage.path)

	if !pathExists {
		err := os.MkdirAll(metaStorage.path, 0700)

		if err != nil {
			return false
		}
	}

	metaStorage.load()
	metaStorage.init() = true

	return true
}

func (metaStorage *LocalRaftMetaStorage) load() {

	metaPath := metaStorage.path + string(os.PathSeparator) + sRaftMeta

	if metaExists, _ := PathExists(metaPath); metaExists {
		data, err := ioutil.ReadFile(metaPath)
		if err != nil {
			panic("raed meta file err")
		}

		meta := &StablePBMeta{}
		err = proto.Unmarshal(data, meta)

		if err != nil {
			panic("load reft meta file error")
		} else {
			metaStorage.term = meta.Term
			metaStorage.votedFor.parse(meta.Votedfor)
		}
	}

}


func (metaStorage *LocalRaftMetaStorage) save() {
	meta := &StablePBMeta{Term:metaStorage.term, Votedfor:metaStorage.votedFor.toString()}
	metaPath := metaStorage.path + string(os.PathSeparator) + sRaftMeta

	data, err := proto.Marshal(meta)

	if err != nil {
		panic("marshal meta error")
	}

	ioutil.WriteFile(metaPath, data, 0700)

}