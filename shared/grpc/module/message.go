package module

const (
	AttrBalancerMessage = "balancer_message"
)

type Event uint8

const (
	Register   Event = iota // 服务正常
	Unregister              // 服务下线
)

func (e Event) IsRegister() bool {
	return e == Register
}

func (e Event) IsUnregister() bool {
	return e == Unregister
}

type ResolverMessage struct {
	Event  Event
	Server string
	Addr   string
}

func NewResolverMessage(server, addr string) *ResolverMessage {
	return &ResolverMessage{
		Event:  Register,
		Server: server,
		Addr:   addr,
	}
}

type BalancerMessage struct {
	Event Event
}

func NewBalancerMessage(event Event) *BalancerMessage {
	return &BalancerMessage{
		Event: event,
	}
}
