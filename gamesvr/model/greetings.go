package model

import "shared/protobuf/pb"

const (
	GreetingTypeCharacter = 1
	GreetingTypeWorldItem = 2
)

type Greetings struct {
	Sequence  []*Greeting `json:"sequence"`
	SendQueue []*Greeting `json:"send_queue"`
}

type GreetingKey struct {
	Id           int32 `json:"id"`
	GreetingType int32 `json:"greeting_type"`
}

type Greeting struct {
	*GreetingKey
	Star      int32 `json:"star"`
	Count     int32 `json:"count"`
	Timestamp int64 `json:"timestamp"`
}

func (gr *Greeting) VOGreetings() *pb.VOGreetings {
	return &pb.VOGreetings{
		GreetingType: gr.GreetingType,
		Id:           gr.Id,
		Count:        gr.Count,
		Star:         gr.Star,
		Timestamp:    gr.Timestamp,
	}
}

func NewGreetings() *Greetings {
	return &Greetings{
		Sequence:  []*Greeting{},
		SendQueue: []*Greeting{},
	}
}

func NewGreeting(star, count int32, timestamp int64, key GreetingKey) *Greeting {
	return &Greeting{
		GreetingKey: &key,
		Star:        star,
		Count:       count,
		Timestamp:   timestamp,
	}
}

func NewGreetingKey(id, gType int32) *GreetingKey {
	return &GreetingKey{
		Id:           id,
		GreetingType: gType,
	}
}

func (g *Greetings) VOGreetings() []*pb.VOGreetings {

	voGreetings := make([]*pb.VOGreetings, 0, len(g.Sequence))

	for _, greeting := range g.Sequence {
		voGreetings = append(voGreetings, greeting.VOGreetings())
	}

	return voGreetings
}

func (g *Greetings) AddSendGreetings(greeting *Greeting) {
	g.SendQueue = append(g.SendQueue, greeting)
}
