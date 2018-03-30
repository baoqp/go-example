package algorithm

import "encoding/binary"

var GROUPIDXLEN = binary.Size(int32(0))
var HEADLEN_LEN = binary.Size(uint16(0))
var CHECKSUM_LEN = binary.Size(uint32(0))

const (
	BroadcastMessage_Type_RunSelf_First = 1
	BroadcastMessage_Type_RunSelf_Final = 2
	BroadcastMessage_Type_RunSelf_None  = 3
)


const (
	PrepareTimer = iota
	AcceptTimer
	LearnerTimer
)

// instance id status
type Status int
const (
	Decided   = iota + 1
	Pending   // not yet decided.
)



