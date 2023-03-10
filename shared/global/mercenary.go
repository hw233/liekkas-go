package global

import (
	"context"
	"encoding/json"
	"shared/common"
	"shared/utility/errors"
	"shared/utility/global"
	"shared/utility/glog"
	"shared/utility/key"
	"time"

	"github.com/go-redis/redis/v8"
)

// todo 设置过期时间

const (
	KeyMercenary           = "Mercenary"
	KeyMercenaryOwn        = "MercenaryOwn"
	KeyMercenaryAvaliable  = "MercenaryAvaliable"
	KeyMercenarySendRecord = "MercenarySendRecord"
	KeyMercenaryApply      = "MercenaryApply"
	KeyMercenaryApplyCount = "MercenaryApplyCount"
	KeyCharacterList       = "MercenaryCharacterList"
	KeyMercenaryLock       = "MercenaryLock"
)

type UserMercenary struct {
	*Mercenary
	*MercenaryOwn
	*MercenaryAvailable
	*MercenaryApply
	// MercenarySendRecord
	*MercenaryApplyCount
	*CharacterList
}

// 玩家的角色数据，所有可借的佣兵
type Mercenary struct {
	key    string
	client *redis.Client
}

// 玩家可借的佣兵是否已经借出
type MercenaryAvailable struct {
	key    string
	client *redis.Client
	global *global.Global
}

// 玩家可借的佣兵的申请人数量
type MercenaryApply struct {
	key    string
	client *redis.Client
	global *global.Global
}

// 玩家可借的佣兵的申请人数量
type MercenaryApplyCount struct {
	key    string
	client *redis.Client
	global *global.Global
}

type CharacterList struct {
	key    string
	client *redis.Client
	global *global.Global
}

func NewUserMercenary(client *redis.Client) *UserMercenary {
	return &UserMercenary{
		Mercenary:          NewMercenary(client),
		MercenaryOwn:       NewMercenaryOwn(client),
		MercenaryAvailable: NewMercenaryAvaliable(client),
		MercenaryApply:     NewMercenaryApply(client),
		// MercenarySendRecord: NewMercenarySendRecord(client),
		MercenaryApplyCount: NewMercenaryApplyCount(client),
		CharacterList:       NewCharacterList(client),
	}
}

func NewMercenary(client *redis.Client) *Mercenary {
	return &Mercenary{
		key:    KeyMercenary,
		client: client,
	}
}

func NewMercenaryApply(client *redis.Client) *MercenaryApply {
	return &MercenaryApply{
		key:    KeyMercenaryApply,
		client: client,
		global: global.NewGlobal(client),
	}
}

func NewMercenaryApplyCount(client *redis.Client) *MercenaryApplyCount {
	return &MercenaryApplyCount{
		key:    KeyMercenaryApplyCount,
		client: client,
		global: global.NewGlobal(client),
	}
}

func NewMercenaryAvaliable(client *redis.Client) *MercenaryAvailable {
	return &MercenaryAvailable{
		key:    KeyMercenaryAvaliable,
		client: client,
		global: global.NewGlobal(client),
	}
}

func NewCharacterList(client *redis.Client) *CharacterList {
	return &CharacterList{
		key:    KeyCharacterList,
		client: client,
		global: global.NewGlobal(client),
	}
}

// -----------------Mercenary----------------

// 获得角色数据
func (um *Mercenary) GetMercenaryCharacter(ctx context.Context, userId int64) (map[int32]*common.MercenaryCharacter, error) {
	characters := map[int32]*common.MercenaryCharacter{}

	values, err := um.client.HGetAll(ctx, key.MakeRedisKey(um.key, userId)).Result()
	if err != nil {
		return nil, err
	}
	for _, val := range values {
		character := &common.MercenaryCharacter{}
		err := json.Unmarshal([]byte(val), character)
		if err != nil {
			return nil, err
		}
		characters[character.Id] = character
	}

	return characters, nil
}

// 每隔十分钟去修改一次，只有自己会去写
func (um *Mercenary) SetMercenaryCharacterData(ctx context.Context, userId int64, characters []*common.MercenaryCharacter) error {
	hashKey := key.MakeRedisKey(um.key, userId)

	for _, character := range characters {
		bytes, err := json.Marshal(character)
		if err != nil {
			return err
		}

		err = um.client.HSet(ctx, hashKey, key.MakeRedisKey(character.Id), bytes).Err()
		if err != nil {
			return errors.WrapTrace(err)
		}
	}

	return nil
}

// -----------------MercenaryOwn------------------
// 玩家已有的借到的佣兵
type MercenaryOwn struct {
	key    string
	client *redis.Client
	Global *global.Global
}

