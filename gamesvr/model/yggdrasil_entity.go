package model

import (
	"context"
	"gamesvr/manager"
	"shared/common"
	"shared/csv/entry"
	"shared/csv/static"
	"shared/protobuf/pb"
	"shared/utility/coordinate"
	"shared/utility/errors"
	"shared/utility/glog"
	"shared/utility/number"
	"shared/utility/servertime"
)

type YggEntities struct {
	//   Builds,CoBuilds,Marks,DiscardGoods,Objects,Messages uid不重复
	Builds       map[int64]*YggBuild        `json:"builds"`        // 建筑
	CoBuilds     map[int64]*YggCoBuild      `json:"co_builds"`     // 联合建筑
	Marks        map[int64]*YggMark         `json:"marks"`         // 标注    标注可以互相重叠吗
	Messages     map[int64]*YggMessage      `json:"messages"`      // 留言    留言和其他东西（建筑，丢弃物品等）可以重合
	DiscardGoods map[int64]*YggDiscardGoods `json:"discard_goods"` // 丢弃物品
	Objects      map[int64]*YggObject       `json:"objects"`       // 互动物

	EnvRemovedBuilds       map[int32]map[int64]*YggBuild        `json:"env_removed_builds"`        // env 暂时移除的建筑  (key subtaskId)
	EnvRemovedObjects      map[int32]map[int64]*YggObject       `json:"env_removed_objects"`       // env 暂时移除的obj  (key subtaskId)
	EnvRemovedDiscardGoods map[int32]map[int64]*YggDiscardGoods `json:"env_removed_discard_goods"` // env 暂时移除的丢弃物  (key subtaskId)

	DestroyedBuildUids         *number.NonRepeatableArrInt64 `json:"destroyed_build_uids"`           // 拆毁的建筑，下次不再匹配
	LastMatchDiscardGoodsUsers *number.NonRepeatableArrInt64 `json:"last_match_discard_goods_users"` //上次掉落物匹配的用户

	NpcHp       map[int32]int32  `json:"npc_hp"`       // 助战npc血量
	EnvTerrains []*YggEnvTerrain `json:"env_terrains"` // 任务环境导致的地形变化 ，需要保持有序

	Blocks            *YggBlocks                     `json:"-"` // 只同步数据用
	EntityFinder      *YggEntityFinder               `json:"-"` // 计算用
	EnvObjectDeleteAt map[int32]map[int64]*YggObject `json:"-"` //缓存用 env创造的obj (key subtaskId 需要在某subtask完成后被删除)     init的时候load
	EnvObjectCreateAt map[int32]map[int64]*YggObject `json:"-"` //缓存用 env创造的obj  (key subtaskId 在某subtask被创建) init的时候load
	PosTypeWithObject map[coordinate.Position]int32  `json:"-"` // 地形 init的时候重新计算 0  可行走区域  1 不可行走区域  2  可行走不可建造  注：为0的不存
	PosHeight         map[coordinate.Position]int32  `json:"-"` // 高度 init的时候重新计算  注：为1的不存
	PosType           map[coordinate.Position]int32  `json:"-"` // 地形 init的时候重新计算 0  可行走区域  1 不可行走区域  2  可行走不可建造  注：为0的不存
}

func NewYggEntities() *YggEntities {
	return &YggEntities{
		Builds:       map[int64]*YggBuild{},
		CoBuilds:     map[int64]*YggCoBuild{},
		Marks:        map[int64]*YggMark{},
		Messages:     map[int64]*YggMessage{},
		DiscardGoods: map[int64]*YggDiscardGoods{},
		Objects:      map[int64]*YggObject{},

		EnvRemovedBuilds:       map[int32]map[int64]*YggBuild{},
		EnvRemovedObjects:      map[int32]map[int64]*YggObject{},
		EnvRemovedDiscardGoods: map[int32]map[int64]*YggDiscardGoods{},

		DestroyedBuildUids:         number.NewNonRepeatableArrInt64(),
		LastMatchDiscardGoodsUsers: number.NewNonRepeatableArrInt64(),

		NpcHp: map[int32]int32{},

		Blocks:            NewYggBlocks(),
		EntityFinder:      NewYggEntityFinder(),
		EnvObjectDeleteAt: map[int32]map[int64]*YggObject{},
		EnvObjectCreateAt: map[int32]map[int64]*YggObject{},
		PosTypeWithObject: map[coordinate.Position]int32{},
		PosHeight:         map[coordinate.Position]int32{},
		PosType:           map[coordinate.Position]int32{},
	}
}

