package main

import (
	"log"

	"shared/protobuf/pb"
)

func (c *Client) GraveyardGetInfo(req *pb.C2SGraveyardGetInfo) (*pb.S2CGraveyardGetInfo, error) {
	gameResp, err := c.Request(1401, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CGraveyardGetInfo{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: GraveyardGetInfo  resp: %+v", resp)

	return resp, nil
}

func (c *Client) GraveyardBuildCreate(req *pb.C2SGraveyardBuildCreate) (*pb.S2CGraveyardBuildCreate, error) {
	gameResp, err := c.Request(1403, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CGraveyardBuildCreate{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: GraveyardBuildCreate resp: %+v", resp)

	return resp, nil
}

func (c *Client) GraveyardBuildLvUp(req *pb.C2SGraveyardBuildLvUp) (*pb.S2CGraveyardBuildLvUp, error) {
	gameResp, err := c.Request(1405, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CGraveyardBuildLvUp{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: GraveyardBuildLvUp resp: %+v", resp)

	return resp, nil
}

func (c *Client) GraveyardBuildStageUp(req *pb.C2SGraveyardBuildStageUp) (*pb.S2CGraveyardBuildStageUp, error) {
	gameResp, err := c.Request(1407, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CGraveyardBuildStageUp{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: GraveyardBuildStageUp resp: %+v", resp)

	return resp, nil
}

func (c *Client) GraveyardRelocation(req *pb.C2SGraveyardRelocation) (*pb.S2CGraveyardRelocation, error) {
	gameResp, err := c.Request(1419, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CGraveyardRelocation{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: GraveyardRelocation resp: %+v", resp)

	return resp, nil
}
func (c *Client) GraveyardOpenCurtain(req *pb.C2SGraveyardOpenCurtain) (*pb.S2CGraveyardOpenCurtain, error) {
	gameResp, err := c.Request(1413, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CGraveyardOpenCurtain{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: GraveyardOpenCurtain resp: %+v", resp)

	return resp, nil
}

func (c *Client) GraveyardHelp(req *pb.C2SGraveyardHelp) (*pb.S2CGraveyardHelp, error) {
	gameResp, err := c.Request(1425, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CGraveyardHelp{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: GraveyardHelp resp: %+v", resp)

	return resp, nil
}

func (c *Client) GraveyardSendHelpRequest(req *pb.C2SGraveyardSendHelpRequest) (*pb.S2CGraveyardSendHelpRequest, error) {
	gameResp, err := c.Request(1427, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CGraveyardSendHelpRequest{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: GenVOAddGraveyardRequest resp: %+v", resp)

	return resp, nil
}

func (c *Client) GraveyardGetHelpRequests(req *pb.C2SGraveyardGetHelpRequests) (*pb.S2CGraveyardGetHelpRequests, error) {
	gameResp, err := c.Request(1429, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CGraveyardGetHelpRequests{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: GraveyardGetHelpRequests resp: %+v", resp)

	return resp, nil
}

func (c *Client) GraveyardProduceStart(req *pb.C2SGraveyardProduceStart) (*pb.S2CGraveyardProduceStart, error) {
	gameResp, err := c.Request(1409, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CGraveyardProduceStart{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: GraveyardProduceStart resp: %+v", resp)

	return resp, nil
}
func (c *Client) GraveyardProductionGet(req *pb.C2SGraveyardProductionGet) (*pb.S2CGraveyardProductionGet, error) {
	gameResp, err := c.Request(1411, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CGraveyardProductionGet{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: GraveyardProductionGet resp: %+v", resp)

	return resp, nil
}

func (c *Client) GraveyardAccelerate(req *pb.C2SGraveyardAccelerate) (*pb.S2CGraveyardAccelerate, error) {
	gameResp, err := c.Request(1423, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CGraveyardAccelerate{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: GraveyardGetHelpRequests resp: %+v", resp)

	return resp, nil
}

func (c *Client) GraveyardReceivePlotReward(req *pb.C2SGraveyardReceivePlotReward) (*pb.S2CGraveyardReceivePlotReward, error) {
	gameResp, err := c.Request(1435, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CGraveyardReceivePlotReward{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: GraveyardReceivePlotReward resp: %+v", resp)

	return resp, nil
}

func (c *Client) GraveyardGetRewardHours(req *pb.C2SGraveyardGetRewardHours) (*pb.S2CGraveyardGetRewardHours, error) {
	gameResp, err := c.Request(1437, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CGraveyardGetRewardHours{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: GraveyardGetRewardHours resp: %+v", resp)

	return resp, nil
}
