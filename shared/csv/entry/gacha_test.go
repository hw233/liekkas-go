package entry

import (
	_ "net/http/pprof"
	"shared/common"
	"shared/utility/rand"
	"testing"
)

func TestGachaEntry_Gacha_1(t *testing.T) {
	pool, _ := CSV.Gacha.GetGachaPool(1)
	poolRecord := common.NewGachaPoolRecord(pool.Id)
	typeRecord := common.NewGachaTypeRecord(pool.Type)

	CSV.Gacha.Gacha(1, pool, 1000000, poolRecord, typeRecord)

	//t.Log(poolRecord)
	//t.Log(typeRecord)

}
func TestGachaEntry_Gacha_2(t *testing.T) {
	pool, _ := CSV.Gacha.GetGachaPool(2)
	poolRecord := common.NewGachaPoolRecord(pool.Id)
	typeRecord := common.NewGachaTypeRecord(pool.Type)

	CSV.Gacha.Gacha(1, pool, 1000, poolRecord, typeRecord)

}

func TestGachaEntry_Gacha_3(t *testing.T) {
	pool, _ := CSV.Gacha.GetGachaPool(3)
	poolRecord := common.NewGachaPoolRecord(pool.Id)
	typeRecord := common.NewGachaTypeRecord(pool.Type)

	CSV.Gacha.Gacha(1, pool, 1000000, poolRecord, typeRecord)

}

func TestGachaEntry_Gacha_4(t *testing.T) {
	pool, _ := CSV.Gacha.GetGachaPool(4)
	poolRecord := common.NewGachaPoolRecord(pool.Id)
	typeRecord := common.NewGachaTypeRecord(pool.Type)

	CSV.Gacha.Gacha(1, pool, 1000000, poolRecord, typeRecord)

}

func TestGachaEntry_Gacha_5(t *testing.T) {
	pool, _ := CSV.Gacha.GetGachaPool(5)
	poolRecord := common.NewGachaPoolRecord(pool.Id)
	typeRecord := common.NewGachaTypeRecord(pool.Type)

	CSV.Gacha.Gacha(1, pool, 1000000, poolRecord, typeRecord)

}

func TestGachaEntry_Rand(t *testing.T) {
	for i := 0; i < 1000; i++ {
		rangeInt32 := rand.RangeInt32(0, 1)
		t.Log(rangeInt32)
	}

}
