package main

import (
	"testing"

	"shared/protobuf/pb"
)

func TestEquipmentRecast(t *testing.T) {
	client := Pool.NewClient()

	_, err := client.Login(&pb.C2SLogin{
		UserId: 20211208201649,
	})
	if err != nil {
		return
	}

	// _, err = client.EquipmentRecastCamp(&pb.C2SEquipmentRecastCamp{
	// 	Id: 20,
	// })
	//
	// if err != nil {
	// 	return
	// }

	respConfirmRecastCamp, err := client.EquipmentConfirmRecastCamp(&pb.C2SEquipmentConfirmRecastCamp{
		Id:      20,
		Confirm: true,
	})

	if err != nil {
		return
	}

	t.Logf("%+v", respConfirmRecastCamp.Equipment)
}