func NewMercenaryOwn(client *redis.Client) *MercenaryOwn {
	return &MercenaryOwn{
		key:    KeyMercenaryOwn,
		client: client,
		Global: global.NewGlobal(client),
	}
}

func (uo *MercenaryOwn) IsMercenaryOwn(ctx context.Context, userId int64, characterId int32) (bool, error) {
	result, err := uo.client.SIsMember(ctx, key.MakeRedisKey(uo.key, userId), characterId).Result()
	if err != nil {
		return false, err
	}
	return result, nil
}

func (uo *MercenaryOwn) GetLengthMercenaryOwn(ctx context.Context, userId int64) (int32, error) {

	count, err := uo.client.SCard(ctx, key.MakeRedisKey(uo.key, userId)).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, errors.WrapTrace(err)
	}

	return int32(count), nil
}

// 返回值为1，代表put失败，该用户已经借到角色或者该用户佣兵已满;返回值为0，代表成功
func (uo *MercenaryOwn) PutMercenaryOwn(ctx context.Context, userId int64, characterId int32, tm time.Time) (bool, error) {
	hashKey := key.MakeRedisKey(uo.key, userId)

	gLock, err := uo.Global.ObtainLock(ctx, hashKey)
	if err != nil {
		return false, errors.WrapTrace(err)
	}
	defer gLock.Release()

	isMember, err := uo.IsMercenaryOwn(ctx, userId, characterId)
	if err != nil {
		return false, err
	}
	count, err := uo.GetLengthMercenaryOwn(ctx, userId)
	if err != nil {
		return false, err
	}
	// todo 常量
	if isMember || (count >= 3) {
		return false, nil
	}

	err = uo.client.SAdd(ctx, hashKey, characterId).Err()
	if err != nil {
		return false, err
	}
	err = uo.client.ExpireAt(ctx, hashKey, tm).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (uo *MercenaryOwn) ClearMercenaryOwn(ctx context.Context, userId int64) error {
	hashKey := key.MakeRedisKey(uo.key, userId)

	gLock, err := uo.Global.ObtainLock(ctx, hashKey)
	if err != nil {
		return errors.WrapTrace(err)
	}
	defer gLock.Release()

	err = uo.client.Del(ctx, hashKey).Err()
	if err != nil {
		return err
	}
	return nil
}

// -----------------MercenaryAvaliable------------------

// 返回值为true，说明可用；返回值为false，说明这个角色已经被人用了
func (ua *MercenaryAvailable) CheckMercenaryAvaliable(ctx context.Context, userId int64, characterId int32) (bool, error) {
	hashKey := key.MakeRedisKey(ua.key, userId)

	IsMember, err := ua.client.SIsMember(ctx, hashKey, characterId).Result()
	if err != nil {
		if err == redis.Nil {
			return true, nil
		}
		return false, errors.WrapTrace(err)
	}

	return !IsMember, nil
}

func (ua *MercenaryAvailable) ChangeMercenaryAvaliable(ctx context.Context, userId int64, characterId int32, tm time.Time) error {
	hashKey := key.MakeRedisKey(ua.key, userId)

	gLock, err := ua.global.ObtainLock(ctx, hashKey)
	if err != nil {
		return errors.WrapTrace(err)
	}
	defer gLock.Release()

	err = ua.client.SAdd(ctx, hashKey, characterId).Err()
	if err != nil {
		return errors.WrapTrace(err)
	}
	err = ua.client.ExpireAt(ctx, hashKey, tm).Err()
	if err != nil {
		return errors.WrapTrace(err)
	}
	return nil
}

func (ua *MercenaryAvailable) Clear(ctx context.Context, userId int64) error {
	hashKey := key.MakeRedisKey(ua.key, userId)

	gLock, err := ua.global.ObtainLock(ctx, hashKey)
	if err != nil {
		return errors.WrapTrace(err)
	}
	defer gLock.Release()

	err = ua.client.Del(ctx, key.MakeRedisKey(ua.key, userId)).Err()
	if err != nil {
		return errors.WrapTrace(err)
	}
	return nil
}

// -----------------MercenaryApplyCount------------------

func (uc *MercenaryApplyCount) MercenaryApplyGet(ctx context.Context, userId int64, characterId int32) (int32, error) {

	count, err := uc.client.HGet(ctx, key.MakeRedisKey(uc.key, userId), key.MakeRedisKey(characterId)).Int()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, errors.WrapTrace(err)
	}

	return int32(count), nil
}

