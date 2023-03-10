package model

import (
	"context"
	"gamesvr/manager"
	"shared/common"
	"shared/protobuf/pb"
	"shared/utility/number"
	"shared/utility/servertime"
	"shared/utility/uid"
)

// YggdrasilMailBox 世界探索异界邮箱
type YggdrasilMailBox struct {
	*YggdrasilOrderedMails
	UID *uid.UID `json:"uid"` // 生成唯一ID
}

func NewYggdrasilMailBox() *YggdrasilMailBox {
	return &YggdrasilMailBox{
		UID:                   uid.NewUID(),
		YggdrasilOrderedMails: NewYggdrasilOrderedMails(),
	}
}

func (y *YggdrasilMailBox) IsMailBoxFull() bool {

	return y.Len() >= manager.CSV.Yggdrasil.GetYggMailMaxCount()
}
func (y *YggdrasilMailBox) AddOne(ctx context.Context, userId int64, fromUserName string, Attachment *common.Rewards) error {

	return y.Put(ctx, userId, NewYggdrasilMail(y.UID.Gen(), fromUserName, Attachment))
}

type YggdrasilMail struct {
	Uid int64 `json:"uid"`
	//FromUserId   int64           `json:"from_user_id"`
	FromUserName string          `json:"from_user_name"`
	Attachment   *common.Rewards `json:"attachment"`
	SendTime     int64           `json:"send_time"`
}

func (m *YggdrasilMail) VOYggdrasilMail() *pb.VOYggdrasilMail {
	return &pb.VOYggdrasilMail{
		Uid:          m.Uid,
		FromUserName: m.FromUserName,
		Attachment:   m.Attachment.MergeVOResource(),
		SendTime:     m.SendTime,
	}
}

func NewYggdrasilMail(uid int64, FromUserName string, Attachment *common.Rewards) *YggdrasilMail {
	return &YggdrasilMail{
		Uid:          uid,
		FromUserName: FromUserName,
		Attachment:   Attachment,
		SendTime:     servertime.Now().Unix(),
	}
}

type YggdrasilOrderedMails struct {
	Keys *number.SortedInt64sSet  `json:"keys"`
	M    map[int64]*YggdrasilMail `json:"m"`
}

func NewYggdrasilOrderedMails() *YggdrasilOrderedMails {
	return &YggdrasilOrderedMails{
		Keys: number.NewSortedInt64sSet(),
		M:    map[int64]*YggdrasilMail{},
	}
}

func (y *YggdrasilOrderedMails) Len() int32 {
	return int32(len(y.M))
}

func (y *YggdrasilOrderedMails) Put(ctx context.Context, userId int64, v *YggdrasilMail) error {
	y.Keys.Add(v.Uid)
	y.M[v.Uid] = v
	return SetYggMailNumInRedis(ctx, userId, len(y.M))

}
func (y *YggdrasilOrderedMails) Get(k int64) (*YggdrasilMail, bool) {
	v, ok := y.M[k]
	return v, ok
}

func (y *YggdrasilOrderedMails) Delete(ctx context.Context, userId int64, k int64) error {
	_, ok := y.Get(k)
	if ok {
		delete(y.M, k)
		y.Keys.Delete(k)
	}
	return SetYggMailNumInRedis(ctx, userId, len(y.M))

}

// PagingSearch 返回key <offset 的n个数据
func (y *YggdrasilOrderedMails) PagingSearch(offset int64, n int) []*YggdrasilMail {

	ids := y.Keys.PagingSearch(offset, n)
	values := make([]*YggdrasilMail, 0, len(ids))
	for _, id := range ids {
		values = append(values, y.M[id])
	}
	return values

}

func (y *YggdrasilOrderedMails) GetAll() map[int64]*YggdrasilMail {

	return y.M

}
