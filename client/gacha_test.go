package main

import (
	"log"
	"net"
	"shared/protobuf/pb"
	"testing"
)

func TestGacha(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:9080")
	if err != nil {
		t.Logf("net.Dial error: %v", err)
		return
	}

	client := NewClient(conn)

	go client.KeepReading()

	_, err = client.Login(&pb.C2SLogin{
		UserId: 12313212315,
	})
	suitUser, err := client.SuitUser(&pb.C2SSuitUser{})
	if err != nil {
		log.Printf("SuitUser error: %v", err)
		return
	}
	log.Printf("SuitUser Characters: %v", suitUser.Characters)

	//info, err := client.GetGachaList(&pb.C2SGetGachaList{})
	//t.Logf("info: %v", info)

	drop, err := client.UserGachaDrop(&pb.C2SUserGachaDrop{
		IsSingle: false,
		GachaId:  1,
	})
	t.Logf("drop: %+v", drop)
}
