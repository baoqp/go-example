package bplus_tree_aop

import (
	"os"

	"fmt"
	"util"
	"github.com/pkg/errors"
)

type Writer struct {
	file         *os.File
	originalName string
	fileName     string
	compactName  string
	fileSize     uint64
	padding      [PADDING]byte // TODO
}

func (t *Tree) creatreWriter(fileName string, isCompact bool) (*Writer, error) {
	var err error
	fileName1 := fmt.Sprintf("%s.%d", fileName, 1)
	fileName2 := fmt.Sprintf("%s.%d", fileName, 2)

	file1Exists, err := util.Exists(fileName1)
	if err != nil {
		return nil, err
	}
	file2Exists, err := util.Exists(fileName2)
	if err != nil {
		return nil, err
	}

	writer := new(Writer)
	writer.originalName = fileName
	if !file1Exists && !file2Exists {
		writer.fileName = fileName1
		writer.compactName = fileName2
	} else if file2Exists && file1Exists {
		return nil, errors.New("db file and compact file is not allowd exists both")
	} else if file1Exists {
		writer.fileName = fileName1
		writer.compactName = fileName2
		if isCompact {
			writer.fileName = fileName2
			writer.compactName = fileName1
		}
	} else {
		writer.fileName = fileName2
		writer.compactName = fileName1
		if isCompact {
			writer.fileName = fileName1
			writer.compactName = fileName2
		}
	}

	writer.file, err = os.OpenFile(writer.fileName,
		os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return nil, EFILE
	}

	fileSize, err := writer.file.Seek(0, os.SEEK_END)
	if err != nil {
		return nil, err
	}
	writer.fileSize = uint64(fileSize)
	for i := range writer.padding {
		writer.padding[i] = 0
	}

	return writer, nil
}

func (t *Tree) destroyWriter() error {
	return t.Writer.file.Close()
}

func (t *Tree) deleteTreeFile() error {
	return os.Remove(t.Writer.fileName)
}

// linux使用 syscall.Fdatasync(int(Tree.file.Fd()))
func (t *Tree) writerFsync() {
	t.Writer.file.Sync()
}

func (t *Tree) writerCompactName() (string, error) {
	w := t.Writer
	compactName := fmt.Sprintf("%s.compact", w.fileName)

	if exits, err := util.Exists(compactName); err != nil {
		return "", err
	} else {
		if exits {
			return "", ECOMPACT_EXISTS
		}
	}

	return compactName, nil
}

// TODO
func writerCompactFinalize(s *Writer, t *Writer) {

}

func (t *Tree) writerRead(compType CompType, offset uint64,
	size *uint64) ([]byte, error) {

	w := t.Writer

	if w.fileSize < offset + *size {
		return nil, EFILEREAD_OOB
	}

	if *size == 0 {
		return nil, nil
	}

	data := make([]byte, *size, *size)

	bytesRead, err := w.file.ReadAt(data, int64(offset))
	if err != nil || bytesRead != int(*size) {
		return nil, EFILEREAD
	}

	if compType == kNotCompressed { // 没有使用压缩
		return data, nil

	} else { // 如果写入时使用了压缩需要解压缩 TODO
		return nil, errors.New("not support yet")

	}

}

// 写入文件末尾，并返回在文件中的offset
func (t *Tree) writerWrite(compType CompType, data []byte,
	offset *uint64, size *uint64) error {

	w := t.Writer
	paddding := uint64(len(w.padding)) - (w.fileSize % uint64(len(w.padding)))

	// writer padding
	if paddding != uint64(len(w.padding)) {
		written, err := w.file.Write(w.padding[:paddding])
		if err != nil || uint64(written) != paddding {
			return EFILEWRITE
		} else {
			w.fileSize += paddding
		}
	}

	if size == nil || *size == 0 {
		if offset != nil {
			*offset = w.fileSize
		}
		return nil
	}

	var written int
	var err error
	// head shouldn't be compressed
	if compType == kNotCompressed {
		written, err = w.file.Write(data)
	} else { // TODO

	}

	if err != nil || uint64(written) != *size {
		return EFILEWRITE
	}

	*offset = w.fileSize
	w.fileSize += *size

	return nil
}

type WriterCallback func(t *Tree, data []byte) error

func (t *Tree) writerFind(compType CompType, size uint64, data []byte,
	seek WriterCallback, miss WriterCallback) error {
	w := t.Writer
	match := false
	// Write padding first
	err := t.writerWrite(kNotCompressed, nil, nil, nil)
	if err != nil {
		return err
	}

	offset := w.fileSize
	sizeTmp := size

	//  Start seeking from bottom of file
	for ; offset >= size; {
		data, err := t.writerRead(compType, offset-size, &sizeTmp)
		if err != nil {
			return err
		}
		// Break if matched
		if seek(t, data) == nil {
			match = true
			break
		}

		offset -= size
	}

	if !match {
		return miss(t, data)
	}

	return nil

}
