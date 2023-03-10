package entry

import (
	"log"
	"shared/common"
	"shared/csv/static"
	"shared/utility/errors"
	"shared/utility/rand"
	"shared/utility/servertime"
	"shared/utility/transfer"
	"sync"
)

const CfgNameGachaDrop = "cfg_gacha_drop"
const CfgNameGachaGlobal = "cfg_gacha_global"
const CfgNameGachaPool = "cfg_gacha_pool"

type GachaEntry struct {
	sync.RWMutex

	gachaMultipleValue    int32
	globalConfigs         map[int32]GachaGlobalConfig
	pools                 map[int32]GachaPool
	dropsMap              map[int32][]GachaDrop
	charaShards           map[int32]CharacterCharaShardInfo
	gachaRecordShowMonth  int
	gachaRecordStoreMonth int
}

func NewGachaEntry() *GachaEntry {
	return &GachaEntry{
		gachaMultipleValue: 0,
		globalConfigs:      map[int32]GachaGlobalConfig{},
		pools:              map[int32]GachaPool{},
		dropsMap:           map[int32][]GachaDrop{},
		charaShards:        map[int32]CharacterCharaShardInfo{},
	}
}

type CharacterCharaShardInfo struct {
	Id         int32
	CharaShard *common.Reward `rule:"reward"`
}

type GachaPool struct {
	Id               int32
	Type             int32
	UnlockCondition  *common.Conditions `rule:"conditions"`
	OpenTime         int64              `rule:"stringToTimeUnix"`
	CloseTime        int64              `rule:"stringToTimeUnix"`
	DailyLimit       int32
	SingleConsume    *common.Reward `rule:"reward"`
	DropN            int32
	DropR            int32
	DropSR           int32
	DropSSR          int32
	DropSafeSSR      int32
	DropUpSSR        int32
	UpGuaranteeCount int32
}

type GachaDrop struct {
	Item     int32
	Count    int32
	Weight   int32
	Rarity   int32 `ignore:"true"`
	ItemType int32 `ignore:"true"`
}

type GachaGlobalConfig struct {
	Id                       int32
	DropWeight               []int32 `src:"DropWeightN,DropWeightR,DropWeightSR" rule:"mergeInt32"` // 依次是N,R,SR的掉落权重
	DropProbSSR              int32
	GuaranteeCountSR         int32
	GuaranteeTrigCountSSR    int32
	GuaranteeTrigProbIncrSSR int32
}

func (g *GachaEntry) Check(config *Config) error {
	return nil
}