// 加载数据到block 和tree 方便数据处理
func (y *YggEntities) init(matchUserIds map[int64]struct{}) {

	var entities []*YggEntity
	for _, v := range y.Builds {
		y.Blocks.AppendBuild(v)
		entities = append(entities, v.YggEntity)
		if !v.IsOwn() {
			matchUserIds[v.MatchUserId] = struct{}{}
		}
	}
	for _, v := range y.CoBuilds {
		y.Blocks.AppendCoBuild(v)
		entities = append(entities, v.YggEntity)
	}
	for _, v := range y.Messages {
		y.Blocks.AppendMessage(v)
		entities = append(entities, v.YggEntity)
		if !v.IsOwn() {
			matchUserIds[v.MatchUserId] = struct{}{}
		}
	}
	for _, v := range y.Objects {
		if v.DeleteAtSubTaskId > 0 {
			m, ok := y.EnvObjectDeleteAt[v.DeleteAtSubTaskId]
			if !ok {
				m = map[int64]*YggObject{}
				y.EnvObjectDeleteAt[v.DeleteAtSubTaskId] = m
			}
			m[v.Uid] = v
		}
		if v.CreateAt > 0 {
			m, ok := y.EnvObjectCreateAt[v.CreateAtSubTaskId]
			if !ok {
				m = map[int64]*YggObject{}
				y.EnvObjectCreateAt[v.CreateAtSubTaskId] = m
			}
			m[v.Uid] = v

		}
		y.Blocks.AppendObject(v)
		entities = append(entities, v.YggEntity)
	}
	for _, v := range y.DiscardGoods {
		y.Blocks.AppendDiscardGoods(v)
		entities = append(entities, v.YggEntity)
		if !v.IsOwn() {
			matchUserIds[v.MatchUserId] = struct{}{}
		}
	}
	for _, v := range y.Marks {
		y.Blocks.AppendMark(v)
		entities = append(entities, v.YggEntity)

	}
	y.EntityFinder.Append(entities...)

	//
	y.RefreshTerrain()

}

// CheckSameTypeBuildingAround 检查周围有无同type建筑
func (y *YggEntities) CheckSameTypeBuildingAround(position coordinate.Position, buildType, r int32) error {
	positions := coordinate.CubeRange(position, r)
	for _, p := range positions {
		tmp, ok := y.FindBuildingsByPos(p)
		if !ok {
			continue
		}
		buildCfg, err := manager.CSV.Yggdrasil.GetYggdrasilBuilding(tmp.BuildId)
		if err != nil {
			return errors.WrapTrace(err)
		}
		if buildCfg.BuildingType == buildType {
			return errors.Swrapf(common.ErrYggdrasilPosTooCloseToSameBuild, position.X, position.Y)
		}
	}
	return nil

}

// AppendBuild 添加建筑
func (y *YggEntities) AppendBuild(ctx context.Context, userId int64, build *YggBuild) error {
	err := AppendYggBuildInRedis(ctx, userId, build)
	if err != nil {
		return errors.WrapTrace(err)
	}
	y.appendBuild(build)
	return nil
}
func (y *YggEntities) appendBuild(build *YggBuild) {
	y.Builds[build.Uid] = build
	y.EntityFinder.Append(build.YggEntity)
	y.Blocks.AppendBuild(build)
}

// AppendMessage 添加留言
func (y *YggEntities) AppendMessage(ctx context.Context, userId int64, msg *YggMessage) error {
	err := AppendYggMessageInRedis(ctx, userId, msg)
	if err != nil {
		return errors.WrapTrace(err)
	}
	y.appendMessage(msg)
	return nil
}

func (y *YggEntities) appendMessage(msg *YggMessage) {
	y.Messages[msg.Uid] = msg
	y.EntityFinder.Append(msg.YggEntity)
	y.Blocks.AppendMessage(msg)
}

// 添加联合建筑
func (y *YggEntities) AppendCoBuild(build *YggCoBuild) {
	y.CoBuilds[build.Uid] = build
	y.EntityFinder.Append(build.YggEntity)
	y.Blocks.AppendCoBuild(build)
}

