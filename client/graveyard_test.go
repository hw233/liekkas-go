package main

import (
	"net"
	"shared/protobuf/pb"
	"testing"
)

func TestGraveyard(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:9080")
	if err != nil {
		t.Logf("net.Dial error: %v", err)
		return
	}

	client := NewClient(conn)

	go client.KeepReading()

	_, err = client.Login(&pb.C2SLogin{
		UserId: 20220216185012,
	})
	info, err := client.GraveyardGetInfo(&pb.C2SGraveyardGetInfo{})
	t.Logf("info: %v", info)
	//build,err:=client.GraveyardOpenCurtain(&pb.C2SGraveyardOpenCurtain{BuildUid: 1})
	//t.Logf("build: %v", build)
	//client.GraveyardBuildCreate(&pb.C2SGraveyardBuildCreate{
	//	BuildId: 2,
	//	Position: &pb.VOPosition{PosX: 0,PosY: -4},
	//})
	//client.GraveyardBuildCreate(&pb.C2SGraveyardBuildCreate{
	//	BuildId: 3,
	//	Position: &pb.VOPosition{PosX: 2,PosY: -9},
	//
	//})

	//client.GraveyardOpenCurtain(&pb.C2SGraveyardOpenCurtain{BuildUid: 2})
	//client.GraveyardOpenCurtain(&pb.C2SGraveyardOpenCurtain{BuildUid: 3})

	//client.GraveyardBuildLvUp(&pb.C2SGraveyardBuildLvUp{
	//	BuildUid: 2,
	//})
	//client.GraveyardBuildLvUp(&pb.C2SGraveyardBuildLvUp{
	//	BuildUid: 3,
	//})

	//client.GraveyardBuildStageUp(&pb.C2SGraveyardBuildStageUp{
	//	BuildUid: 3,
	//})

	// client.GraveyardRelocation(&pb.C2SGraveyardRelocation{
	//	Relocations: []*pb.VOBuildRelocation{
	//		{BuildUid: 1, Position: &pb.VOPosition{PosX: -6, PosY: -12}},
	//		{BuildUid: 3, Position: &pb.VOPosition{PosX: -8, PosY: -21}},
	//		{BuildUid: 2, Position: &pb.VOPosition{PosX: 5, PosY: -7}},
	//	},
	//})
	//test, err := client.Test(&pb.C2STest{})
	//t.Logf("info: %v", test)
	//client.GraveyardAskForHelp(&pb.C2SGraveyardAskForHelp{})
	//client.GenVOAddGraveyardRequest(&pb.C2SGraveyardSendHelpRequest{
	//	BuildUid: 1,
	//})
	//requests, err := client.GraveyardGetHelpRequests(&pb.C2SGraveyardGetHelpRequests{})
	//t.Logf("requests: %v", requests)

	//client.HeartBeatReq(&pb.C2SHeartBeatReq{})
	//
	//client.GraveyardProduceStart(&pb.C2SGraveyardProduceStart{
	//	BuildUid:   5,
	//	ProduceNum: 20,
	//})
	//
	//client.GraveyardAccelerate(&pb.C2SGraveyardAccelerate{
	//	BuildUid: 5,
	//	Consumes: []*pb.VOResource{{ItemId: 1610001, Count: 5}},
	//})
	//get, err := client.GraveyardProductionGet(&pb.C2SGraveyardProductionGet{
	//	BuildUidList: []int64{5},
	//})
	//t.Logf("get: %v", get)
	//reward, err := client.GraveyardReceivePlotReward(&pb.C2SGraveyardReceivePlotReward{})
	//t.Logf("reward: %v", reward)

	productionGet, err := client.GraveyardProductionGet(&pb.C2SGraveyardProductionGet{
		BuildUidList: []int64{8},
	})
	t.Log(productionGet)
	t.Log(err)
}
