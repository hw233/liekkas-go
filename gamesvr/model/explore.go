package model

import (
	"shared/common"
	"shared/protobuf/pb"
	"shared/utility/servertime"
)

type ExploreMap struct {
	ChapterId   int32        `json:"chapter_id"`
	CurPosition *common.Vec2 `json:"cur_position"`
}

type ExploreEventPoint struct {
	Id            int32 `json:"id"`
	InteractTime  int64 `json:"interact_time"`
	InteractTimes int32 `json:"interact_times"`
}

type ExploreResourcePoint struct {
	Id               int32 `json:"id"`
	Collecting       bool  `json:"collecting"`
	CollectStartTime int64 `json:"collect_start_time"`
	CollectTimes     int32 `json:"collect_times"`
	HasMonster       bool  `json:"has_monster"`
}

type ExploreInfo struct {
	Maps           map[int32]*ExploreMap           `json:"maps"`
	ResourcePoints map[int32]*ExploreResourcePoint `json:"resource_points"`
	EventPoints    map[int32]*ExploreEventPoint    `json:"event_point"`
	FogUnlocks     map[int32]bool                  `json:"fog_unlocks"`
	TransportGates map[int32]*ExploreTransportGate `json:"transport_gates"`

	CurChapter int32 `json:"cur_chapter"`
}

type ExploreTransportGate struct {
	Id            int32 `json:"id"`
	UseTimes      int32 `json:"user_times"`
	DestroyShowed bool  `json:"destroy_showed"`
}

func NewExploreInfo() *ExploreInfo {
	return &ExploreInfo{
		Maps:           map[int32]*ExploreMap{},
		ResourcePoints: map[int32]*ExploreResourcePoint{},
		EventPoints:    map[int32]*ExploreEventPoint{},
		FogUnlocks:     map[int32]bool{},
		TransportGates: map[int32]*ExploreTransportGate{},

		CurChapter: 0,
	}
}

func NewExploreMap(chapterId int32, initPos *common.Vec2) *ExploreMap {
	return &ExploreMap{
		ChapterId:   chapterId,
		CurPosition: common.NewVec2ByVec2(initPos),
	}
}

func NewExploreEventPoint(id int32) *ExploreEventPoint {
	return &ExploreEventPoint{
		Id:            id,
		InteractTime:  0,
		InteractTimes: 0,
	}
}

func NewExploreResourcePoint(id int32) *ExploreResourcePoint {
	return &ExploreResourcePoint{
		Id:               id,
		Collecting:       false,
		CollectStartTime: 0,
		CollectTimes:     0,
		HasMonster:       true,
	}
}

func NewExploreTransportGate(id int32) *ExploreTransportGate {
	return &ExploreTransportGate{
		Id:            id,
		UseTimes:      0,
		DestroyShowed: false,
	}
}

func (ei *ExploreInfo) SetCurChapter(id int32) {
	ei.CurChapter = id
}

func (ei *ExploreInfo) GetCurChapter() int32 {
	return ei.CurChapter
}

func (ei *ExploreInfo) AddMap(chapterId int32, initPos *common.Vec2) *ExploreMap {
	exploreMap, ok := ei.Maps[chapterId]
	if ok {
		return exploreMap
	}

	exploreMap = NewExploreMap(chapterId, initPos)

	ei.Maps[chapterId] = exploreMap

	return exploreMap
}

func (ei *ExploreInfo) GetMap(chapterId int32) (*ExploreMap, bool) {
	exploreMap, ok := ei.Maps[chapterId]
	return exploreMap, ok
}

func (ei *ExploreInfo) AddEventPoint(id int32) *ExploreEventPoint {
	eventPoint := NewExploreEventPoint(id)
	ei.EventPoints[id] = eventPoint
	return eventPoint
}

func (ei *ExploreInfo) GetEventPoint(id int32) (*ExploreEventPoint, bool) {
	eventPoint, ok := ei.EventPoints[id]
	return eventPoint, ok
}

func (ei *ExploreInfo) AddResourcePoint(id int32) *ExploreResourcePoint {
	resourcePoint := NewExploreResourcePoint(id)
	ei.ResourcePoints[id] = resourcePoint

	return resourcePoint
}

func (ei *ExploreInfo) GetResourcePoint(id int32) (*ExploreResourcePoint, bool) {
	resourcePoint, ok := ei.ResourcePoints[id]
	return resourcePoint, ok
}

func (ei *ExploreInfo) UnlockFog(id int32) {
	ei.FogUnlocks[id] = true
}

func (ei *ExploreInfo) IsFogUnlocked(id int32) bool {
	unlock, ok := ei.FogUnlocks[id]
	if !ok {
		return ok
	}

	return unlock
}

