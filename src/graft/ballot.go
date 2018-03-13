package graft

type PosHint struct {
	pos0 int
	pos1 int
}

type UnfoundPeerId struct {
	peedId PeerId
	found  bool
}

// 投票
type Ballot struct {
	peers     []UnfoundPeerId
	quorum    int
	oldPeers  []UnfoundPeerId
	oldQuorum int
}

// TODO
func (ballot *Ballot) init(conf *Configuration, oldConf *Configuration) {

}

// TODO
func (ballot *Ballot) grant(peer *PeerId, hint PosHint) PosHint {
	return PosHint{}
}


func (ballot *Ballot) granted() bool {
	return ballot.quorum<=0 && ballot.oldQuorum<=0
}



type BallotBox struct {

}