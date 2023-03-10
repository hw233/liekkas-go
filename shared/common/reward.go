package common

import (
	"math"
	"sync"

	"shared/protobuf/pb"
	"shared/utility/rand"
)

var rewardsType = map[int32]int32{}
var lock sync.RWMutex

func SetRewardsType(rst map[int32]int32) {
	lock.Lock()
	defer lock.Unlock()

	rewardsType = rst
}

func rewardType(id int32) int32 {
	lock.RLock()
	defer lock.RUnlock()

	return rewardsType[id]
}

type Reward struct {
	ID   int32
	Type int32
	Num  int32
}

func NewReward(id, num int32) *Reward {
	return &Reward{
		ID:   id,
		Type: rewardType(id),
		Num:  num,
	}
}

func (r *Reward) VOResource() *pb.VOResource {
	return &pb.VOResource{
		ItemId: r.ID,
		Count:  r.Num,
	}
}

type Rewards [][]Reward

func NewRewards() *Rewards {
	return (*Rewards)(&[][]Reward{})
}

func ParseFromVOConsume(list []*pb.VOResource) (*Rewards, error) {
	rewards := NewRewards()

	for _, consume := range list {
		if consume == nil {
			return nil, ErrParamError
		}
		if consume.Count <= 0 {
			return nil, ErrParamError
		}
		rewards.AddReward(NewReward(consume.ItemId, consume.Count))
	}
	return rewards, nil

}

func (rs *Rewards) AddReward(reward *Reward) {
	if len(*rs) == 0 {
		*rs = append(*rs, []Reward{})
	}

	lastIndex := len(*rs) - 1
	(*rs)[lastIndex] = append((*rs)[lastIndex], *reward)
}

//func (rs *Rewards) AddRewards(rewards ...*Reward) {
//	for _, reward := range rewards {
//		lastIndex := len(*rs) - 1
//
//		(*rs)[lastIndex] = append((*rs)[lastIndex], *reward)
//	}
//}

func (rs *Rewards) AddRewards(rewards *Rewards) {
	for _, vs := range *rewards {
		*rs = append(*rs, vs)
	}
}

func (rs *Rewards) Cut() {
	*rs = append(*rs, []Reward{})
}

func (rs *Rewards) Append(rewards *Rewards) {
	*rs = append(*rs, *rewards...)
}

func (rs *Rewards) Clone() *Rewards {
	cloneRewards := NewRewards()

	for _, rewards := range *rs {
		*cloneRewards = append(*cloneRewards, []Reward{})
		idx := len(*cloneRewards) - 1
		(*cloneRewards)[idx] = append((*cloneRewards)[idx], rewards...)
	}

	return cloneRewards
}

func (rs *Rewards) ContainsID(id int32) bool {
	for _, vs := range *rs {
		for _, v := range vs {
			if v.ID == id {
				return true
			}
		}
	}

	return false
}

func (rs *Rewards) ContainsType(typ int32) bool {
	for _, vs := range *rs {
		for _, v := range vs {
			if v.Type == typ {
				return true
			}
		}
	}

	return false
}

// Shuffle 打乱顺序
func (rs *Rewards) Shuffle() *Rewards {
	shuffleRewards := rs.Value()
	for i := len(shuffleRewards) - 1; i > 0; i-- {
		num := rand.RangeInt32(0, int32(i))
		shuffleRewards[i], shuffleRewards[num] = shuffleRewards[num], shuffleRewards[i]
	}
	rewards := NewRewards()

	for _, reward := range shuffleRewards {
		rewards.AddReward(&reward)
	}
	return rewards
}

func (rs *Rewards) Value() []Reward {
	var rewards []Reward

	for _, vs := range *rs {
		rewards = append(rewards, vs...)
	}

	return rewards
}

func (rs *Rewards) MergeValue() []Reward {
	var rewards []Reward
	// map[type]map[id]value
	record := map[int32]*Reward{}

	for _, vs := range *rs {
		for _, v := range vs {
			r, hasId := record[v.ID]
			if hasId {
				r.Num += v.Num
			} else {
				reward := NewReward(v.ID, v.Num)
				record[v.ID] = reward
			}
		}
	}

	for _, reward := range record {
		rewards = append(rewards, *reward)
	}
	return rewards
}

func (rs *Rewards) MultiValue() [][]Reward {
	return *rs
}

func (rs *Rewards) Multiple(multi int32) *Rewards {
	result := NewRewards()
	for _, vs := range *rs {
		for _, v := range vs {
			result.AddReward(NewReward(v.ID, v.Num*multi))
		}
	}
	return result
}

// MultipleFloor 向下取整
func (rs *Rewards) MultipleFloor(multi float64) *Rewards {
	result := NewRewards()
	for _, vs := range *rs {
		for _, v := range vs {
			result.AddReward(NewReward(v.ID, int32(math.Floor(float64(v.Num)*multi))))
		}
	}
	return result
}

// MultipleCeil 向上取整
func (rs *Rewards) MultipleCeil(multi float64) *Rewards {
	result := NewRewards()
	for _, vs := range *rs {
		for _, v := range vs {
			result.AddReward(NewReward(v.ID, int32(math.Ceil(float64(v.Num)*multi))))
		}
	}
	return result
}
func (rs *Rewards) VOResourceMultiple() []*pb.VOResourceMultiple {
	multi := make([]*pb.VOResourceMultiple, 0, len(*rs))

	for _, vs := range *rs {
		resource := make([]*pb.VOResource, 0, len(vs))
		for _, v := range vs {
			resource = append(resource, v.VOResource())
		}
		multi = append(multi, &pb.VOResourceMultiple{
			Resources: resource,
		})
	}

	return multi
}

func (rs *Rewards) VOMergedResourceMultiple() []*pb.VOResourceMultiple {
	multi := make([]*pb.VOResourceMultiple, 0, len(*rs))

	mergedRewards := rs.MergeValue()

	resource := make([]*pb.VOResource, 0, len(mergedRewards))
	for _, reward := range mergedRewards {
		resource = append(resource, reward.VOResource())
	}
	multi = append(multi, &pb.VOResourceMultiple{
		Resources: resource,
	})

	return multi
}

func (rs *Rewards) MergeVOResource() []*pb.VOResource {
	value := rs.MergeValue()
	result := make([]*pb.VOResource, 0, len(value))
	for _, reward := range value {
		result = append(result, reward.VOResource())
	}
	return result
}

func (rs *Rewards) IsEmpty() bool {
	return len(*rs) == 0
}

type RandReward struct {
	ID   int32
	Type int32
	Low  int32
	Up   int32
}

func NewRandReward(id, low, up int32) *RandReward {
	return &RandReward{
		ID:   id,
		Type: rewardType(id),
		Low:  low,
		Up:   up,
	}
}

func (rr *RandReward) NewReward() *Reward {
	return &Reward{
		ID:   rr.ID,
		Type: rr.Type,
		Num:  rand.RangeInt32(rr.Low, rr.Up),
	}
}
