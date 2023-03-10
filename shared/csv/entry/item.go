package entry

import (
	"shared/csv/base"
	"sort"
	"sync"

	"shared/common"
	"shared/csv/static"
	"shared/utility/errors"
	"shared/utility/transfer"
)

const CfgItemData = "cfg_item_data"

type Item struct {
	sync.RWMutex

	items                    map[int32]ItemBase
	characterExpItems        map[int32]CharacterExpItem
	usePramNilItems          map[int32]UsePramNilItem
	graveyardAccItems        map[int32]GraveyardAccItem
	graveyardGetProductItems map[int32]GraveyardGetProduct
	giftSelectOneItems       map[int32]GiftSelectOne
	giftRandomDropItems      map[int32]GiftRandomDrop
	graveyardBuffItems       map[int32]GraveyardBuffItem
	energyItems              map[int32]EnergyItem
	expItems                 map[int32]int32
	characterAdaptiveGifts   map[int32]CharacterAdaptiveGift
}

func NewItem() *Item {
	return &Item{
		items:                    map[int32]ItemBase{},
		characterExpItems:        map[int32]CharacterExpItem{},
		usePramNilItems:          map[int32]UsePramNilItem{},
		graveyardAccItems:        map[int32]GraveyardAccItem{},
		graveyardGetProductItems: map[int32]GraveyardGetProduct{},
		giftSelectOneItems:       map[int32]GiftSelectOne{},
		giftRandomDropItems:      map[int32]GiftRandomDrop{},
		graveyardBuffItems:       map[int32]GraveyardBuffItem{},
		energyItems:              map[int32]EnergyItem{},
		expItems:                 map[int32]int32{},
		characterAdaptiveGifts:   map[int32]CharacterAdaptiveGift{},
	}
}

func (i *Item) Check(configManager *Config) error {
	return nil
}

func (i *Item) Reload(config *Config) error {
	i.Lock()
	defer i.Unlock()
	items := map[int32]ItemBase{}

	characterExpItems := map[int32]CharacterExpItem{}
	usePramNilItems := map[int32]UsePramNilItem{}
	graveyardAccItems := map[int32]GraveyardAccItem{}
	graveyardGetProductItems := map[int32]GraveyardGetProduct{}
	giftSelectOneItems := map[int32]GiftSelectOne{}
	giftRandomDropItems := map[int32]GiftRandomDrop{}
	graveyardBuffItems := map[int32]GraveyardBuffItem{}
	energyItems := map[int32]EnergyItem{}
	expItems := map[int32]int32{}
	characterAdaptiveGifts := map[int32]CharacterAdaptiveGift{}

	for _, v := range config.CfgItemDataConfig.GetAllData() {

		itemBase := &ItemBase{}
		err := transfer.Transfer(v, itemBase)
		if err != nil {
			return errors.WrapTrace(err)
		}
		items[v.Id] = *itemBase

		switch v.ItemType {
		case static.ItemTypeCharExp:
			characterExpItems[v.Id] = CharacterExpItem{AddExp: v.UseParam[0]}
		case static.ItemTypeCharSkill, static.ItemTypeCharPiece, static.ItemTypeGacha:
			usePramNilItems[v.Id] = UsePramNilItem{}
		case static.ItemTypeGraveyardAccelerate:
			if len(v.UseParam) != 3 {
				return errors.Swrapf(common.ErrCSVFormatInvalid, CfgItemData, v.Id)
			}
			graveyardAccItems[v.Id] = NewGraveyardAccItem(v.Id, v.UseParam[0], v.UseParam[1], v.UseParam[2])
		case static.ItemTypeGraveyardGetProduct:
			if len(v.UseParam) != 4 {
				return errors.Swrapf(common.ErrCSVFormatInvalid, CfgItemData, v.Id)
			}
			graveyardGetProductItems[v.Id] = NewGraveyardGetProduct(v.UseParam[0], v.UseParam[1],
				v.UseParam[2], v.UseParam[3])
		case static.ItemTypeGiftSelectOne:
			if len(v.UseParam) != 1 {
				return errors.Swrapf(common.ErrCSVFormatInvalid, CfgItemData, v.Id)
			}
			dropId := v.UseParam[0]
			cfgDropDataConfig, ok := config.CfgDropDataConfig.Find(dropId)
			if !ok {
				return errors.Swrapf(common.ErrNotFoundInCSV, CfgDropDataConfig, dropId)

			}
			var dropGroups []*base.CfgDropGroup
			for _, group := range config.CfgDropGroupConfig.GetAllData() {
				if group.DropGroup == cfgDropDataConfig.DropGroup {
					dropGroups = append(dropGroups, group)
				}
			}
			if len(dropGroups) == 0 {
				return errors.Swrapf(common.ErrNotFoundInCSV, CfgDropDataConfig, dropId)
			}
			// 按dropGroup表的id升序
			less := func(i, j int) bool {
				return dropGroups[i].Id < dropGroups[j].Id
			}
			sort.Slice(dropGroups, less)
			rewards := make([]common.Reward, 0, len(dropGroups))
			for _, group := range dropGroups {
				rewards = append(rewards, *common.NewReward(group.DropItem, group.DropNumber[0]))
			}
			giftSelectOneItems[v.Id] = NewGiftSelectOne(rewards)

		case static.ItemTypeGiftRandomDrop:
			giftRandomDropItems[v.Id] = NewGiftRandomDrop(v.UseParam)
		case static.ItemTypeGraveyardProduceBuff:
			graveyardBuffItems[v.Id] = NewGraveyardBuffItem(v.UseParam)
		case static.ItemTypeEnergyItem:
			if len(v.UseParam) != 1 {
				return errors.Swrapf(common.ErrCSVFormatInvalid, CfgItemData, v.Id)
			}
			energyItems[v.Id] = NewEnergyItem(v.UseParam[0])
		case static.ItemTypeEquipExp, static.ItemTypeWorldItemExp:
			if len(v.UseParam) != 1 {
				return errors.Swrapf(common.ErrCSVFormatInvalid, CfgItemData, v.Id)
			}

			expItems[v.Id] = v.UseParam[0]
		case static.ItemTypeCharacterAdaptiveGift:
			if len(v.UseParam)%2 != 0 {
				return errors.Swrapf(common.ErrCSVFormatInvalid, CfgItemData, v.Id)
			}
			characterAdaptiveGifts[v.Id] = NewCharacterAdaptiveGift(v.UseParam...)
		}
	}
	i.items = items
	i.characterExpItems = characterExpItems
	i.usePramNilItems = usePramNilItems
	i.graveyardAccItems = graveyardAccItems
	i.graveyardGetProductItems = graveyardGetProductItems
	i.giftSelectOneItems = giftSelectOneItems
	i.giftRandomDropItems = giftRandomDropItems
	i.graveyardBuffItems = graveyardBuffItems
	i.energyItems = energyItems
	i.expItems = expItems
	i.characterAdaptiveGifts = characterAdaptiveGifts
	return nil
}

