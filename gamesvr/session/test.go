package session

import (
	"context"

	"shared/protobuf/pb"
)

func (s *Session) Test(ctx context.Context, req *pb.TestReq) (*pb.TestResp, error) {
	// log.Printf("session: name: %s uid: %d \n", req.Name, req.Uid)

	return &pb.TestResp{
		// Name: s.Name,
		Uid: s.User.ID,
	}, nil
}
