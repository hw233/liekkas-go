package entry

import (
	"shared/common"
	"shared/utility/errors"
	"shared/utility/transfer"
	"sync"
)

type PatchCfg struct {
	Channel         string   `json:"channel"`
	AppVersion      string   `json:"app_version"`
	ResourceUrl     []string `json:"resource_url"`
	ResourceVersion string   `json:"resource_version"`
}

type Patches struct {
	sync.RWMutex

	Patches map[string]map[string]*PatchCfg
}

func NewPatches() *Patches {
	return &Patches{}
}

func (p *Patches) Reload(config *Config) error {
	p.Lock()
	defer p.Unlock()

	patches := map[string]map[string]*PatchCfg{}

	for _, patchCSV := range config.PatchListConfig.GetAllData() {
		patchCfg := &PatchCfg{}

		err := transfer.Transfer(patchCSV, patchCfg)
		if err != nil {
			return errors.WrapTrace(err)
		}

		channel := patchCfg.Channel
		channelPatches, ok := patches[channel]
		if !ok {
			channelPatches = map[string]*PatchCfg{}
			patches[channel] = channelPatches
		}
		channelPatches[patchCfg.AppVersion] = patchCfg
	}

	p.Patches = patches

	return nil
}

func (p *Patches) GetPatchCfg(channel, appVersion string) (*PatchCfg, error) {
	channelPatches, ok := p.Patches[channel]
	if !ok {
		return nil, common.ErrPatchCfgNotFound
	}

	patchCfg, ok := channelPatches[appVersion]
	if !ok {
		return nil, common.ErrPatchCfgNotFound

	}

	return patchCfg, nil
}
