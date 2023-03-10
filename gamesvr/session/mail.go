package session

import (
	"context"
	"gamesvr/manager"
	"shared/protobuf/pb"
)

func (s *Session) MailInfo(ctx context.Context, req *pb.C2SMailInfo) (*pb.S2CMailInfo, error) {

	s.User.MailBatchNotify()

	return &pb.S2CMailInfo{
		// Mails: s.User.MailInfo.VOMailInfo(),
	}, nil
}

func (s *Session) MailMarkRead(ctx context.Context, req *pb.C2SMailMarkRead) (*pb.S2CMailMarkRead, error) {
	err := s.User.ReadMail(req.MailId)

	if err != nil {
		return nil, err
	}

	mail, _ := s.User.MailInfo.GetMail(req.MailId)

	return &pb.S2CMailMarkRead{
		Mail: mail.VOMail(),
	}, nil
}

func (s *Session) MailReceiveAttachment(ctx context.Context, req *pb.C2SMailReceiveAttachment) (*pb.S2CMailReceiveAttachment, error) {
	err := s.User.ReceiveMailAttachment(req.MailId)

	if err != nil {
		return nil, err
	}

	mail, _ := s.User.MailInfo.GetMail(req.MailId)

	return &pb.S2CMailReceiveAttachment{
		Mail:           mail.VOMail(),
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

func (s *Session) MailRemove(ctx context.Context, req *pb.C2SMailRemove) (*pb.S2CMailRemove, error) {
	err := s.User.RemoveMail(req.MailId)

	if err != nil {
		return nil, err
	}

	return &pb.S2CMailRemove{
		MailId: req.MailId,
	}, nil
}

func (s *Session) MailReceiveAll(ctx context.Context, req *pb.C2SMailReceiveAll) (*pb.S2CMailReceiveAll, error) {
	s.User.ReceiveAllMailAttachment()

	return &pb.S2CMailReceiveAll{
		Mails:          []*pb.VOMail{},
		ResourceResult: s.User.VOMergedResourceResult(),
	}, nil
}

func (s *Session) MailRemoveReadAndReceived(ctx context.Context, req *pb.C2SMailRemoveReadAndReceived) (*pb.S2CMailRemoveReadAndReceived, error) {
	mailIds := s.User.RemoveReadMails()

	return &pb.S2CMailRemoveReadAndReceived{
		RemoveMails: mailIds,
	}, nil
}

//----------------------------------------
//Notify
//----------------------------------------
func (s *Session) TryPushMailNotify() {
	notifies := s.User.PopMailNotifies()

	for _, notify := range notifies {
		s.push(manager.CSV.Protocol.Pushes.MailNotify, notify)
	}
}
