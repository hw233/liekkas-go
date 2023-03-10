package main

import (
	"log"
	"shared/protobuf/pb"
)

func (c *Client) GetGachaList(req *pb.C2SGetGachaList) (*pb.S2CGetGachaList, error) {
	gameResp, err := c.Request(1201, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CGetGachaList{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: GetGachaList  resp: %+v", resp)

	return resp, nil
}

func (c *Client) UserGachaDrop(req *pb.C2SUserGachaDrop) (*pb.S2CUserGachaDrop, error) {
	gameResp, err := c.Request(1203, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CUserGachaDrop{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: UserGachaDrop resp: %+v", resp)

	return resp, nil
}
