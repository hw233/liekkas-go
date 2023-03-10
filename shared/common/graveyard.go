package common

import (
	"math"
	"shared/protobuf/pb"
	"shared/utility/coordinate"
	"shared/utility/servertime"
	"shared/utility/uid"
	"sort"
)

const (
	CurtainTypeCreate  = 1
	CurtainTypeLvUp    = 2
	CurtainTypeStageUp = 3

	GraveyardCountBuffType = 1
	GraveyardTimeBuffType  = 2

	RequestStateBefore = 1
	RequestStateDuring = 2
)

type UserGraveyardBuild struct {
	*UserGraveyardBuildBase  `json:"base"`       // 位置，等级阶级等信息
	*UserGraveyardProduce    `json:"produce"`    // 如果是产出型建筑 则有此数据， 主堡和装饰建筑无此数据
	*UserGraveyardTransition `json:"transition"` // 如果在 建造，升级升阶等状态时，有此数据，揭幕后置为空
	*UserGraveyardRequest    `json:"request"`    // 公会互助相关
}

func (u *UserGraveyardBuild) FetchLevel() int32 {
	// 如果是建造状态中 返回0级
	if u.UserGraveyardTransition != nil && u.UserGraveyardTransition.CurtainType == CurtainTypeCreate {
		return 0
	}
	// 如果是升级状态中 返回上一级
	if u.UserGraveyardTransition != nil && u.UserGraveyardTransition.CurtainType == CurtainTypeLvUp {
		return u.Lv - 1
	}
	return u.Lv
}
func (u *UserGraveyardBuild) FetchStage() int32 {
	// 如果是建造状态中 返回0阶
	if u.UserGraveyardTransition != nil && u.UserGraveyardTransition.CurtainType == CurtainTypeCreate {
		return 0
	}
	// 如果是升阶状态中 返回上一阶
	if u.UserGraveyardTransition != nil && u.UserGraveyardTransition.CurtainType == CurtainTypeStageUp {
		return u.Stage - 1
	}
	return u.Stage
}

// FetchCharacters 查询建筑入驻的角色
func (u *UserGraveyardBuild) FetchCharacters() (*CharacterPositions, bool) {
	//只有产出型建筑可以入驻角色
	if u.UserGraveyardProduce != nil {

		return u.UserGraveyardProduce.Characters, true
	}
	return nil, false
}

// SetCharacters 设置建筑入驻的角色
func (u *UserGraveyardBuild) SetCharacters(characters *CharacterPositions) {
	//只有产出型建筑可以入驻角色
	if u.UserGraveyardProduce != nil {
		u.UserGraveyardProduce.Characters = characters
	}
}

// NewUserGraveyardMainTower 主堡
func NewUserGraveyardMainTower(mainTowerBuildId int32, position *coordinate.Position) *UserGraveyardBuild {
	return &UserGraveyardBuild{
		UserGraveyardBuildBase:  NewUserGraveyardBuildBase(mainTowerBuildId, 1, 0, position),
		UserGraveyardTransition: NewUserGraveyardTransition(0, CurtainTypeCreate),
		UserGraveyardRequest:    NewUserGraveyardRequest(),
	}
}

// NewUserGraveyardBuild 产出型建筑
func NewUserGraveyardBuild(buildId int32, position *coordinate.Position, sec int32) *UserGraveyardBuild {
	return &UserGraveyardBuild{
		UserGraveyardBuildBase:  NewUserGraveyardBuildBase(buildId, 1, 0, position),
		UserGraveyardTransition: NewUserGraveyardTransition(sec, CurtainTypeCreate),
		UserGraveyardProduce:    NewUserGraveyardProduce(),
		UserGraveyardRequest:    NewUserGraveyardRequest(),
	}
}

// NewUserGraveyardAdornment 装饰物
func NewUserGraveyardAdornment(buildId int32, position *coordinate.Position) *UserGraveyardBuild {
	return &UserGraveyardBuild{
		// 前端需要装饰物的等级给1级
		UserGraveyardBuildBase: NewUserGraveyardBuildBase(buildId, 1, 0, position),
	}
}

type UserGraveyardTransition struct {
	StartAt           int64 `json:"start_at"`     // 开始时间(打点用)
	EndAt             int64 `json:"end_at"`       // 剩余时间
	TransitionNeedSec int32 `json:"need_sec"`     // 需要多场时间
	CurtainType       int32 `json:"curtain_type"` // 1,2,3分别表示建造，升级，升阶的幕布
}

