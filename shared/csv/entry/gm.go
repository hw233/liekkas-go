package entry

import (
	"sync"

	"shared/common"
	"shared/utility/errors"
)

type GM struct {
	sync.RWMutex

	gmCodes map[int32]string
}

func NewGM() *GM {
	return &GM{
		gmCodes: map[int32]string{},
	}
}

func (p *GM) Check(config *Config) error {
	gmCodes := map[int32]bool{}

	for _, gm := range config.GmConfig.GetAllData() {
		if gmCodes[gm.Id] {
			// repeat
			return errors.New("repeat data")
		}

		gmCodes[gm.Id] = true
	}

	return nil
}

func (p *GM) Reload(config *Config) error {
	p.Lock()
	defer p.Unlock()

	gmCodes := map[int32]string{}

	for _, gm := range config.GmConfig.GetAllData() {
		gmCodes[gm.Id] = gm.Code
	}

	p.gmCodes = gmCodes

	return nil
}

func (p *GM) GetGMCode(i int32) (string, error) {
	p.RLock()
	defer p.RUnlock()

	code, ok := p.gmCodes[i]
	if !ok {
		return "", errors.Swrapf(common.ErrNotFoundInCSV, "gm.csv", i)
	}

	return code, nil
}
