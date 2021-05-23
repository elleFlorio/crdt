package main

import (
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/elleFlorio/crdt/network"
)

type GSetElement struct {
	Id    uint32
	Value interface{}
}

type gSet struct {
	id      uint32
	replica string
	state   map[string][]GSetElement
	buffer  map[string][]GSetElement
	bufLock sync.Mutex
	net     network.Overlay
	ch      chan network.Message
}

func NewGSet(net network.Overlay, syncTimeMs int) *gSet {
	return newGSet(nil, net, syncTimeMs)
}

func ConnectGSet(id uint32, net network.Overlay, syncTimeMs int) *gSet {
	return newGSet(&id, net, syncTimeMs)
}

func newGSet(id *uint32, net network.Overlay, syncTimeMs int) *gSet {
	var gsId uint32
	if id == nil {
		gsId = uuid.New().ID()
	} else {
		gsId = *id
	}
	gsChan := make(chan network.Message, 10)
	gsState := make(map[string][]GSetElement)
	gsBuf := make(map[string][]GSetElement)
	gs := &gSet{
		id:      gsId,
		replica: net.GetLocalAddr(),
		state:   gsState,
		buffer:  gsBuf,
		net:     net,
		ch:      gsChan,
	}

	net.Listen(gsChan)
	gs.listen()
	gs.synchronize(syncTimeMs)

	return gs
}

func (gs *gSet) listen() {
	go func() {
		for msg := range gs.ch {
			if msg.Id == gs.id {
				received := msg.Payload.(gSetDState)
				dState := gs.getDelta(received)
				gs.store(dState)
			}
		}
	}()
}

func (gs *gSet) store(dState gSetDState) {
	if len(dState.DState) > 0 {
		gs.state[dState.Replica] = append(gs.state[dState.Replica], dState.DState...)
		gs.bufLock.Lock()
		gs.buffer[dState.Replica] = append(gs.buffer[dState.Replica], dState.DState...)
		gs.bufLock.Unlock()
	}
}

func (gs *gSet) getDelta(dState gSetDState) gSetDState {
	local := gs.state[gs.replica]
	delta := arrayDif(local, dState.DState)

	return gSetDState{Replica: dState.Replica, DState: delta}
}

func (gs *gSet) synchronize(intevalMs int) {
	ticker := time.NewTicker(time.Duration(intevalMs) * time.Millisecond)

	go func() {
		for range ticker.C {
			replicas := gs.net.GetNodes()
			gs.bufLock.Lock()
			for replicaBuf, dState := range gs.buffer {
				for _, replica := range replicas {
					if replica != replicaBuf && replica != gs.replica {
						msg := network.Message{
							Id:      gs.id,
							Payload: gSetDState{replicaBuf, dState},
						}

						gs.net.Send(msg, replica)
					}
				}
			}

			gs.buffer = make(map[string][]GSetElement)
			gs.bufLock.Unlock()
		}

	}()
}

func (gs *gSet) Id() uint32 {
	return gs.id
}

func (gs *gSet) Add(element GSetElement) {
	dState := gSetDState{gs.replica, []GSetElement{element}}
	delta := gs.getDelta(dState)
	gs.store(delta)
}

func (gs *gSet) Read() []GSetElement {
	state := make([]GSetElement, 0)
	for _, elements := range gs.state {
		state = append(state, elements...)
	}

	return state
}