type CharacterExpItem struct {
	AddExp int32
}
type UsePramNilItem struct {
}
type GraveyardAccItem struct {
	Id  int32
	Acc GraveyardAcc
}

func NewGraveyardAccItem(Id, AccType, BuildId, Sec int32) GraveyardAccItem {
	return GraveyardAccItem{
		Id:  Id,
		Acc: NewGraveyardAcc(AccType, BuildId, Sec),
	}
}

type GraveyardAcc struct {
	AccType int32 // 加速
	BuildId int32 // 为0表示都可以用
	Sec     int32 // 加速时间
}

func NewGraveyardAcc(AccType, BuildId, Sec int32) GraveyardAcc {
	return GraveyardAcc{
		AccType: AccType,
		BuildId: BuildId,
		Sec:     Sec,
	}
}

type GraveyardGetProduct struct {
	BuildId int32
	Sec     int32
	Low     int32
	Up      int32
}

func NewGraveyardGetProduct(BuildId, Sec, Low, Up int32) GraveyardGetProduct {
	return GraveyardGetProduct{
		BuildId: BuildId,
		Sec:     Sec,
		Low:     Low,
		Up:      Up,
	}
}

type ItemBase struct {
	ItemType    int32
	LimitNumber int32
	SellPrice   int32
	Rarity      int32
	Cname       string
}

type GiftSelectOne struct {
	Rewards []common.Reward
}

func NewGiftSelectOne(Rewards []common.Reward) GiftSelectOne {
	return GiftSelectOne{
		Rewards: Rewards,
	}
}

type GiftRandomDrop struct {
	DropIds []int32
}

func NewGiftRandomDrop(DropIds []int32) GiftRandomDrop {
	return GiftRandomDrop{
		DropIds: DropIds,
	}
}

type GraveyardBuffItem struct {
	BuffIds []int32
}

func NewGraveyardBuffItem(BuffIds []int32) GraveyardBuffItem {
	return GraveyardBuffItem{
		BuffIds: BuffIds,
	}
}

type EnergyItem struct {
	AddNum int32
}

func NewEnergyItem(AddNum int32) EnergyItem {
	return EnergyItem{
		AddNum: AddNum,
	}
}

type CharacterAdaptiveGift struct {
	DropMap map[int32]int32
}

func NewCharacterAdaptiveGift(args ...int32) CharacterAdaptiveGift {
	DropMap := map[int32]int32{}
	for i := 0; i < len(args); i += 2 {
		DropMap[args[i]] = args[i+1]
	}
	return CharacterAdaptiveGift{
		DropMap: DropMap,
	}
}
func (i *Item) GetAllCharExpItem() map[int32]CharacterExpItem {
	i.RLock()
	defer i.RUnlock()
	return i.characterExpItems
}

func (i *Item) GetCharExpItem(itemId int32) (CharacterExpItem, bool) {
	i.RLock()
	defer i.RUnlock()
	item, ok := i.characterExpItems[itemId]
	return item, ok
}