func (y *YggEntities) UpdateMessage(ctx context.Context, userId int64, msg *YggMessage, m string) error {
	msg.Comment = m
	msg.UpdateAt = servertime.Now().Unix()
	err := AppendYggMessageInRedis(ctx, userId, msg)
	if err != nil {
		return errors.WrapTrace(err)
	}
	y.EntityFinder.Remove(msg.YggEntity)
	y.EntityFinder.Append(msg.YggEntity)
	y.Blocks.RemoveMessage(msg)
	y.Blocks.AppendMessage(msg)
	y.Messages[msg.Uid] = msg
	return nil
}

// AppendObject 添加互动物
func (y *YggEntities) AppendObject(object *YggObject) {
	y.Objects[object.Uid] = object
	y.EntityFinder.Append(object.YggEntity)
	y.Blocks.AppendObject(object)

	if object.DeleteAtSubTaskId > 0 {
		m, ok := y.EnvObjectDeleteAt[object.DeleteAtSubTaskId]
		if !ok {
			m = map[int64]*YggObject{}
			y.EnvObjectDeleteAt[object.DeleteAtSubTaskId] = m
		}
		m[object.Uid] = object
	}
	if object.CreateAt > 0 {
		m, ok := y.EnvObjectCreateAt[object.CreateAtSubTaskId]
		if !ok {
			m = map[int64]*YggObject{}
			y.EnvObjectCreateAt[object.CreateAtSubTaskId] = m
		}
		m[object.Uid] = object
	}
	y.RefreshTerrainForObjectChange(*object.Position)

	state, err := manager.CSV.Yggdrasil.GetObjectState(object.ObjectId, object.State)
	if err == nil {
		if state.ObjectType == static.YggObjectTypeFollowbattlenpc {
			y.AddBattleNpc(state.ObjectParam)
		}
	}
}

// AppendMark 添加标注
func (y *YggEntities) AppendMark(mark *YggMark) {
	y.Marks[mark.Uid] = mark
	y.EntityFinder.Append(mark.YggEntity)
	y.Blocks.AppendMark(mark)
}

func (y *YggEntities) CalPosType(position coordinate.Position) int32 {

	initType, ok := y.PosType[position]
	if !ok {
		initType = entry.PosTypeWalkable
	}
	objs, ok := y.FindObjectsByPos(position)
	if !ok {
		return initType
	}

	/**
	PosTypeWalkable    = 0 // 可行走区域
	PosTypeUnWalkable  = 1 // 不可行走区域
	PosTypeUnBuildable = 2 // 可行走不可建造

	重叠的情况下1优先级最高 其次是2 然后是0
	*/
	var posType int32 = entry.PosTypeWalkable
	for _, obj := range objs {
		state, err := manager.CSV.Yggdrasil.GetObjectState(obj.ObjectId, obj.State)
		if err != nil {
			glog.Errorf("CalPosType GetYggdrasilObjectConfig err%+v:", err)
			continue
		}
		if entry.PosTypePriority[state.PosType] > entry.PosTypePriority[posType] {
			posType = state.PosType
		}
	}
	return posType
}

// AppendDiscardGoods 添加丢弃物
func (y *YggEntities) AppendDiscardGoods(ctx context.Context, userId int64, discardGoods *YggDiscardGoods) error {
	err := AppendYggDiscardGoodsInRedis(ctx, userId, discardGoods)
	if err != nil {
		return errors.WrapTrace(err)
	}
	y.appendDiscardGoods(discardGoods)
	return nil
}

func (y *YggEntities) appendDiscardGoods(discardGoods *YggDiscardGoods) {
	y.DiscardGoods[discardGoods.Uid] = discardGoods
	y.EntityFinder.Append(discardGoods.YggEntity)
	y.Blocks.AppendDiscardGoods(discardGoods)
}

// RemoveBuild 移除建筑
func (y *YggEntities) RemoveBuild(ctx context.Context, userId int64, build *YggBuild) error {
	err := DelYggBuildInRedis(ctx, userId, build)
	if err != nil {
		return errors.WrapTrace(err)
	}
	y.removeBuild(build)
	return nil
}

func (y *YggEntities) removeBuild(build *YggBuild) {
	delete(y.Builds, build.Uid)
	y.EntityFinder.Remove(build.YggEntity)
	y.Blocks.RemoveBuild(build)
	y.DestroyedBuildUids.Append(build.Uid)
}

func (y *YggEntities) RemoveMessage(ctx context.Context, userId int64, msg *YggMessage) error {
	err := DelYggMessageInRedis(ctx, userId, msg)
	if err != nil {
		return errors.WrapTrace(err)
	}
	y.removeMessage(msg)
	return nil
}

