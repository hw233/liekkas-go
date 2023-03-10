package mail

import (
	"context"

	"mailserver/manager"
	"shared/common"
	"shared/protobuf/pb"
)

func NewGroupMailNotify(ctx context.Context, groupMail *common.ServerGroupMail) error {
	req := &pb.NewGroupMailReq{
		Mail: groupMail.VOServerGroupMail(),
	}

	_, err := manager.RPCGameClient.NewGroupMailNotify(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func NewPersonalMailNotify(ctx context.Context, personalMail *common.ServerPersonalMail) error {
	_, err := manager.RPCGameClient.NewPersonalMailNotify(
		ctx,
		&pb.NewPersonalMailReq{
			Mail: personalMail.VOServerPersonalMail(),
		},
	)

	if err != nil {
		return err
	}

	return nil
}
