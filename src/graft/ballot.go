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
	peers    []UnfoundPeerId
	oldPeers []UnfoundPeerId
}
