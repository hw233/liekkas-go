package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"shared/protobuf/pb"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9080")
	if err != nil {
		log.Printf("net.Dial error: %v", err)
		return
	}
	defer conn.Close()

	client := NewClient(conn)

	go client.KeepReading()

	_, err = client.Login(&pb.C2SLogin{
		UserId: 111117,
	})
	if err != nil {
		log.Printf("login error: %v", err)
	}

	suitUserResp, err := client.SuitUser(&pb.C2SSuitUser{})
	if err != nil {
		log.Printf("SuitUser error: %v", err)
	}

	if suitUserResp != nil {
		for _, v := range suitUserResp.Items {
			log.Printf("ItemId: %d, Amount: %d", v.ItemId, v.Amount)
		}
	}

	// s2CCharacterLvUp, err := client.CharacterLvUp(&pb.C2SCharacterLvUp{
	// 	CharacterId: 1002,
	// 	Costs: []*pb.VOConsume{
	// 		{ItemId: 121004, Count: 1},
	// 	},
	// })
	// if err != nil {
	// 	log.Printf("CharacterLvUp error: %v", err)
	// 	return
	// }
	// log.Printf("s2CCharacterLvUp:%v", s2CCharacterLvUp.Result)
	// s2CCharacterSkillLvUp, err := client.CharacterSkillLvUp(&pb.C2SCharacterSkillLvUp{
	// 	CharacterId: 1002,
	// 	SkillId:     100201,
	// 	LvUpAmount:  1,
	// })
	// if err != nil {
	// 	log.Printf("CharacterSkillLvUp error: %v", err)
	// 	return
	// }
	// log.Printf("s2CCharacterSkillLvUp:%v", s2CCharacterSkillLvUp.ResourceResult)

	getQuestInfoResp, err := client.GetQuestInfo(&pb.C2SGetQuestInfo{})
	if err != nil {
		log.Printf("get quest req fail, error: %v\n", err)
	}

	fmt.Printf("getQuestInfoResp: %v\n", getQuestInfoResp)

	time.Sleep(5 * time.Second)
}
