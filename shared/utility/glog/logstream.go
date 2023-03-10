package glog

import (
	"bytes"
	"os"
	"path"
)

type LogFileStream struct {
	dir        string
	namePrefix string
	file       *os.File
	spliter    LogFileSpliter
}

func NewLogFileStream(dir, namePrefix string, spliter LogFileSpliter) *LogFileStream {
	return &LogFileStream{
		dir:        dir,
		namePrefix: namePrefix,
		spliter:    spliter,
	}
}

func (lfs *LogFileStream) Write(bytes []byte) (int, error) {
	err := lfs.ensureFile()
	if err != nil {
		return 0, err
	}

	writeByte, err := lfs.file.Write(bytes)
	lfs.spliter.OnWrite(writeByte)

	return writeByte, err
}

func (lfs *LogFileStream) Sync() error {
	if lfs.file != nil {
		return lfs.file.Sync()
	}

	return os.ErrInvalid
}

func (lfs *LogFileStream) ensureFile() error {
	if lfs.file == nil {
		err := lfs.init()
		if err != nil {
			return err
		}
	} else {
		if lfs.spliter.TrySplit() {
			err := lfs.loadFile()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (lfs *LogFileStream) loadFile() error {
	dirBuf := bytes.NewBufferString(lfs.dir)
	nameBuf := bytes.NewBufferString(lfs.namePrefix)

	lfs.spliter.AppendDirAndName(dirBuf, nameBuf)

	if nameBuf.Len() <= 0 {
		nameBuf.WriteString("log")
	} else {
		nameBuf.WriteString(".log")
	}

	return lfs.openFile(dirBuf.String(), nameBuf.String())
}

func (lfs *LogFileStream) openFile(dir, filename string) error {
	os.MkdirAll(dir, os.ModeDir|os.ModePerm)

	filePath := path.Join(dir, filename)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModeAppend|os.ModePerm)
	if err != nil {
		return err
	}

	if lfs.file != nil {
		lfs.file.Close()
	}

	lfs.file = file

	return nil
}

func (lfs *LogFileStream) init() error {
	dirBuf := bytes.NewBufferString(lfs.dir)
	nameBuf := bytes.NewBufferString(lfs.namePrefix)
	lfs.spliter.MoveLast(dirBuf, nameBuf)

	for {
		err := lfs.loadFile()
		if err != nil {
			return err
		}

		fileInfo, err := lfs.file.Stat()
		if err != nil {
			return err
		}

		lfs.spliter.OnWrite(int(fileInfo.Size()))

		if !lfs.spliter.TrySplit() {
			break
		}
	}

	return nil
}
