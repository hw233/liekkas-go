package main

import (
	"log"
	"shared/protobuf/pb"
)

func (c *Client) CharacterLvUp(req *pb.C2SCharacterLvUp) (*pb.S2CCharacterLvUp, error) {
	gameResp, err := c.Request(1503, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CCharacterLvUp{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: character lv up resp: %+v", resp)

	return resp, nil
}

func (c *Client) CharacterSkillLvUp(req *pb.C2SCharacterSkillLvUp) (*pb.S2CCharacterSkillLvUp, error) {
	gameResp, err := c.Request(1509, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CCharacterSkillLvUp{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: character skill lv up resp: %+v", resp)

	return resp, nil
}
