package method

import (
	"context"
	"gm/manager"
	"gm/param"
	"shared/protobuf/pb"
	"shared/utility/errors"
)

func (p *HttpPostHandler) MaintainSwitch(ctx context.Context, m *param.MaintainSwitch) error {

	err := manager.Global.SetMaintainSwitch(ctx, m.Switch)
	if err != nil {
		return errors.WrapTrace(err)
	}
	_, err = manager.RPCGameClient.ReloadMaintain(ctx, &pb.ReloadMaintainReq{})
	if err != nil {
		return errors.WrapTrace(err)
	}

	_, err = manager.RPCForeplayClient.ReloadMaintain(ctx, &pb.ReloadMaintainReq{})
	if err != nil {
		return errors.WrapTrace(err)
	}

	return nil
}
