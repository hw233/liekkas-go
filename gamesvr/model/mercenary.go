package model

import (
	"shared/common"
	"shared/protobuf/pb"
	"shared/utility/errors"
	"shared/utility/servertime"
	"time"
)

// todo 超过统一归还时间的话，清除

const (
	MercenaryStatusNotApplied    = 0
	MercenaryStatusApplied       = 1
	MercenaryStatusCancelApplied = 2
)

const (
	MercenaryRelationFriend = 1
	MercenaryRelationGuild  = 2
	MercenaryRelationBoth   = 3
)

type Mercenaries struct {
	*DailyRefreshChecker
	Mercenaries            []*Mercenary                        `json:"mercenaries"`              // 已经借来的佣兵，最大个数三个
	MercenaryReceivedApply map[int32]map[int64]*MercenaryApply `json:"mercenary_received_apply"` // 收到的申请
	MercenarySentApply     []*MercenarySend                    `json:"mercenary_sent_apply"`     // 已经发送的申请
	MercenaryRecords       map[int32]*MercenaryRecord          `json:"mercenary_records"`        // 记录
	ExpireTime             int64                               `json:"expiretime"`
}

type Mercenary struct {
	ID         int64               `json:"id"`
	Relation   int8                `json:"relation"`
	Owner      string              `json:"owner"`
	Character  *Character          `json:"character"`
	UseCount   map[int32]int32     `json:"use_count"`
	Equipments []*common.Equipment `json:"equipments"`
	WorldItem  *common.WorldItem   `json:"worlditem"`
}

type MercenarySend struct {
	IfDelete    bool  `json:"if_delete"`
	IsCanceled  bool  `json:"is_canceled"`
	CharacterId int32 `json:"character_id"`
	SentTo      int64 `json:"sent_to"`
	ExpireTime  int64 `json:"expire_time"`
}

type MercenaryApply struct {
	CharacterId   int32  `json:"character_id"`
	Relation      int32  `json:"relation"`
	ApplicantId   int64  `json:"applicant_id"`
	ApplyTime     int64  `json:"apply_time"`
	ApplicantName string `json:"applicant_name"`
}

type MercenaryRecord struct {
	CharacterId   int32  `json:"character_id"`
	Star          int32  `json:"star"`
	Level         int32  `json:"level"`
	Power         int32  `json:"power"`
	Relation      int32  `json:"relation"`
	ApplicantId   int64  `json:"applicant_id"`
	SendTime      int64  `json:"send_time"`
	ApplicantName string `json:"applicant_name"`
}

type MercenaryUser struct {
	Uid      int64  `json:"uid"`
	UserName string `json:"user_name"`
	// Intimacy   int32
	Relation   int32                         `json:"relation"`
	Characters map[int32]*MercenaryCharacter `json:"characters"`
}

func NewMercenaries() *Mercenaries {
	return &Mercenaries{
		DailyRefreshChecker:    NewDailyRefreshChecker(),
		Mercenaries:            make([]*Mercenary, 0, 3),
		MercenaryReceivedApply: make(map[int32]map[int64]*MercenaryApply),
		MercenarySentApply:     make([]*MercenarySend, 0, 5),
		MercenaryRecords:       make(map[int32]*MercenaryRecord),
		ExpireTime:             ThisWeekRefreshTime().Add(time.Hour * 24 * 7).Unix(),
	}
}

func NewMercenary(uid int64, name string, relation int8, c *common.CharacterData) *Mercenary {
	character := NewCharacter(c.ID)
	character.Level = c.Level
	character.Power = c.Power
	character.Rarity = c.Rarity
	character.Stage = c.Stage
	character.Star = c.Star
	character.Skills = c.Skills
	if c.WorldItem != nil && c.WorldItem.ID > 0 {
		character.WorldItem = c.WorldItem.ID
	} else {
		c.WorldItem = nil
	}
	for _, equip := range c.Equipments {
		character.Equipments = append(character.Equipments, equip.ID)
	}

	return &Mercenary{
		ID:         uid,
		Relation:   relation,
		Owner:      name,
		Character:  character,
		UseCount:   map[int32]int32{},
		Equipments: c.Equipments,
		WorldItem:  c.WorldItem,
	}
}

func NewMercenarySend(userTo int64, characterId int32) *MercenarySend {
	return &MercenarySend{
		IfDelete:    false,
		IsCanceled:  false,
		CharacterId: characterId,
		SentTo:      userTo,
		ExpireTime:  0,
	}
}

func NewMercenaryApply(cid, relation int32, applicantId int64, name string, timestamp int64) *MercenaryApply {
	return &MercenaryApply{
		CharacterId:   cid,
		Relation:      relation,
		ApplicantId:   applicantId,
		ApplyTime:     timestamp,
		ApplicantName: name,
	}
}

// 还需要增加星级战力等信息
func NewMercenaryRecord(cid, relation int32, applicantId int64, name string, star, level, power int32) *MercenaryRecord {
	return &MercenaryRecord{
		CharacterId:   cid,
		Star:          star,
		Level:         level,
		Power:         power,
		Relation:      relation,
		ApplicantId:   applicantId,
		SendTime:      0,
		ApplicantName: name,
	}
}

