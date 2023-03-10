package balancer

import (
	"errors"
	"fmt"

	"shared/grpc/module"
	"shared/utility/glog"

	googleBalancer "google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/resolver"
)

type builder struct {
	name          string
	pickerBuilder base.PickerBuilder
	config        base.Config
}

func RegisterPickerBuilder(name string, pickerBuilder base.PickerBuilder) {
	builder := &builder{
		name:          name,
		pickerBuilder: pickerBuilder,
		config:        base.Config{HealthCheck: true},
	}

	googleBalancer.Register(builder)
}

func (b *builder) Build(cc googleBalancer.ClientConn, opt googleBalancer.BuildOptions) googleBalancer.Balancer {
	bal := &balancer{
		service:       b.name,
		cc:            cc,
		pickerBuilder: b.pickerBuilder,
		picker:        b.pickerBuilder.Build(base.PickerBuildInfo{}),

		subConns: newServerMap(),
		scStates: make(map[googleBalancer.SubConn]connectivity.State),
		csEvltor: &googleBalancer.ConnectivityStateEvaluator{},
		config:   b.config,
	}

	// Initialize picker to a picker that always returns
	// ErrNoSubConnAvailable, because when state of a SubConn changes, we
	// may call UpdateState with this picker.
	bal.picker = NewErrPicker(googleBalancer.ErrNoSubConnAvailable)

	return bal
}

func (b *builder) Name() string {
	return b.name
}

type balancer struct {
	service       string
	cc            googleBalancer.ClientConn
	pickerBuilder base.PickerBuilder

	csEvltor *googleBalancer.ConnectivityStateEvaluator
	state    connectivity.State

	// subConns *resolver.AddressMap // `attributes` is stripped from the keys of this map (the addresses)
	subConns *serverMap
	scStates map[googleBalancer.SubConn]connectivity.State
	picker   googleBalancer.Picker
	config   base.Config

	resolverErr error // the last error reported by the resolver; cleared on successful resolution
	connErr     error // the last connection error; cleared upon leaving TransientFailure
}

func (b *balancer) ResolverError(err error) {
	b.resolverErr = err

	if b.subConns.Len() == 0 {
		b.state = connectivity.TransientFailure
	}

	if b.state != connectivity.TransientFailure {
		// The picker will not change since the balancer does not currently
		// report an error.
		return
	}

	b.regeneratePicker()

	b.cc.UpdateState(googleBalancer.State{
		ConnectivityState: b.state,
		Picker:            b.picker,
	})
}

func (b *balancer) UpdateClientConnState(s googleBalancer.ClientConnState) error {
	// Successful resolution; clear resolver error and ensure we return nil.
	b.resolverErr = nil

	for _, a := range s.ResolverState.Addresses {
		// get balancer message
		balancerMessage := a.Attributes.Value(module.AttrBalancerMessage).(*module.BalancerMessage)

		if _, ok := b.subConns.Get(a); !ok {
			if balancerMessage.Event.IsRegister() {
				// new server connect
				glog.Infof("+ service [%s] register server [%s, %s]", b.service, a.ServerName, a.Addr)
				sc, err := b.cc.NewSubConn([]resolver.Address{a}, googleBalancer.NewSubConnOptions{HealthCheckEnabled: b.config.HealthCheck})
				if err != nil {
					glog.Warnf("failed to create new SubConn: %v", err)
					continue
				}
				b.subConns.Set(a, sc)
				b.scStates[sc] = connectivity.Idle
				b.csEvltor.RecordTransition(connectivity.Shutdown, connectivity.Idle)
				sc.Connect()
			}
		} else {
			if balancerMessage.Event.IsRegister() { // special condition, addr changed
				// remove server connect
				glog.Infof("- service [%s] unregister server [%s]", b.service, a.ServerName)
				sci, _ := b.subConns.Get(a)
				sc := sci.(googleBalancer.SubConn)
				b.cc.RemoveSubConn(sc)
				b.subConns.Delete(a)

				// new server connect
				glog.Infof("+ service [%s] register server [%s, %s]", b.service, a.ServerName, a.Addr)
				sc, err := b.cc.NewSubConn([]resolver.Address{a}, googleBalancer.NewSubConnOptions{HealthCheckEnabled: b.config.HealthCheck})
				if err != nil {
					glog.Warnf("failed to create new SubConn: %v", err)
					continue
				}
				b.subConns.Set(a, sc)
				b.scStates[sc] = connectivity.Idle
				b.csEvltor.RecordTransition(connectivity.Shutdown, connectivity.Idle)
				sc.Connect()
			} else if balancerMessage.Event.IsUnregister() {
				glog.Infof("- service [%s] unregister Server [%s]", b.service, a.ServerName)
				// remove server connect
				sci, _ := b.subConns.Get(a)
				sc := sci.(googleBalancer.SubConn)
				b.cc.RemoveSubConn(sc)
				b.subConns.Delete(a)
			}
		}
	}

	// If resolver state contains no addresses, return an error so ClientConn
	// will trigger re-resolve. Also records this as an resolver error, so when
	// the overall state turns transient failure, the error message will have
	// the zero address information.
	if len(s.ResolverState.Addresses) == 0 {
		b.ResolverError(errors.New("produced zero addresses"))
		return googleBalancer.ErrBadResolverState
	}

	return nil
}