// 修改申请人数
func (uc *MercenaryApplyCount) MercenaryApplyAdd(ctx context.Context, userId int64, characterId int32, val int32, tm time.Time) error {
	hashKey := key.MakeRedisKey(uc.key, userId)

	gLock, err := uc.global.ObtainLock(ctx, hashKey)
	if err != nil {
		return errors.WrapTrace(err)
	}
	defer gLock.Release()

	count, err := uc.MercenaryApplyGet(ctx, userId, characterId)
	if err != nil {
		return errors.WrapTrace(err)
	}

	count += val
	if count < 0 { // 申请人数不能为负
		count = 0
	}

	err = uc.client.HSet(ctx, hashKey, key.MakeRedisKey(characterId), count).Err()
	if err != nil {
		return errors.WrapTrace(err)
	}
	err = uc.client.ExpireAt(ctx, hashKey, tm).Err()
	if err != nil {
		return errors.WrapTrace(err)
	}

	return nil
}

// 将申请人数置0
func (uc *MercenaryApplyCount) MercenaryApplyReset(ctx context.Context, userId int64, characterId int32) error {
	hashKey := key.MakeRedisKey(uc.key, userId)

	gLock, err := uc.global.ObtainLock(ctx, hashKey)
	if err != nil {
		return errors.WrapTrace(err)
	}
	defer gLock.Release()

	return uc.client.HSet(ctx, hashKey, key.MakeRedisKey(characterId), 0).Err()
}

// -----------CharacterList----------------
func (cl *CharacterList) CharacterListPush(ctx context.Context, userIdFrom int64, name string, userIdTo int64, relation int8, c *common.CharacterData, tm time.Time) error {
	hashKey := key.MakeRedisKey(cl.key, userIdTo)

	gLock, err := cl.global.ObtainLock(ctx, hashKey)
	if err != nil {
		return errors.WrapTrace(err)
	}
	defer gLock.Release()

	characterData := common.NewMercenaryData(userIdFrom, name, relation, c)

	bytes, err := json.Marshal(characterData)
	if err != nil {
		return errors.WrapTrace(err)
	}
	_, err = cl.client.RPush(ctx, hashKey, bytes).Result()
	if err != nil {
		return errors.WrapTrace(err)
	}
	err = cl.client.ExpireAt(ctx, hashKey, tm).Err()
	if err != nil {
		return errors.WrapTrace(err)
	}

	return nil
}

// 使用前加分布式锁，用后调用clear清理list
func (cl *CharacterList) CharacterListGetAndClear(ctx context.Context, userId int64) []*common.MercenaryData {
	hashKey := key.MakeRedisKey(cl.key, userId)

	var characters []*common.MercenaryData

	gLock, err := cl.global.ObtainLock(ctx, hashKey)
	if err != nil {
		glog.Errorf("MercenaryCharacterList ObtainLock err:%+v", err)
		return characters
	}
	defer gLock.Release()

	results, err := cl.client.LRange(ctx, hashKey, 0, -1).Result() //读list所有元素
	if err != nil {
		glog.Errorf("CharacterListPop LRange err:%+v", err)
		return characters
	}
	for _, r := range results {
		var character = &common.MercenaryData{}
		err := json.Unmarshal([]byte(r), character)
		if err != nil {
			glog.Errorf("CharacterListPop json Unmarshal err:%+v", err)
			continue
		}
		characters = append(characters, character)
	}
	cl.CharacterListClear(ctx, userId)
	// fmt.Println("======characterlist=====afterclear==========")

	return characters
}

func (cl *CharacterList) CharacterListClear(ctx context.Context, userId int64) {
	hashKey := key.MakeRedisKey(cl.key, userId)
	_, err := cl.client.Del(ctx, hashKey).Result()
	if err != nil {
		glog.Errorf("CharacterListPop Del err:%v", errors.WrapTrace(err))
		return
	}
}

