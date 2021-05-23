package main

import (
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/elleFlorio/crdt/network"
)

type gCounter struct {
	id      uint32
	replica string
	state   map[string]uint64
	buffer  map[string]uint64
	bufLock sync.Mutex
	net     network.Overlay
	ch      chan network.Message
}

func NewGCounter(net network.Overlay, syncTimeMs int) *gCounter {
	return newGCounter(nil, net, syncTimeMs)
}

func ConnectGCounter(id uint32, net network.Overlay, syncTimeMs int) *gCounter {
	return newGCounter(&id, net, syncTimeMs)
}

func newGCounter(id *uint32, net network.Overlay, syncTimeMs int) *gCounter {
	var gcntId uint32
	if id == nil {
		gcntId = uuid.New().ID()
	} else {
		gcntId = *id
	}
	gcntChan := make(chan network.Message, 10)
	gcntState := make(map[string]uint64)
	gcntBuf := make(map[string]uint64)
	cnt := &gCounter{
		id:      gcntId,
		replica: net.GetLocalAddr(),
		state:   gcntState,
		buffer:  gcntBuf,
		net:     net,
		ch:      gcntChan,
	}

	net.Listen(gcntChan)
	cnt.listen()
	cnt.synchronize(syncTimeMs)

	return cnt
}

func (gc *gCounter) listen() {
	go func() {
		for msg := range gc.ch {
			if msg.Id == gc.id {
				received := msg.Payload.(gCounterDState)
				dState := gc.getDelta(received)
				gc.store(dState)
			}
		}
	}()
}

func (gc *gCounter) store(dState gCounterDState) {
	gc.state[dState.Replica] = dState.DState
	gc.bufLock.Lock()
	gc.buffer[dState.Replica] = dState.DState
	gc.bufLock.Unlock()
}

func (gc *gCounter) getDelta(dState gCounterDState) gCounterDState {
	// No need to compute the minimum delta
	return dState
}

func (gc *gCounter) synchronize(intevalMs int) {
	ticker := time.NewTicker(time.Duration(intevalMs) * time.Millisecond)

	go func() {
		for range ticker.C {
			replicas := gc.net.GetNodes()
			gc.bufLock.Lock()
			for replicaBuf, dState := range gc.buffer {
				for _, replica := range replicas {
					if replica != replicaBuf && replica != gc.replica {
						msg := network.Message{
							Id:      gc.id,
							Payload: gCounterDState{replicaBuf, dState},
						}

						gc.net.Send(msg, replica)
					}
				}
			}

			gc.buffer = make(map[string]uint64)
			gc.bufLock.Unlock()
		}

	}()
}

func (gc *gCounter) Id() uint32 {
	return gc.id
}

func (gc *gCounter) Inc() {
	gc.state[gc.replica] += 1
	delta := gc.state[gc.replica]
	dState := gCounterDState{gc.replica, delta}
	gc.store(dState)
}

func (gc *gCounter) Read() uint64 {
	var total uint64 = 0
	for _, counter := range gc.state {
		total += counter
	}

	return total
}
