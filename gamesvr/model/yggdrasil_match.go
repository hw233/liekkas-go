package model

import (
	"context"
	"encoding/json"
	"gamesvr/manager"
	"math"
	"shared/common"
	"shared/csv/static"
	"shared/utility/errors"
	"shared/utility/glog"
	"shared/utility/rand"
	"sort"
)

type YggMatchEnablePairs []*YggMatchEnablePair

func NewYggMatchEnablePairs() *YggMatchEnablePairs {
	return (*YggMatchEnablePairs)(&[]*YggMatchEnablePair{})
}

func (y *YggMatchEnablePairs) Append(userId int64, entities ...YggMatchEnable) {
	for _, entity := range entities {
		*y = append(*y, NewYggMatchEnablePair(userId, entity))
	}
}

// Shuffle 打乱顺序
func (y *YggMatchEnablePairs) Shuffle() {
	for i := len(*y) - 1; i > 0; i-- {
		num := rand.RangeInt32(0, int32(i))
		(*y)[i], (*y)[num] = (*y)[num], (*y)[i]
	}
}

type YggMatchEnablePair struct {
	userId int64
	entity YggMatchEnable
}

func NewYggMatchEnablePair(userId int64, entity YggMatchEnable) *YggMatchEnablePair {
	return &YggMatchEnablePair{
		userId: userId,
		entity: entity,
	}
}

type MatchPool map[int64]*UserMatchEntities

func NewMatchPool() *MatchPool {
	return (*MatchPool)(&map[int64]*UserMatchEntities{})
}

type UserMatchEntities struct {
	YggdrasilMailNum int32                      `json:"mail_num"`      // 邮件数量满了不进入匹配池
	Builds           map[int64]*YggBuild        `json:"Builds"`        // 建筑
	Messages         map[int64]*YggMessage      `json:"Messages"`      // 留言
	DiscardGoods     map[int64]*YggDiscardGoods `json:"discard_goods"` // 丢弃物品
}

func NewUserMatchEntities(allBuild, allMessage, allDiscardGoods map[string]string, mailNum int) (*UserMatchEntities, error) {
	Builds := map[int64]*YggBuild{}
	Messages := map[int64]*YggMessage{}
	DiscardGoods := map[int64]*YggDiscardGoods{}

	for _, s := range allBuild {

		build := &YggBuild{}
		err := json.Unmarshal([]byte(s), build)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		Builds[build.Uid] = build
	}

	for _, s := range allMessage {
		message := &YggMessage{}
		err := json.Unmarshal([]byte(s), message)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		Messages[message.Uid] = message

	}

	for _, s := range allDiscardGoods {
		discardGoods := &YggDiscardGoods{}
		err := json.Unmarshal([]byte(s), discardGoods)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		DiscardGoods[discardGoods.Uid] = discardGoods

	}

	return &UserMatchEntities{
		Builds:           Builds,
		Messages:         Messages,
		DiscardGoods:     DiscardGoods,
		YggdrasilMailNum: int32(mailNum),
	}, nil
}

func (u *UserMatchEntities) IsMailBoxFull() bool {
	return u.YggdrasilMailNum >= manager.CSV.Yggdrasil.GetYggMailMaxCount()
}

func (m *MatchPool) filterBuilds(areaId int32, area *common.Area, buildType int32) map[int64][]YggMatchEnable {
	ret := map[int64][]YggMatchEnable{}
	for userId, userMatchEntities := range *m {
		for _, build := range userMatchEntities.Builds {
			if build.AreaId != areaId {
				continue
			}
			if !area.Contains(*build.Position) {
				continue
			}
			cfg, err := manager.CSV.Yggdrasil.GetYggdrasilBuilding(build.BuildId)
			if err != nil {
				glog.Errorf("BuildTypeFilter filter GetYggdrasilBuilding err:%+v", err)
				continue
			}
			if cfg.BuildingType != buildType {
				continue
			}
			ret[userId] = append(ret[userId], build)
		}
	}

	return ret
}

