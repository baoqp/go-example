package node

// 实现InsideSM接口
// TODO
type SystemVSM struct {

}


func(systemVM *SystemVSM) GetCheckpointBuffer() (string, error){
	return "", nil
}

func(systemVM *SystemVSM) UpdateByCheckpoint(systemVariables []byte) (bool, error)  {
	return true, nil
}