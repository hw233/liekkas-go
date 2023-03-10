package model

import (
	"shared/common"
	"shared/csv/static"
	"shared/protobuf/pb"
	"shared/utility/servertime"
)

type Mail struct {
	Id          int64           `json:"id"`
	TemplateId  int32           `json:"template_id"`
	State       int32           `json:"state"`
	Title       string          `json:"title"`
	TitleArgs   []string        `json:"title_args"`
	Content     string          `json:"content"`
	ContentArgs []string        `json:"content_args"`
	Attachment  *common.Rewards `json:"attachment"`
	Sender      string          `json:"sender"`
	SendTime    int64           `json:"send_time"`
	ExpireTime  int64           `json:"expire_time"`
}

type MailInfo struct {
	Mails             map[int64]*Mail `json:"mails"`
	RecvGroupMails    map[int64]bool  `json:"recv_group_mails"`
	RecvPersonalMails map[int64]bool  `json:"recv_personal_mails"`
	AllocatedId       int64           `json:"allocated_id"`
}

func NewMail(id int64, templateId int32, title string, titleArgs []string,
	content string, contentArgs []string, attachment *common.Rewards,
	sender string, sendTime, expireTime int64) *Mail {
	return &Mail{
		Id:          id,
		TemplateId:  templateId,
		State:       static.MailStateUnread,
		Title:       title,
		TitleArgs:   append([]string{}, titleArgs...),
		Content:     content,
		ContentArgs: append([]string{}, contentArgs...),
		Attachment:  attachment,
		Sender:      sender,
		SendTime:    sendTime,
		ExpireTime:  expireTime,
	}
}

func NewMailInfo() *MailInfo {
	return &MailInfo{
		Mails:             map[int64]*Mail{},
		RecvGroupMails:    map[int64]bool{},
		RecvPersonalMails: map[int64]bool{},
		AllocatedId:       0,
	}
}

//----------------------------------------
//MailInfo
//----------------------------------------
func (mi *MailInfo) AddMail(templateId int32, title string, titleArgs []string,
	content string, contentArgs []string, attachment *common.Rewards,
	sender string, sendTime, expireTime int64) *Mail {
	mailId := mi.genMailId()
	mail := NewMail(mailId, templateId, title, titleArgs, content, contentArgs,
		attachment, sender, sendTime, expireTime)

	mi.Mails[mailId] = mail

	return mail
}

func (mi *MailInfo) GetMail(id int64) (*Mail, bool) {
	mail, ok := mi.Mails[id]

	return mail, ok
}

func (mi *MailInfo) GetMails() map[int64]*Mail {
	return mi.Mails
}

func (mi *MailInfo) MailsCount() int32 {
	return int32(len(mi.Mails))
}

func (mi *MailInfo) GetEarliest() *Mail {
	var earliestMail *Mail = nil

	for _, mail := range mi.Mails {
		if earliestMail == nil {
			earliestMail = mail
			continue
		}

		if earliestMail.SendTime > mail.SendTime {
			earliestMail = mail
		}
	}

	return earliestMail
}

func (mi *MailInfo) RemoveMail(id int64) {
	delete(mi.Mails, id)
}

func (mi *MailInfo) RecordGroupMail(id int64) {
	mi.RecvGroupMails[id] = true
}

func (mi *MailInfo) IsGroupMailRecveived(id int64) bool {
	_, ok := mi.RecvGroupMails[id]
	return ok
}

func (mi *MailInfo) RecordPersonalMail(id int64) {
	mi.RecvPersonalMails[id] = true
}

func (mi *MailInfo) IsPersonalMailRecveived(id int64) bool {
	_, ok := mi.RecvPersonalMails[id]
	return ok
}

func (mi *MailInfo) VOMailInfo() []*pb.VOMail {
	mails := make([]*pb.VOMail, 0, len(mi.Mails))

	for _, mail := range mi.Mails {
		mails = append(mails, mail.VOMail())
	}

	return mails
}

func (mi *MailInfo) genMailId() int64 {
	mi.AllocatedId = mi.AllocatedId + 1
	return mi.AllocatedId
}

//----------------------------------------
//Mail
//----------------------------------------
func (m *Mail) SetRead() {
	m.State = static.MailStateRead
}

func (m *Mail) SetAttachmentReceived() {
	m.State = static.MailStateReceived
}

func (m *Mail) HasAttachment() bool {
	return m.Attachment != nil
}

func (m *Mail) GetAttachment() *common.Rewards {
	return m.Attachment
}

func (m *Mail) IsExpired() bool {
	return servertime.Now().Unix() > m.ExpireTime
}

func (m *Mail) IsRemovable() bool {
	if m.Attachment == nil {
		return m.State == static.MailStateRead
	} else {
		return m.State == static.MailStateReceived
	}
}

func (m *Mail) VOMail() *pb.VOMail {
	voMail := &pb.VOMail{
		Id:          m.Id,
		TemplateId:  m.TemplateId,
		State:       m.State,
		Title:       m.Title,
		TitleArgs:   m.TitleArgs,
		ContentArgs: m.ContentArgs,
		Content:     m.Content,
		SendTime:    m.SendTime,
		ExpireTime:  m.ExpireTime,
	}

	if m.Attachment != nil {
		voMail.Attachment = m.Attachment.MergeVOResource()
	}

	return voMail
}
