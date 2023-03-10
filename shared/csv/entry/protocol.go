package entry

import (
	"reflect"
	"sync"

	"google.golang.org/protobuf/proto"

	"shared/common"
	"shared/utility/errors"
)

type PushCommands struct {
	ExploreEventNotify    int32
	ExploreResourceNotify int32
	QuestUpdateNotify     int32
	TowerUpdateNotify     int32
	LevelNotify           int32
	DailyRefreshNotify    int32
	MailNotify            int32
	ScorePassNotify       int32
	StoreDailyRefresh     int32
}

type Protocol struct {
	sync.RWMutex

	protocols map[int32]string

	Pushes            *PushCommands
	fullNameToCommand map[string]int32
}

func NewProtocol() *Protocol {
	return &Protocol{
		protocols:         map[int32]string{},
		Pushes:            &PushCommands{},
		fullNameToCommand: map[string]int32{},
	}
}

func (p *Protocol) Check(config *Config) error {
	protocols := map[int32]bool{}

	for _, protocol := range config.ProtocolConfig.GetAllData() {
		if protocols[protocol.Id] {
			// repeat
			return errors.New("repeat data")
		}

		protocols[protocol.Id] = true
	}

	return nil
}

func (p *Protocol) Reload(config *Config) error {
	p.Lock()
	defer p.Unlock()

	protocols := map[int32]string{}
	nameToCommand := map[string]int32{}
	fullNameToCommand := map[string]int32{}

	pushes := &PushCommands{}

	for _, protocol := range config.ProtocolConfig.GetAllData() {
		protocols[protocol.Id] = protocol.Function
		nameToCommand[protocol.Function] = protocol.Id
		if protocol.Type == 0 {
			fullNameToCommand["C2S"+protocol.Function] = protocol.Id
			fullNameToCommand["S2C"+protocol.Function] = protocol.Id + 1
		} else if protocol.Type == 1 {
			fullNameToCommand["S2C"+protocol.Function] = protocol.Id
		}
	}

	pushesType := reflect.TypeOf(pushes).Elem()
	pushesValue := reflect.ValueOf(pushes).Elem()

	for i := 0; i < pushesType.NumField(); i++ {
		fieldType := pushesType.Field(i)
		fieldValue := pushesValue.Field(i)

		pushName := fieldType.Name
		cmd, ok := nameToCommand[pushName]
		if !ok {
			return errors.WrapTrace(errors.Swrapf(common.ErrProtocolCmdNotFound, pushName))
		}

		fieldValue.Set(reflect.ValueOf(cmd))

	}

	p.protocols = protocols
	p.Pushes = pushes
	p.fullNameToCommand = fullNameToCommand
	return nil
}

func (p *Protocol) Protocols() map[int32]string {
	p.RLock()
	defer p.RUnlock()

	return p.protocols
}

func (p *Protocol) GetCmdByProtoName(msg proto.Message) (int32, error) {
	p.RLock()
	defer p.RUnlock()
	fullName := string(proto.MessageName(msg).Name())
	cmd, ok := p.fullNameToCommand[fullName]
	if !ok {
		return 0, errors.WrapTrace(errors.Swrapf(common.ErrProtocolCmdNotFound, fullName))
	}
	return cmd, nil
}