// mergeErrors builds an error from the last connection error and the last
// resolver error.  Must only be called if b.state is TransientFailure.
func (b *balancer) mergeErrors() error {
	// connErr must always be non-nil unless there are no SubConns, in which
	// case resolverErr must be non-nil.
	if b.connErr == nil {
		return fmt.Errorf("last resolver error: %v", b.resolverErr)
	}
	if b.resolverErr == nil {
		return fmt.Errorf("last connection error: %v", b.connErr)
	}
	return fmt.Errorf("last connection error: %v; last resolver error: %v", b.connErr, b.resolverErr)
}

// regeneratePicker takes a snapshot of the balancer, and generates a picker
// from it. The picker is
//  - errPicker if the balancer is in TransientFailure,
//  - built by the pickerBuilder with all READY SubConns otherwise.
func (b *balancer) regeneratePicker() {
	if b.state == connectivity.TransientFailure {
		b.picker = NewErrPicker(b.mergeErrors())
		return
	}
	readySCs := make(map[googleBalancer.SubConn]base.SubConnInfo)

	// Filter out all ready SCs from full subConn map.
	for _, addr := range b.subConns.Keys() {
		sci, _ := b.subConns.Get(addr)
		sc := sci.(googleBalancer.SubConn)
		if st, ok := b.scStates[sc]; ok && st == connectivity.Ready {
			readySCs[sc] = base.SubConnInfo{Address: addr}
		}
	}
	b.picker = b.pickerBuilder.Build(base.PickerBuildInfo{ReadySCs: readySCs})
}

func (b *balancer) UpdateSubConnState(sc googleBalancer.SubConn, state googleBalancer.SubConnState) {
	s := state.ConnectivityState
	oldS, ok := b.scStates[sc]
	if !ok {
		glog.Infof("base.balancer: got state changes for an unknown SubConn: %p, %v", sc, s)
		return
	}
	if oldS == connectivity.TransientFailure && s == connectivity.Connecting {
		// Once a subconn enters TRANSIENT_FAILURE, ignore subsequent
		// CONNECTING transitions to prevent the aggregated state from being
		// always CONNECTING when many backends exist but are all down.
		return
	}
	b.scStates[sc] = s
	switch s {
	case connectivity.Idle:
		sc.Connect()
	case connectivity.Shutdown:
		// When an address was removed by resolver, b called RemoveSubConn but
		// kept the sc's state in scStates. Remove state for this sc here.
		delete(b.scStates, sc)
	case connectivity.TransientFailure:
		// Save error to be reported via picker.
		b.connErr = state.ConnectionError
	}

	b.state = b.csEvltor.RecordTransition(oldS, s)

	// Regenerate picker when one of the following happens:
	//  - this sc entered or left ready
	//  - the aggregated state of balancer is TransientFailure
	//    (may need to update error message)
	if (s == connectivity.Ready) != (oldS == connectivity.Ready) ||
		b.state == connectivity.TransientFailure {
		b.regeneratePicker()
	}

	b.cc.UpdateState(googleBalancer.State{ConnectivityState: b.state, Picker: b.picker})
}

// Stop is a nop because base balancer doesn't have internal state to clean up,
// and it doesn't need to call RemoveSubConn for the SubConns.
func (b *balancer) Close() {
}

// NewErrPicker returns a picker that always returns err on Pick().
func NewErrPicker(err error) googleBalancer.Picker {
	return &errPicker{err: err}
}

type errPicker struct {
	err error // Pick() always returns this err.
}

func (p *errPicker) Pick(info googleBalancer.PickInfo) (googleBalancer.PickResult, error) {
	return googleBalancer.PickResult{}, p.err
}
