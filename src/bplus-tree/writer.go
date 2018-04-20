package bplus_tree

import (
	"os"

)

type Writer struct {
	file *os.File
	fileName string
	fileSize uint64
	padding   [PADDING]byte // TODO
}


func writeCreate(fileName string) (*Writer, error) {

	writer := new(Writer)

	writer.fileName = fileName

	var err error
	writer.file, err = os.OpenFile(fileName, os.O_RDWR | os.O_APPEND | os.O_CREATE,0666) // TODO file op permits
	if err != nil {
		return nil, EFILE
	}

	// TODO




	return writer, nil
}