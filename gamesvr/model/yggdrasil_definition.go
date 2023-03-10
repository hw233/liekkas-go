package model

import (
	"context"
	"encoding/json"
	"gamesvr/manager"
	"shared/common"
	"shared/protobuf/pb"
	"shared/utility/coordinate"
	"shared/utility/errors"
	"shared/utility/global"
	"shared/utility/glog"
	"shared/utility/rand"
	"shared/utility/servertime"
)

const (
	YggdrasilPackUnitCount int32 = 1
)

type YggArea struct {
	AreaId                      int32           `json:"area_id"`
	ExploredPosCount            int32           `json:"-"` // 点亮坐标总数 init的时候重新计算
	IsTaskDone                  bool            `json:"is_task_done"`
	ExploredProgressRewardIndex int             `json:"explored_progress_reward_index"` // 探索进度已领奖的index
	BuildCreateCount            map[int32]int32 `json:"-"`                              // 建造数量 init的时候重新计算
	MessageCreateCount          int32           `json:"-"`                              // 标记数量 init的时候重新计算
}

// 登陆时候初始化
func (y *YggArea) init(yggdrasil *Yggdrasil) {
	// 有些数值可以计算得到，所以不存，每次登录时候重新计算
	y.BuildCreateCount = map[int32]int32{}
	for _, build := range yggdrasil.Entities.Builds {
		if build.AreaId != y.AreaId {
			continue
		}
		if !build.IsOwn() {
			continue
		}
		y.BuildCreateCount[build.BuildId] = y.BuildCreateCount[build.BuildId] + 1
	}

	for _, m := range yggdrasil.Entities.EnvRemovedBuilds {
		for _, build := range m {
			if build.AreaId != y.AreaId {
				continue
			}
			if !build.IsOwn() {
				continue
			}
			y.BuildCreateCount[build.BuildId] = y.BuildCreateCount[build.BuildId] + 1
		}
	}
	for _, message := range yggdrasil.Entities.Messages {
		if message.AreaId != y.AreaId {
			continue
		}
		if !message.IsOwn() {
			continue
		}
		y.MessageCreateCount++
	}
	area, err := manager.CSV.Yggdrasil.GetArea(y.AreaId)
	if err == nil {
		cut := area.Area.CutArea(yggdrasil.UnlockArea)
		y.ExploredPosCount = cut.Count()
	}
}

func NewYggArea(areaId int32) *YggArea {
	return &YggArea{
		AreaId:                      areaId,
		ExploredPosCount:            0,
		IsTaskDone:                  false,
		ExploredProgressRewardIndex: -1,
		BuildCreateCount:            map[int32]int32{},
		MessageCreateCount:          0,
	}
}
func (y *YggArea) VOYggdrasilArea() *pb.VOYggdrasilArea {
	var buildCount []*pb.VOYggdrasilBuildCount

	for buildId, count := range y.BuildCreateCount {
		buildCount = append(buildCount, &pb.VOYggdrasilBuildCount{
			BuildId: buildId,
			Count:   count,
		})
	}

	return &pb.VOYggdrasilArea{
		AreaId:                      y.AreaId,
		ExploredPosCount:            y.ExploredPosCount,
		IsTaskDone:                  y.IsTaskDone,
		ExploredProgressRewardIndex: int32(y.ExploredProgressRewardIndex),
		BuildCount:                  buildCount,
		MessageCreateCount:          y.MessageCreateCount,
	}
}

func (y *YggArea) FetchPrestige(u *User) int32 {
	area, err := manager.CSV.Yggdrasil.GetArea(y.AreaId)
	if err != nil {
		return 0
	}
	number, ok := (*u.ItemPack)[area.PrestigeItemID]
	if !ok {
		return 0
	}
	return number.Value()

}

type YggCharacter struct {
	CharacterId int32 `json:"character_id"`
	HpPercent   int32 `json:"hp_percent"`
}

type YggMark struct {
	*YggEntity `json:"entity"`
	MarkId     int32 `json:"mark_id"`
}

func NewYggMark(ctx context.Context, markId int32, position coordinate.Position) (*YggMark, error) {
	entity, err := NewYggEntity(ctx, position)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &YggMark{
		YggEntity: entity,
		MarkId:    markId,
	}, nil
}

