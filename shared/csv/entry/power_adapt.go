package entry

import (
	"shared/common"
	"shared/utility/errors"
	"shared/utility/transfer"
	"sync"
)

type PowerAdaptEntry struct {
	sync.RWMutex

	CareerParms map[int32]*CareerParm
}

type CareerParm struct {
	Parms []*Parm
}

type Parm struct {
	HpAdaption           float64
	PhyAtkAdaption       float64
	MagAtkAdaption       float64
	PhyDfsAdaption       float64
	MagDfsAdaption       float64
	CritAtkRatioAdaption float64
}

type CareerData struct {
	Id                   int32
	CareerId             int32
	Symbol               int32
	HpAdaption           float64
	PhyAtkAdaption       float64
	MagAtkAdaption       float64
	PhyDfsAdaption       float64
	MagDfsAdaption       float64
	CritAtkRatioAdaption float64
}

func NewPowerAdaptEntry() *PowerAdaptEntry {
	return &PowerAdaptEntry{}
}

func NewCareerParm() *CareerParm {
	return &CareerParm{
		Parms: make([]*Parm, 2, 2),
	}
}

func (p *PowerAdaptEntry) Check(config *Config) error {
	return nil
}

func (p *PowerAdaptEntry) Reload(config *Config) error {
	p.Lock()
	defer p.Unlock()

	careerParms := map[int32]*CareerParm{}

	for _, careerCsv := range config.CfgCombatPowerAdaptConfig.GetAllData() {
		careerCfg := &CareerData{}

		err := transfer.Transfer(careerCsv, careerCfg)
		if err != nil {
			return errors.WrapTrace(err)
		}

		careerParm, ok := careerParms[careerCfg.CareerId]
		if !ok {
			careerParm = NewCareerParm()
			careerParms[careerCfg.CareerId] = careerParm
		}

		parm := &Parm{}

		parm.HpAdaption = careerCfg.HpAdaption
		parm.PhyAtkAdaption = careerCfg.PhyAtkAdaption
		parm.MagAtkAdaption = careerCfg.MagAtkAdaption
		parm.PhyDfsAdaption = careerCfg.PhyDfsAdaption
		parm.MagDfsAdaption = careerCfg.MagDfsAdaption
		parm.CritAtkRatioAdaption = careerCfg.CritAtkRatioAdaption

		if careerCfg.Symbol < 1 || careerCfg.Symbol > 2 {
			return errors.Swrapf(common.ErrPowerAdaptWrongNumberForSymbol, careerCfg.Symbol)
		}

		careerParm.Parms[careerCfg.Symbol-1] = parm
		careerParms[careerCfg.CareerId] = careerParm
	}

	p.CareerParms = careerParms

	return nil
}

// 根据职业id和乘除方向返回对应的参数，symbol为1代表乘法，为2代表除法
func (p *PowerAdaptEntry) GetCareerParm(careerID, symbol int32) (*Parm, error) {
	careerParm, ok := p.CareerParms[careerID]
	if !ok {
		return nil, errors.Swrapf(common.ErrPowerAdaptWrongCareerIDForParm, careerID)
	}

	if symbol != 1 && symbol != 2 {
		return nil, errors.Swrapf(common.ErrPowerAdaptWrongSymbolForParm, symbol)
	}

	return careerParm.Parms[symbol-1], nil
}
