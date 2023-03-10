package method

import (
	"context"
	"gm/manager"
	"gm/param"
	"shared/common"
	"shared/protobuf/pb"
	"shared/utility/errors"
	"shared/utility/servertime"
)

const SystemMailTemplateId int32 = 999

func (p *HttpPostHandler) SendMail(ctx context.Context, sendMailParam *param.SendMailParam) error {
	if len(sendMailParam.UserIds) == 0 {
		return errors.New("userIds can not be nil")

	}
	userIds := make([]int64, 0, len(sendMailParam.UserIds))
	for _, id := range sendMailParam.UserIds {
		userIds = append(userIds, id.Id)
	}
	startTime := servertime.Now().Unix()
	if sendMailParam.StartTime != "" {
		t, err := servertime.ParseTime(sendMailParam.StartTime)
		if err != nil {
			return errors.WrapTrace(err)
		}
		startTime = t
	}
	rewards := common.NewRewards()
	for _, reward := range sendMailParam.Attachment {
		if reward.ID > 0 && reward.Num <= 0 {
			return errors.New("wrong Num of reward in attachment")
		}
		rewards.AddReward(common.NewReward(reward.ID, reward.Num))
	}
	if sendMailParam.ExpireDay <= 0 {
		return errors.New("expire can not be smaller than 0")
	}
	expireDay := sendMailParam.ExpireDay
	if len(userIds) == 1 {
		req := &pb.SendPersonalMailReq{
			TemplateId:  SystemMailTemplateId,
			Title:       sendMailParam.Title,
			TitleArgs:   []string{},
			Content:     sendMailParam.Content,
			ContentArgs: []string{},
			Attachment:  rewards.MergeVOResource(),
			Sender:      "",
			StartTime:   startTime,
			EndTime:     servertime.Now().Unix() + servertime.SecondPerDay*int64(expireDay),
			ExpireTime:  servertime.Now().Unix() + servertime.SecondPerDay*int64(expireDay),
			Users:       userIds,
		}
		// _ = req
		_, err := manager.RPCMailClient.SendPersonalMail(ctx, req)
		if err != nil {
			return errors.WrapTrace(err)
		}
	} else {
		req := &pb.SendGroupMailReq{
			TemplateId:  SystemMailTemplateId,
			Title:       sendMailParam.Title,
			TitleArgs:   []string{},
			Content:     sendMailParam.Content,
			ContentArgs: []string{},
			Attachment:  rewards.MergeVOResource(),
			Sender:      "",
			StartTime:   startTime,
			EndTime:     servertime.Now().Unix() + servertime.SecondPerDay*int64(expireDay),
			ExpireTime:  servertime.Now().Unix() + servertime.SecondPerDay*int64(expireDay),
			Users:       userIds,
		}
		// _ = req

		_, err := manager.RPCMailClient.SendGroupMail(ctx, req)
		if err != nil {
			return errors.WrapTrace(err)
		}

	}

	return nil
}

func (p *HttpPostHandler) SendWholeServerMail(ctx context.Context, sendMailParam *param.SendWholeServerMailParam) error {

	startTime := servertime.Now().Unix()
	if sendMailParam.StartTime != "" {
		t, err := servertime.ParseTime(sendMailParam.StartTime)
		if err != nil {
			return errors.WrapTrace(err)
		}
		startTime = t
	}
	rewards := common.NewRewards()
	for _, reward := range sendMailParam.Attachment {
		if reward.ID > 0 && reward.Num <= 0 {
			return errors.New("wrong Num of reward in attachment")
		}
		rewards.AddReward(common.NewReward(reward.ID, reward.Num))
	}
	if sendMailParam.ExpireDay <= 0 {
		return errors.New("expire can not be smaller than 0")
	}
	expireDay := sendMailParam.ExpireDay

	req := &pb.SendWholeServerMailReq{
		TemplateId:  SystemMailTemplateId,
		Title:       sendMailParam.Title,
		TitleArgs:   []string{},
		Content:     sendMailParam.Content,
		ContentArgs: []string{},
		Attachment:  rewards.MergeVOResource(),
		Sender:      "",
		StartTime:   startTime,
		EndTime:     servertime.Now().Unix() + servertime.SecondPerDay*int64(expireDay),
		ExpireTime:  servertime.Now().Unix() + servertime.SecondPerDay*int64(expireDay),
	}
	// _ = req
	_, err := manager.RPCMailClient.SendWholeServerMail(ctx, req)
	if err != nil {
		return errors.WrapTrace(err)
	}
	return nil
}