func (ei *ExploreInfo) AddTransportGate(id int32) *ExploreTransportGate {
	transportGate := NewExploreTransportGate(id)
	ei.TransportGates[id] = transportGate

	return transportGate
}

func (ei *ExploreInfo) GetTransportGate(id int32) (*ExploreTransportGate, bool) {
	transportGate, ok := ei.TransportGates[id]
	return transportGate, ok
}

func (ei *ExploreInfo) VOExploreInfo() *pb.VOExploreInfo {
	info := &pb.VOExploreInfo{
		Maps:       make([]*pb.VOExploreMap, 0, len(ei.Maps)),
		Events:     make([]*pb.VOExploreEvent, 0, len(ei.EventPoints)),
		Resources:  make([]*pb.VOExploreResource, 0, len(ei.ResourcePoints)),
		UnlockFogs: make([]int32, 0, len(ei.FogUnlocks)),
		TpGate:     make([]*pb.VOExploreTransportGate, 0, len(ei.TransportGates)),
		CurChapter: ei.CurChapter,
	}

	for _, exploreMap := range ei.Maps {
		info.Maps = append(info.Maps, exploreMap.VOExploreMap())
	}

	for _, event := range ei.EventPoints {
		info.Events = append(info.Events, event.VOExploreEvent())
	}

	for _, resource := range ei.ResourcePoints {
		info.Resources = append(info.Resources, resource.VOExploreResource())
	}

	for fogId := range ei.FogUnlocks {
		info.UnlockFogs = append(info.UnlockFogs, fogId)
	}

	for _, gate := range ei.TransportGates {
		info.TpGate = append(info.TpGate, gate.VOExploreTransportGate())
	}

	return info
}

//----------------------------------------
// Map
//----------------------------------------
func (em *ExploreMap) SetCurPosition(position *common.Vec2) {
	em.CurPosition.Set(position)
}

func (em *ExploreMap) VOExploreMap() *pb.VOExploreMap {
	return &pb.VOExploreMap{
		ChapterId: em.ChapterId,
		CurPosition: &pb.VOPosition{
			PosX: int32(em.CurPosition.X),
			PosY: int32(em.CurPosition.Y),
		},
	}
}

//----------------------------------------
// ExploreEventPoint
//----------------------------------------
func (ei *ExploreEventPoint) RecordInteration() {
	ei.InteractTime = servertime.Now().Unix()
	ei.InteractTimes = ei.InteractTimes + 1
}

func (ei *ExploreEventPoint) GetInteractTime() int64 {
	return ei.InteractTime
}

func (ei *ExploreEventPoint) GetInteractTimes() int32 {
	return ei.InteractTimes
}

func (ei *ExploreEventPoint) VOExploreEvent() *pb.VOExploreEvent {
	exploreEvent := &pb.VOExploreEvent{
		Id:            ei.Id,
		InteractTime:  ei.InteractTime,
		InteractTimes: ei.InteractTimes,
	}

	return exploreEvent
}

//----------------------------------------
// ExploreResourcePoint
//----------------------------------------
func (er *ExploreResourcePoint) ClearMonster() {
	er.HasMonster = false
}

func (er *ExploreResourcePoint) IsMonsterExist() bool {
	return er.HasMonster
}

func (er *ExploreResourcePoint) StartCollect() {
	er.Collecting = true
	er.CollectStartTime = servertime.Now().Unix()
}

func (er *ExploreResourcePoint) FinishCollect() {
	er.Collecting = false
	er.CollectTimes = er.CollectTimes + 1
}

func (er *ExploreResourcePoint) IsCollecting() bool {
	return er.Collecting
}

func (er *ExploreResourcePoint) GetCollectTimes() int32 {
	return er.CollectTimes
}

func (er *ExploreResourcePoint) GetCollectStartTime() int64 {
	return er.CollectStartTime
}

func (er *ExploreResourcePoint) VOExploreResource() *pb.VOExploreResource {
	resource := &pb.VOExploreResource{
		Id:               er.Id,
		Collecting:       er.Collecting,
		CollectStartTime: er.CollectStartTime,
		CollectTimes:     er.CollectTimes,
		HasMonster:       er.HasMonster,
	}

	return resource
}

//----------------------------------------
// ExploreTransportGate
//----------------------------------------
func (et *ExploreTransportGate) GetUseTimes() int32 {
	return et.UseTimes
}

func (et *ExploreTransportGate) RecordUsed() {
	et.UseTimes = et.UseTimes + 1
}

func (et *ExploreTransportGate) SetDestroyShowed() {
	et.DestroyShowed = true
}

func (et *ExploreTransportGate) VOExploreTransportGate() *pb.VOExploreTransportGate {
	return &pb.VOExploreTransportGate{
		Id:        et.Id,
		UseTimes:  et.UseTimes,
		Destroyed: et.DestroyShowed,
	}
}
