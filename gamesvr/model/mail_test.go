package model

import (
	"context"
	"fmt"
	"testing"
)

func TestAddMail(t *testing.T) {
	user := NewUser(1001)
	user.InitForCreate(context.Background())

	mail, err := addMail(user, 0)
	if err != nil {
		t.Errorf("failed, %s", err.Error())
		return
	}

	fmt.Printf("mail: %v\n", mail)
}

func TestReceiveMailAttachment(t *testing.T) {
	user := NewUser(1001)
	user.InitForCreate(context.Background())

	user.RewardsResult.Clear()
	mail, _ := addMail(user, 50001)

	err := user.ReceiveMailAttachment(mail.Id)
	if err != nil {
		t.Errorf("failed, %s", err.Error())
		return
	}

	fmt.Printf("mail: %v\n", mail)
	fmt.Printf("RewardsResult: %v\n", user.RewardsResult.VOResourceResult())
}

func TestRemoveMail(t *testing.T) {
	user := NewUser(1001)
	user.InitForCreate(context.Background())

	mail, _ := addMail(user, 50001)

	fmt.Printf("beforce remove mail count: %v\n", user.MailInfo.MailsCount())

	err := user.RemoveMail(mail.Id)
	if err != nil {
		t.Errorf("failed, %s", err.Error())
		return
	}

	fmt.Printf("after remove mail count: %v\n", user.MailInfo.MailsCount())
}

func addMail(user *User, dropId int32) (*Mail, error) {
	//var rewards *common.Rewards = nil
	//var err error
	//if dropId != 0 {
	//	rewards, err = manager.CSV.Drop.DropRewards(dropId)
	//	if err != nil {
	//		return nil, err
	//	}
	//}
	//
	//templateCfg, err := manager.CSV.Mail.GetTemplate(entry.MailTemplateGM)
	//if err != nil {
	//	return nil, err
	//}
	//
	//title := fmt.Sprintf(templateCfg.Title)
	//content := fmt.Sprintf(templateCfg.Content, user.GetUserId())
	//sendTime := servertime.Now().Unix()
	//expireTime := sendTime + int64(templateCfg.ExpireDays)*servertime.SecondPerDay
	//mail, err := user.AddMail(entry.MailTemplateGM, title, nil, content, ,ninil, rewards, sendTime, expireTime)
	//
	//return mail, err
	return nil, nil
}
