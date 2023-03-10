package controller

import (
	"context"
	"mailserver/mail"
	"shared/common"
	"shared/protobuf/pb"
	"shared/utility/glog"
)

type ServiceHandler struct {
	*pb.UnimplementedMailServer
}

func NewServiceHandler() *ServiceHandler {
	return &ServiceHandler{
		UnimplementedMailServer: &pb.UnimplementedMailServer{},
	}
}

func (sh *ServiceHandler) Fetch(ctx context.Context, req *pb.MailFetchReq) (*pb.MailFetchResp, error) {

	userId := req.UserId
	groupMails, personalMails, err := mail.MailMgr.FetchMail(ctx, userId)
	if err != nil {
		return nil, err
	}

	voGroupMails := make([]*pb.VOServerMail, 0, len(groupMails))
	for _, mail := range groupMails {
		if mail.HasUser(userId) {
			voGroupMails = append(voGroupMails, mail.VOServerMail())
		}
	}

	voPersonalMails := make([]*pb.VOServerMail, 0, len(personalMails))
	for _, mail := range personalMails {
		voPersonalMails = append(voPersonalMails, mail.VOServerMail())
	}

	return &pb.MailFetchResp{
		UserId:        userId,
		GroupMails:    voGroupMails,
		PersonalMails: voPersonalMails,
	}, nil
}

func (sh *ServiceHandler) SendWholeServerMail(ctx context.Context, req *pb.SendWholeServerMailReq) (*pb.SendWholeServerMailResp, error) {
	attachment, err := common.ParseFromVOConsume(req.Attachment)
	if err != nil {
		return nil, err
	}

	err = mail.MailMgr.SendGroupMail(ctx, []int64{}, true, req.TemplateId, req.Title, req.Content,
		req.TitleArgs, req.ContentArgs, attachment, req.Sender, req.StartTime, req.ExpireTime, req.EndTime)

	if err != nil {
		return nil, err
	}

	return &pb.SendWholeServerMailResp{}, nil
}

func (sh *ServiceHandler) SendGroupMail(ctx context.Context, req *pb.SendGroupMailReq) (*pb.SendGroupMailResp, error) {
	rewardList := req.Attachment
	if rewardList == nil {
		rewardList = []*pb.VOResource{}
	}

	attachment, err := common.ParseFromVOConsume(rewardList)
	if err != nil {
		return nil, err
	}

	err = mail.MailMgr.SendGroupMail(ctx, req.Users, false, req.TemplateId, req.Title, req.Content,
		req.TitleArgs, req.ContentArgs, attachment, req.Sender, req.StartTime, req.ExpireTime, req.EndTime)

	if err != nil {
		return nil, err
	}

	return &pb.SendGroupMailResp{}, nil
}

func (sh *ServiceHandler) SendPersonalMail(ctx context.Context, req *pb.SendPersonalMailReq) (*pb.SendPersonalMailResp, error) {
	rewardList := req.Attachment
	if rewardList == nil {
		rewardList = []*pb.VOResource{}
	}

	attachment, err := common.ParseFromVOConsume(rewardList)
	if err != nil {
		return nil, err
	}

	for _, userId := range req.Users {
		err = mail.MailMgr.SendPersonalMail(ctx, userId, req.TemplateId, req.Title, req.Content,
			req.TitleArgs, req.ContentArgs, attachment, req.Sender, req.StartTime, req.ExpireTime, req.EndTime)

		if err != nil {
			glog.Errorf("send personal mail faild, user: %d\n%s", userId, err.Error())
		}
	}

	return &pb.SendPersonalMailResp{}, nil
}
