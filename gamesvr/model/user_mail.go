package model

import (
	"gamesvr/manager"
	"log"
	"shared/common"
	"shared/statistic/bilog"
	"shared/statistic/logreason"
	"shared/utility/errors"
)

func (u *User) MailBatchNotify() {
	for mailId := range u.MailInfo.Mails {
		u.AddMailNotify(mailId)
	}
}

func (u *User) MailOnLine() {
	u.FetchMail()
}

func (u *User) FetchMail() {
	groupMails, personalMails, err := u.RPCMailFetch()
	if err != nil {
		log.Printf("fetch mail failure, err: %+v\n", err)
		return
	}

	for _, voMail := range groupMails {
		if u.MailInfo.IsGroupMailRecveived(voMail.Id) {
			continue
		}

		mail, err := common.ParseFromVOServerMail(voMail)
		if err != nil {
			log.Printf("fetch mail failure, err: %+v\n", err)
			continue
		}
		u.addMailByServerMail(mail)
		u.MailInfo.RecordGroupMail(mail.Id)
	}

	for _, voMail := range personalMails {
		if u.MailInfo.IsPersonalMailRecveived(voMail.Id) {
			continue
		}

		mail, err := common.ParseFromVOServerMail(voMail)
		if err != nil {
			log.Printf("fetch mail failure, err: %+v\n", err)
			continue
		}

		u.addMailByServerMail(mail)
		u.MailInfo.RecordPersonalMail(mail.Id)
	}
}

func (u *User) AddMail(templateId int32, title string, titleArgs []string,
	content string, contentArgs []string, attachment *common.Rewards,
	sender string, sendTime, expireTime int64) (*Mail, error) {

	mail := u.MailInfo.AddMail(templateId, title, titleArgs, content, contentArgs,
		attachment, sender, sendTime, expireTime)

	if u.MailInfo.MailsCount() > manager.CSV.GlobalEntry.MailMaxCount {
		removeMail := u.MailInfo.GetEarliest()
		u.removeMail(removeMail.Id)

		u.AddMainRemoveNotify(removeMail.Id)
	}

	u.AddMailNotify(mail.Id)

	u.BIMail(mail, 0, bilog.MailOpReceive)

	return mail, nil

}

func (u *User) ReadMail(mailId int64) error {
	mail, ok := u.MailInfo.GetMail(mailId)
	if !ok {
		return errors.Swrapf(common.ErrMailNotFound, mailId)
	}

	mail.SetRead()

	u.BIMail(mail, 0, bilog.MailOpRead)

	return nil
}

func (u *User) ReceiveMailAttachment(mailId int64) error {
	mail, ok := u.MailInfo.GetMail(mailId)
	if !ok {
		return errors.Swrapf(common.ErrMailNotFound, mailId)
	}

	if !mail.HasAttachment() {
		return errors.Swrapf(common.ErrMailHasNotAttachment, mailId)
	}

	u.receiveMailAttachment(mail)

	return nil
}

func (u *User) ReceiveAllMailAttachment() []int64 {
	receiveList := []int64{}
	mails := u.MailInfo.GetMails()
	for _, mail := range mails {
		if mail.HasAttachment() && !mail.IsRemovable() {
			u.receiveMailAttachment(mail)
			receiveList = append(receiveList, mail.Id)
			u.AddMailNotify(mail.Id)
		}
	}

	return receiveList
}

func (u *User) RemoveMail(mailId int64) error {
	_, ok := u.MailInfo.GetMail(mailId)
	if !ok {
		return errors.Swrapf(common.ErrMailNotFound, mailId)
	}

	u.removeMail(mailId)

	return nil
}

func (u *User) RemoveReadMails() []int64 {
	mails := u.MailInfo.GetMails()

	removeList := []int64{}
	for _, mail := range mails {
		if mail.IsRemovable() {
			removeList = append(removeList, mail.Id)
		}
	}

	for _, mailId := range removeList {
		u.removeMail(mailId)
	}

	return removeList
}

func (u *User) RemoveExpireMails() {
	mails := u.MailInfo.GetMails()

	removeList := []int64{}
	for _, mail := range mails {
		if mail.IsExpired() {
			removeList = append(removeList, mail.Id)
		}
	}

	for _, mailId := range removeList {
		u.removeMail(mailId)
		u.AddMainRemoveNotify(mailId)
	}
}

func (u *User) AddServerGroupMail(mail *common.ServerMail) {
	u.addMailByServerMail(mail)
	u.MailInfo.RecordGroupMail(mail.Id)
}

func (u *User) AddServerPersonalMail(mail *common.ServerMail) {
	u.addMailByServerMail(mail)
	u.MailInfo.RecordPersonalMail(mail.Id)
}

func (u *User) addMailByServerMail(mail *common.ServerMail) {
	u.AddMail(mail.TemplateId, mail.Title, mail.TitleArgs, mail.Content, mail.ContentArgs,
		mail.Attachment, mail.Sender, mail.SendTime, mail.ExpireTime)
}

func (u *User) receiveMailAttachment(mail *Mail) {
	reason := logreason.NewReason(logreason.ReceiveMail)
	u.addRewards(mail.GetAttachment(), reason)
	mail.SetAttachmentReceived()

	u.BIMail(mail, 0, bilog.MailOpReceiveAttachment)

}

func (u *User) removeMail(mailId int64) {
	mail, _ := u.MailInfo.GetMail(mailId)

	u.MailInfo.RemoveMail(mailId)

	u.BIMail(mail, 0, bilog.MailOpRemove)
}
