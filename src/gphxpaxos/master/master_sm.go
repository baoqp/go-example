package master


// 实现InsideSM接口
type MasterStateMachine struct{

}


// TODO
func(masterStateMachine *MasterStateMachine) GetCheckpointBuffer() (string, error){
	return "", nil
}


func(masterStateMachine *MasterStateMachine) UpdateByCheckpoint(systemVariables []byte) (bool, error)  {
return true, nil
}