func (i *Item) GetSortedCharExpItem(rewards *common.Rewards) ([]AddExp, error) {
	//判断是否是经验道具
	var items []AddExp
	for _, reward := range rewards.Value() {
		itemBase, ok := i.GetItem(reward.ID)
		if !ok {
			return nil, errors.WrapTrace(common.ErrParamError)
		}
		item, ok := i.GetCharExpItem(reward.ID)
		if !ok {
			return nil, errors.WrapTrace(common.ErrParamError)
		}
		items = append(items, AddExp{CharacterExpItem: item, ItemId: reward.ID, Count: reward.Num, Rarity: itemBase.Rarity})
	}
	// 按稀有度排序

	less := func(i, j int) bool {
		return items[i].Rarity < items[j].Rarity
	}
	sort.Slice(items, less)
	return items, nil
}

type AddExp struct {
	CharacterExpItem
	ItemId int32
	Count  int32
	Rarity int32
}

func (i *Item) CalRealConsume(items []AddExp, totalExp int32) (*common.Rewards, int32, error) {
	if totalExp == 0 {
		return nil, 0, common.ErrCharacterLevelMax
	}

	// 累计增加的经验
	var addExp int32
	realConsume := common.NewRewards()
	for _, item := range items {
		if totalExp <= 0 {
			break
		}
		var consumeNum int32
		for i := 0; i < int(item.Count); i++ {
			if totalExp <= 0 {
				break
			}
			totalExp -= item.AddExp
			addExp += item.AddExp
			consumeNum++
		}
		realConsume.AddReward(common.NewReward(item.ItemId, consumeNum))
	}
	return realConsume, addExp, nil
}

func (i *Item) GetAllUsePramNilItems() map[int32]UsePramNilItem {
	i.RLock()
	defer i.RUnlock()
	return i.usePramNilItems
}

func (i *Item) GetGraveyardAcc(itemId int32) (*GraveyardAccItem, bool) {
	i.RLock()
	defer i.RUnlock()
	item, ok := i.graveyardAccItems[itemId]
	if !ok {
		return nil, false
	}
	return &item, ok
}
func (i *Item) GetGraveyardGetProductItem(itemId int32) (*GraveyardGetProduct, bool) {
	i.RLock()
	defer i.RUnlock()
	item, ok := i.graveyardGetProductItems[itemId]
	if !ok {
		return nil, false
	}
	return &item, ok
}

func (i *Item) GetItem(itemId int32) (*ItemBase, bool) {
	i.RLock()
	defer i.RUnlock()
	item, ok := i.items[itemId]
	if !ok {
		return nil, false
	}
	return &item, ok
}

func (i *Item) GetGiftSelectOne(itemId int32) (*GiftSelectOne, bool) {
	i.RLock()
	defer i.RUnlock()
	item, ok := i.giftSelectOneItems[itemId]
	if !ok {
		return nil, false
	}
	return &item, ok
}

func (i *Item) GetGiftRandomDrop(itemId int32) (*GiftRandomDrop, bool) {
	i.RLock()
	defer i.RUnlock()
	item, ok := i.giftRandomDropItems[itemId]
	if !ok {
		return nil, false
	}
	return &item, ok
}

func (i *Item) GetGraveyardBuffItem(itemId int32) (*GraveyardBuffItem, bool) {
	i.RLock()
	defer i.RUnlock()
	item, ok := i.graveyardBuffItems[itemId]
	if !ok {
		return nil, false
	}
	return &item, ok
}

func (i *Item) GetEnergyItem(itemId int32) (*EnergyItem, bool) {
	i.RLock()
	defer i.RUnlock()
	item, ok := i.energyItems[itemId]
	if !ok {
		return nil, false
	}
	return &item, ok
}

func (i *Item) CalculateTotalEquipmentEXP(items []int32) int32 {
	i.RLock()
	defer i.RUnlock()

	if len(items) != len(common.EquipmentEXPItems) {
		return 0
	}

	var totalEXP int32 = 0

	for index, v := range common.EquipmentEXPItems {
		totalEXP += i.expItems[v] * items[index]
	}

	return totalEXP
}

func (i *Item) CalculateTotalWorldItemEXP(items []int32) int32 {
	i.RLock()
	defer i.RUnlock()

	if len(items) != len(common.WorldItemEXPItems) {
		return 0
	}

	var totalEXP int32 = 0

	for index, v := range common.WorldItemEXPItems {
		totalEXP += i.expItems[v] * items[index]
	}

	return totalEXP
}

func (i *Item) GetCharacterAdaptiveGift(itemId int32) (*CharacterAdaptiveGift, bool) {
	i.RLock()
	defer i.RUnlock()
	item, ok := i.characterAdaptiveGifts[itemId]
	if !ok {
		return nil, false
	}
	return &item, ok
}

func (i *Item) GetLimit(itemId int32) int32 {
	itemCfg, ok := i.GetItem(itemId)
	if !ok {
		return 0
	}

	return itemCfg.LimitNumber
}