func NewUserGraveyardTransition(sec int32, curtainType int32) *UserGraveyardTransition {
	return &UserGraveyardTransition{
		StartAt:           servertime.Now().Unix(),
		EndAt:             servertime.Now().Unix() + int64(sec),
		TransitionNeedSec: sec,
		CurtainType:       curtainType,
	}
}

// 获得倒计时秒
func getCountDownSec(EndAt int64) int32 {
	return int32(math.Max(0, float64(EndAt-servertime.Now().Unix())))
}

// 获得计时秒
func getSec(startAt int64) int32 {
	sec := int32(servertime.Now().Unix() - startAt)
	// 往前调时间的容错
	if sec > 0 {
		return sec
	}
	return 0
}

func (u *UserGraveyardTransition) VOGraveyardTransition() *pb.VOGraveyardTransition {
	if u == nil {
		return nil
	}
	return &pb.VOGraveyardTransition{
		DiffSecs:    getCountDownSec(u.EndAt),
		CurtainType: u.CurtainType,
	}
}

type UserGraveyardBuildBase struct {
	BuildId   int32                `json:"build_id"`
	Position  *coordinate.Position `json:"position"` // position为nil的时候，表示在背包
	Lv        int32                `json:"lv"`
	Stage     int32                `json:"stage"`
	InBagTime int64                `json:"in_bag_time"` // position为nil的时候，表示入背包时间
}

func NewUserGraveyardBuildBase(buildId, lv, stage int32, position *coordinate.Position) *UserGraveyardBuildBase {
	return &UserGraveyardBuildBase{
		BuildId:  buildId,
		Position: position,
		Lv:       lv,
		Stage:    stage,
	}
}

func (u *UserGraveyardBuildBase) InBag() bool {
	return u.Position == nil
}

func (u *UserGraveyardBuild) VOGraveyardBase() *pb.VOGraveyardBase {
	var voPosition *pb.VOPosition
	if u.Position != nil {
		voPosition = u.Position.VOPosition()
	}
	return &pb.VOGraveyardBase{
		BuildId:   u.BuildId,
		Position:  voPosition,
		Level:     u.FetchLevel(),
		Stage:     u.FetchStage(),
		InBagTime: u.InBagTime,
	}
}

// ResetConsumeProduce 重置消耗型产出
func (u *UserGraveyardBuild) ResetConsumeProduce() {
	u.Productions = NewRewards()

	// 改成可以一件一件领了，
	if u.CurrProduceNum == u.CurrReceiveProduceNum {
		u.CurrReceiveProduceNum = 0
		u.CurrProduceNum = 0
		u.ReachLimit = false

		u.AccRecords.clear()
		// 重置产出的时 重置buff
		u.InvalidateBuff()
		// 请求状态重置
		u.ResetRequest()

	}

}

// ResetContinuousProduce 重置持续生产产出
func (u *UserGraveyardBuild) ResetContinuousProduce() {
	u.Productions = NewRewards()
	u.CurrProduceNum = 0
	u.ReachLimit = false
	u.AccRecords.clear()
	// 重置产出的时 重置buff
	u.InvalidateBuff()

}

// StartProduce 持续性建筑开始产出
func (u *UserGraveyardBuild) StartProduce() {
	if u.UserGraveyardProduce == nil {
		return
	}
	u.ProduceStartAt = servertime.Now().Unix()
}

// StartConsumeProduce 手动型建筑开始产出
func (u *UserGraveyardBuild) StartConsumeProduce(produceNum int32) {
	u.CurrProduceNum = produceNum
	u.ProduceStartAt = servertime.Now().Unix()
	if u.CountBuff != nil {
		// 开始产出的时候restCount --
		u.CountBuff.RestCount--
	}

}

func (u *UserGraveyardBuild) PopulationSet(population int32) {
	produce := u.UserGraveyardProduce
	if produce != nil {
		produce.PopulationCount = population
	}

}

func (u *UserGraveyardBuild) LvUpTransition(toLv, levelUpTime int32) {
	u.Lv = toLv
	u.UserGraveyardTransition = NewUserGraveyardTransition(levelUpTime, CurtainTypeLvUp)
}

func (u *UserGraveyardBuild) StageUpTransition(toStage int32, stageUpTime int32) {
	u.Stage = toStage

	u.UserGraveyardTransition = NewUserGraveyardTransition(stageUpTime, CurtainTypeStageUp)
}

