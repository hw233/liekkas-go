/*
 *
 * Copyright 2017 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package roundrobin defines a roundrobin balancer. Roundrobin balancer is
// installed as one of the default balancers in gRPC, users don't need to
// explicitly install this balancer.
package round_robin

import (
	"math/rand"
	"sync"
	"time"

	"shared/grpc/balancer"

	googleBalancer "google.golang.org/grpc/balancer"
	googleBase "google.golang.org/grpc/balancer/base"
)

// Name is the name of round_robin balancer.
const Name = "round_robin"

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type rrPickerBuilder struct{}

func NewRRPickerBuilder() *rrPickerBuilder {
	return &rrPickerBuilder{}
}

func (*rrPickerBuilder) Build(info googleBase.PickerBuildInfo) googleBalancer.Picker {
	if len(info.ReadySCs) == 0 {
		return balancer.NewErrPicker(googleBalancer.ErrNoSubConnAvailable)
	}
	scs := make([]googleBalancer.SubConn, 0, len(info.ReadySCs))
	for sc := range info.ReadySCs {
		scs = append(scs, sc)
	}
	return &rrPicker{
		subConns: scs,
		// Start at a random index, as the same RR balancer rebuilds a new
		// picker when SubConn states change, and we don't want to apply excess
		// load to the first server in the list.

		next: rand.Intn(len(scs)),
	}
}

type rrPicker struct {
	// subConns is the snapshot of the roundrobin balancer when this picker was
	// created. The slice is immutable. Each GetBalance() will do a round robin
	// selection from it and return the selected SubConn.
	subConns []googleBalancer.SubConn

	mu   sync.Mutex
	next int
}

func (p *rrPicker) Pick(googleBalancer.PickInfo) (googleBalancer.PickResult, error) {
	p.mu.Lock()
	sc := p.subConns[p.next]
	p.next = (p.next + 1) % len(p.subConns)
	p.mu.Unlock()
	return googleBalancer.PickResult{SubConn: sc}, nil
}