func (g *GachaEntry) Reload(config *Config) error {
	g.Lock()
	defer g.Unlock()

	globalConfigs := map[int32]GachaGlobalConfig{}
	pools := map[int32]GachaPool{}
	dropsMap := map[int32][]GachaDrop{}
	charaShards := map[int32]CharacterCharaShardInfo{}

	for _, v := range config.CfgGachaGlobalConfig.GetAllData() {
		gachaGlobalConfig := &GachaGlobalConfig{}
		err := transfer.Transfer(v, gachaGlobalConfig)
		if err != nil {
			return errors.WrapTrace(err)
		}
		globalConfigs[v.Id] = *gachaGlobalConfig
	}
	for _, v := range config.CfgGachaPoolConfig.GetAllData() {
		gachaPool := &GachaPool{}
		err := transfer.Transfer(v, gachaPool)
		if err != nil {
			return errors.WrapTrace(err)
		}
		pools[v.Id] = *gachaPool
	}

	for _, v := range config.CfgGachaDropConfig.GetAllData() {
		gachaDrop := &GachaDrop{}
		err := transfer.Transfer(v, gachaDrop)
		if err != nil {
			return errors.WrapTrace(err)
		}
		item, ok := config.CfgItemDataConfig.Find(gachaDrop.Item)
		if !ok {
			return errors.Swrapf(common.ErrNotFoundInCSV, CfgItemData, gachaDrop.Item)
		}
		gachaDrop.Rarity = item.Rarity
		if item.ItemType != static.ItemTypeCharacter && item.ItemType != static.ItemTypeWorldItem {
			return errors.New("CfgGachaDropConfig item type err ,id:%d", v.Id)
		}

		gachaDrop.ItemType = item.ItemType

		gachaDrops := dropsMap[v.DropId]
		gachaDrops = append(gachaDrops, *gachaDrop)
		dropsMap[v.DropId] = gachaDrops
	}

	for _, v := range config.CfgCharacterConfig.GetAllData() {
		if v.Visible {
			c := &CharacterCharaShardInfo{}
			err := transfer.Transfer(v, c)
			if err != nil {
				return errors.WrapTrace(err)
			}
			charaShards[v.Id] = *c
		}

	}
	g.globalConfigs = globalConfigs
	g.pools = pools
	g.dropsMap = dropsMap
	g.gachaMultipleValue = config.GachaMultipleValue
	g.charaShards = charaShards
	g.gachaRecordShowMonth = int(config.GachaRecordShowMonth)
	g.gachaRecordStoreMonth = int(config.GachaRecordStoreMonth)

	// check------------------

	for _, pool := range pools {
		if isTypeWithUp(pool.Type) {
			if pool.DropUpSSR == 0 || pool.UpGuaranteeCount == 0 {
				// 限定池子 必须有限定的item且UpGuaranteeCount不为0
				return errors.WrapTrace(common.ErrParamError)
			}
		} else {
			if pool.DropUpSSR != 0 {
				// 非限定池子 必须没有限定的item
				return errors.WrapTrace(common.ErrParamError)
			}
		}

		// check 各掉落包稀有度
		err := CheckDropRarity(dropsMap, pool.DropN, static.RarityN, config)
		if err != nil {
			return errors.WrapTrace(err)
		}
		err = CheckDropRarity(dropsMap, pool.DropR, static.RarityR, config)
		if err != nil {
			return errors.WrapTrace(err)
		}
		err = CheckDropRarity(dropsMap, pool.DropSR, static.RaritySr, config)
		if err != nil {
			return errors.WrapTrace(err)
		}
		err = CheckDropRarity(dropsMap, pool.DropSSR, static.RaritySsr, config)
		if err != nil {
			return errors.WrapTrace(err)
		}
		err = CheckDropRarity(dropsMap, pool.DropSafeSSR, static.RaritySsr, config)
		if err != nil {
			return errors.WrapTrace(err)
		}
		err = CheckDropRarity(dropsMap, pool.DropUpSSR, static.RaritySsr, config)
		if err != nil {
			return errors.WrapTrace(err)
		}
	}
	return nil
}

func CheckDropRarity(dropsMap map[int32][]GachaDrop, dropId, rarity int32, config *Config) error {
	if dropId == 0 {
		return nil
	}
	drops, ok := dropsMap[dropId]
	if !ok {
		return errors.Swrapf(common.ErrNotFoundInCSV, CfgNameGachaDrop, dropId)
	}
	if len(drops) == 0 {
		return errors.Swrapf(common.ErrNotFoundInCSV, CfgNameGachaDrop, dropId)
	}
	for _, drop := range drops {
		if drop.Rarity != rarity {
			// 掉落稀有度不正确
			return errors.WrapTrace(common.ErrParamError)
		}
		if drop.ItemType == static.ItemTypeCharacter {
			find, ok := config.CfgCharacterConfig.Find(drop.Item)
			if !ok {
				return errors.Swrapf(common.ErrNotFoundInCSV, CfgCharacter, drop.Item)
			}
			if find.Visible == false {
				return errors.New("Gacha CfgCharacter Visible cannot be false,characterId:%d", drop.Item)
			}
			if find.Type != 1 {
				return errors.New("Gacha CfgCharacter Type error,characterId:%d", drop.Item)
			}
		} else if drop.ItemType == static.ItemTypeWorldItem {
			_, ok := config.CfgWorldItemDataConfig.Find(drop.Item)
			if !ok {
				return errors.Swrapf(common.ErrNotFoundInCSV, CfgWorldItemData, drop.Item)
			}
		}
	}
	return nil
}

func (g *GachaEntry) GetGachaGlobalConfig(Type int32) (*GachaGlobalConfig, error) {
	g.RLock()
	defer g.RUnlock()
	config, ok := g.globalConfigs[Type]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameGachaGlobal, Type)
	}
	return &config, nil
}

func (g *GachaEntry) GetGachaPool(poolId int32) (*GachaPool, error) {
	g.RLock()
	defer g.RUnlock()
	config, ok := g.pools[poolId]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameGachaPool, poolId)
	}
	return &config, nil
}