func (ym *YggMark) VOYggdrasilMark() *pb.VOYggdrasilMark {
	return &pb.VOYggdrasilMark{
		Uid:      ym.Uid,
		Position: ym.VOPosition(),
		MarkId:   ym.MarkId,
	}
}

type YggMessage struct {
	*YggMatchEntity `json:"match_entity"`
	Comment         string `json:"comment"`
}

func (y *YggMessage) MarshalBinary() (data []byte, err error) {
	return json.Marshal(y)
}
func NewYggMessage(ctx context.Context, msg string, position coordinate.Position) (*YggMessage, error) {
	entity, err := NewYggEntity(ctx, position)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	matchEntity := NewYggMatchEntity(entity, nil)

	return &YggMessage{
		YggMatchEntity: matchEntity,
		Comment:        msg,
	}, nil
}

func (y *YggMessage) CanMatch(entities *YggEntities) error {

	p := *y.Position
	err := entities.WalkablePos(p)
	if err != nil {
		return errors.WrapTrace(err)
	}
	if !entities.IsBlankPos(p) {
		return errors.Swrapf(common.ErrYggdrasilPosAlreadyHasEntity, p.X, p.Y)
	}
	//
	//if _, ok := entities.FindMessageByPos(p); ok {
	//	return errors.WrapTrace(errors.Swrapf(common.ErrYggdrasilMessageRepeated, p.X, p.Y))
	//}

	return nil
}
func (y *YggMessage) MatchLog() {
	glog.Debugf("Match YggMessage success uid:%d,pos, %+v,value:%+v", y.Uid, y.Position, y)
}

func NewMatchYggMessage(matchUserId int64, from *YggMessage) (*YggMessage, error) {
	entity, err := NewYggEntityWithUid(from.Uid, *from.Position)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &YggMessage{
		YggMatchEntity: NewYggMatchEntity(entity, NewYggMatchInfo(matchUserId)),
		Comment:        from.Comment,
	}, nil
}

func (y *YggMessage) VOYggdrasilMessage() *pb.VOYggdrasilMessage {
	var MatchUserId int64
	if !y.IsOwn() {
		MatchUserId = y.MatchUserId
	}
	return &pb.VOYggdrasilMessage{
		Comment:     y.Comment,
		Position:    y.VOPosition(),
		MessageUid:  y.Uid,
		MatchUserId: MatchUserId,
	}
}

type YggBuild struct {
	*YggMatchEntity `json:"match_entity"`
	BuildId         int32 `json:"build_id"`
	UseCount        int32 `json:"use_count"`
}

func (y *YggBuild) MarshalBinary() (data []byte, err error) {
	return json.Marshal(y)
}
func NewMatchYggBuild(matchUserId int64, from *YggBuild) (*YggBuild, error) {
	yggEntity, err := NewYggEntityWithUid(from.Uid, *from.Position)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &YggBuild{
		YggMatchEntity: NewYggMatchEntity(yggEntity, NewYggMatchInfo(matchUserId)),
		BuildId:        from.BuildId,
		UseCount:       0,
	}, nil
}

func NewYggBuild(ctx context.Context, buildID int32, position coordinate.Position) (*YggBuild, error) {

	entity, err := NewYggEntity(ctx, position)
	if err != nil {
		return nil, err
	}
	//  匹配相关为nil
	matchEntity := NewYggMatchEntity(entity, nil)

	return &YggBuild{
		YggMatchEntity: matchEntity,
		BuildId:        buildID,
		UseCount:       0,
	}, nil
}

func (y *YggBuild) CanMatch(entities *YggEntities) error {
	if entities.DestroyedBuildUids.Contains(y.Uid) {
		return errors.Swrapf(common.ErrYggdrasilBuildDestroyedBefore, y.Uid)
	}
	p := *y.Position
	err := entities.WalkablePos(p)
	if err != nil {
		return errors.WrapTrace(err)
	}
	if !entities.BuildablePos(p) {
		return errors.Swrapf(common.ErrYggdrasilPosTypeForbidBuild, p.X, p.Y)
	}
	buildCfg, err := manager.CSV.Yggdrasil.GetYggdrasilBuilding(y.BuildId)
	if err != nil {
		return errors.WrapTrace(err)
	}
	err = entities.CheckSameTypeBuildingAround(p, buildCfg.BuildingType, buildCfg.BuildingR)
	if err != nil {
		return errors.WrapTrace(err)
	}
	err = entities.CheckSameTypeBuildingAround(p, buildCfg.BuildingType, buildCfg.MatchParam)
	if err != nil {
		return errors.WrapTrace(err)
	}

	return nil
}