func (m *MatchPool) filterDiscardGoods(areaId int32) map[int64][]YggMatchEnable {
	ret := map[int64][]YggMatchEnable{}
	for userId, userMatchEntities := range *m {
		// 邮件数量满了不进入匹配池
		if userMatchEntities.IsMailBoxFull() {
			continue
		}
		for _, DiscardGoods := range userMatchEntities.DiscardGoods {
			if DiscardGoods.AreaId != areaId {
				continue
			}
			ret[userId] = append(ret[userId], DiscardGoods)
		}
	}
	return ret

}

func (m *MatchPool) filterMessages(areaId int32) map[int64][]YggMatchEnable {
	ret := map[int64][]YggMatchEnable{}
	for userId, userMatchEntities := range *m {
		for _, message := range userMatchEntities.Messages {
			if message.AreaId != areaId {
				continue
			}
			ret[userId] = append(ret[userId], message)
		}
	}
	return ret
}

func (m *MatchPool) contains(uid int64) bool {
	for _, matchEntities := range *m {
		_, ok := matchEntities.Builds[uid]
		if ok {
			return true
		}
		_, ok = matchEntities.Messages[uid]
		if ok {
			return true
		}
		_, ok = matchEntities.DiscardGoods[uid]
		if ok {
			return true
		}
	}
	return false
}

var buildTypes []int32

func init() {
	// todo:以后新增建筑type要在这里加
	buildTypes = append(buildTypes, static.YggBuildingTypeApSpring, static.YggBuildingTypeHpSpring, static.YggBuildingTypeMagicPeeping, static.YggBuildingTypeMagicTransport,
		static.YggBuildingTypeMagicBuff, static.YggBuildingTypeStepladder)

}

type UserIdIntimacyPair struct {
	userId   int64
	intimacy int32
}

// GenUserIdSortByIntimacy 按亲密度排序 且上次匹配的人在最后
func (y *YggEntities) GenUserIdSortByIntimacy(ctx context.Context, guildId, ownUserId int64, pool *MatchPool) []int64 {
	// 按亲密度从大到小排序 且上次匹配的人在最后

	var pairs []*UserIdIntimacyPair
	for userId := range *pool {
		if y.LastMatchDiscardGoodsUsers.Contains(userId) {
			continue
		}
		intimacy, err := manager.Global.GetIntimacy(ctx, guildId, userId, ownUserId)
		if err != nil {
			glog.Errorf("YggEntities GenUserIdSortByIntimacy GetIntimacy err:%+v", err)
		}

		pairs = append(pairs, &UserIdIntimacyPair{userId: userId, intimacy: intimacy})
	}
	// 按亲密度从大到小排序
	less := func(i, j int) bool {
		return pairs[i].intimacy > pairs[j].intimacy
	}
	sort.Slice(pairs, less)
	// 排序的userid
	var sortedUserIds []int64
	for _, pair := range pairs {
		sortedUserIds = append(sortedUserIds, pair.userId)
	}

	// 上次匹配的人在最后
	for _, userId := range y.LastMatchDiscardGoodsUsers.ReverseValues() {
		if _, ok := (*pool)[userId]; !ok {
			continue
		}
		sortedUserIds = append(sortedUserIds, userId)

	}
	return sortedUserIds
}

// 处理上次匹配残余得数据
func (y *YggEntities) handleLastMatch(pool *MatchPool) {
	// 移除上次上次掉落物匹配的用户
	y.LastMatchDiscardGoodsUsers.Clear()

	for _, build := range y.Builds {
		if build.IsOwn() {
			continue
		}
		// 源建筑已经被删除,或者玩家已经退公会
		if !pool.contains(build.Uid) {
			y.removeBuild(build)
		}
	}
	for _, message := range y.Messages {
		if message.IsOwn() {
			continue
		}
		// 源留言已经被删除,或者玩家已经退公会
		if !pool.contains(message.Uid) {
			y.removeMessage(message)
		}
	}
	for _, DiscardGoods := range y.DiscardGoods {
		if DiscardGoods.IsOwn() {
			continue
		}
		// 删除上次匹配的丢弃物
		y.removeDiscardGoods(DiscardGoods)
	}
}

