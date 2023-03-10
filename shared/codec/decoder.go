package codec

import (
	"bytes"
	"encoding/binary"
)

type Decoder struct {
	Head      uint16
	Version   uint8
	Length    int32
	Command   uint16
	Serial    uint32
	Timestamp int64
	Content   []byte
}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (td *Decoder) Decode(bs []byte) error {
	buf := bytes.NewBuffer(bs)

	// 包头 2字节
	err := binary.Read(bytes.NewBuffer(buf.Next(2)), binary.BigEndian, &td.Head)
	if err != nil {
		return err
	}

	// 版本号 1字节
	err = binary.Read(bytes.NewBuffer(buf.Next(1)), binary.BigEndian, &td.Version)
	if err != nil {
		return err
	}

	// 协议号 2字节
	err = binary.Read(bytes.NewBuffer(buf.Next(2)), binary.BigEndian, &td.Command)
	if err != nil {
		return err
	}

	// 序列号 4字节
	err = binary.Read(bytes.NewBuffer(buf.Next(4)), binary.BigEndian, &td.Serial)
	if err != nil {
		return err
	}

	// 时间戳 8字节
	err = binary.Read(bytes.NewBuffer(buf.Next(8)), binary.BigEndian, &td.Timestamp)
	if err != nil {
		return err
	}

	// 包内容长度 4字节
	err = binary.Read(bytes.NewBuffer(buf.Next(4)), binary.BigEndian, &td.Length)
	if err != nil {
		return err
	}

	// 包内容
	td.Content = buf.Next(int(td.Length))

	return nil
}