func (y *YggBuild) MatchLog() {
	glog.Debugf("Match YggBuild success uid:%d,pos, %+v,value:%+v", y.Uid, y.Position, y)
}
func (y *YggBuild) VOYggdrasilBuild(ctx context.Context) (*pb.VOYggdrasilBuild, error) {
	var MatchUserId int64
	if !y.IsOwn() {
		MatchUserId = y.MatchUserId
	}
	cnt, err := manager.Global.FetchBuildUseNum(ctx, y.Uid)
	if err != nil && err != global.ErrNil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.VOYggdrasilBuild{
		Position:      y.VOPosition(),
		BuildUid:      y.Uid,
		UseCount:      y.UseCount,
		BuildId:       y.BuildId,
		TotalUseCount: int32(cnt),
		MatchUserId:   MatchUserId,
	}, nil
}

func (y *YggBuild) AfterUse(ctx context.Context) error {
	_, err := manager.Global.IncrBuildUseNum(ctx, y.Uid)
	if err != nil {
		return errors.WrapTrace(err)
	}
	y.UseCount += 1
	y.UpdateAt = servertime.Now().Unix()

	return nil
}

type YggMatchEnable interface {
	CanMatch(entities *YggEntities) error
	MatchLog()
	GetMatchUserId() int64
}

type YggMatchEntity struct {
	*YggEntity    `json:"entity"`
	*YggMatchInfo `json:"match_info"`
}

func (y *YggMatchEntity) GetMatchUserId() int64 {
	if y.YggMatchInfo != nil {
		return y.YggMatchInfo.MatchUserId
	}
	return 0
}

func NewYggMatchEntity(yggEntity *YggEntity, yggMatchInfo *YggMatchInfo) *YggMatchEntity {
	return &YggMatchEntity{
		YggEntity:    yggEntity,
		YggMatchInfo: yggMatchInfo,
	}
}

func (y *YggMatchEntity) CanMatch(entities *YggEntities) error {

	p := *y.Position
	err := entities.WalkablePos(p)
	if err != nil {
		return errors.WrapTrace(err)
	}
	if !entities.IsBlankPos(p) {
		return errors.Swrapf(common.ErrYggdrasilPosAlreadyHasEntity, p.X, p.Y)
	}
	return nil
}

type YggMatchInfo struct { // 匹配物品的uid=原物品的uid
	MatchUserId int64 `json:"match_user_id"` // 匹配的userId
}

func (y *YggMatchInfo) IsOwn() bool {
	return y == nil
}

func NewYggMatchInfo(MatchUserId int64) *YggMatchInfo {
	return &YggMatchInfo{
		MatchUserId: MatchUserId,
	}
}

type YggEntity struct {
	*coordinate.Position `json:"position"`
	Uid                  int64 `json:"uid"`
	AreaId               int32 `json:"area_id"`
	CreateAt             int64 `json:"create_at"`
	UpdateAt             int64 `json:"update_at"`
}

func (e *YggEntity) Clone() *YggEntity {
	return &YggEntity{
		Uid:      e.Uid,
		Position: coordinate.NewPosition(e.X, e.Y),
		AreaId:   e.AreaId,
		CreateAt: e.CreateAt,
		UpdateAt: e.UpdateAt,
	}
}
func NewYggEntityWithUid(uid int64, position coordinate.Position) (*YggEntity, error) {
	areaId, err := manager.CSV.Yggdrasil.GetPosAreaId(position)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	now := servertime.Now().Unix()
	return &YggEntity{
		Uid:      uid,
		Position: &position, // 不传指针 防止指针修改
		AreaId:   areaId,
		CreateAt: now,
		UpdateAt: now,
	}, nil
}