func (u *UserGraveyardBuild) EndTransition() {
	u.UserGraveyardTransition = nil
	// 请求状态重置
	u.ResetRequest()
}

type UserGraveyardProduce struct {
	*UserGraveyardProduceBuffs `json:"produce_buff"`
	PopulationCount            int32                `json:"population_count"`         // 人口数
	ProduceStartAt             int64                `json:"produce_start_at"`         // 何时开始生产
	AccRecords                 *GraveyardAccRecords `json:"acc_records"`              // 产出加速记录
	Characters                 *CharacterPositions  `json:"characters"`               // 派遣的侍从（派遣侍从星级会降低手动产出建筑单位时间）
	Productions                *Rewards             `json:"productions"`              // 已产出的物品
	TotalProduceNum            int32                `json:"total_produce_num"`        // 累计生产数量
	CurrProduceNum             int32                `json:"curr_produce_num"`         // 如果建筑是持续性产出的话表示，此次生产已生产x次，否则表示此次产出目标生产x次
	CurrReceiveProduceNum      int32                `json:"curr_receive_produce_num"` // 如果建筑是非持续性产出的话表示，当前生产已领取次数
	ReachLimit                 bool                 `json:"reach_limit"`              // 持续型产出，如果达到storeLimit后，置ReachLimit为true
}

type GraveyardAccRecords struct {
	Records []*GraveyardAccRecord `json:"records"`
}

func NewGraveyardAccRecords() *GraveyardAccRecords {
	return &GraveyardAccRecords{}
}

func (g *GraveyardAccRecords) GetProduceAccSec() int32 {
	var total int32
	for _, record := range g.Records {
		total += record.AccSec
	}
	return total
}

func (g *GraveyardAccRecords) AddRecord(sec int32, accAt int64) {
	g.Records = append(g.Records, &GraveyardAccRecord{
		AccAt:  accAt,
		AccSec: sec,
	})
	AccAtLess := func(i, j int) bool {
		return g.Records[i].AccAt < g.Records[j].AccAt
	}
	sort.Slice(g.Records, AccAtLess)
}

func (g *GraveyardAccRecords) clear() {
	g.Records = nil
}

type GraveyardAccRecord struct {
	AccSec int32 `json:"acc_sec"` //加速秒数
	AccAt  int64 `json:"acc_at"`  //加速时间戳
}

func NewUserGraveyardProduce() *UserGraveyardProduce {
	return &UserGraveyardProduce{
		UserGraveyardProduceBuffs: NewUserGraveyardProduceBuffs(),
		PopulationCount:           0,
		ProduceStartAt:            servertime.Now().Unix(),
		AccRecords:                NewGraveyardAccRecords(),
		Characters:                NewCharacterPositions(nil),
		Productions:               NewRewards(),
		TotalProduceNum:           0,
		CurrProduceNum:            0,
		CurrReceiveProduceNum:     0,
		ReachLimit:                false,
	}

}

// CalProduceSecs 计算已经生产多少时间(秒)
func (u *UserGraveyardProduce) CalProduceSecs() int32 {
	return getSec(u.ProduceStartAt) + u.AccRecords.GetProduceAccSec()
}

func (u *UserGraveyardProduce) VOGraveyardProduce(build *UserGraveyardBuild) *pb.VOGraveyardProduce {
	if u == nil {
		return nil
	}
	Characters := make([]*pb.VOCharacterPos, 0, len(u.Characters.Positions))
	for k, v := range u.Characters.Positions {
		Characters = append(Characters, &pb.VOCharacterPos{
			CharacterId: k,
			Position:    v,
		})
	}
	buff, ok := u.DuringBuff()
	var voBuff *pb.VOGraveyardBuff
	if ok {
		voBuff = buff.VOBuff(build)
	}
	return &pb.VOGraveyardProduce{
		PopulationCount: u.PopulationCount,
		ProduceSecs:     u.CalProduceSecs(),
		Characters:      Characters,
		Productions:     u.Productions.MergeVOResource(),
		ProduceNum:      u.CurrProduceNum,
		TotalProduceNum: u.TotalProduceNum,
		Buff:            voBuff,
		ReceiveNum:      u.CurrReceiveProduceNum,
	}
}

type CharacterPositions struct {
	Positions map[int32]int32
}

