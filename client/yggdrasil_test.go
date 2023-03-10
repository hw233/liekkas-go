package main

import (
	"shared/protobuf/pb"
	"testing"
)

func TestClient_YggdrasilGetMain(t *testing.T) {
	client := Pool.NewClient()

	client.Login(&pb.C2SLogin{
		UserId: 123,
	})
	info, err := client.YggdrasilGetMain(&pb.C2SYggdrasilGetMain{})
	t.Logf("info: %v", info)
	t.Logf("err: %v", err)
}
