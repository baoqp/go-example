package bplus_tree

import (
	"os"

	"fmt"
	"util"
)

type Writer struct {
	file     *os.File
	fileName string
	fileSize uint64
	padding  [PADDING]byte // TODO
}

func writerCreate(writer *Writer, fileName string) error {

	writer.fileName = fileName

	var err error
	writer.file, err = os.OpenFile(fileName,
		os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return EFILE
	}

	fileSize, err := writer.file.Seek(0, os.SEEK_END)
	if err != nil {
		return err
	}
	writer.fileSize = uint64(fileSize)
	for i := range writer.padding {
		writer.padding[i] = 0
	}

	return nil
}

func writerDestroy(w *Writer) {
	w.file.Close()
}

// linux使用 syscall.Fdatasync(int(db.file.Fd()))
func writerFsync(w *Writer) {
	w.file.Sync()
}

func writerCompactName(w *Writer) (string, error) {
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

func writerRead(w *Writer, compType CompType, offset uint64,
	size *uint64) ([]byte, error) {

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

	} else { // 如果写入时使用了压缩需要解压缩 TODO

	}

	return data, nil
}

func writerWrite(w *Writer, compType CompType, data []byte,
	offset *uint64, size *uint64) error {

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

type WriterCallback func(w *Writer, data []byte) error

func writerFind(w *Writer, compType CompType, size uint64, data []byte,
	seek WriterCallback, miss WriterCallback) error {

	match := false
	// Write padding first
	err := writerWrite(w, kNotCompressed, nil, nil, nil)
	if err != nil {
		return err
	}

	offset := w.fileSize
	sizeTmp := offset

	//  Start seeking from bottom of file
	for ; offset >= size; {
		data, err := writerRead(w, compType, offset-size, &sizeTmp)
		if err != nil {
			return err
		}
		// Break if matched
		if seek(w, data) == nil {
			match = true
			break
		}

		offset -= size
	}

	if !match {
		return miss(w, data)
	}

	return nil

}
