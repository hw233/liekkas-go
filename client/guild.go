package main

import (
	"log"

	"shared/protobuf/pb"
)

func (c *Client) GuildInfo(req *pb.C2SGuildInfo) (*pb.S2CGuildInfo, error) {
	gameResp, err := c.Request(10501, req)
	if err != nil {
		log.Printf("RESP ERROR: GuildInfo request error: %+v", err)
		return nil, err
	}

	resp := &pb.S2CGuildInfo{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		log.Printf("RESP ERROR: GuildInfo handle error: %+v", err)
		return nil, err
	}

	log.Printf("RESP SUCCESS: GuildInfo  resp: %+v", resp)

	return resp, nil
}

func (c *Client) GuildCreate(req *pb.C2SGuildCreate) (*pb.S2CGuildCreate, error) {
	gameResp, err := c.Request(10503, req)
	if err != nil {
		log.Printf("RESP ERROR: GuildCreate request error: %+v", err)
		return nil, err
	}

	resp := &pb.S2CGuildCreate{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		log.Printf("RESP ERROR: GuildCreate handle error: %+v", err)
		return nil, err
	}

	log.Printf("RESP SUCCESS: GuildCreate  resp: %+v", resp)

	return resp, nil
}

func (c *Client) GuildDissolve(req *pb.C2SGuildDissolve) (*pb.S2CGuildDissolve, error) {
	gameResp, err := c.Request(10505, req)
	if err != nil {
		log.Printf("RESP ERROR: GuildDissolve request error: %+v", err)
		return nil, err
	}

	resp := &pb.S2CGuildDissolve{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		log.Printf("RESP ERROR: GuildDissolve handle error: %+v", err)
		return nil, err
	}

	log.Printf("RESP SUCCESS: GuildDissolve  resp: %+v", resp)

	return resp, nil
}

func (c *Client) GuildApply(req *pb.C2SGuildApply) (*pb.S2CGuildApply, error) {
	gameResp, err := c.Request(10509, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CGuildApply{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		log.Printf("RESP ERROR: GuildApply handle error: %+v", err)
		return nil, err
	}

	log.Printf("RESP SUCCESS: GuildApply  resp: %+v", resp)

	return resp, nil
}

func (c *Client) GuildHandleApplied(req *pb.C2SGuildHandleApplied) (*pb.S2CGuildHandleApplied, error) {
	gameResp, err := c.Request(10513, req)
	if err != nil {
		log.Printf("RESP ERROR: GuildHandleApplied request error: %+v", err)
		return nil, err
	}

	resp := &pb.S2CGuildHandleApplied{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		log.Printf("RESP ERROR: GuildHandleApplied handle error: %+v", err)
		return nil, err
	}

	log.Printf("RESP SUCCESS: GuildHandleApplied  resp: %+v", resp)

	return resp, nil
}

func (c *Client) GuildPromotion(req *pb.C2SGuildPromotion) (*pb.S2CGuildPromotion, error) {
	gameResp, err := c.Request(10519, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CGuildPromotion{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		log.Printf("RESP ERROR: GuildPromotion handle error: %+v", err)
		return nil, err
	}

	log.Printf("RESP SUCCESS: GuildPromotion  resp: %+v", resp)

	return resp, nil
}

func (c *Client) GuildMercenaryList(req *pb.C2SGetMercenaryList) (*pb.S2CGetMercenaryList, error) {
	gameResp, err := c.Request(10601, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CGetMercenaryList{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		log.Printf("RESP ERROR:GetMercenaryList handle error: %+v", err)
		return nil, err
	}

	log.Printf("RESP SUCCESS: GetMercenaryList  resp: %+v", resp)

	return resp, nil
}

func (c *Client) GuildGetList(req *pb.C2SGuildGetList) (*pb.S2CGuildGetList, error) {
	gameResp, err := c.Request(10529, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CGuildGetList{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		log.Printf("RESP ERROR:GuildGetListhandle error: %+v", err)
		return nil, err
	}

	log.Printf("RESP SUCCESS: GuildGetListresp: %+v", resp)

	return resp, nil
}

func (c *Client) GuildListInfo(req *pb.C2SGuildListInfo) (*pb.S2CGuildListInfo, error) {
	gameResp, err := c.Request(10535, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CGuildListInfo{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		log.Printf("RESP ERROR:GuildListInfo handle error: %+v", err)
		return nil, err
	}

	log.Printf("RESP SUCCESS: GuildListInfo resp: %+v", resp)

	return resp, nil
}

func (c *Client) MercenarySendApply(req *pb.C2SMercenarySendApply) (*pb.S2CMercenarySendApply, error) {
	gameResp, err := c.Request(10603, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CMercenarySendApply{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		log.Printf("RESP ERROR:MercenarySendApply handle error: %+v", err)
		return nil, err
	}

	log.Printf("RESP SUCCESS: MercenarySendApply  resp: %+v", resp)

	return resp, nil
}

func (c *Client) MercenaryHandleApply(req *pb.C2SMercenaryHandleApply) (*pb.S2CMercenaryHandleApply, error) {
	gameResp, err := c.Request(10607, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CMercenaryHandleApply{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		log.Printf("RESP ERROR:MercenaryhandleApply handle error: %+v", err)
		return nil, err
	}

	log.Printf("RESP SUCCESS: MercenaryhandleApply  resp: %+v", resp)

	return resp, nil
}

func (c *Client) MercenaryManagement(req *pb.C2SMercenaryManagement) (*pb.S2CMercenaryManagement, error) {
	gameResp, err := c.Request(10605, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CMercenaryManagement{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		log.Printf("RESP ERROR:MercenaryManagement handle error: %+v", err)
		return nil, err
	}

	log.Printf("RESP SUCCESS: MercenaryManagement  resp: %+v", resp)

	return resp, nil
}

func (c *Client) MercenaryRecord(req *pb.C2SMercenaryRecord) (*pb.S2CMercenaryRecord, error) {
	gameResp, err := c.Request(10609, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CMercenaryRecord{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		log.Printf("RESP ERROR:MercenaryRecord handle error: %+v", err)
		return nil, err
	}

	log.Printf("RESP SUCCESS: MercenaryRecord  resp: %+v", resp)

	return resp, nil
}

func (c *Client) GuildModify(req *pb.C2SGuildModify) (*pb.S2CGuildModify, error) {
	gameResp, err := c.Request(10549, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CGuildModify{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		log.Printf("RESP ERROR:MGuildModify handle error: %+v", err)
		return nil, err
	}

	log.Printf("RESP SUCCESS:GuildModify resp: %+v", resp)

	return resp, nil
}

func (c *Client) GuildChat(req *pb.C2SGuildChat) (*pb.S2CGuildChat, error) {
	gameResp, err := c.Request(10517, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CGuildChat{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		log.Printf("RESP ERROR:GuildChat handle error: %+v", err)
		return nil, err
	}

	log.Printf("RESP SUCCESS:GuildChat resp: %+v", resp)

	return resp, nil
}

func (c *Client) GuildApplyList(req *pb.C2SGuildApplyList) (*pb.S2CGuildApplyList, error) {
	gameResp, err := c.Request(10537, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CGuildApplyList{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		log.Printf("RESP ERROR:GuildApplyListhandle error: %+v", err)
		return nil, err
	}

	log.Printf("RESP SUCCESS:GuildApplyList resp: %+v", resp)

	return resp, nil
}
