package entry

import (
	"shared/common"
	"shared/utility/errors"
	"shared/utility/transfer"
	"sync"
)

type MailTemplate struct {
	Id         int32
	ExpireDays int32
	Title      string
	Content    string
}

type MailEntry struct {
	sync.RWMutex

	MailTemplates map[int32]*MailTemplate
}

const (
	MailTemplateGM int32 = 101
)

func NewMailEntry() *MailEntry {
	return &MailEntry{
		MailTemplates: map[int32]*MailTemplate{},
	}
}

func (me *MailEntry) Check(config *Config) error {
	return nil
}

func (me *MailEntry) Reload(config *Config) error {
	me.Lock()
	defer me.Unlock()

	err := me.checkTemplete(config)
	if err != nil {
		return errors.WrapTrace(err)
	}

	for _, templateCSV := range config.CfgMailTemplateConfig.GetAllData() {
		template := &MailTemplate{}

		err := transfer.Transfer(templateCSV, template)
		if err != nil {
			return errors.WrapTrace(err)
		}

		me.MailTemplates[template.Id] = template
	}

	return nil
}

func (me *MailEntry) GetTemplate(id int32) (*MailTemplate, error) {
	template, ok := me.MailTemplates[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrMailTemplateConfigNotFound, id)
	}

	return template, nil
}

func (me *MailEntry) checkTemplete(config *Config) error {
	templateIds := []int32{
		MailTemplateGM,
	}

	for _, id := range templateIds {
		_, ok := config.CfgMailTemplateConfig.Find(id)
		if !ok {
			return errors.Swrapf(common.ErrMailTemplateConfigNotFound, id)
		}
	}

	return nil
}