func NewYggEntity(ctx context.Context, position coordinate.Position) (*YggEntity, error) {
	uid, err := manager.Global.GenYggdrasilEntityUid(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	areaId, err := manager.CSV.Yggdrasil.GetPosAreaId(position)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	now := servertime.Now().Unix()
	return &YggEntity{
		Uid:      uid,
		Position: &position, // 不传指针 防止指针修改
		AreaId:   areaId,
		CreateAt: now,
		UpdateAt: now,
	}, nil
}

type YggObject struct {
	*YggEntity        `json:"entity"`
	ObjectId          int32                `json:"object_id"`
	State             int32                `json:"state"`
	CreateAtSubTaskId int32                `json:"create_at_sub_task_id"`
	DeleteAtSubTaskId int32                `json:"delete_at_sub_task_id"` // 在指定的subtask完成后被移除
	OrgPos            *coordinate.Position `json:"org_pos"`               //初始位置
}

func NewYggObject(ctx context.Context, position coordinate.Position, objectId, state int32) (*YggObject, error) {
	entity, err := NewYggEntity(ctx, position)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &YggObject{
		ObjectId:  objectId,
		State:     state,
		YggEntity: entity,
		OrgPos:    &position,
	}, nil
}

func (y *YggObject) Clone() *YggObject {

	return &YggObject{
		ObjectId:          y.ObjectId,
		State:             y.State,
		CreateAtSubTaskId: y.CreateAtSubTaskId,
		DeleteAtSubTaskId: y.DeleteAtSubTaskId,
		YggEntity:         y.YggEntity.Clone(),
	}
}
func (y *YggObject) SetCreateAndDeleteAt(CreateAt int32, DeleteAt int32) {
	y.CreateAtSubTaskId = CreateAt
	y.DeleteAtSubTaskId = DeleteAt
}
func (y *YggObject) VOYggdrasilObject() *pb.VOYggdrasilObject {
	return &pb.VOYggdrasilObject{
		ObjectUid: y.Uid,
		ObjectId:  y.ObjectId,
		State:     y.State,
		Position:  y.VOPosition(),
		//EnvId:     y.EnvId,
	}
}

type YggGoods struct {
	ItemId int32 `json:"item_id"`
}

func NewYggGoods(ItemId int32) *YggGoods {
	return &YggGoods{
		ItemId: ItemId,
	}
}

func (y *YggGoods) VOResource() *pb.VOResource {
	return &pb.VOResource{
		ItemId: y.ItemId,
		Count:  YggdrasilPackUnitCount,
	}
}

type YggDiscardGoods struct {
	*YggMatchEntity `json:"match_entity"`
	*YggGoods       `json:"goods"`
}

func (y *YggDiscardGoods) MarshalBinary() ([]byte, error) {
	return json.Marshal(y)
}
func (y *YggDiscardGoods) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, y)
}

func (y *YggDiscardGoods) MatchLog() {
	glog.Debugf("Match YggDiscardGoods success uid:%d,pos, %+v,value:%+v", y.Uid, y.Position, y)
}

type YggPack map[int64]*YggPackGoods

func NewYggPack() *YggPack {
	return (*YggPack)(&map[int64]*YggPackGoods{})
}
func (y *YggPack) Values() []*YggPackGoods {
	values := make([]*YggPackGoods, 0, len(*y))
	for _, goods := range *y {
		values = append(values, goods)
	}
	return values
}

func (y *YggPack) IsFull(lv int32) bool {
	return len(*y) >= int(manager.CSV.Yggdrasil.GetYggBagAllCount(lv))
}

func (y *YggPack) Get(goodsId int64) (*YggPackGoods, bool) {
	goods, ok := (*y)[goodsId]
	return goods, ok
}

func (y *YggPack) Add(goods *YggPackGoods) {
	(*y)[goods.Uid] = goods
}

func (y *YggPack) Delete(goodsId int64) {
	delete(*y, goodsId)
}

// GetAllOwnGoods 获得被宝贝自己的所有物品
func (y *YggPack) GetAllOwnGoods() *YggPack {
	ret := NewYggPack()
	for _, goods := range *y {
		if !goods.IsOwn() {
			continue
		}
		ret.Add(goods)
	}
	return ret
}

// CostReward 返回剩余要消耗的和要删除的goodsId
func (y *YggPack) CostReward(reward common.Reward) (*common.Reward, []int64) {
	var ret []int64
	num := reward.Num
	for _, goods := range *y {
		if goods.ItemId == reward.ID {
			num--
			ret = append(ret, goods.Uid)
			if num == 0 {
				break
			}

		}
	}
	return common.NewReward(reward.ID, num), ret
}

