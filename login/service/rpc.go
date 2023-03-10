package service

import (
	"context"

	"google.golang.org/grpc/status"

	"login/model"
	"shared/utility/errors"
	"shared/utility/glog"

	"shared/protobuf/pb"
)

type RPCLoginHandler struct {
	pb.UnimplementedLoginServer
}

func (RPCLoginHandler) CheckToken(ctx context.Context, req *pb.CheckTokenReq) (*pb.CheckTokenResp, error) {
	// login, err := model.GetLoginAccountByUserID(ctx, req.UserID)
	// if err != nil {
	// 	return nil, err
	// }

	token, err := model.GetToken(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	return &pb.CheckTokenResp{
		OK: token == req.Token,
	}, nil
}

func (RPCLoginHandler) Test(ctx context.Context, req *pb.LoginTestReq) (*pb.LoginTestResp, error) {
	if req.Id == 0 {
		return nil, status.Errorf(2100, "test error!")
	}

	if req.Id == 1 {
		errors.SetDefaultCode(1000)
		err := errors.NewCode(2200, "test error2!")

		if s, ok := err.(interface {
			GRPCStatus() *status.Status
		}); ok {
			glog.Infof("code :%d", s.GRPCStatus().Code())
		}

		s, ok := status.FromError(err)
		if !ok {
			glog.Errorf("xxxxxxxxxx!ok!ok!ok")
		}

		glog.Infof("code :%d, msg: %s", s.Code(), s.Message())

		err = errors.WrapTrace(err)

		err = errors.WrapText(err, "123")
		err = errors.WrapTrace(err)
		err = errors.Wrap(err, "456")

		err = errors.WrapTrace(err)
		return nil, err
	}

	return &pb.LoginTestResp{
		Msg: req.Msg,
		Id:  req.Id,
	}, nil
}
