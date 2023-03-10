package common

import "shared/protobuf/pb"

type VisitingCardShow struct {
	CharactersSet   bool                          `json:"characters_set"`
	Characters      map[int32]int32               `json:"characters"`
	CharactersCache map[int32]*pb.VOUserCharacter `json:"characters_cache"` //他人访问主页时候的cache

	WorldItemUId int64 `json:"world_item_u_id"`
	WorldItemId  int32 `json:"world_item_id"`
}

func NewVisitingCardShow() *VisitingCardShow {
	return &VisitingCardShow{
		CharactersSet:   false,
		Characters:      map[int32]int32{},
		CharactersCache: map[int32]*pb.VOUserCharacter{},
		WorldItemUId:    0,
		WorldItemId:     0,
	}
}
