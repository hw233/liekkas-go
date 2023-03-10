package model

import (
	"context"
	"time"

	"gamesvr/manager"
	"shared/protobuf/pb"
	"shared/utility/errors"
)

func (u *User) RPCMailFetch() ([]*pb.VOServerMail, []*pb.VOServerMail, error) {
	userId := u.GetUserId()

	req := &pb.MailFetchReq{
		UserId: userId,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	resp, err := manager.RPCMailClient.Fetch(ctx, req)
	if err != nil {
		return nil, nil, errors.WrapTrace(err)
	}

	return resp.GroupMails, resp.PersonalMails, nil
}