func (y *YggEntities) removeMessage(msg *YggMessage) {
	delete(y.Messages, msg.Uid)
	y.EntityFinder.Remove(msg.YggEntity)
	y.Blocks.RemoveMessage(msg)
}

// RemoveObject 移除互动物
func (y *YggEntities) RemoveObject(object *YggObject) {
	delete(y.Objects, object.Uid)
	y.EntityFinder.Remove(object.YggEntity)
	y.Blocks.RemoveObject(object)

	if object.DeleteAtSubTaskId > 0 && y.EnvObjectDeleteAt[object.DeleteAtSubTaskId] != nil {
		delete(y.EnvObjectDeleteAt[object.DeleteAtSubTaskId], object.Uid)

	}
	if object.CreateAt > 0 && y.EnvObjectCreateAt[object.CreateAtSubTaskId] != nil {
		delete(y.EnvObjectCreateAt[object.CreateAtSubTaskId], object.Uid)

	}
	y.RefreshTerrainForObjectChange(*object.Position)

	state, err := manager.CSV.Yggdrasil.GetObjectState(object.ObjectId, object.State)
	if err == nil {
		if state.ObjectType == static.YggObjectTypeFollowbattlenpc {
			y.RemoveBattleNpc(state.ObjectParam)
		}
	}
}

func (y *YggEntities) SetObjectNextState(object *YggObject, nextState int32) {
	// battleNpc处理
	before, err := manager.CSV.Yggdrasil.GetObjectState(object.ObjectId, object.State)
	if err == nil {
		if before.ObjectType == static.YggObjectTypeFollowbattlenpc {
			y.RemoveBattleNpc(before.ObjectParam)
		}
	}
	// 设置状态
	object.State = nextState
	// 重新设置地形
	y.RefreshTerrainForObjectChange(*object.Position)
	// battleNpc处理
	after, err := manager.CSV.Yggdrasil.GetObjectState(object.ObjectId, object.State)
	if err == nil {
		if after.ObjectType == static.YggObjectTypeFollowbattlenpc {
			y.AddBattleNpc(after.ObjectParam)
		}
	}

}

// RemoveDiscardGoods 移除丢弃物
func (y *YggEntities) RemoveDiscardGoods(ctx context.Context, userId int64, discardGoods *YggDiscardGoods) error {
	err := DelYggDiscardGoodsInRedis(ctx, userId, discardGoods)
	if err != nil {
		return errors.WrapTrace(err)
	}
	y.removeDiscardGoods(discardGoods)
	return nil
}
func (y *YggEntities) removeDiscardGoods(discardGoods *YggDiscardGoods) {
	delete(y.DiscardGoods, discardGoods.Uid)
	y.EntityFinder.Remove(discardGoods.YggEntity)
	y.Blocks.RemoveDiscardGoods(discardGoods)
}

// RemoveMark 移除标记
func (y *YggEntities) RemoveMark(mark *YggMark) {
	delete(y.Marks, mark.Uid)
	y.EntityFinder.Remove(mark.YggEntity)
	y.Blocks.RemoveMark(mark)
}

func (y *YggEntities) IsBlankPos(position coordinate.Position) bool {
	//如果不是区域的点则为false
	_, err := manager.CSV.Yggdrasil.GetPosAreaId(position)
	if err != nil {
		return false
	}
	posType := y.getPosType(position)
	if posType == entry.PosTypeUnWalkable {
		return false
	}
	uidList, ok := y.EntityFinder.Find(position)
	if !ok {
		return true
	}
	for _, uid := range uidList {
		_, ok = y.Builds[uid]
		if ok {
			return false
		}
		_, ok = y.CoBuilds[uid]
		if ok {
			return false
		}
		_, ok = y.DiscardGoods[uid]
		if ok {
			return false
		}
		object, ok := y.Objects[uid]
		if ok {
			// 跟隨npc所在位置算空
			state, err := manager.CSV.Yggdrasil.GetObjectState(object.ObjectId, object.State)
			if err != nil {
				glog.Errorf("YggEntities IsBlankPos err:%v", err)
				continue
			}
			if state.ObjectType == static.YggObjectTypeFollownpc || state.ObjectType == static.YggObjectTypeFollowbattlenpc {
				continue
			}
			return false
		}
	}

	return true
}
func (y *YggEntities) GetPosHeight(position coordinate.Position) int32 {
	height, ok := y.PosHeight[position]
	if !ok {
		height = entry.DefaultHeight
	}

	return height
}

