package controller

import (
	"shared/protobuf/pb"
	"shared/utility/errors"
	"shared/utility/servertime"

	"google.golang.org/protobuf/proto"
)

func newGameResp(req *pb.GameReq, content []byte) *pb.GameResp {
	return &pb.GameResp{
		Command:   req.Command + 1,
		Serial:    req.Serial,
		Timestamp: servertime.Now().Unix(),
		Content:   content,
	}
}

func newErrResp(req *pb.GameReq, err error) *pb.GameResp {
	content, _ := proto.Marshal(&pb.S2CError{
		Error: int32(errors.Code(err)),
		Arg:   "",
	})

	return &pb.GameResp{
		Command:   23,
		Serial:    req.Serial,
		Timestamp: servertime.Now().Unix(),
		Content:   content,
	}
}

func newConnectResp(uid int64, gameResp *pb.GameResp) *pb.ConnectResp {
	return &pb.ConnectResp{
		GameResp: gameResp,
		Uid:      uid,
	}
}

func newTransferResp(gameResp *pb.GameResp) *pb.TransferResp {
	return &pb.TransferResp{
		GameResp: gameResp,
	}
}
