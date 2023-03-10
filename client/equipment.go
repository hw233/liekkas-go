package main

import (
	"log"

	"shared/protobuf/pb"
)

func (c *Client) EquipmentRecastCamp(req *pb.C2SEquipmentRecastCamp) (*pb.S2CEquipmentRecastCamp, error) {
	gameResp, err := c.Request(1607, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CEquipmentRecastCamp{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: EquipmentRecastCamp  resp: %+v", resp)

	return resp, nil
}

func (c *Client) EquipmentConfirmRecastCamp(req *pb.C2SEquipmentConfirmRecastCamp) (*pb.S2CEquipmentConfirmRecastCamp, error) {
	gameResp, err := c.Request(1609, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CEquipmentConfirmRecastCamp{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: EquipmentConfirmRecastCamp  resp: %+v", resp)

	return resp, nil
}
