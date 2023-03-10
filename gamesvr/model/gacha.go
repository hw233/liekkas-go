package model

import (
	"encoding/json"
	"gamesvr/manager"
	"shared/common"
	"shared/csv/entry"
	"shared/csv/static"
	"shared/protobuf/pb"
	"shared/utility/errors"
	"shared/utility/glog"
)

type GachaRecords struct {
	*DailyRefreshChecker
	NewPlayerDrew   bool                              `json:"new_player_drew"`   // 是否抽过新手池
	DestinyChild    int32                             `json:"new_player_ssr_id"` // 新手池抽到的ssr
	PoolRecords     map[int32]*common.GachaPoolRecord `json:"pool_records"`      // 保底规则计算用
	TypeRecords     map[int32]*common.GachaTypeRecord `json:"type_records"`      // 保底规则计算用
	WorldItemRecord *common.GachaResultRecords        `json:"world_item_record"` // 世界级道具抽卡记录
	CharacterRecord *common.GachaResultRecords        `json:"character_record"`  // 角色抽卡记录
	TotalTimes      int32                             `json:"total_times"`       // 总抽卡次数
}

func NewGachaRecords() *GachaRecords {
	return &GachaRecords{
		DailyRefreshChecker: NewDailyRefreshChecker(),
		PoolRecords:         map[int32]*common.GachaPoolRecord{},
		TypeRecords:         map[int32]*common.GachaTypeRecord{},
		WorldItemRecord:     common.NewGachaResultRecords(),
		CharacterRecord:     common.NewGachaResultRecords(),
	}
}

func (g *GachaRecords) GetByPoolIdCreate(pool *entry.GachaPool) (*common.GachaPoolRecord, *common.GachaTypeRecord) {
	record, ok := g.GetByPoolId(pool.Id)
	if !ok {
		record = common.NewGachaPoolRecord(pool.Id)
		g.PoolRecords[pool.Id] = record
	}
	typeRecord, ok := g.GetByPoolType(pool.Type)
	if !ok {
		typeRecord = common.NewGachaTypeRecord(pool.Type)
		g.TypeRecords[pool.Type] = typeRecord
	}
	return record, typeRecord

}

func (g *GachaRecords) GetByPoolId(poolId int32) (*common.GachaPoolRecord, bool) {
	record, ok := g.PoolRecords[poolId]
	return record, ok
}
func (g *GachaRecords) GetByPoolType(poolId int32) (*common.GachaTypeRecord, bool) {
	record, ok := g.TypeRecords[poolId]
	return record, ok
}

func (g *GachaRecords) Marshal() ([]byte, error) {
	return json.Marshal(g)
}

func (g *GachaRecords) Unmarshal(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	err := json.Unmarshal(data, &g)
	if err != nil {
		glog.Errorf("json.Unmarshal error: %v", errors.WrapTrace(err))
		return errors.WrapTrace(err)
	}

	return nil
}

func (g *GachaRecords) removePoolsIfOutOfDate() {
	for i := range g.PoolRecords {
		if !manager.CSV.Gacha.CheckPoolInTime(i) {
			g.removePool(i)
		}
	}
}

func (g *GachaRecords) removePool(poolId int32) {
	delete(g.PoolRecords, poolId)
}

func (g *GachaRecords) Gacha(userId int64, pool *entry.GachaPool, num int32, poolRecord *common.GachaPoolRecord, typeRecord *common.GachaTypeRecord) ([]*common.GachaReward, error) {

	rewards, err := manager.CSV.Gacha.Gacha(userId, pool, num, poolRecord, typeRecord)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	// 记录抽卡结果
	g.RecordResult(pool.Type, poolRecord.PoolId, rewards)

	return rewards, nil

}

// RecordNewPlayerPool 记录新手池一些数据
func (g *GachaRecords) RecordNewPlayerPool(rewards []*common.GachaReward) {
	g.NewPlayerDrew = true
	for _, reward := range rewards {
		item, ok := manager.CSV.Item.GetItem(reward.ID)
		if !ok {
			continue
		}
		if item.Rarity == static.RaritySsr {
			g.DestinyChild = reward.ID
		}
	}
}

func (g *GachaRecords) VOGachaInfo(pool entry.GachaPool) *pb.VOGachaInfo {

	poolRecord, _ := g.GetByPoolId(pool.Id)
	typeRecord, _ := g.GetByPoolType(pool.Type)

	var TodayGachaCount, SSRProbIncr int32
	if poolRecord != nil {
		TodayGachaCount = poolRecord.TodayGachaCount
	}
	if typeRecord != nil {
		SSRProbIncr = manager.CSV.Gacha.CalSSRProbIncr(typeRecord.SSRMissCount, typeRecord.Type)
	}
	return &pb.VOGachaInfo{
		GachaId:         pool.Id,
		TodayGachaCount: TodayGachaCount,
		SSRProbIncr:     SSRProbIncr,
	}
}

func (g *GachaRecords) RecordResult(poolType, poolId int32, rewards []*common.GachaReward) {

	switch poolType {
	case static.GachaPoolTypeNewPlayer,
		static.GachaPoolTypeCharacterCommon,
		static.GachaPoolTypeCharacterUp:
		if poolType == static.GachaPoolTypeNewPlayer {
			g.RecordNewPlayerPool(rewards)

		}
		g.CharacterRecord.Add(poolId, rewards)

	case static.GachaPoolTypeWorldItemCommon,
		static.GachaPoolTypeWorldItemUp:
		g.WorldItemRecord.Add(poolId, rewards)

	}
}

// 清除过期的抽卡记录
func (g *GachaRecords) removeRecordIfOutOfDate() {
	g.WorldItemRecord.RemoveIfOutOfDate(manager.CSV.Gacha.GetGachaRecordStoreMonth())
	g.CharacterRecord.RemoveIfOutOfDate(manager.CSV.Gacha.GetGachaRecordStoreMonth())
}

func (g *GachaRecords) RecordTotalTimes(times int32) {
	g.TotalTimes = g.TotalTimes + times
}

func (g *GachaRecords) GetTotalTimes() int32 {
	return g.TotalTimes
}