func (y *YggEntities) getPosType(position coordinate.Position) int32 {
	t, ok := y.PosTypeWithObject[position]
	if !ok {
		t = entry.PosTypeWalkable
	}
	return t
}

func (y *YggEntities) BuildablePos(position coordinate.Position) bool {
	_, err := manager.CSV.Yggdrasil.GetPosAreaId(position)
	if err != nil {
		glog.Errorf("CanPosBuild GetPosAreaId err:%+v", err)
		return false
	}
	if y.getPosType(position) == entry.PosTypeUnBuildable {
		return false
	}
	return true
}

func (y *YggEntities) WalkablePos(position coordinate.Position) error {
	_, err := manager.CSV.Yggdrasil.GetPosAreaId(position)
	if err != nil {
		return errors.WrapTrace(err)
	}

	if y.getPosType(position) == entry.PosTypeUnWalkable {
		return common.ErrYggdrasilCannotMovePosUnWalkable
	}
	return nil
}

func (y *YggEntities) FindDiscardGoodsByPos(position coordinate.Position) (*YggDiscardGoods, bool) {
	uidList, ok := y.EntityFinder.Find(position)
	if !ok {
		return nil, false
	}
	for _, i := range uidList {
		goods, ok := y.DiscardGoods[i]
		if ok {
			return goods, true
		}
	}
	return nil, false
}

func (y *YggEntities) FindBuildingsByPos(position coordinate.Position) (*YggBuild, bool) {
	uidList, ok := y.EntityFinder.Find(position)
	if !ok {
		return nil, false
	}
	for _, i := range uidList {
		build, ok := y.Builds[i]
		if ok {
			return build, true
		}
	}

	return nil, false
}

func (y *YggEntities) FindCoBuildingsByPos(position coordinate.Position) (*YggCoBuild, bool) {
	uidList, ok := y.EntityFinder.Find(position)
	if !ok {
		return nil, false
	}
	for _, i := range uidList {
		coBuild, ok := y.CoBuilds[i]
		if ok {
			return coBuild, true
		}
	}

	return nil, false
}

func (y *YggEntities) FindMessageByPos(position coordinate.Position) (*YggMessage, bool) {
	uidList, ok := y.EntityFinder.Find(position)
	if !ok {
		return nil, false
	}
	for _, i := range uidList {
		m, ok := y.Messages[i]
		if ok {
			return m, true
		}
	}

	return nil, false
}

func checkObjectState(object *YggObject, stateTypes ...int32) (*entry.YggdrasilObjectState, error) {
	config, err := manager.CSV.Yggdrasil.GetYggdrasilObjectConfig(object.ObjectId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	state, ok := config.Sates[object.State]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, entry.CfgNameYggdrasilObjectState, object.State)
	}

	pass := false
	for _, stateType := range stateTypes {
		if stateType == state.ObjectType {
			pass = true
		}
	}
	if !pass {
		return nil, errors.WrapTrace(errors.Swrapf(common.ErrYggdrasilObjectTypeError, stateTypes))
	}

	return &state, nil
}

func (y *YggEntities) FindObjectByPosAndType(position coordinate.Position, stateTypes ...int32) (*YggObject, error) {
	objs, ok := y.FindObjectsByPos(position)
	if !ok {
		return nil, errors.WrapTrace(common.ErrYggdrasilObjectNotFound)
	}
	for _, obj := range objs {
		_, err := checkObjectState(obj, stateTypes...)
		if err != nil {
			continue
		}
		return obj, err
	}

	return nil, errors.WrapTrace(common.ErrYggdrasilObjectNotFound)
}

// obj会有重叠
func (y *YggEntities) FindObjectsByPos(position coordinate.Position) ([]*YggObject, bool) {
	uidList, ok := y.EntityFinder.Find(position)
	if !ok {
		return nil, false
	}
	var objs []*YggObject
	for _, i := range uidList {
		object, ok := y.Objects[i]
		if ok {
			objs = append(objs, object)
		}
	}
	if len(objs) == 0 {
		return nil, false
	}
	return objs, true
}

func (y *YggEntities) FindObjectByUid(uid int64) (*YggObject, bool) {
	object, ok := y.Objects[uid]
	return object, ok
}

