package main

import (
	"testing"
	"time"

	"shared/protobuf/pb"
)

func TestGMAddItem(t *testing.T) {
	client := Pool.NewClient()
	if client == nil {
		return
	}

	_, err := client.Login(&pb.C2SLogin{
		UserId: 20211208201649,
	})
	if err != nil {
		return
	}

	for i := 0; i < 100; i++ {
		time.Sleep(time.Second * 1)
		_, err = client.ExchangeCDKey(&pb.C2SExchangeCDKey{
			Code: "#GM:addItem(11,1);addItem(11,2)",
		})
		if err != nil {
			return
		}
	}

}

func TestGMConfigTask(t *testing.T) {
	client := Pool.NewClient()
	if client == nil {
		return
	}

	_, err := client.Login(&pb.C2SLogin{
		UserId: 20211208201649,
	})
	if err != nil {
		return
	}

	_, err = client.ExchangeCDKey(&pb.C2SExchangeCDKey{
		Code: "#GM:configTask(1)",
	})

	if err != nil {
		return
	}
}

func TestGMPassGuide(t *testing.T) {
	client := Pool.NewClient()
	if client == nil {
		return
	}

	_, err := client.Login(&pb.C2SLogin{
		UserId: 20211208201649,
	})
	if err != nil {
		return
	}

	_, err = client.ExchangeCDKey(&pb.C2SExchangeCDKey{
		Code: "#GM:PassGuide(60)",
	})

	if err != nil {
		return
	}
}
