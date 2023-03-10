package codec

import (
	"bytes"
	"encoding/binary"
)

type Codec interface {
	Decode(in []byte) error
	Encode() ([]byte, error)
}

type PortalCodec struct {
	Head      uint16
	Version   uint8
	Length    int32
	Command   uint16
	Serial    uint32
	Timestamp int64
	Content   []byte
}

func (c *PortalCodec) Decode(in []byte) error {
	buf := bytes.NewBuffer(in)

	// 包头 2字节
	err := binary.Read(bytes.NewBuffer(buf.Next(2)), binary.BigEndian, &c.Head)
	if err != nil {
		return err
	}

	// 版本号 1字节
	err = binary.Read(bytes.NewBuffer(buf.Next(1)), binary.BigEndian, &c.Version)
	if err != nil {
		return err
	}

	// 协议号 2字节
	err = binary.Read(bytes.NewBuffer(buf.Next(2)), binary.BigEndian, &c.Command)
	if err != nil {
		return err
	}

	// 序列号 4字节
	err = binary.Read(bytes.NewBuffer(buf.Next(4)), binary.BigEndian, &c.Serial)
	if err != nil {
		return err
	}

	// 时间戳 8字节
	err = binary.Read(bytes.NewBuffer(buf.Next(8)), binary.BigEndian, &c.Timestamp)
	if err != nil {
		return err
	}

	// 包内容长度 4字节
	err = binary.Read(bytes.NewBuffer(buf.Next(4)), binary.BigEndian, &c.Length)
	if err != nil {
		return err
	}

	// 包内容
	c.Content = buf.Next(int(c.Length))

	return nil
}

func (c *PortalCodec) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)

	// 包头 2字节
	err := binary.Write(buf, binary.BigEndian, c.Head)
	if err != nil {
		return nil, err
	}

	// 版本号 1字节
	err = binary.Write(buf, binary.BigEndian, &c.Version)
	if err != nil {
		return nil, err
	}

	// 协议号 2字节
	err = binary.Write(buf, binary.BigEndian, &c.Command)
	if err != nil {
		return nil, err
	}

	// 序列号 4字节
	err = binary.Write(buf, binary.BigEndian, &c.Serial)
	if err != nil {
		return nil, err
	}

	// 时间戳 8字节
	err = binary.Write(buf, binary.BigEndian, &c.Timestamp)
	if err != nil {
		return nil, err
	}

	// 包内容长度 4字节
	err = binary.Write(buf, binary.BigEndian, &c.Length)
	if err != nil {
		return nil, err
	}

	// 包内容
	err = binary.Write(buf, binary.BigEndian, &c.Content)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