func NewCharacterPositions(positions map[int32]int32) *CharacterPositions {
	return &CharacterPositions{
		Positions: positions,
	}
}
func (c *CharacterPositions) GetCharacters() []int32 {
	var characters []int32

	for k := range c.Positions {
		characters = append(characters, k)
	}
	return characters
}

func (c *CharacterPositions) RemoveIfExist(other *CharacterPositions) (*CharacterPositions, bool) {
	newPositions := map[int32]int32{}
	for k, v := range c.Positions {
		newPositions[k] = v
	}
	remove := false
	for k := range other.Positions {
		if _, ok := newPositions[k]; ok {
			delete(newPositions, k)
			remove = true
		}
	}
	if remove {
		return NewCharacterPositions(newPositions), remove
	}
	return nil, remove
}

func VOGraveyardBuild(uid int64, build *UserGraveyardBuild) *pb.VOGraveyardBuild {

	return &pb.VOGraveyardBuild{
		BuildUid:     uid,
		BuildBase:    build.VOGraveyardBase(),
		Transition:   build.VOGraveyardTransition(),
		Produce:      build.VOGraveyardProduce(build),
		RequestState: build.VOGraveyardRequestState(),
	}
}

type UserGraveyardProduceBuffs struct {
	CountBuff *GraveyardProduceCountBuff  `json:"count_buff"`
	TimeBuffs []*GraveyardProduceTimeBuff `json:"time_buffs"`
}

func NewUserGraveyardProduceBuffs() *UserGraveyardProduceBuffs {
	return &UserGraveyardProduceBuffs{}
}
func (u UserGraveyardProduceBuffs) Clone() *UserGraveyardProduceBuffs {
	buffs := NewUserGraveyardProduceBuffs()
	if u.CountBuff != nil {
		buffs.CountBuff = NewGraveyardProduceCountBuff(u.CountBuff.BuffId, u.CountBuff.StartAt, u.CountBuff.RestCount)
	}
	for _, buff := range u.TimeBuffs {
		buffs.TimeBuffs = append(buffs.TimeBuffs, NewGraveyardProduceTimeBuff(buff.BuffId, buff.StartAt, buff.EndAt))
	}
	StartAtLess := func(i, j int) bool {
		return buffs.TimeBuffs[i].StartAt < buffs.TimeBuffs[j].StartAt
	}
	sort.Slice(buffs.TimeBuffs, StartAtLess)
	return buffs
}

type GraveyardProduceCountBuff struct {
	BuffId    int32 `json:"buff_id"`
	StartAt   int64 `json:"start_at"`   // buff开始时间
	RestCount int32 `json:"rest_count"` // 如果是次数类buff有 剩余使用次数
}

func NewGraveyardProduceCountBuff(Id int32, startAt int64, RestCount int32) *GraveyardProduceCountBuff {
	return &GraveyardProduceCountBuff{
		BuffId:    Id,
		StartAt:   startAt,
		RestCount: RestCount,
	}
}

type GraveyardProduceTimeBuff struct {
	BuffId  int32 `json:"buff_id"`
	StartAt int64 `json:"start_at"` // buff开始时间
	EndAt   int64 `json:"end_at"`   // buff截止时间
}

func NewGraveyardProduceTimeBuff(Id int32, startAt, endAt int64) *GraveyardProduceTimeBuff {
	return &GraveyardProduceTimeBuff{
		BuffId:  Id,
		StartAt: startAt,
		EndAt:   endAt,
	}
}

// InvalidateCountBuff 次数类buff失效
func (u *UserGraveyardBuild) InvalidateCountBuff() {
	if u.CountBuff != nil {
		if u.CountBuff.RestCount <= 0 {
			u.CountBuff = nil
		}

	}
}

// InvalidateTimeBuff 过期的时间类buff失效
func (u *UserGraveyardBuild) InvalidateTimeBuff() {
	now := servertime.Now().Unix()
	if len(u.TimeBuffs) != 0 {
		var newTimeBuffs []*GraveyardProduceTimeBuff
		for _, buff := range u.TimeBuffs {
			if buff.EndAt > now {
				// 未过期的buff保留
				newTimeBuffs = append(newTimeBuffs, buff)
			}
		}
		u.TimeBuffs = newTimeBuffs

	}
}