func (g *GachaEntry) GetGachaDrops(dropId int32) ([]GachaDrop, error) {
	g.RLock()
	defer g.RUnlock()
	config, ok := g.dropsMap[dropId]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameGachaDrop, dropId)
	}
	return config, nil
}

func (g *GachaEntry) GetCharaShards(character int32) (*CharacterCharaShardInfo, error) {
	g.RLock()
	defer g.RUnlock()
	config, ok := g.charaShards[character]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgCharacter, character)
	}
	return &config, nil
}

func (g *GachaEntry) CheckPoolInTime(poolId int32) bool {
	pool, err := g.GetGachaPool(poolId)
	if err != nil {
		return false
	}
	now := servertime.Now().Unix()
	return pool.CheckPoolInTime(now)
}

func (g *GachaPool) CheckPoolInTime(timeUnix int64) bool {
	if g.Type == static.GachaPoolTypeNewPlayer {
		return true
	}
	return timeUnix >= g.OpenTime && timeUnix <= g.CloseTime
}

// CalSSRProbIncr 连续n次未抽到ssr 计算ssr概率增量
func (g *GachaEntry) CalSSRProbIncr(SSRMissCount int32, poolType int32) int32 {
	config, err := g.GetGachaGlobalConfig(poolType)
	if err != nil {
		log.Printf("CalSSRProbIncr GetGachaGlobalConfig err:%+v", err)
		return 0
	}

	if SSRMissCount < config.GuaranteeTrigCountSSR {
		return 0
	} else {
		prob := (SSRMissCount + 1 - config.GuaranteeTrigCountSSR) * config.GuaranteeTrigProbIncrSSR
		if prob+config.DropProbSSR > tenThousand {
			return tenThousand - config.DropProbSSR
		}
		return prob
	}

}

// GetCurrentPools 获得当前可抽的卡池
func (g *GachaEntry) GetCurrentPools(newPlayerDrew bool) map[int32]GachaPool {
	g.RLock()
	defer g.RUnlock()
	result := map[int32]GachaPool{}
	now := servertime.Now().Unix()

	for i, pool := range g.pools {
		if pool.Type == static.GachaPoolTypeNewPlayer && newPlayerDrew {
			continue
		}
		if pool.CheckPoolInTime(now) {
			result[i] = pool
		}
	}
	return result

}

func (g *GachaEntry) GetGachaMultipleValue() int32 {
	g.RLock()
	defer g.RUnlock()
	return g.gachaMultipleValue
}

func (g *GachaEntry) GetGachaRecordShowMonth() int {
	g.RLock()
	defer g.RUnlock()
	return g.gachaRecordShowMonth
}

func (g *GachaEntry) GetGachaRecordStoreMonth() int {
	g.RLock()
	defer g.RUnlock()
	return g.gachaRecordStoreMonth
}

