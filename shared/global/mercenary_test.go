package global

import (
	"context"
	"shared/common"
	"shared/utility/number"
	"shared/utility/servertime"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func TestMercenarySetAndGet(t *testing.T) {
	MercenaryGlobal := NewMercenary(redis.NewClient(&redis.Options{
		Username: "root",
		Password: "",
		Addr:     ":6379",
	}))

	ctx := context.Background()

	characters := make([]*common.MercenaryCharacter, 0, 2)
	first := common.NewMercenaryCharacter(1005, 11, 2, 700)
	second := common.NewMercenaryCharacter(1006, 10, 3, 510)
	characters = append(characters, first)
	characters = append(characters, second)
	err := MercenaryGlobal.SetMercenaryCharacterData(ctx, 1737373, characters)
	if err != nil {
		t.Fatal(err)
	}

	cs, err := MercenaryGlobal.GetMercenaryCharacter(ctx, 1737373)
	if err != nil {
		t.Fatal(err)
	}

	for id, c := range cs {
		t.Log(id, *c)
	}

	// 修改之后
	charactersAfter := make([]*common.MercenaryCharacter, 0, 3)
	third := common.NewMercenaryCharacter(1005, 15, 1, 160)
	fourth := common.NewMercenaryCharacter(1006, 12, 4, 2160)
	fifth := common.NewMercenaryCharacter(1007, 10, 4, 1440)
	charactersAfter = append(charactersAfter, third)
	charactersAfter = append(charactersAfter, fourth)
	charactersAfter = append(charactersAfter, fifth)
	err = MercenaryGlobal.SetMercenaryCharacterData(ctx, 1737373, charactersAfter)
	if err != nil {
		t.Fatal(err)
	}

	cas, err := MercenaryGlobal.GetMercenaryCharacter(ctx, 1737373)
	if err != nil {
		t.Fatal(err)
	}

	for id, c := range cas {
		t.Log(id, *c)
	}
}

func TestMercenaryOwn(t *testing.T) {
	MercenaryGlobal := NewMercenaryOwn(redis.NewClient(&redis.Options{
		Username: "root",
		Password: "",
		Addr:     ":6379",
	}))

	ctx := context.Background()

	isSuccess, err := MercenaryGlobal.PutMercenaryOwn(ctx, 1, 1009, servertime.Now().Add(1*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(isSuccess)

	isSuccess, err = MercenaryGlobal.PutMercenaryOwn(ctx, 1, 1009, servertime.Now().Add(1*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(isSuccess)

	isMercenaryOwn, err := MercenaryGlobal.IsMercenaryOwn(ctx, 1, 1009)
	if err != nil {
		t.Fatal(err)
	}
	len, err := MercenaryGlobal.GetLengthMercenaryOwn(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(isMercenaryOwn, len)

	time.Sleep(1 * time.Second)

	isMercenaryOwn, err = MercenaryGlobal.IsMercenaryOwn(ctx, 1, 1009)
	if err != nil {
		t.Fatal(err)
	}
	len, err = MercenaryGlobal.GetLengthMercenaryOwn(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(isMercenaryOwn, len)
}

func TestMercenaryAvailable(t *testing.T) {
	MercenaryGlobal := NewMercenaryAvaliable(redis.NewClient(&redis.Options{
		Username: "root",
		Password: "",
		Addr:     ":6379",
	}))

	ctx := context.Background()
	isAvaliable, err := MercenaryGlobal.CheckMercenaryAvaliable(ctx, 1, 1010)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(isAvaliable)

	err = MercenaryGlobal.ChangeMercenaryAvaliable(ctx, 1, 1010, servertime.Now().Add(1*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	isAvaliable, err = MercenaryGlobal.CheckMercenaryAvaliable(ctx, 1, 1010)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(isAvaliable)

	time.Sleep(1 * time.Second)
	isAvaliable, err = MercenaryGlobal.CheckMercenaryAvaliable(ctx, 1, 1010)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(isAvaliable)
}

func TestMercenaryApplyCount(t *testing.T) {
	MercenaryGlobal := NewMercenaryApplyCount(redis.NewClient(&redis.Options{
		Username: "root",
		Password: "",
		Addr:     ":6379",
	}))

	ctx := context.Background()

	count, err := MercenaryGlobal.MercenaryApplyGet(ctx, 1, 1010)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(count)

	err = MercenaryGlobal.MercenaryApplyAdd(ctx, 1, 1010, 1, servertime.Now().Add(time.Second*1))
	if err != nil {
		t.Fatal(err)
	}

	count, err = MercenaryGlobal.MercenaryApplyGet(ctx, 1, 1010)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(count)

	time.Sleep(1 * time.Second)
	count, err = MercenaryGlobal.MercenaryApplyGet(ctx, 1, 1010)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(count)
}

func TestCharacterListPush(t *testing.T) {
	MercenaryGlobal := NewCharacterList(redis.NewClient(&redis.Options{
		Username: "root",
		Password: "",
		Addr:     ":6379",
	}))

	ctx := context.Background()

	c := &common.CharacterData{
		ID:         1010,
		Level:      1,
		Star:       2,
		Stage:      3,
		Skills:     map[int32]int32{},
		Power:      1000,
		Rarity:     3,
		Equipments: []*common.Equipment{},
		WorldItem: &common.WorldItem{
			ID:     1001,
			WID:    0,
			CID:    0,
			EXP:    number.NewCalNumber(0),
			Level:  number.NewCalNumber(1),
			Stage:  number.NewCalNumber(0),
			Rarity: 0,
			IsLock: false,
			CTime:  0,
		},
	}

	err := MercenaryGlobal.CharacterListPush(ctx, 1, "me", 2, 1, c, servertime.Now().Add(1*time.Second))
	if err != nil {
		t.Fatal(err)
	}

	data := MercenaryGlobal.CharacterListGetAndClear(ctx, 2)
	for _, d := range data {
		t.Log(d.Uid, d.Name, *d.Data)
	}

	err = MercenaryGlobal.CharacterListPush(ctx, 1, "me", 2, 1, c, servertime.Now().Add(1*time.Second))
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	data = MercenaryGlobal.CharacterListGetAndClear(ctx, 2)
	t.Log(len(data))
}

func TestMercenaryApplySend(t *testing.T) {
	MercenaryGlobal := NewMercenaryApply(redis.NewClient(&redis.Options{
		Username: "root",
		Password: "",
		Addr:     ":6379",
	}))

	ctx := context.Background()

	FirstSend := common.NewMercenarySend(1, "from", 1010)
	SecondSend := common.NewMercenarySend(1, "from", 1011)
	thirdSend := common.NewMercenarySend(1, "from", 1012)
	err := MercenaryGlobal.MercenaryApplySend(ctx, 2, FirstSend, servertime.Now().Add(2*time.Second))
	if err != nil {
		t.Fatal(err)
	}

	err = MercenaryGlobal.MercenaryApplySend(ctx, 2, SecondSend, servertime.Now().Add(2*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	all, err := MercenaryGlobal.MercenaryApplyReceive(ctx, 2)
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range all {
		t.Log(*a)
	}

	err = MercenaryGlobal.MercenaryApplySend(ctx, 2, thirdSend, servertime.Now().Add(2*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Second)

	all, err = MercenaryGlobal.MercenaryApplyReceive(ctx, 2)
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range all {
		t.Log(*a)
	}
}