func (y *YggEntities) FindObjectById(objectId int32) ([]*YggObject, bool) {
	// 怪物的obj的id可以重复，但是比如任务宝箱那些还是不重的
	var objs []*YggObject
	for _, object := range y.Objects {
		if object.ObjectId == objectId {
			objs = append(objs, object)
		}
	}
	if len(objs) == 0 {
		return nil, false
	}
	return objs, true
}

func (y *YggEntities) FindMarkByUid(markUid int64) (*YggMark, bool) {
	mark, ok := y.Marks[markUid]
	return mark, ok
}

func (y *YggEntities) AppendMatchEntity(entity YggMatchEnable, yggdrasil *Yggdrasil) {

	switch entity.(type) {
	case *YggMessage:
		y.appendMessage(entity.(*YggMessage))
	case *YggDiscardGoods:
		y.appendDiscardGoods(entity.(*YggDiscardGoods))
	case *YggBuild:
		y.appendBuild(entity.(*YggBuild))
	}
	yggdrasil.MatchUserIds[entity.GetMatchUserId()] = struct{}{}
}

// AddBattleNpc 增加助战npc
func (y *YggEntities) AddBattleNpc(npcId int32) {
	_, ok := y.NpcHp[npcId]
	if !ok {
		y.NpcHp[npcId] = 10000
	}
}

// RemoveBattleNpc 移除助战npc
func (y *YggEntities) RemoveBattleNpc(npcId int32) {
	delete(y.NpcHp, npcId)
}

// RemoveEnvTerrainCreateAt 移除env产生的地形
func (y *YggEntities) RemoveEnvTerrainCreateAt(createAt int32, yggdrasil *Yggdrasil) {
	var replace []*YggEnvTerrain
	for _, terrain := range y.EnvTerrains {
		if terrain.CreateAt == createAt {
			yggdrasil.EntityChange.RemoveEnvTerrain = append(yggdrasil.EntityChange.RemoveEnvTerrain, terrain.EnvId)
		} else {
			replace = append(replace, terrain)
		}
	}
	y.EnvTerrains = replace
	y.RefreshTerrain()

}

// RemoveEnvTerrainDeleteAt 移除env产生的地形
func (y *YggEntities) RemoveEnvTerrainDeleteAt(deleteAt int32, yggdrasil *Yggdrasil) {
	var replace []*YggEnvTerrain
	for _, terrain := range y.EnvTerrains {
		if terrain.DeleteAt == deleteAt {
			yggdrasil.EntityChange.RemoveEnvTerrain = append(yggdrasil.EntityChange.RemoveEnvTerrain, terrain.EnvId)
		} else {
			replace = append(replace, terrain)
		}
	}
	y.EnvTerrains = replace
	y.RefreshTerrain()

}

// CreateEnvTerrain 创建env产生的地形
func (y *YggEntities) CreateEnvTerrain(createAt, deleteAt, envId int32, yggdrasil *Yggdrasil) {
	y.EnvTerrains = append(y.EnvTerrains, NewYggEnvTerrain(envId, createAt, deleteAt))

	yggdrasil.EntityChange.AddEnvTerrain = append(yggdrasil.EntityChange.AddEnvTerrain, envId)
	y.RefreshTerrain()

}

// RefreshTerrain 刷新地形
func (y *YggEntities) RefreshTerrain() {
	y.RefreshEnvTerrain()
	// obj会改变初始地形
	for _, object := range y.Objects {
		y.PosTypeWithObject[*object.Position] = y.CalPosType(*object.Position)
	}
}

func (y *YggEntities) RefreshTerrainForObjectChange(p coordinate.Position) {
	glog.Debugf("RefreshTerrainForObjectChange pos:%v,posTypeBefore:%d", p, y.getPosType(p))
	y.PosTypeWithObject[p] = y.CalPosType(p)
	glog.Debugf("RefreshTerrainForObjectChange pos:%v,posTypeAfter:%d", p, y.getPosType(p))

}

func (y *YggEntities) RefreshEnvTerrain() {
	y.PosHeight = map[coordinate.Position]int32{}
	y.PosType = map[coordinate.Position]int32{}
	initHeight, initType := manager.CSV.Yggdrasil.GetPosInitHeightType()
	// 初始高度
	for position, i := range initHeight {
		y.PosHeight[position] = i
	}
	// 初始地形
	for position, i := range initType {
		y.PosType[position] = i
	}

	for _, terrain := range y.EnvTerrains {
		env, err := manager.CSV.Yggdrasil.GetEnv(terrain.EnvId)
		if err != nil {
			glog.Errorf("RefreshEnvTerrain err:%v", err)
			continue
		}
		if env.Terrains == nil {
			continue

		}
		for _, heightTerrain := range env.Terrains.HeightTerrains {
			for _, position := range heightTerrain.Points() {
				y.PosHeight[position] = heightTerrain.PosHeight
			}
		}
		for _, typeTerrain := range env.Terrains.TypeTerrains {
			for _, position := range typeTerrain.Points() {
				y.PosType[position] = typeTerrain.PosType
			}
		}
	}

}

