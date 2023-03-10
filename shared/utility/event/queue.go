package event

import (
	"context"
	"encoding/json"
	"fmt"

	"shared/utility/errors"
	"shared/utility/global"
	"shared/utility/glog"
	"shared/utility/param"

	"github.com/go-redis/redis/v8"
)

const (
	eventUserIdKey = "e:u:%d"
)

type EventQueue struct {
	client *redis.Client
	global *global.Global
}

func NewEventQueue(client *redis.Client) *EventQueue {
	return &EventQueue{
		client: client,
		global: global.NewGlobal(client),
	}
}

// ExecuteEventsInQueue 执行所有队列中的event后删除list
func (e *EventQueue) ExecuteEventsInQueue(ctx context.Context, userId int64) {
	// 分布式锁
	lock, err := e.global.ObtainLock(ctx, userId)
	if err != nil {
		glog.Errorf("EventQueue ExecuteEventsInQueue ObtainLock err:%+v:", err)
		return
	}
	defer lock.Release()

	for _, event := range e.eventsInQueue(ctx, userId) {
		handler, ok := UserEventHandler.GetHandler(userId, event.Type)
		if !ok {
			glog.Errorf("EventQueue ExecuteEventsInQueue Unregistied handler ,event type %d", event.Type)
			continue
		}
		glog.Infof("EventQueue handle event, userid:%d,event:%+v,", userId, event)
		err := handler(param.NewParam(event.Params))
		if err != nil {
			glog.Errorf("EventQueue handle event error occur, userid:%d,event:%+v,err:%+v", userId, event, err)
			continue
		}
	}
	e.clear(ctx, userId)

}
func (e *EventQueue) makeKey(userId int64) string {
	return fmt.Sprintf(eventUserIdKey, userId)
}

func (e *EventQueue) eventsInQueue(ctx context.Context, userId int64) []*Event {

	results, err := e.client.LRange(ctx, e.makeKey(userId), 0, -1).Result() //读list所有元素
	var events []*Event
	if err != nil {
		glog.Errorf("EventQueue eventsInQueue LRange err:%+v", err)
		return events
	}
	for _, r := range results {
		var event = &Event{}
		err := json.Unmarshal([]byte(r), event)
		if err != nil {
			glog.Errorf("EventQueue eventsInQueue json Unmarshal err:%+v", err)
			continue
		}
		events = append(events, event)
	}
	return events
}

// Push 给某user添加一个事件
func (e *EventQueue) Push(ctx context.Context, userId int64, event *Event) {

	// 分布式锁
	lock, err := e.global.ObtainLock(ctx, userId)
	if err != nil {
		glog.Errorf("EventQueue Push ObtainLock err:%+v:", err)
		return
	}
	defer lock.Release()
	glog.Debugf("EventQueue Push userId:%d,event:%+v:", userId, event)

	bytes, err := json.Marshal(event)
	if err != nil {
		glog.Errorf("EventQueue json Marshal err:%v", errors.WrapTrace(err))
		return
	}
	_, err = e.client.RPush(ctx, e.makeKey(userId), bytes).Result()
	if err != nil {
		glog.Errorf("EventQueue rpush err:%v", errors.WrapTrace(err))
		return
	}
}

func (e *EventQueue) clear(ctx context.Context, userId int64) {
	_, err := e.client.Del(ctx, e.makeKey(userId)).Result()
	if err != nil {
		glog.Errorf("EventQueue Del err:%v", errors.WrapTrace(err))
		return
	}
}