func (g *GachaEntry) Gacha(userId int64, pool *GachaPool, num int32, poolRecord *common.GachaPoolRecord, typeRecord *common.GachaTypeRecord) ([]*common.GachaReward, error) {
	config, err := g.GetGachaGlobalConfig(pool.Type)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	/**
	1. 保底 ssr  --|
	              |->如果是ssr 判断是否要替换保底 up (up逻辑上必须是ssr)
	2. 随机 ssr  --|
	3. 保底 sr
	4. 随机 n/r/sr
	*/
	// todo:  record remove
	dropRecord := map[int]string{}
	mapIndex := 0
	ret := make([]*common.GachaReward, 0, num)
	for i := 0; i < int(num); i++ {
		mapIndex++
		SSRProb := g.SSRProb(typeRecord.SSRMissCount, config)
		randProb := rand.RangeInt32(1, tenThousand)
		//log.Printf("num:%d,SSRProb:%d,randProb:%d", i, SSRProb, randProb)
		// 是否掉落ssr
		var SSR bool
		dropId := pool.DropN
		if SSRProb >= tenThousand {
			//保底掉落ssr
			dropId = pool.DropSafeSSR
			dropRecord[mapIndex] = "保底 ssr"
			SSR = true
		} else if SSRProb >= randProb {
			// 优先判断掉落ssr
			dropId = pool.DropSSR
			dropRecord[mapIndex] = "随机 ssr"
			SSR = true
		}
		if SSR {
			// up池子判断是否保底up了
			if isTypeWithUp(pool.Type) && poolRecord.UpMissCount >= pool.UpGuaranteeCount {
				dropId = pool.DropUpSSR
				dropRecord[mapIndex] = "保底 up ssr"
			}
		} else {
			// 未掉落ssr

			if typeRecord.SRMissCount >= config.GuaranteeCountSR {
				// 保底sr
				dropId = pool.DropSR
				dropRecord[mapIndex] = "保底 sr"

			} else {
				// 没有任何保底,从N,R,SR中随机掉落
				perm := rand.SinglePerm(config.DropWeight)
				switch perm {
				case 0:
					dropRecord[mapIndex] = "随机 n"
					dropId = pool.DropN
				case 1:
					dropRecord[mapIndex] = "随机 r"
					dropId = pool.DropR
				case 2:
					dropRecord[mapIndex] = "随机 sr"
					dropId = pool.DropSR
				}

			}
		}
		drop, err := g.GachaDrop(dropId)
		if err != nil {
			return nil, err
		}

		// 更新卡池记录
		updateGachaRecord(g.isUpItem(pool, drop.Item), drop.Rarity, poolRecord, typeRecord)

		ret = append(ret, common.NewGachaReward(common.NewReward(drop.Item, drop.Count), typeRecord.SSRMissCount))

		log.Printf("gachaRecord,userId:%d,poolId:%d,todayGachaCount:%d, ssrProb:%d,randProb:%d,dropMode:%s,itemId:%d,itemRarity:%d,isUpItem:%v,UpMissCount:%d,SSRMissCount:%d,SRMissCount:%d", userId, poolRecord.PoolId, poolRecord.TodayGachaCount, SSRProb, randProb, dropRecord[mapIndex], drop.Item, drop.Rarity, g.isUpItem(pool, drop.Item), poolRecord.UpMissCount, typeRecord.SSRMissCount, typeRecord.SRMissCount)
		//dropRecord[mapIndex] = fmt.Sprintf(dropRecord[mapIndex]+",isUpItem:%v", g.isUpItem(pool, drop.Item))

	}

	////for i := 1; i < mapIndex; i++ {
	////	log.Printf(dropRecord[i])
	////}
	//
	//result := map[string]int32{}
	//for _, s := range dropRecord {
	//	num, ok := result[s]
	//	if ok {
	//		result[s] = num + 1
	//	} else {
	//		result[s] = 1
	//	}
	//}
	//for s, i := range result {
	//	log.Printf("%s %d", s, i)
	//
	//}

	return ret, nil
}

func (g *GachaEntry) isUpItem(pool *GachaPool, itemId int32) bool {
	if pool.DropUpSSR == 0 {
		return false
	}
	drops, err := g.GetGachaDrops(pool.DropUpSSR)
	if err != nil {
		return false
	}
	for _, drop := range drops {
		if drop.Item == itemId {
			return true
		}
	}
	return false
}
func isTypeWithUp(poolType int32) bool {
	return poolType == static.GachaPoolTypeCharacterUp || poolType == static.GachaPoolTypeWorldItemUp
}

// 更新卡池记录
func updateGachaRecord(isUpItem bool, rarity int32, poolRecord *common.GachaPoolRecord, typeRecord *common.GachaTypeRecord) {
	poolRecord.Drop(rarity)
	typeRecord.Drop(rarity)
	// 更新卡池记录
	if isUpItem {
		poolRecord.ClearUpMissCount()
	}
}

func (g *GachaEntry) SSRProb(SSRMissCount int32, config *GachaGlobalConfig) int32 {

	return g.CalSSRProbIncr(SSRMissCount, config.Id) + config.DropProbSSR
}

func (g *GachaEntry) GachaDrop(dropId int32) (*GachaDrop, error) {

	gachaDrops, err := g.GetGachaDrops(dropId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	weights := make([]int32, 0, len(gachaDrops))

	for _, gachaDrop := range gachaDrops {
		weights = append(weights, gachaDrop.Weight)
	}
	return &(gachaDrops)[rand.SinglePerm(weights)], nil
}

func (g *GachaEntry) GachaConsume(singleConsume *common.Reward, num int32) *common.Rewards {
	rewards := common.NewRewards()
	rewards.AddReward(singleConsume)
	return rewards.Multiple(num)
}