func (y *YggEntities) GetEnvTerrain() []int32 {
	ret := make([]int32, 0, len(y.EnvTerrains))
	for _, envTerrain := range y.EnvTerrains {
		ret = append(ret, envTerrain.EnvId)
	}
	return ret
}

// CheckElevatorCanBuild 位置是否可以建造电梯
func (y *YggEntities) CheckElevatorCanBuild(pos coordinate.Position) error {
	cubeRing := coordinate.CubeRing(pos, 1)
	height := y.GetPosHeight(pos)
	for _, position := range cubeRing {
		posType := y.getPosType(pos)
		if posType == entry.PosTypeUnWalkable {
			continue
		}
		if y.GetPosHeight(position)-height >= 1 {
			return nil
		}

	}
	return errors.WrapTrace(common.ErrYggdrasilCannotMoveTerrainDiff)
}

// YggEntityFinder 查找坐标得到entities
type YggEntityFinder struct {
	Entities map[coordinate.Position]*number.NonRepeatableArrInt64 // 记录position 和uid的关系
}

func NewYggEntityFinder() *YggEntityFinder {
	return &YggEntityFinder{
		Entities: map[coordinate.Position]*number.NonRepeatableArrInt64{},
	}
}

func (y *YggEntityFinder) Find(position coordinate.Position) ([]int64, bool) {
	v, ok := y.Entities[position]
	if !ok {
		return nil, false
	}
	return v.Values(), true
}
func (y *YggEntityFinder) Append(entities ...*YggEntity) {
	for _, entity := range entities {

		list, ok := y.Entities[*entity.Position]
		if !ok {
			list = number.NewNonRepeatableArrInt64()
			y.Entities[*entity.Position] = list
		}
		list.Append(entity.Uid)
	}
}

func (y *YggEntityFinder) Remove(entities ...*YggEntity) {
	for _, entity := range entities {
		list, ok := y.Entities[*entity.Position]
		if ok {
			list.Remove(entity.Uid)
			if list.IsEmpty() {
				delete(y.Entities, *entity.Position)
			}
		}
	}
}

type YggBlock struct {
	Corner       *coordinate.Position //左下角坐标
	Areas        *number.NonRepeatableArr
	Builds       map[int64]*YggBuild
	CoBuilds     map[int64]*YggCoBuild
	Messages     map[int64]*YggMessage
	DiscardGoods map[int64]*YggDiscardGoods
	Objects      map[int64]*YggObject
	Marks        map[int64]*YggMark
}

func (y *YggBlock) VOYggdrasilBlock(ctx context.Context, unlockArea *common.Area) *pb.VOYggdrasilBlock {
	h, v := manager.CSV.Yggdrasil.GetYggBlockLengthAndWidth()
	area := unlockArea.CutRectangle(*y.Corner, h, v)
	var UnlockPosList []*pb.VOPosition

	for _, position := range area.Points() {
		UnlockPosList = append(UnlockPosList, position.VOPosition())
	}
	Builds := make([]*pb.VOYggdrasilBuild, 0, len(y.Builds))
	for _, build := range y.Builds {
		yggdrasilBuild, err := build.VOYggdrasilBuild(ctx)
		if err != nil {
			glog.Errorf("VOYggdrasilBlock VOYggdrasilBuild  err:%+v", err)
			continue
		}

		Builds = append(Builds, yggdrasilBuild)
	}

	DiscardGoods := make([]*pb.VOYggdrasilDiscardGoods, 0, len(y.DiscardGoods))
	for _, good := range y.DiscardGoods {
		DiscardGoods = append(DiscardGoods, good.VOYggdrasilDiscardGoods())
	}
	Objects := make([]*pb.VOYggdrasilObject, 0, len(y.Objects))
	for _, object := range y.Objects {
		Objects = append(Objects, object.VOYggdrasilObject())
	}
	Marks := make([]*pb.VOYggdrasilMark, 0, len(y.Marks))
	for _, mark := range y.Marks {
		Marks = append(Marks, mark.VOYggdrasilMark())
	}
	Messages := make([]*pb.VOYggdrasilMessage, 0, len(y.Messages))
	for _, message := range y.Messages {
		Messages = append(Messages, message.VOYggdrasilMessage())
	}
	CoBuilds := make([]*pb.VOYggdrasilCoBuildBase, 0, len(y.CoBuilds))
	for _, coBuild := range y.CoBuilds {
		CoBuilds = append(CoBuilds, coBuild.VOYggdrasilCoBuildBase(coBuild.BuildId, []int64{}, coBuild.TotalUseCount))
	}

	return &pb.VOYggdrasilBlock{
		Position:      y.Corner.VOPosition(),
		UnlockPosList: UnlockPosList,
		Builds:        Builds,
		DiscardGoods:  DiscardGoods,
		Objects:       Objects,
		Marks:         Marks,
		Messages:      Messages,
		CoBuilds:      CoBuilds,
	}

}