// Match 匹配
func (y *YggEntities) Match(ctx context.Context, yggdrasil *Yggdrasil, guildId, userId int64, pool *MatchPool) {
	// 按亲密度排序 且上次匹配的人在最后
	sortedUserIds := y.GenUserIdSortByIntimacy(ctx, guildId, userId, pool)
	// 处理上次匹配残余得数据
	y.handleLastMatch(pool)
	// 按区域匹配
	for areaId, config := range manager.CSV.Yggdrasil.GetAllAreaConfig() {
		area, ok := yggdrasil.Areas.get(areaId)
		// 区域还没初始化，不匹配
		if !ok {
			continue
		}
		// 区域任务还没完成，不匹配
		if !area.IsTaskDone {
			continue
		}
		areaTotal := config.Area
		areaUnlocked := areaTotal.CutArea(yggdrasil.UnlockArea)
		areaLocked := areaTotal.MinusArea(areaUnlocked)

		// 匹配
		for _, buildType := range buildTypes {
			y.commonMatch(NewBuildMatchHandler(areaId, areaUnlocked, buildType), pool.filterBuilds(areaId, areaUnlocked, buildType), yggdrasil)
			y.commonMatch(NewBuildMatchHandler(areaId, areaLocked, buildType), pool.filterBuilds(areaId, areaLocked, buildType), yggdrasil)
		}
		y.discardGoodsMatch(NewDiscardGoodsMatchHandler(areaId), pool.filterDiscardGoods(areaId), sortedUserIds, yggdrasil)
		y.commonMatch(NewMessageMatchHandler(areaId), pool.filterMessages(areaId), yggdrasil)

	}
}

//
func (y *YggEntities) commonMatch(handler matchHandler, userEntities map[int64][]YggMatchEnable, yggdrasil *Yggdrasil) {
	limit := handler.getMatchLimit()
	nowMatchCount := handler.getNowMatchCount(y)
	needMatchCount := limit - nowMatchCount
	if needMatchCount <= 0 {
		return
	}
	pairs := NewYggMatchEnablePairs()
	for userId, entities := range userEntities {
		pairs.Append(userId, entities...)
	}
	// 打乱顺序
	pairs.Shuffle()
	y.matchByPairs(pairs, needMatchCount, yggdrasil)

}

// 物品匹配走抢红包算法
func (y *YggEntities) discardGoodsMatch(handler matchHandler, userEntities map[int64][]YggMatchEnable, sortedUserIds []int64, yggdrasil *Yggdrasil) {
	limit := handler.getMatchLimit()
	nowMatchCount := handler.getNowMatchCount(y)
	needMatchCount := limit - nowMatchCount
	glog.Debugf("discardGoodsMatch before needMatchCount %d", needMatchCount)
	if needMatchCount <= 0 {
		return
	}
	first := true
	for _, userId := range sortedUserIds {
		entities, ok := userEntities[userId]
		if !ok {
			continue
		}
		// 抢红包算法
		grabNum := rand.GrabRedPacket(needMatchCount, manager.CSV.Yggdrasil.GetYggMatchGrabAlgorithmParam(), first)
		glog.Debugf("discardGoodsMatch during userid:%d needMatchCount %d,isFirst %v,grabNum:%d,matchEnableNum:%d", userId, needMatchCount, first, grabNum, len(entities))
		if grabNum <= 0 {
			continue
		}
		pairs := NewYggMatchEnablePairs()
		pairs.Append(userId, entities...)
		// 打乱顺序
		pairs.Shuffle()
		matchCount := y.matchByPairs(pairs, grabNum, yggdrasil)

		if matchCount > 0 {
			// 第一个实际匹配到丢弃物的人
			first = false
			// 加到上次匹配列表中
			y.LastMatchDiscardGoodsUsers.Append(userId)
		}
		needMatchCount -= matchCount
		if needMatchCount <= 0 {
			break
		}

	}
}