func (y *YggPack) init(matchUserIds map[int64]struct{}) {
	for _, v := range y.Values() {
		if !v.IsOwn() {
			matchUserIds[v.MatchUserId] = struct{}{}
		}
	}

}

func NewMatchYggDiscardGoods(matchUserId int64, from *YggDiscardGoods) (*YggDiscardGoods, error) {
	entity, err := NewYggEntityWithUid(from.Uid, *from.Position)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &YggDiscardGoods{
		YggMatchEntity: NewYggMatchEntity(entity, NewYggMatchInfo(matchUserId)),
		YggGoods:       NewYggGoods(from.ItemId),
	}, nil

}

func NewYggDiscardGoods(position coordinate.Position, goods *YggPackGoods) (*YggDiscardGoods, error) {
	entity, err := NewYggEntityWithUid(goods.Uid, position)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	matchEntity := NewYggMatchEntity(entity, goods.YggMatchInfo)

	return &YggDiscardGoods{
		YggMatchEntity: matchEntity,
		YggGoods:       NewYggGoods(goods.ItemId),
	}, nil
}

func NewYggDiscardGoodsByReward(ctx context.Context, position coordinate.Position, itemId int32) (*YggDiscardGoods, error) {
	entity, err := NewYggEntity(ctx, position)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	matchEntity := NewYggMatchEntity(entity, nil)

	return &YggDiscardGoods{
		YggMatchEntity: matchEntity,
		YggGoods:       NewYggGoods(itemId),
	}, nil
}

func (y *YggDiscardGoods) VOYggdrasilDiscardGoods() *pb.VOYggdrasilDiscardGoods {
	var MatchUserId int64
	if y.YggMatchInfo != nil {
		MatchUserId = y.MatchUserId
	}
	return &pb.VOYggdrasilDiscardGoods{
		Uid:      y.Uid,
		Position: y.VOPosition(),

		Resource:    y.VOResource(),
		MatchUserId: MatchUserId,
	}
}

type YggPackGoods struct {
	Uid           int64 `json:"uid"`
	*YggGoods     `json:"goods"`
	*YggMatchInfo `json:"match_info"`
}

func NewYggPackGoodsByDiscard(discardGoods *YggDiscardGoods) (*YggPackGoods, error) {

	return &YggPackGoods{
		Uid:          discardGoods.Uid,
		YggGoods:     NewYggGoods(discardGoods.ItemId),
		YggMatchInfo: discardGoods.YggMatchInfo,
	}, nil
}