func NewYggBlock(corner coordinate.Position) *YggBlock {

	h, v := manager.CSV.Yggdrasil.GetYggBlockLengthAndWidth()
	Areas := number.NewNonRepeatableArr()

	// 计算block 包含的区域
	for i := corner.X; i < corner.X+h; i++ {
		for j := corner.Y; j < corner.Y+v; j++ {
			areaId, err := manager.CSV.Yggdrasil.GetPosAreaId(*coordinate.NewPosition(i, j))
			if err == nil {
				Areas.Append(areaId)
			}
		}
	}

	return &YggBlock{
		Corner:       &corner,
		Areas:        Areas,
		Builds:       map[int64]*YggBuild{},
		CoBuilds:     map[int64]*YggCoBuild{},
		Messages:     map[int64]*YggMessage{},
		DiscardGoods: map[int64]*YggDiscardGoods{},
		Objects:      map[int64]*YggObject{},
		Marks:        map[int64]*YggMark{},
	}
}

type YggBlocks map[coordinate.Position]*YggBlock

func NewYggBlocks() *YggBlocks {
	return (*YggBlocks)(&map[coordinate.Position]*YggBlock{})
}

func (y *YggBlocks) getByCreate(position coordinate.Position) *YggBlock {

	corner := calCorner(position)
	block, ok := (*y)[corner]
	if !ok {
		block = NewYggBlock(corner)
		(*y)[corner] = block
	}
	return block
}

func (y *YggBlocks) AppendBuild(build *YggBuild) {
	block := y.getByCreate(*build.Position)
	block.Builds[build.Uid] = build
}
func (y *YggBlocks) AppendMessage(msg *YggMessage) {
	block := y.getByCreate(*msg.Position)
	block.Messages[msg.Uid] = msg

}
func (y *YggBlocks) AppendCoBuild(coBuild *YggCoBuild) {
	block := y.getByCreate(*coBuild.Position)
	block.CoBuilds[coBuild.Uid] = coBuild
}

func (y *YggBlocks) AppendObject(object *YggObject) {
	block := y.getByCreate(*object.Position)
	block.Objects[object.Uid] = object
}

func (y *YggBlocks) AppendDiscardGoods(discardGoods *YggDiscardGoods) {
	block := y.getByCreate(*discardGoods.Position)
	block.DiscardGoods[discardGoods.Uid] = discardGoods

}
func (y *YggBlocks) AppendMark(mark *YggMark) {
	block := y.getByCreate(*mark.Position)
	block.Marks[mark.Uid] = mark
}

func (y *YggBlocks) RemoveBuild(build *YggBuild) {
	block := y.getByCreate(*build.Position)
	delete(block.Builds, build.Uid)
}
func (y *YggBlocks) RemoveMessage(msg *YggMessage) {
	block := y.getByCreate(*msg.Position)
	delete(block.Messages, msg.Uid)
}
func (y *YggBlocks) RemoveObject(object *YggObject) {
	block := y.getByCreate(*object.Position)
	delete(block.Objects, object.Uid)
}

func (y *YggBlocks) RemoveDiscardGoods(discardGoods *YggDiscardGoods) {
	block := y.getByCreate(*discardGoods.Position)
	delete(block.DiscardGoods, discardGoods.Uid)
}

func (y *YggBlocks) RemoveMark(mark *YggMark) {
	block := y.getByCreate(*mark.Position)
	delete(block.Marks, mark.Uid)
}