func (u *UserGraveyardBuild) InvalidateBuff() {
	u.InvalidateCountBuff()
	u.InvalidateTimeBuff()

}
func (u *UserGraveyardBuild) AddBuff(Id, Type, TypeContent int32, startAt int64) {
	duringBuff, ok := u.DuringBuff()
	if ok {
		duringBuff.buffExtend(TypeContent)
	} else {
		if Type == GraveyardCountBuffType {
			u.CountBuff = NewGraveyardProduceCountBuff(Id, startAt, TypeContent)
		} else if Type == GraveyardTimeBuffType {
			u.TimeBuffs = append(u.TimeBuffs, NewGraveyardProduceTimeBuff(Id, startAt, startAt+int64(TypeContent)))
			StartAtLess := func(i, j int) bool {
				return u.TimeBuffs[i].StartAt < u.TimeBuffs[j].StartAt
			}
			sort.Slice(u.TimeBuffs, StartAtLess)
		}
	}

}

// DuringBuff 返回现在正在使用的buff
func (u *UserGraveyardProduce) DuringBuff() (GraveyardBuff, bool) {
	return u.DuringBuffWhen(servertime.Now().Unix())
}

// DuringBuffWhen 返回某一时间点使用的buff
func (u *UserGraveyardProduce) DuringBuffWhen(tickTime int64) (GraveyardBuff, bool) {
	if u == nil {
		return nil, false
	}

	if u.CountBuff != nil {
		return u.CountBuff, true

	}
	if len(u.TimeBuffs) != 0 {
		for _, buff := range u.TimeBuffs {
			// 如果buff持续时间只有1秒，则只生效在第一秒
			if buff.EndAt > tickTime && buff.StartAt <= tickTime {
				return buff, true
			}
		}

	}
	return nil, false
}

func (u *UserGraveyardBuild) BuildComplete() bool {
	return u.FetchLevel() > 0

}

func (u *UserGraveyardBuild) CloneBuffFrom(from *UserGraveyardBuild) {
	u.UserGraveyardProduceBuffs = from.UserGraveyardProduceBuffs.Clone()
}

// BuildAcc 建造，升级，升阶加速
func (u *UserGraveyardBuild) BuildAcc(sec int32) {
	if u.UserGraveyardTransition == nil {
		return
	}
	u.EndAt -= int64(sec)
}

func (u *UserGraveyardBuild) GetRequestId() int64 {
	return u.RequestId
}

func (g *GraveyardProduceTimeBuff) GetBuffId() int32 {
	return g.BuffId
}
func (g *GraveyardProduceTimeBuff) IsEffective(build *UserGraveyardBuild) bool {
	return true
}
func (g *GraveyardProduceTimeBuff) buffExtend(TypeContent int32) {
	g.EndAt += int64(TypeContent)
}
func (g *GraveyardProduceTimeBuff) VOBuff(build *UserGraveyardBuild) *pb.VOGraveyardBuff {
	return &pb.VOGraveyardBuff{
		BuffId:      g.BuffId,
		IsEffective: g.IsEffective(build),
		EndAt:       g.EndAt,
	}
}

func (g *GraveyardProduceCountBuff) GetBuffId() int32 {
	return g.BuffId
}
func (g *GraveyardProduceCountBuff) IsEffective(build *UserGraveyardBuild) bool {
	// 次数类buff如果在生产中使用，则下次才生效
	return build.ProduceStartAt > g.StartAt
}
func (g *GraveyardProduceCountBuff) buffExtend(TypeContent int32) {
	g.RestCount += TypeContent
}
func (g *GraveyardProduceCountBuff) VOBuff(build *UserGraveyardBuild) *pb.VOGraveyardBuff {
	return &pb.VOGraveyardBuff{
		BuffId:      g.BuffId,
		IsEffective: g.IsEffective(build),
		RestCount:   g.RestCount,
	}
}

type UserGraveyardRequest struct {
	RequestId            int64 `json:"request_id"`              // 非0时候表示该建筑发请求帮助的请求id，生产建造升级升阶结束的时候清为0并去公会服删掉此request_id的请求
	RequestState         int32 `json:"request_state"`           // 1.未发送帮助请求，2.正在请求中，3.请求已经结束
	ReceiveHelpCount     int32 `json:"receive_help_count"`      // 如果是 RequestState=2则表示本次生产获得帮助次数
	LastReceiveHelpCount int32 `json:"last_receive_help_count"` // 上次cd过程中接受的帮助次数
}

