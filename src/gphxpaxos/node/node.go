package node

import "gphxpaxos/comm"

type Node struct {
}

//Base function.
func (node *Node) Propose(groupIdx int, value []byte, instanceId uint64) error { return nil }

func (node *Node) ProposeWithCtx(groupIdx int, value []byte, instanceId uint64, smCtx *SMCtx) error { return nil }

func (node *Node) GetNowInstanceID(groupIdx int) error { return nil }

func (node *Node) GetMinChosenInstanceID(groupIdx int) error { return nil }

func (node *Node) GetMyNodeID() uint64 { return 0 }

//Batch propose.

//Only set options::bUserBatchPropose as true can use this batch API.
//Warning: BatchProposal will have same llInstanceID returned but different iBatchIndex.
//Batch values's execute order in StateMachine is certain, the return value iBatchIndex
//means the execute order index, start from 0.
func (node *Node) BatchPropose(groupIdx int, value []byte, instanceId uint64, batchIndex uint32) error { return nil }

func (node *Node) BatchProposeWithCtx(groupIdx int, value []byte, instanceId uint64, batchIndex uint32, smCtx *SMCtx) error {
	return nil
}

//PhxPaxos will batch proposal while waiting proposals count reach to BatchCount, 
//or wait time reach to BatchDelayTimeMs.
func (node *Node) SetBatchCount(groupIdx int, batchCount int) error { return nil }

func (node *Node) SetBatchDelayTimeMs(groupIdx int, iBatchDelayTimeMs int) error { return nil }

//State machine.

//This function will add state machine to all group.
func (node *Node) AddStateMachineToAllGroup(sm *StateMachine) error { return nil }

func (node *Node) AddStateMachine(groupIdx int, sm *StateMachine) error { return nil }

//Timeout control.
func (node *Node) SetTimeoutMs(timeoutMs int) error {
	return nil
}

//Checkpoint

//Set the number you want to keep paxoslog's count.
//We will only delete paxoslog before checkpoint instanceid.
//If llHoldCount < 300, we will set it to 300. Not suggest too small holdcount.
func (node *Node) SetHoldPaxosLogCount(llHoldCount uint64) error { return nil }

//Replayer is to help sm make checkpoint.
//Checkpoint replayer default is paused, if you not use this, ignord this function.
//If sm use ExecuteForCheckpoint to make checkpoint, you need to run replayer(you can run in any time).

//Pause checkpoint replayer.
func (node *Node) PauseCheckpointReplayer() error { return nil }

//Continue to run replayer
func (node *Node) ContinueCheckpointReplayer() error { return nil }

//Paxos log cleaner working for deleting paxoslog before checkpoint instanceid.
//Paxos log cleaner default is pausing.

//pause paxos log cleaner.
func (node *Node) PausePaxosLogCleaner() error { return nil }

//Continue to run paxos log cleaner.
func (node *Node) ContinuePaxosLogCleaner() error { return nil }

//Membership

//Show now membership.
func (node *Node) ShowMembership(groupIdx int, ndeInfoList *comm.NodeInfoList) error { return nil }

//Add a paxos node to membership.
func (node *Node) AddMember(groupIdx int, ndeInfoList *comm.NodeInfoList) error { return nil }

//Remove a paxos node from membership.
func (node *Node) RemoveMember(groupIdx int, ndeInfoList *comm.NodeInfoList) error { return nil }

//Change membership by one node to another node.
func (node *Node) ChangeMember(groupIdx int, fromNode *comm.NodeInfo, toNode *comm.NodeInfo) error { return nil }

//Master

//Check who is master.
func (node *Node) GetMaster(groupIdx int) (*comm.NodeInfo, error) { return nil, nil }

//Check who is master and get version.
func (node *Node) GetMasterWithVersion(groupIdx int, version uint64) (*comm.NodeInfo, error) { return nil, nil }

//Check is i'm master.
func (node *Node) IsIMMaster(groupIdx int) (bool, error) { return false, nil }

func (node *Node) SetMasterLease(groupIdx int, leaseTimeMs int) error {
	return nil
}

func (node *Node) DropMaster(groupIdx int) error { return nil }

//Qos

//If many threads propose same group, that some threads will be on waiting status.
//Set max hold threads, and we will reject some propose request to avoid to many threads be holded.
//Reject propose request will get retcode(PaxosTryCommitRet_TooManyThreadWaiting_Reject), check on def.h.
func (node *Node) SetMaxHoldThreads(groupIdx int, iMaxHoldThreads int) error { return nil }

//To avoid threads be holded too long time, we use this threshold to reject some propose to control thread's wait time.
func (node *Node) SetProposeWaitTimeThresholdMS(groupIdx int, iWaitTimeThresholdMS int) error { return nil }

//write disk
func (node *Node) SetLogSync(groupIdx int, logSync bool) error { return nil }

//Not suggest to use this function
//pair: value,smid.
//Because of BatchPropose, a InstanceID maybe include multi-value.
func (node *Node) GetInstanceValue(groupIdx int, instanceId uint64, vecValues *[]map[string]int) error { return nil }
