package model

import "shared/protobuf/pb"

type YggdrasilResourceResult struct {
	DiscardGoods       map[int64]*YggDiscardGoods
	PackGoods          map[int64]*YggPackGoods
	DeletePackGoodsIds []int64
	DeleteDiscardGoods []int64
	TaskItems          map[int32]int32
}

func NewYggdrasilResourceResult() *YggdrasilResourceResult {
	return &YggdrasilResourceResult{
		DiscardGoods:       map[int64]*YggDiscardGoods{},
		PackGoods:          map[int64]*YggPackGoods{},
		DeletePackGoodsIds: nil,
		DeleteDiscardGoods: nil,
		TaskItems:          map[int32]int32{},
	}
}

func (y *YggdrasilResourceResult) VOYggdrasilResourceResult() *pb.VOYggdrasilResourceResult {

	ret := &pb.VOYggdrasilResourceResult{}

	DiscardGoods := make([]*pb.VOYggdrasilDiscardGoods, 0, len(y.DiscardGoods))
	PackGoods := make([]*pb.VOYggdrasilPackGoods, 0, len(y.PackGoods))
	DeletePackGoodsIds := make([]int64, 0, len(y.DeletePackGoodsIds))
	DeleteDiscardGoods := make([]int64, 0, len(y.DeleteDiscardGoods))
	TaskPackChange := make([]*pb.VOResource, 0, len(y.TaskItems))

	for _, good := range y.DiscardGoods {
		DiscardGoods = append(DiscardGoods, good.VOYggdrasilDiscardGoods())
	}
	for _, good := range y.PackGoods {
		PackGoods = append(PackGoods, good.VOYggdrasilPackGoods())
	}
	for _, id := range y.DeleteDiscardGoods {
		DeleteDiscardGoods = append(DeleteDiscardGoods, id)
	}
	for _, id := range y.DeletePackGoodsIds {
		DeletePackGoodsIds = append(DeletePackGoodsIds, id)
	}
	for k, v := range y.TaskItems {
		TaskPackChange = append(TaskPackChange, &pb.VOResource{ItemId: k, Count: v})
	}
	ret.DiscardGoods = DiscardGoods
	ret.PackGoods = PackGoods
	ret.DeletePackGoodsIds = DeletePackGoodsIds
	ret.DeleteDiscardGoods = DeleteDiscardGoods
	ret.TaskPackChange = TaskPackChange
	y.Clear()

	return ret
}

func (y *YggdrasilResourceResult) AppendDeleteDiscardGoods(deleteDiscardGoods int64) {
	y.DeleteDiscardGoods = append(y.DeleteDiscardGoods, deleteDiscardGoods)

}

func (y *YggdrasilResourceResult) AppendDeletePackGoods(deletePackGoods int64) {
	y.DeletePackGoodsIds = append(y.DeletePackGoodsIds, deletePackGoods)

}

func (y *YggdrasilResourceResult) AppendDiscardGoods(discardGoods *YggDiscardGoods) {
	y.DiscardGoods[discardGoods.Uid] = discardGoods

}

func (y *YggdrasilResourceResult) AppendPackGoods(packGoods *YggPackGoods) {
	y.PackGoods[packGoods.Uid] = packGoods
}

func (y *YggdrasilResourceResult) TaskItemChange(itemId, num int32) {
	y.TaskItems[itemId] = num

}

func (y *YggdrasilResourceResult) TaskItemChanges(idNum map[int32]int32) {
	for itemId, num := range idNum {
		y.TaskItemChange(itemId, num)
	}

}
func (y *YggdrasilResourceResult) Clear() {
	y.DiscardGoods = map[int64]*YggDiscardGoods{}
	y.PackGoods = map[int64]*YggPackGoods{}
	y.DeletePackGoodsIds = nil
	y.DeleteDiscardGoods = nil
	y.TaskItems = map[int32]int32{}
}
