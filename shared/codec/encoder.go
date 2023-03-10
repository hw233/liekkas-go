package codec

import (
	"bytes"
	"encoding/binary"
	"log"
	"time"
)

type Encoder struct {
	Head      uint16
	Version   uint8
	Length    int32
	Command   uint16
	Serial    uint32
	Timestamp int64
	Content   []byte
}

func NewEncoder(cmd uint16, serial uint32, content []byte) *Encoder {
	return &Encoder{
		Head:      1,
		Version:   0,
		Length:    int32(len(content)),
		Command:   cmd,
		Serial:    serial,
		Timestamp: time.Now().Unix(),
		Content:   content,
	}
}

func (e *Encoder) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)

	// 包头 2字节
	err := binary.Write(buf, binary.BigEndian, e.Head)
	if err != nil {
		return nil, err
	}

	log.Printf("buf: %v", buf)

	// 版本号 1字节
	err = binary.Write(buf, binary.BigEndian, &e.Version)
	if err != nil {
		return nil, err
	}

	// 协议号 2字节
	err = binary.Write(buf, binary.BigEndian, &e.Command)
	if err != nil {
		return nil, err
	}

	// 序列号 4字节
	err = binary.Write(buf, binary.BigEndian, &e.Serial)
	if err != nil {
		return nil, err
	}

	// 时间戳 8字节
	err = binary.Write(buf, binary.BigEndian, &e.Timestamp)
	if err != nil {
		return nil, err
	}

	// 包内容长度 4字节
	err = binary.Write(buf, binary.BigEndian, &e.Length)
	if err != nil {
		return nil, err
	}

	// 包内容
	err = binary.Write(buf, binary.BigEndian, &e.Content)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
