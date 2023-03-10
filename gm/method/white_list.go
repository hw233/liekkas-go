package method

import (
	"context"
	"gm/manager"
	"gm/param"
	"shared/protobuf/pb"
	"shared/utility/errors"
	"shared/utility/whitelist"
	"strconv"
)

func (g *HttpGetHandler) FetchWhiteList(ctx context.Context) (map[whitelist.WhiteListType][]interface{}, error) {
	list, err := manager.Global.GetIdWhiteList(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	m := map[whitelist.WhiteListType][]interface{}{}
	m[whitelist.Id] = []interface{}{}
	for _, option := range *list {
		l := m[option.T]
		l = append(l, option.V)
		m[option.T] = l
	}
	return m, nil
}

func (p *HttpPostHandler) AddWhiteList(ctx context.Context, item *param.WhiteListItem) error {

	switch item.Type {
	case whitelist.Id:
		uid, err := strconv.ParseInt(item.Value, 10, 64)

		if err != nil {
			return errors.WrapTrace(err)
		}
		err = manager.Global.AddIdWhiteList(ctx, uid)
		if err != nil {
			return errors.WrapTrace(err)
		}
	}
	_, err := manager.RPCGameClient.ReloadWhiteList(ctx, &pb.ReloadWhiteListReq{})
	if err != nil {
		return errors.WrapTrace(err)
	}

	return nil
}

func (p *HttpPostHandler) DelWhiteList(ctx context.Context, item *param.WhiteListItem) error {

	switch item.Type {
	case whitelist.Id:
		uid, err := strconv.ParseInt(item.Value, 10, 64)
		if err != nil {
			return errors.WrapTrace(err)
		}

		err = manager.Global.DelIdWhiteList(ctx, uid)
		if err != nil {
			return errors.WrapTrace(err)
		}
	}
	_, err := manager.RPCGameClient.ReloadWhiteList(ctx, &pb.ReloadWhiteListReq{})
	if err != nil {
		return errors.WrapTrace(err)
	}
	return nil
}

func (p *HttpPostHandler) WhiteListSwitch(ctx context.Context, w *param.MaintainSwitch) error {

	err := manager.Global.SetWLSwitch(ctx, w.Switch)
	if err != nil {
		return errors.WrapTrace(err)
	}
	_, err = manager.RPCGameClient.ReloadWhiteList(ctx, &pb.ReloadWhiteListReq{})
	if err != nil {
		return errors.WrapTrace(err)
	}
	return nil
}