// 返回匹配成功个数
func (y *YggEntities) matchByPairs(entities *YggMatchEnablePairs, needMatchCount int32, yggdrasil *Yggdrasil) int32 {
	var ret int32 = 0
	for _, pair := range *entities {
		err := pair.entity.CanMatch(y)
		if err != nil {
			glog.Debugf("Match pass cause CanMatch err:%+v", err)
			continue
		}

		newYggMatchEnable, err := NewYggMatchEnable(pair.userId, pair.entity)
		if err != nil {
			glog.Debugf("Match pass cause NewYggMatchEnable err:%+v", err)
			continue
		}
		needMatchCount--
		ret++
		y.AppendMatchEntity(newYggMatchEnable, yggdrasil)
		pair.entity.MatchLog()
		if needMatchCount <= 0 {
			break
		}
	}
	return ret
}
func NewYggMatchEnable(matchUserId int64, entity YggMatchEnable) (YggMatchEnable, error) {
	switch entity.(type) {
	case *YggMessage:
		return NewMatchYggMessage(matchUserId, entity.(*YggMessage))
	case *YggDiscardGoods:
		return NewMatchYggDiscardGoods(matchUserId, entity.(*YggDiscardGoods))
	case *YggBuild:
		return NewMatchYggBuild(matchUserId, entity.(*YggBuild))
	}
	return nil, errors.Swrapf(common.ErrYggdrasilUnknownMatchEntity, entity)
}

type matchHandler interface {
	getMatchLimit() int32
	getNowMatchCount(y *YggEntities) int32
}

type buildMatchHandler struct {
	areaId    int32
	matchArea *common.Area
	buildType int32
}

func NewBuildMatchHandler(areaId int32, matchArea *common.Area, buildType int32) *buildMatchHandler {
	return &buildMatchHandler{
		areaId:    areaId,
		matchArea: matchArea,
		buildType: buildType,
	}
}

func (b *buildMatchHandler) getMatchLimit() int32 {
	area, err := manager.CSV.Yggdrasil.GetArea(b.areaId)
	if err != nil {
		glog.Errorf("buildMatchHandler getMatchLimit GetArea err:%v", err)
		return 0
	}
	count := area.GetBuildMatchCount(b.buildType)

	totalArea := area.Area
	// 匹配数量= 匹配区域格子数/总格子数*总可匹配数并向下取整
	ret := math.Floor(float64(b.matchArea.Count()) / float64(totalArea.Count()) * float64(count))
	return int32(ret)
}

func (b *buildMatchHandler) getNowMatchCount(y *YggEntities) int32 {
	var ret int32 = 0
	for _, build := range y.Builds {
		if build.IsOwn() {
			continue
		}
		if build.AreaId != b.areaId {
			continue
		}
		if b.matchArea.Contains(*build.Position) {
			ret++
		}
	}
	return ret
}

type discardGoodsMatchHandler struct {
	areaId int32
}

func NewDiscardGoodsMatchHandler(areaId int32) *discardGoodsMatchHandler {
	return &discardGoodsMatchHandler{
		areaId: areaId,
	}
}

func (b *discardGoodsMatchHandler) getMatchLimit() int32 {
	area, err := manager.CSV.Yggdrasil.GetArea(b.areaId)
	if err != nil {
		glog.Errorf("buildMatchHandler getMatchLimit GetArea err:%v", err)
		return 0
	}
	return area.ItemMaxCount
}

func (b *discardGoodsMatchHandler) getNowMatchCount(y *YggEntities) int32 {

	var ret int32 = 0

	for _, DiscardGoods := range y.DiscardGoods {
		if DiscardGoods.IsOwn() {
			continue
		}
		if DiscardGoods.AreaId != b.areaId {
			continue
		}
		ret++

	}
	return ret

}

type messageMatchHandler struct {
	areaId int32
}

func NewMessageMatchHandler(areaId int32) *messageMatchHandler {
	return &messageMatchHandler{
		areaId: areaId,
	}
}
func (b *messageMatchHandler) getMatchLimit() int32 {
	area, err := manager.CSV.Yggdrasil.GetArea(b.areaId)
	if err != nil {
		glog.Errorf("messageMatchHandler getMatchLimit GetArea err:%v", err)
		return 0
	}
	return area.MessageMaxMatchCount
}

func (b *messageMatchHandler) getNowMatchCount(y *YggEntities) int32 {
	var ret int32 = 0

	for _, message := range y.Messages {
		if message.IsOwn() {
			continue
		}
		if message.AreaId != b.areaId {
			continue
		}
		ret++

	}
	return ret
}
