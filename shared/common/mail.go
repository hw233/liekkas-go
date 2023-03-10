package common

import (
	"shared/protobuf/pb"
	"shared/utility/servertime"
)

type ServerMail struct {
	Id          int64    `json:"id"`
	TemplateId  int32    `json:"template_id"`
	Title       string   `json:"title"`
	TitleArgs   []string `json:"title_args"`
	Content     string   `json:"content"`
	ContentArgs []string `json:"content_args"`
	Attachment  *Rewards `json:"attachment"`
	Sender      string   `json:"sender"`
	SendTime    int64    `json:"send_time"`
	ExpireTime  int64    `json:"expire_time"`
	EndTime     int64    `json:"end_time"`
}

type ServerPersonalMail struct {
	*ServerMail
	UserId int64 `json:"user_id"`
}

type ServerGroupMail struct {
	*ServerMail
	Users   map[int64]bool `json:"users"`
	SendAll bool           `json:"send_all"`
}

//----------------------------------------
//ServerMail
//----------------------------------------
func ParseFromVOServerMail(voMail *pb.VOServerMail) (*ServerMail, error) {
	mail := &ServerMail{
		Id:          voMail.Id,
		TemplateId:  voMail.TemplateId,
		Title:       voMail.Title,
		TitleArgs:   voMail.TitleArgs,
		Content:     voMail.Content,
		ContentArgs: voMail.ContentArgs,
		Sender:      voMail.Sender,
		SendTime:    voMail.SendTime,
		ExpireTime:  voMail.ExpireTime,
		EndTime:     voMail.EndTime,
	}

	attachment, err := ParseFromVOConsume(voMail.Attachment)
	if err != nil {
		return nil, err
	}

	mail.Attachment = attachment

	return mail, nil
}

func (sm *ServerMail) CanBeSend() bool {
	if sm.SendTime > servertime.Now().Unix() {
		return false
	}

	return !sm.IsExpired()
}

func (sm *ServerMail) IsExpired() bool {
	now := servertime.Now().Unix()
	if now >= sm.ExpireTime {
		return true
	}

	if now >= sm.EndTime {
		return true
	}

	return false
}

func (m *ServerMail) VOServerMail() *pb.VOServerMail {
	voMail := &pb.VOServerMail{
		Id:          m.Id,
		TemplateId:  m.TemplateId,
		Title:       m.Title,
		TitleArgs:   m.TitleArgs,
		Content:     m.Content,
		ContentArgs: m.ContentArgs,
		Sender:      m.Sender,
		SendTime:    m.SendTime,
		ExpireTime:  m.ExpireTime,
		EndTime:     m.EndTime,
	}

	if m.Attachment != nil {
		voMail.Attachment = m.Attachment.MergeVOResource()
	}

	return voMail
}

//----------------------------------------
//ServerGroupMail
//----------------------------------------
func (sgm *ServerGroupMail) HasUser(userId int64) bool {
	if sgm.SendAll {
		return true
	}

	_, has := sgm.Users[userId]
	return has
}

func (sgm *ServerGroupMail) VOServerGroupMail() *pb.VOServerGroupMail {
	voMail := &pb.VOServerGroupMail{
		Mail:    sgm.VOServerMail(),
		SendAll: sgm.SendAll,
		Users:   make([]int64, 0, len(sgm.Users)),
	}

	for userId := range sgm.Users {
		voMail.Users = append(voMail.Users, userId)
	}

	return voMail
}

//----------------------------------------
//ServerPersonalMail
//----------------------------------------
func (spm *ServerPersonalMail) VOServerPersonalMail() *pb.VOServerPersonalMail {
	voMail := &pb.VOServerPersonalMail{
		Mail: spm.VOServerMail(),
		User: spm.UserId,
	}

	return voMail
}
