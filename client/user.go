package main

import (
	"log"

	"shared/protobuf/pb"
)

func (c *Client) Login(req *pb.C2SLogin) (*pb.S2CLogin, error) {
	gameResp, err := c.Request(1001, req)
	if err != nil {
		log.Printf("RESP ERROR: login error: %+v", err)
		return nil, err
	}

	resp := &pb.S2CLogin{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		log.Printf("RESP ERROR: login error: %+v", err)
		return nil, err
	}

	log.Printf("RESP SUCCESS: login resp: %+v", resp)

	return resp, nil
}

func (c *Client) SuitUser(req *pb.C2SSuitUser) (*pb.S2CSuitUser, error) {
	gameResp, err := c.Request(1005, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CSuitUser{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: suit user resp: %+v", resp)

	return resp, nil
}

func (c *Client) GetQuestInfo(req *pb.C2SGetQuestInfo) (*pb.S2CGetQuestInfo, error) {
	gameResp, err := c.Request(2001, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CGetQuestInfo{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: get quest info resp: %+v", resp)

	return resp, nil
}

func (c *Client) Test(req *pb.C2STest) (*pb.S2CTest, error) {
	gameResp, err := c.Request(1, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CTest{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: Test resp: %+v", resp)

	return resp, nil
}

func (c *Client) HeartBeatReq(req *pb.C2SHeartBeatReq) (*pb.S2CHeartBeatReq, error) {
	gameResp, err := c.Request(101, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CHeartBeatReq{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: HeartBeatReq resp: %+v", resp)

	return resp, nil
}

func (c *Client) ExchangeCDKey(req *pb.C2SExchangeCDKey) (*pb.S2CExchangeCDKey, error) {
	gameResp, err := c.Request(1033, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CExchangeCDKey{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: ExchangeCDKey resp: %+v", resp)

	return resp, nil
}

func (c *Client) Mail(req *pb.C2SMailInfo) (*pb.S2CMailInfo, error) {
	gameResp, err := c.Request(3601, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CMailInfo{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: MailInfo resp: %+v", resp)

	return resp, nil
}

func (c *Client) StarUp(req *pb.C2SCharacterStarUp) (*pb.S2CCharacterStarUp, error) {
	gameResp, err := c.Request(1505, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CCharacterStarUp{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: CharacterStarUp resp: %+v", resp)

	return resp, nil
}

func (c *Client) SkillUp(req *pb.C2SCharacterSkillLvUp) (*pb.S2CCharacterSkillLvUp, error) {
	gameResp, err := c.Request(1509, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CCharacterSkillLvUp{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: CharacterSkillLvUp resp: %+v", resp)

	return resp, nil
}

func (c *Client) GM(req *pb.C2SGM) (*pb.S2CGM, error) {
	gameResp, err := c.Request(1029, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CGM{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: gm resp: %+v", resp)

	return resp, nil
}

func (c *Client) StoreGetGoods(req *pb.C2SStoreGetGoods) (*pb.S2CStoreGetGoods, error) {
	gameResp, err := c.Request(3401, req)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CStoreGetGoods{}
	err = c.Handle(gameResp, resp)
	if err != nil {
		return nil, err
	}

	log.Printf("RESP SUCCESS: StoreGetGoods resp: %+v", resp)

	return resp, nil
}
