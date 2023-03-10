package entry

import (
	"shared/common"
	"shared/utility/errors"
	"shared/utility/transfer"
	"sync"
)

type GreetingsEntry struct {
	sync.RWMutex

	Greetings map[int32]GreetingsConfig
}

type GreetingsConfig struct {
	Id     int32
	DropId int32
}

func NewGreetingsEntry() *GreetingsEntry {
	return &GreetingsEntry{}
}

func (g *GreetingsEntry) Check(config *Config) error {

	return nil
}
func (g *GreetingsEntry) Reload(config *Config) error {
	g.Lock()
	defer g.Unlock()

	greetings := map[int32]GreetingsConfig{}

	for _, greetingsCsv := range config.CfgGuildGreetingsConfig.GetAllData() {
		greetingsCfg := &GreetingsConfig{}

		err := transfer.Transfer(greetingsCsv, greetingsCfg)
		if err != nil {
			return errors.WrapTrace(err)
		}

		greetings[greetingsCsv.Id] = *greetingsCfg
	}

	g.Greetings = greetings

	return nil
}

func (g *GreetingsEntry) GetDropId(star int32) (int32, error) {
	g.Lock()
	defer g.Unlock()

	result, ok := g.Greetings[star]
	if !ok {
		return 0, errors.Swrapf(common.ErrGreetingsNotFoundInCsv, star)
	}
	return result.DropId, nil
}