func NewYggPackGoods(ctx context.Context, itemId int32) (*YggPackGoods, error) {
	uid, err := manager.Global.GenYggdrasilEntityUid(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &YggPackGoods{
		Uid:          uid,
		YggGoods:     NewYggGoods(itemId),
		YggMatchInfo: nil,
	}, nil
}

func (y *YggPackGoods) VOYggdrasilPackGoods() *pb.VOYggdrasilPackGoods {
	var MatchUserId int64
	if y.YggMatchInfo != nil {
		MatchUserId = y.MatchUserId
	}

	return &pb.VOYggdrasilPackGoods{
		GoodsId: y.Uid,
		Resource: &pb.VOResource{
			ItemId: y.ItemId,
			Count:  YggdrasilPackUnitCount,
		},
		OriDiscardUserId: MatchUserId,
	}
}

type YggAreas struct {
	Areas map[int32]*YggArea `json:"areas"`
}

func NewYggAreas() *YggAreas {
	return &YggAreas{
		Areas: map[int32]*YggArea{},
	}
}

func (y *YggAreas) VOYggdrasilArea(areaId int32) *pb.VOYggdrasilArea {
	area, ok := y.Areas[areaId]
	if ok {
		return area.VOYggdrasilArea()
	}

	return nil

}
func (y *YggAreas) VOYggdrasilAreas() []*pb.VOYggdrasilArea {
	vos := make([]*pb.VOYggdrasilArea, 0, len(y.Areas))
	for _, area := range y.Areas {
		vos = append(vos, area.VOYggdrasilArea())
	}
	return vos

}
func (y *YggAreas) contains(areaId int32) bool {
	_, ok := y.Areas[areaId]
	return ok
}

func (y *YggAreas) get(areaId int32) (*YggArea, bool) {
	area, ok := y.Areas[areaId]
	return area, ok
}
func (y *YggAreas) getByCreate(ctx context.Context, yggdrasil *Yggdrasil, areaId int32) *YggArea {

	area, ok := y.Areas[areaId]
	if !ok {
		area = NewYggArea(areaId)
		// 初始化区域内互动物
		yggdrasil.InitAreaInitObj(ctx, areaId)
		y.Areas[areaId] = area
	}
	return area
}

func (y *YggAreas) init(yggdrasil *Yggdrasil) {
	for _, area := range y.Areas {
		area.init(yggdrasil)
	}
}

type YggTaskPack struct {
	Items map[int32]*ItemPack `json:"items"` // key subtaskId
}

func NewYggTaskPack() *YggTaskPack {
	return &YggTaskPack{
		Items: map[int32]*ItemPack{},
	}
}

func (p *YggTaskPack) GetBySubTaskId(subtaskId int32) (*ItemPack, bool) {
	pack, ok := p.Items[subtaskId]
	return pack, ok
}

func (p *YggTaskPack) Add(itemId int32, num int32, subtaskId int32) int32 {
	subPack, ok := p.Items[subtaskId]
	if !ok {
		subPack = NewItemPack()
		p.Items[subtaskId] = subPack
	}
	subPack.Add(itemId, num)

	return (*p.merge())[itemId].Value()
}

func (p *YggTaskPack) Remove(subtaskId int32) map[int32]int32 {
	subPack, ok := p.Items[subtaskId]
	idNum := map[int32]int32{}
	if !ok {
		return idNum
	}
	delete(p.Items, subtaskId)
	merge := *p.merge()
	for id := range *subPack {
		idNum[id] = merge[id].Value()
	}

	return idNum

}

func (p *YggTaskPack) CostReward(reward common.Reward) (*common.Reward, int32) {
	for _, pack := range p.Items {
		restNum := CostItemPack(pack, reward.ID, reward.Num)
		reward.Num = restNum
		if restNum == 0 {
			break
		}
	}
	eventNumber, ok := (*p.merge())[reward.ID]
	if ok {
		return &reward, eventNumber.Value()
	} else {
		return &reward, 0

	}
}

// CostItemPack 消耗并返回还剩多少个未消耗
func CostItemPack(itemPack *ItemPack, id, num int32) int32 {
	var tmp int32 = 0
	enough := itemPack.Enough(id, num)
	if enough {
		itemPack.Minus(id, num)
	} else {
		eventNumber, ok := (*itemPack)[id]
		if ok {
			value := eventNumber.Value()
			itemPack.Minus(id, value)
			tmp = num - value
		} else {
			tmp = num
		}
	}

	return tmp
}

func (p *YggTaskPack) merge() *ItemPack {
	tmp := NewItemPack()
	for _, pack := range p.Items {
		for id, count := range *pack {
			tmp.Add(id, count.Value())
		}
	}
	return tmp
}

func (p *YggTaskPack) VOTaskPackInfo() []*pb.VOResource {

	var vos []*pb.VOResource

	for k, v := range *p.merge() {
		vos = append(vos, &pb.VOResource{
			ItemId: k,
			Count:  v.Value(),
		})
	}

	return vos
}

type YggdrasilEntityChange struct {
	Objects             map[int64]*YggObject
	Builds              map[int64]*YggBuild
	RemovedBuildingUids []int64
	RemovedObjectUids   []int64
	AddEnvTerrain       []int32 // 增加的env地形
	RemoveEnvTerrain    []int32 // 删除的env地形
	*YggdrasilResourceResult
}

func NewYggdrasilEntityChange() *YggdrasilEntityChange {
	return &YggdrasilEntityChange{
		Objects:                 map[int64]*YggObject{},
		Builds:                  map[int64]*YggBuild{},
		RemovedObjectUids:       nil,
		RemovedBuildingUids:     nil,
		YggdrasilResourceResult: NewYggdrasilResourceResult(),
	}
}

func (y *YggdrasilEntityChange) VOYggdrasilEntityChange(ctx context.Context, yggdrasil *Yggdrasil) *pb.VOYggdrasilEntityChange {
	ObjectVOs := make([]*pb.VOYggdrasilObject, 0, len(y.Objects))
	BuildVOs := make([]*pb.VOYggdrasilBuild, 0, len(y.Builds))

	for _, object := range y.Objects {
		_, ok := yggdrasil.Entities.FindObjectByUid(object.Uid)
		if ok {
			// 兼容一些改变状态，但又删除的操作
			ObjectVOs = append(ObjectVOs, object.VOYggdrasilObject())
		}

	}
	for _, build := range y.Builds {
		yggdrasilBuild, err := build.VOYggdrasilBuild(ctx)
		if err != nil {
			glog.Errorf("VOYggdrasilEntityChange VOYggdrasilBuild  err:%+v", err)
			continue
		}
		BuildVOs = append(BuildVOs, yggdrasilBuild)
	}

	ret := &pb.VOYggdrasilEntityChange{
		Objects:           ObjectVOs,
		Builds:            BuildVOs,
		RemovedObjectUids: y.RemovedObjectUids,
		RemovedBuildUids:  y.RemovedBuildingUids,
		YggdrasilResource: y.VOYggdrasilResourceResult(),
		AddEnvTerrain:     y.AddEnvTerrain,
		RemoveEnvTerrain:  y.RemoveEnvTerrain,
	}
	y.clear()

	return ret

}
func (y *YggdrasilEntityChange) clear() {
	y.Objects = map[int64]*YggObject{}
	y.Builds = map[int64]*YggBuild{}
	y.RemovedBuildingUids = nil
	y.RemovedObjectUids = nil
	y.AddEnvTerrain = nil
	y.RemoveEnvTerrain = nil
}

func (y *YggdrasilEntityChange) AppendObject(object *YggObject) {
	y.Objects[object.Uid] = object
}

func (y *YggdrasilEntityChange) AppendRemovedBuildingUid(uid int64) {
	y.RemovedBuildingUids = append(y.RemovedBuildingUids, uid)

}

func (y *YggdrasilEntityChange) AppendRemovedObjectUid(uid int64) {
	y.RemovedObjectUids = append(y.RemovedObjectUids, uid)

}

func (y *YggdrasilEntityChange) AppendBuild(build *YggBuild) {
	y.Builds[build.Uid] = build

}

type YggdrasilDailyMonsters map[int32]*YggdrasilDailyMonster

func NewYggdrasilDailyMonsters() *YggdrasilDailyMonsters {
	return (*YggdrasilDailyMonsters)(&map[int32]*YggdrasilDailyMonster{})
}
func (y *YggdrasilDailyMonsters) init(ctx context.Context, yggdrasil *Yggdrasil) {
	for _, cfg := range manager.CSV.Yggdrasil.GetAllDailyMonster() {
		_, ok := (*y)[cfg.Id]
		if !ok {
			(*y)[cfg.Id] = NewYggdrasilDailyMonster(cfg.Id)
			(*y)[cfg.Id].GenNewMonster(ctx, yggdrasil)
		}
	}

}
func (y *YggdrasilDailyMonsters) DailyRefresh(ctx context.Context, yggdrasil *Yggdrasil) {
	for _, monster := range *y {
		monster.DailyRefresh(ctx, yggdrasil)
	}
}

type YggdrasilDailyMonster struct {
	CfgId       int32 `json:"cfg_id"`
	ObjectUid   int64 `json:"object_uid"`
	IntervalDay int32 `json:"interval_day"` // 持续天数
}

func NewYggdrasilDailyMonster(CfgId int32) *YggdrasilDailyMonster {
	return &YggdrasilDailyMonster{
		CfgId:       CfgId,
		ObjectUid:   -1,
		IntervalDay: 0,
	}
}

func (y *YggdrasilDailyMonster) GenNewMonster(ctx context.Context, yggdrasil *Yggdrasil) {
	monster, err := manager.CSV.Yggdrasil.GetDailyMonster(y.CfgId)
	if err != nil {
		glog.Errorf("YggdrasilDailyMonsters Refresh GetDailyMonster err:%+v", err)
		return
	}
	// 移除老的obj
	object, ok := yggdrasil.Entities.FindObjectByUid(y.ObjectUid)
	if ok {
		yggdrasil.RemoveObject(object)
	}
	// 生成新的obj
	objConfig, err := manager.CSV.Yggdrasil.GetYggdrasilObjectConfig(monster.ObjectId)
	if err != nil {
		glog.Errorf("YggdrasilDailyMonsters Refresh GetYggdrasilObjectConfig err:%+v", err)
		return
	}

	cubeRange := coordinate.CubeRange(*monster.Pos, monster.Radius)
	//打乱顺序
	for i := len(cubeRange) - 1; i > 0; i-- {
		num := rand.RangeInt32(0, int32(i))
		(cubeRange)[i], (cubeRange)[num] = (cubeRange)[num], (cubeRange)[i]
	}
	for _, position := range cubeRange {
		if yggdrasil.Entities.IsBlankPos(position) {
			yggObject, err := NewYggObject(ctx, position, objConfig.Id, objConfig.DefaultState)
			if err != nil {
				glog.Errorf("YggdrasilDailyMonsters Refresh NewYggObject err:%+v", err)
				return
			}
			yggdrasil.AppendObject(ctx, yggObject)
			break
		}

	}
	y.IntervalDay = 0 // 如果是每天刷新间隔天数配1
}

func (y *YggdrasilDailyMonster) DailyRefresh(ctx context.Context, yggdrasil *Yggdrasil) {
	y.IntervalDay++
	monster, err := manager.CSV.Yggdrasil.GetDailyMonster(y.CfgId)
	if err != nil {
		glog.Errorf("YggdrasilDailyMonsters Refresh GetDailyMonster err:%+v", err)
		return
	}
	// 达到刷新天数
	if monster.IntervalDay >= y.IntervalDay {
		y.GenNewMonster(ctx, yggdrasil)
	}

}

type YggdrasilSpecialStatics struct {
	Complete  bool  `json:"complete"`
	MoveCount int32 `json:"move_count"`
}

func NewYggdrasilSpecialStatics() *YggdrasilSpecialStatics {
	return &YggdrasilSpecialStatics{
		Complete:  false,
		MoveCount: 0,
	}
}
func (y *YggdrasilSpecialStatics) VOYggSpecialStatics() *pb.VOYggSpecialStatics {
	return &pb.VOYggSpecialStatics{
		MoveCount: y.MoveCount,
	}
}

func (y *YggdrasilSpecialStatics) TriggerMoveCount() {
	if y.Complete {
		return
	}
	y.MoveCount++
}
func (y *YggdrasilSpecialStatics) IsMoveCountReach(moveCount int32) error {
	if y.MoveCount >= moveCount {
		y.Complete = true
		return nil
	}

	return errors.New("move count not reach:%d", y.MoveCount)
}

type YggdrasilBattleNpc struct {
	NpcId int32 // 对应 cfg_battle_npc表的id
	Hp    int32 // 血量（万分比）
}

func NewYggdrasilBattleNpc(NpcId, Hp int32) *YggdrasilBattleNpc {
	return &YggdrasilBattleNpc{
		NpcId: NpcId,
		Hp:    Hp,
	}
}

type YggdrasilInitObjects map[int32]map[int32]struct{}

func NewYggdrasilInitObjects() *YggdrasilInitObjects {
	return (*YggdrasilInitObjects)(&map[int32]map[int32]struct{}{})

}

func (y *YggdrasilInitObjects) Append(position coordinate.Position) {
	m, ok := (*y)[position.X]
	if !ok {
		m = map[int32]struct{}{}
		(*y)[position.X] = m
	}
	m[position.Y] = struct{}{}
}
func (y *YggdrasilInitObjects) Contains(position coordinate.Position) bool {
	m, ok := (*y)[position.X]
	if !ok {
		return false
	}
	_, ok = m[position.Y]
	return ok
}

type YggEnvTerrain struct {
	EnvId    int32 `json:"env_id"`
	CreateAt int32 `json:"create_at"`
	DeleteAt int32 `json:"delete_at"`
}

func NewYggEnvTerrain(EnvId, CreateAt, DeleteAt int32) *YggEnvTerrain {
	return &YggEnvTerrain{
		EnvId:    EnvId,
		CreateAt: CreateAt,
		DeleteAt: DeleteAt,
	}
}