func (ma *MercenaryApply) VOMercenaryApply() *pb.VOMercenaryApply {
	return &pb.VOMercenaryApply{
		CharacterId:   ma.CharacterId,
		Relation:      ma.Relation,
		ApplicantId:   ma.ApplicantId,
		ApplicantName: ma.ApplicantName,
		Timestamp:     ma.ApplyTime,
	}
}

func (mr *MercenaryRecord) VOMercenaryRecord() *pb.VOMercenaryRecord {
	return &pb.VOMercenaryRecord{
		CharacterId:   mr.CharacterId,
		Relation:      mr.Relation,
		Power:         mr.Power,
		Star:          mr.Star,
		Level:         mr.Level,
		ApplicantId:   mr.ApplicantId,
		ApplicantName: mr.ApplicantName,
	}
}

func NewMercenaryUser(id int64, name string, relation int32) *MercenaryUser {
	return &MercenaryUser{
		Uid:        id,
		UserName:   name,
		Relation:   relation,
		Characters: map[int32]*MercenaryCharacter{},
	}
}

func (mu *MercenaryUser) VOMercenaries() *pb.VOMercenaries {
	voCharacters := make([]*pb.VOMercenaryCharac, 0, len(mu.Characters))
	for _, c := range mu.Characters {
		voCharacters = append(voCharacters, c.VOMercenaryCharac())
	}

	return &pb.VOMercenaries{
		UserId:     mu.Uid,
		UserName:   mu.UserName,
		Relation:   mu.Relation,
		Characters: voCharacters,
	}
}

type MercenaryCharacter struct {
	Id         int32
	Level      int32
	Power      int32
	ApplyCount int32
	Star       int32
	Status     int32
	ExpireTime int64
}

func NewMercenaryCharacter(id, level, power, applyCount, start int32) *MercenaryCharacter {
	return &MercenaryCharacter{
		Id:         id,
		Level:      level,
		Power:      power,
		ApplyCount: applyCount,
		Star:       start,
		Status:     0,
		ExpireTime: 0,
	}
}

func (mc *MercenaryCharacter) VOMercenaryCharac() *pb.VOMercenaryCharac {
	// todo status 注意记录user自己发出的申请
	return &pb.VOMercenaryCharac{
		CharacterId: mc.Id,
		Power:       mc.Power,
		ApplyCount:  mc.ApplyCount,
		Star:        mc.Star,
		Status:      mc.Status,
		Level:       mc.Level,
		Timestamp:   mc.ExpireTime,
	}
}

// ---------------ReceivedApply-----------------------

// 删除其他所有玩家对于指定角色id的申请
func (m *Mercenaries) DeleteApply(id int32) {
	delete(m.MercenaryReceivedApply, id)
}

func (m *Mercenaries) DeleteSpecificApply(characterId int32, uid int64) {
	delete(m.MercenaryReceivedApply[characterId], uid)
}

// ---------------SendApply-----------------------

func (m *Mercenaries) SendApplyCountWithoutCancel() int32 {
	var result int32
	for _, send := range m.MercenarySentApply {
		if !send.IsCanceled {
			result++
		}
	}
	return result
}

func (m *Mercenaries) AddSendApply(uid int64, cid int32) error {
	for _, send := range m.MercenarySentApply {
		if send.SentTo == uid && send.CharacterId == cid {
			// fmt.Println(send.IsCanceled, send.ExpireTime, servertime.Now().Unix())
			// 重复的请求以及虽然被撤销但是还没结束冷却期的请求都会报错
			if !send.IsCanceled || send.ExpireTime >= servertime.Now().Unix() {
				return errors.Swrapf(common.ErrMercenaryRepeatedSendApply, send.SentTo, send.CharacterId)
			}
		}
	}

	apply := NewMercenarySend(uid, cid)
	m.MercenarySentApply = append(m.MercenarySentApply, apply)
	return nil
}

// 删除某个已经发送的申请
func (m *Mercenaries) DeleteSend(uid int64, cid int32) {
	for _, send := range m.MercenarySentApply {
		if send.SentTo == uid && send.CharacterId == cid {
			send.IfDelete = true
		}
	}
}

func (m *Mercenaries) UpdateMercenarySend() {
	afterFilter := make([]*MercenarySend, 0, len(m.MercenarySentApply))
	for _, send := range m.MercenarySentApply {
		if !send.IfDelete {
			newSend := NewMercenarySend(send.SentTo, send.CharacterId)
			newSend.IsCanceled = send.IsCanceled
			newSend.ExpireTime = send.ExpireTime
			afterFilter = append(afterFilter, newSend)
		}
	}
	m.MercenarySentApply = afterFilter
}

func (m *Mercenaries) DailyRefresh() {
	tm := servertime.Now()
	if m.ExpireTime < tm.Unix() {
		m = NewMercenaries()
	}
}

// ---------------Mercenary----------------

func (m *Mercenaries) CheckMercenaryOwn(characterId int32) bool {
	for _, mercenary := range m.Mercenaries {
		if characterId == mercenary.Character.ID {
			return true
		}
	}
	return false
}

// -----------------mercenary record-----------------
func (m *Mercenaries) CheckMercenaryRecord(cid int32) bool {
	_, ok := m.MercenaryRecords[cid]
	return ok
}

func (m *Mercenaries) AddMercenaryRecord(record *MercenaryRecord) {
	m.MercenaryRecords[record.CharacterId] = record
}
