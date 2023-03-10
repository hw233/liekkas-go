package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"shared/protobuf/pb"
	"sync"

	"github.com/golang/protobuf/proto"
)

var (
	bufferInitSize = 1024
	pushCmdMap     map[int32]proto.Message
)

func init() {
	pushCmdMap = make(map[int32]proto.Message)
	pushCmdMap[1] = &pb.C2STest{}
	pushCmdMap[2] = &pb.S2CTest{}

	pushCmdMap[1401] = &pb.C2SGraveyardGetInfo{}
	pushCmdMap[1402] = &pb.S2CGraveyardGetInfo{}
}

type ClientPool struct {
	sync.RWMutex
	pool      []*Client
	newClient func() *Client
}

func NewClientPool(addr string) *ClientPool {
	return &ClientPool{
		pool: []*Client{},
		newClient: func() *Client {
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				return nil
			}

			return NewClient(conn)
		},
	}
}

func (cp *ClientPool) NewClient() *Client {
	client := cp.newClient()
	if client == nil {
		return nil
	}
	go client.KeepReading()
	cp.pool = append(cp.pool, client)
	return client
}

func (cp *ClientPool) Close() {
	for _, c := range cp.pool {
		c.Close()
	}
}

type Client struct {
	sync.RWMutex

	net.Conn
	serial   int32
	receiver map[int32]chan *pb.GameResp
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		Conn:     conn,
		serial:   1,
		receiver: make(map[int32]chan *pb.GameResp),
	}
}

func (c *Client) Request(cmd int32, message proto.Message) (*pb.GameResp, error) {
	content, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}

	req := &pb.GameReq{
		Command: cmd,
		Serial:  c.serial,
		Content: content,
	}

	r, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}

	err = binary.Write(buf, binary.BigEndian, int32(len(r)))
	if err != nil {
		return nil, err
	}

	buf.Write(r)

	_, err = c.Conn.Write(buf.Bytes())
	if err != nil {
		return nil, err
	}

	receive := make(chan *pb.GameResp, 1)

	c.Lock()
	c.receiver[c.serial] = receive
	c.serial++
	c.Unlock()

	select {
	case ret := <-receive:
		return ret, nil
	}
}

func (c *Client) KeepReading() {
	defer c.Close()

	for {
		r, err := c.read()
		if err != nil {
			log.Printf("ERROR: read error, %v", err)
			continue
		}
		// // 未读到消息
		// if len(r) < 4 {
		// 	continue
		// }

		resp := &pb.GameResp{}

		buf := bytes.NewBuffer(r)

		var l int32

		err = binary.Read(bytes.NewBuffer(buf.Next(4)), binary.BigEndian, &l)
		if err != nil {
			log.Printf("ERROR: binary read error: %v", err)
			break
		}
		err = proto.Unmarshal(buf.Next(int(l)), resp)
		if err != nil {
			log.Printf("ERROR: proto unmarshal error: %v", err)
			continue
		}

		// 推送
		if resp.Serial == 0 {
			handlePush(resp)
		} else {
			c.Lock()
			if receive, ok := c.receiver[resp.Serial]; ok {
				receive <- resp
				delete(c.receiver, resp.Serial)
			}
			c.Unlock()
		}
	}
}

func (c *Client) read() ([]byte, error) {
	buf := &bytes.Buffer{}

	for {
		var temp = make([]byte, bufferInitSize)

		n, err := c.Conn.Read(temp)
		if err != nil {
			// 读完了
			if err != io.EOF {
				return nil, err
			}

			break
		}

		if n < bufferInitSize {
			buf.Write(temp[:n+1])
			break
		}

		buf.Write(temp)
	}

	return buf.Bytes(), nil
}

func (c *Client) Error(resp *pb.GameResp) (*pb.S2CError, bool) {
	if resp.Command == 23 {
		e := &pb.S2CError{}
		err := proto.Unmarshal(resp.Content, e)
		if err != nil {
			return nil, false
		}

		return e, true
	}

	return nil, false
}

func (c *Client) Unmarshal(resp *pb.GameResp, message proto.Message) error {
	err := proto.Unmarshal(resp.Content, message)
	if err != nil {
		return err
	}

	return nil

}

func (c *Client) Handle(gameResp *pb.GameResp, message proto.Message) error {
	errResp, ok := c.Error(gameResp)
	if ok {
		log.Printf("RESP ERROR: %d, %s", errResp.Error, errResp.Arg)
		return errors.New("resp error")
	}

	err := c.Unmarshal(gameResp, message)
	if err != nil {
		return err
	}

	return nil
}

func handlePush(resp *pb.GameResp) {
	message, ok := pushCmdMap[resp.Command]
	if !ok {
		log.Printf("unknown push cmd: %d", resp.Command)
		return
	}
	err := proto.Unmarshal(resp.Content, message)
	if err != nil {
		log.Printf("unmarshal push error,cmd: %d ,err: %v", resp.Command, err)
	}
	log.Printf("push resp,cmd: %d ,content: %+v", resp.Command, pushCmdMap[resp.Command])
}
