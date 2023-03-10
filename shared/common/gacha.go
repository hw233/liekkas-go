package common

import (
	"math"
	"shared/csv/static"
	"shared/protobuf/pb"
	"shared/utility/number"
	"shared/utility/servertime"
	"shared/utility/uid"
)

type GachaPoolRecord struct {
	PoolId          int32 `json:"pool_id"`
	TodayGachaCount int32 `json:"today_gacha_count"`
	UpMissCount     int32 `json:"up_miss_count"` // x次抽到ssr但未抽到当期up
}

func NewGachaPoolRecord(poolId int32) *GachaPoolRecord {
	return &GachaPoolRecord{
		PoolId:          poolId,
		TodayGachaCount: 0,
		UpMissCount:     0,
	}
}

func (g *GachaPoolRecord) Drop(rarity int32) {
	switch rarity {
	case static.RarityN:
	case static.RarityR:
	case static.RaritySr:
	case static.RaritySsr:
		g.UpMissCount++
	}
	g.TodayGachaCount++

}
func (g *GachaPoolRecord) ClearUpMissCount() {
	g.UpMissCount = 0
}

type GachaTypeRecord struct {
	Type         int32 `json:"type"`
	SRMissCount  int32 `json:"sr_miss_count"`  // 连续x次未抽到sr
	SSRMissCount int32 `json:"ssr_miss_count"` // 连续x次未抽到ssr
}

func NewGachaTypeRecord(Type int32) *GachaTypeRecord {
	return &GachaTypeRecord{
		Type:         Type,
		SRMissCount:  0,
		SSRMissCount: 0,
	}
}

func (g *GachaTypeRecord) Drop(rarity int32) {
	switch rarity {
	case static.RarityN,
		static.RarityR:
		g.SRMissCount++
		g.SSRMissCount++
	case static.RaritySr:
		g.SSRMissCount++
		g.SRMissCount = 0

	case static.RaritySsr:
		g.SRMissCount = 0
		g.SSRMissCount = 0
	}

}

type GachaResultRecords struct {
	Uid     *uid.UID                     `json:"uid"`
	Keys    *number.SortedInt64sSet      `json:"keys"`
	Records map[int64]*GachaResultRecord `json:"records"`
}

func NewGachaResultRecords() *GachaResultRecords {
	return &GachaResultRecords{
		Uid:     uid.NewUID(),
		Keys:    number.NewSortedInt64sSet(),
		Records: map[int64]*GachaResultRecord{},
	}
}

func (g *GachaResultRecords) PagingSearch(offset int64, n int) []*GachaResultRecord {
	ids := g.Keys.PagingSearch(offset, n)
	values := make([]*GachaResultRecord, 0, len(ids))
	for _, id := range ids {
		values = append(values, g.Records[id])
	}
	return values
}

func (g *GachaResultRecords) Put(v *GachaResultRecord) {
	g.Delete(v.Uid)
	g.Keys.Add(v.Uid)
	g.Records[v.Uid] = v

}
func (g *GachaResultRecords) Delete(k int64) {
	delete(g.Records, k)
	g.Keys.Delete(k)
}

func (g *GachaResultRecords) Add(poolId int32, rewards []*GachaReward) {
	for _, reward := range rewards {
		g.Put(NewGachaResultRecord(g.Uid.Gen(), poolId, reward.ID))
	}
}

func (g *GachaResultRecords) RemoveIfOutOfDate(month int) {
	time := servertime.Now().AddDate(0, -month, 0).Unix()

	search := g.PagingSearch(math.MaxInt64, math.MaxInt64)
	for i := len(search) - 1; i >= 0; i-- {
		record := search[i]
		if record.CreateAt <= time {
			g.Delete(record.Uid)
		} else {
			break
		}
	}

}

type GachaResultRecord struct {
	Uid      int64 `json:"uid"`
	PoolId   int32 `json:"pool_id"`
	CreateAt int64 `json:"create_at"`
	ItemId   int32 `json:"item_id"`
}

func (g *GachaResultRecord) VOGachaResultRecord() *pb.VOGachaResultRecord {
	return &pb.VOGachaResultRecord{
		Uid:      g.Uid,
		PoolId:   g.PoolId,
		ItemId:   g.ItemId,
		CreateAt: g.CreateAt,
	}
}
func NewGachaResultRecord(uid int64, poolId, itemId int32) *GachaResultRecord {
	return &GachaResultRecord{
		Uid:      uid,
		PoolId:   poolId,
		ItemId:   itemId,
		CreateAt: servertime.Now().Unix(),
	}
}

type GachaReward struct {
	*Reward
	SSRMissCount int32
}

func NewGachaReward(reward *Reward, SSRMissCount int32) *GachaReward {
	return &GachaReward{
		Reward:       reward,
		SSRMissCount: SSRMissCount,
	}
}
