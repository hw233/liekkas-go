package common

import (
	"shared/utility/event"
	"shared/utility/servertime"
)

// 事件类型
const (
	EventTypeGraveyardHelp           = 1
	EventTypeYggdrasilMail           = 2
	EventTypeYggdrasilIntimacyChange = 3
	EventTypeMercenaryApply          = 5
)

//----------------------------------------------- 模拟经营互助参数

const (
	GraveyardHelpTypeBuild   = 1
	GraveyardHelpTypeLvUp    = 2
	GraveyardHelpTypeStageUp = 3
	GraveyardHelpTypeProduct = 4
)

func NewGraveyardHelpEvent(helpType int, buildUid int64, sec int32) *event.Event {
	return event.NewEvent(EventTypeGraveyardHelp, helpType, buildUid, sec, servertime.Now())
}

//----------------------------------------------- 世界探索异界邮箱参数

func NewYggdrasilMailEvent(fromUserName string, goodsUid int64) *event.Event {

	return event.NewEvent(EventTypeYggdrasilMail, fromUserName, goodsUid)
}

//----------------------------------------------- 世界探索亲密度变化

func NewYggdrasilIntimacyChangeEvent(userId int64, intimacy, totalIntimacy int32) *event.Event {

	return event.NewEvent(EventTypeYggdrasilIntimacyChange, userId, intimacy, totalIntimacy)
}

// -----------------------佣兵申请参数--------------
func NewMercenaryApplyEvent(applicantId int64, characterId int32) *event.Event {
	return event.NewEvent(EventTypeMercenaryApply, applicantId, characterId, servertime.Now())
}

// func NewMercenaryCharacterEvent(c *CharacterData, e []*Equipment, w []*WorldItem) *event.Event {
// 	return event.NewEvent(EventTypeMercenaryCharacter)
// }
