package main

import (
	"log"
	"shared/protobuf/pb"
)

func (c *Client) YggdrasilGetMain(req *pb.C2SYggdrasilGetMain) (*pb.S2CYggdrasilGetMain, error) {
	gameResp, err := c.Request(2801, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CYggdrasilGetMain{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: YggdrasilGetMain  resp: %+v", resp)

	return resp, nil
}