// ------------MercenaryApply---------------
func (ma *MercenaryApply) MercenaryApplySend(ctx context.Context, userTo int64, apply *common.MercenarySend, tm time.Time) error {
	hashKey := key.MakeRedisKey(ma.key, userTo)

	gLock, err := ma.global.ObtainLock(ctx, hashKey)
	if err != nil {
		return errors.WrapTrace(err)
	}
	defer gLock.Release()

	bytes, err := json.Marshal(apply)
	if err != nil {
		return errors.WrapTrace(err)
	}

	_, err = ma.client.RPush(ctx, hashKey, bytes).Result()
	if err != nil {
		return errors.WrapTrace(err)
	}
	err = ma.client.ExpireAt(ctx, hashKey, tm).Err()
	if err != nil {
		return errors.WrapTrace(err)
	}

	return nil
}
func (ma *MercenaryApply) MercenaryApplyCancel(ctx context.Context, userFrom int64, name string, userTo int64, characterId int32, timestamp int64) error {
	hashKey := key.MakeRedisKey(ma.key, userTo)

	gLock, err := ma.global.ObtainLock(ctx, hashKey)
	if err != nil {
		return errors.WrapTrace(err)
	}
	defer gLock.Release()

	send := common.NewMercenarySend(userFrom, name, characterId)
	send.SendTime = timestamp
	bytes, err := json.Marshal(send)
	if err != nil {
		return errors.WrapTrace(err)
	}

	err = ma.client.LRem(ctx, hashKey, 0, bytes).Err()
	if err != nil {
		return errors.WrapTrace(err)
	}

	return nil
}

func (ma *MercenaryApply) MercenaryApplyReceive(ctx context.Context, userId int64) ([]*common.MercenarySend, error) {
	hashKey := key.MakeRedisKey(ma.key, userId)

	var send []*common.MercenarySend

	gLock, err := ma.global.ObtainLock(ctx, hashKey)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	defer gLock.Release()

	results, err := ma.client.LRange(ctx, hashKey, 0, -1).Result() //读list所有元素
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	for _, r := range results {
		var apply = &common.MercenarySend{}
		err := json.Unmarshal([]byte(r), apply)
		if err != nil {
			glog.Errorf("MercenaryApplyReceive json Unmarshal err:%+v", err)
			continue
		}
		send = append(send, apply)
	}

	_, err = ma.client.Del(ctx, hashKey).Result()
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return send, nil
}

// --------------MercenarySendRecord----------------
// type MercenarySendRecord struct {
// 	key    string
// 	client *redis.Client
// }

// func NewMercenarySendRecord(client *redis.Client) *MercenarySendRecord {
// 	return &MercenarySendRecord{
// 		key:    KeyMercenarySendRecord,
// 		client: client,
// 	}
// }

// func (m *MercenarySendRecord) MercenarySendRecordAdd(ctx context.Context, userFrom int64, record *common.MercenarySendRecord, tm time.Time) error {
// 	hashKey := key.MakeRedisKey(m.key, userFrom)

// 	bytes, err := json.Marshal(record)
// 	if err != nil {
// 		return errors.WrapTrace(err)
// 	}

// 	// isMember, err := m.MercenarySendRecrodCheck(ctx, userFrom, record)
// 	// if err != nil {
// 	// 	return errors.WrapTrace(err)
// 	// }

// 	count, err := m.client.SCard(ctx, hashKey).Result()
// 	if err != nil {
// 		return errors.WrapTrace(err)
// 	}
// 	// todo 常量
// 	// if isMember || count >= 5 {
// 	// 	return errors.Swrapf(common.ErrMercenaryApplyExceedLimit)
// 	// }
// 	if count >= 5 {
// 		return errors.Swrapf(common.ErrMercenaryApplyExceedLimit)
// 	}

// 	err = m.client.SAdd(ctx, hashKey, bytes).Err()
// 	if err != nil {
// 		return errors.WrapTrace(err)
// 	}

// 	err = m.client.ExpireAt(ctx, hashKey, tm).Err()
// 	if err != nil {
// 		return errors.WrapTrace(err)
// 	}

// 	return nil
// }

// func (m *MercenarySendRecord) MercenarySendRecordDelete(ctx context.Context, userFrom int64, record *common.MercenarySendRecord) error {
// 	hashKey := key.MakeRedisKey(m.key, userFrom)

// 	bytes, err := json.Marshal(record)
// 	if err != nil {
// 		return errors.WrapTrace(err)
// 	}
// 	err = m.client.SRem(ctx, hashKey, bytes).Err()
// 	if err != nil {
// 		return errors.WrapTrace(err)
// 	}

// 	return nil
// }

// // 返回true，说明有这个申请；返回值为false，说明没有这个申请
// func (m *MercenarySendRecord) MercenarySendRecrodCheck(ctx context.Context, userFrom int64, record *common.MercenarySendRecord) (bool, error) {
// 	hashKey := key.MakeRedisKey(m.key, userFrom)

// 	bytes, err := json.Marshal(record)
// 	if err != nil {
// 		return false, errors.WrapTrace(err)
// 	}

// 	result, err := m.client.SIsMember(ctx, hashKey, bytes).Result()
// 	if err != nil {
// 		return false, errors.WrapTrace(err)
// 	}

// 	return result, nil
// }