func NewUserGraveyardRequest() *UserGraveyardRequest {
	return &UserGraveyardRequest{
		RequestState: RequestStateBefore,
		RequestId:    0,
	}
}
func (g *UserGraveyardRequest) VOGraveyardRequestState() *pb.VOGraveyardRequestState {
	if g == nil {
		return nil
	}
	return &pb.VOGraveyardRequestState{
		RequestState: g.RequestState,
		HelpCount:    g.LastReceiveHelpCount,
	}
}
func (g *UserGraveyardRequest) ResetRequest() {
	g.RequestId = 0
	g.RequestState = RequestStateBefore
	g.LastReceiveHelpCount = g.ReceiveHelpCount
	g.ReceiveHelpCount = 0
}

func (g *UserGraveyardRequest) SetRequestId(requestId int64) {
	g.RequestId = requestId
	g.RequestState = RequestStateDuring
	g.ReceiveHelpCount = 0
}

func (g *UserGraveyardRequest) IncreaseReceiveHelpCount() {
	g.ReceiveHelpCount++
}

type GraveyardBuff interface {
	GetBuffId() int32
	// IsEffective 次数类buff如果在生产中使用，则下次才生效
	IsEffective(build *UserGraveyardBuild) bool
	// 延长buff的使用时间或者次数
	buffExtend(TypeContent int32)
	VOBuff(build *UserGraveyardBuild) *pb.VOGraveyardBuff
}

type GraveyardHelpRequest struct {
	UserId        int64              `json:"user_id"` // 请求来源自
	BuildUid      int64              `json:"build_uid"`
	HelpType      int                `json:"help_type"`
	Sec           int32              `json:"sec"`
	HelpedCount   int32              `json:"helped_count"` // 已被帮助次数
	BuildId       int32              `json:"build_id"`
	BuildLv       int32              `json:"build_lv"`
	ExpireAt      int64              `json:"expire_at"`       // 生产升级升阶结束后该请求就消失
	HelpedUserIds map[int64]struct{} `json:"helped_user_ids"` // 已经帮助过此请求的玩家
}

func NewGraveyardHelpRequest(userId int64, buildUid int64, HelpType int, sec int32, helpedCount int32, buildId int32, buildLv int32, ExpireAt int64) *GraveyardHelpRequest {
	return &GraveyardHelpRequest{
		UserId:        userId,
		HelpType:      HelpType,
		BuildUid:      buildUid,
		Sec:           sec,
		HelpedCount:   helpedCount,
		BuildId:       buildId,
		BuildLv:       buildLv,
		ExpireAt:      ExpireAt,
		HelpedUserIds: map[int64]struct{}{},
	}
}
func (g *GraveyardHelpRequest) Expired() bool {
	return g.ExpireAt < servertime.Now().Unix()

}

func (g *GraveyardHelpRequest) VOGraveyardHelpRequest(m map[int64]*pb.VOGuildMember) *pb.VOGraveyardHelpRequest {
	return &pb.VOGraveyardHelpRequest{
		Member:      m[g.UserId],
		BuildId:     g.BuildId,
		BuildLv:     g.BuildLv,
		HelpedCount: g.HelpedCount,
		HelpType:    int32(g.HelpType),
	}

}

type GraveyardHelpRequests struct {
	Requests map[int64]*GraveyardHelpRequest `json:"requests"`
	UID      *uid.UID                        `json:"uid"`
}

func NewGraveyardRequests() *GraveyardHelpRequests {
	return &GraveyardHelpRequests{
		Requests: map[int64]*GraveyardHelpRequest{},
		UID:      uid.NewUID(),
	}
}

// GetRequestsExcept 获得除了userId其他人的请求
func (g *GraveyardHelpRequests) GetRequestsExcept(userId int64, m map[int64]struct{}, maxHelpCount int32) []*GraveyardHelpRequest {

	var result []*GraveyardHelpRequest
	for k, request := range g.Requests {
		if request.Expired() {
			delete(g.Requests, k)
			continue
		}
		// 已经不是公会成员
		if _, ok := m[request.UserId]; !ok {
			delete(g.Requests, k)
			continue
		}

		if request.HelpedCount >= maxHelpCount {
			delete(g.Requests, k)
			continue
		}
		if request.UserId == userId {
			continue
		}
		if _, ok := request.HelpedUserIds[userId]; ok {
			continue
		}
		result = append(result, request)

	}

	return result
}

func (g *GraveyardHelpRequests) Add(request *GraveyardHelpRequest) int64 {

	requestId := g.UID.Gen()
	g.Requests[requestId] = request
	return requestId
}
