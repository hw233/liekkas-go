package mail

import (
	"context"
	"mailserver/manager"
	"shared/common"
	"shared/utility/glog"
	"sync"
	"time"
)

var (
	MailMgr *MailManager
)

func Init() error {
	MailMgr = NewMailManager()
	return MailMgr.Init()
}

type MailManager struct {
	GroupMails      map[int64]*common.ServerGroupMail
	LastGroupMailId int64

	sync.RWMutex
}

func NewMailManager() *MailManager {
	return &MailManager{
		GroupMails: map[int64]*common.ServerGroupMail{},
	}
}

func (mm *MailManager) Init() error {
	err := mm.loadGroupMail()
	if err != nil {
		return err
	}

	return nil
}

func (mm *MailManager) SendGroupMail(ctx context.Context, users []int64, sendAll bool,
	templateId int32, title, content string, titleArgs, contentArgs []string,
	attachments *common.Rewards, sender string, startTime, expireTime, endTime int64) error {

	if !sendAll {
		if len(users) <= 0 {
			//todo error
			return common.ErrParamError
		}
	}

	id, err := genGroupMailId(ctx)
	if err != nil {
		return err
	}

	if titleArgs == nil {
		titleArgs = []string{}
	}

	if contentArgs == nil {
		contentArgs = []string{}
	}

	mail := &common.ServerGroupMail{
		ServerMail: &common.ServerMail{
			Id:          id,
			TemplateId:  templateId,
			Title:       title,
			TitleArgs:   titleArgs,
			Content:     content,
			ContentArgs: contentArgs,
			Attachment:  attachments,
			Sender:      sender,
			SendTime:    startTime,
			ExpireTime:  expireTime,
			EndTime:     endTime,
		},
		Users:   map[int64]bool{},
		SendAll: sendAll,
	}

	if !sendAll {
		for _, userId := range users {
			mail.Users[userId] = true
		}
	}

	mm.AddGroupMail(mail)

	err = DBAddGroupMail(ctx, mail)
	if err != nil {
		return err
	}

	err = NewGroupMailNotify(ctx, mail)
	if err != nil {
		return err
	}

	return nil
}

func (mm *MailManager) SendPersonalMail(ctx context.Context, userId int64, templateId int32,
	title, content string, titleArgs, contentArgs []string, attachments *common.Rewards,
	sender string, startTime, expireTime, endTime int64) error {

	id, err := genPersonalMailId(ctx)
	if err != nil {
		return err
	}

	if titleArgs == nil {
		titleArgs = []string{}
	}

	if contentArgs == nil {
		contentArgs = []string{}
	}

	mail := &common.ServerPersonalMail{
		ServerMail: &common.ServerMail{
			Id:          id,
			TemplateId:  templateId,
			Title:       title,
			TitleArgs:   titleArgs,
			Content:     content,
			ContentArgs: contentArgs,
			Attachment:  attachments,
			Sender:      sender,
			SendTime:    startTime,
			ExpireTime:  expireTime,
			EndTime:     endTime,
		},
		UserId: userId,
	}

	err = DBAddPersonalMail(ctx, mail)
	if err != nil {
		glog.Error(err.Error())
		return err
	}

	NewPersonalMailNotify(ctx, mail)

	return nil
}

func (mm *MailManager) AddGroupMail(mail *common.ServerGroupMail) {
	mm.Lock()
	defer mm.Unlock()

	mm.addGroupMail(mail)
}

func (mm *MailManager) FetchMail(ctx context.Context, userId int64) ([]*common.ServerGroupMail, []*common.ServerPersonalMail, error) {
	mm.removeExpireMails()

	lastGroupMailId, err := getLastGroupMailId(ctx)
	if err != nil {
		return nil, nil, err
	}

	if lastGroupMailId > mm.LastGroupMailId {
		err := mm.loadGroupMail()
		if err != nil {
			return nil, nil, err
		}
	}

	var groupMails []*common.ServerGroupMail

	mm.Lock()
	for _, mail := range mm.GroupMails {
		if mail.HasUser(userId) && mail.CanBeSend() {
			groupMails = append(groupMails, mail)
		}
	}
	mm.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	personalMails, err := DBLoadPersonalMail(ctx)
	if err != nil {
		return nil, nil, err
	}

	return groupMails, personalMails, nil
}

func (mm *MailManager) loadGroupMail() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	groupMails, err := DBLoadGroupMail(ctx)
	if err != nil {
		return err
	}

	mm.Lock()
	for _, mail := range groupMails {
		mm.GroupMails[mail.Id] = mail
	}
	mm.Unlock()

	return nil
}

func (mm *MailManager) addGroupMail(mail *common.ServerGroupMail) {
	mm.GroupMails[mail.Id] = mail
}

func (mm *MailManager) removeGroupMail(id int64) {
	mm.Lock()
	defer mm.Unlock()
	delete(mm.GroupMails, id)
}

func (mm *MailManager) removeExpireMails() {
	var removeList []int64
	for mailId, mail := range mm.GroupMails {
		if mail.IsExpired() {
			removeList = append(removeList, mailId)
		}
	}

	for _, mailId := range removeList {
		mm.removeGroupMail(mailId)
	}
}

func genGroupMailId(ctx context.Context) (int64, error) {
	id, err := manager.Global.GenID(ctx, "group_mail")
	if err != nil {
		return 0, err
	}

	return id, nil
}

func genPersonalMailId(ctx context.Context) (int64, error) {
	id, err := manager.Global.GenID(ctx, "personal_mail")
	if err != nil {
		return 0, err
	}

	return id, nil
}

func getLastGroupMailId(ctx context.Context) (int64, error) {
	id, err := manager.Global.GetInt64(ctx, "group_mail")
	return id, err
}
