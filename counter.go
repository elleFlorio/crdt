package main

import (
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/elleFlorio/crdt/network"
)

type counter struct {
	id        uint32
	replica   string
	stateInc  map[string]int64
	stateDec  map[string]int64
	bufferInc map[string]int64
	bufferDec map[string]int64
	bufLock   sync.Mutex
	net       network.Overlay
	ch        chan network.Message
}

func NewCounter(net network.Overlay, syncTimeMs int) *counter {
	return newCounter(nil, net, syncTimeMs)
}

func ConnectCounter(id uint32, net network.Overlay, syncTimeMs int) *counter {
	return newCounter(&id, net, syncTimeMs)
}

func newCounter(id *uint32, net network.Overlay, syncTimeMs int) *counter {
	var cntId uint32
	if id == nil {
		cntId = uuid.New().ID()
	} else {
		cntId = *id
	}
	cntChan := make(chan network.Message, 10)
	cntStateInc := make(map[string]int64)
	cntStateDec := make(map[string]int64)
	cntBufInc := make(map[string]int64)
	cntBufDec := make(map[string]int64)
	cnt := &counter{
		id:        cntId,
		replica:   net.GetLocalAddr(),
		stateInc:  cntStateInc,
		stateDec:  cntStateDec,
		bufferInc: cntBufInc,
		bufferDec: cntBufDec,
		net:       net,
		ch:        cntChan,
	}

	net.Listen(cntChan)
	cnt.listen()
	cnt.synchronize(syncTimeMs)

	return cnt
}

func (c *counter) listen() {
	go func() {
		for msg := range c.ch {
			if msg.Id == c.id {
				received := msg.Payload.(counterDState)
				dState := c.getDelta(received)
				c.store(dState)
			}
		}
	}()
}

func (c *counter) store(dState counterDState) {
	c.bufLock.Lock()
	if dState.DStateInc > 0 {
		c.stateInc[dState.Replica] = dState.DStateInc
		c.bufferInc[dState.Replica] = dState.DStateInc
	}

	if dState.DStateDec > 0 {
		c.stateDec[dState.Replica] = dState.DStateDec
		c.bufferDec[dState.Replica] = dState.DStateDec
	}

	c.bufLock.Unlock()
}

func (c *counter) getDelta(dState counterDState) counterDState {
	// No need to compute the minimum delta
	return dState
}

func (c *counter) synchronize(intevalMs int) {
	ticker := time.NewTicker(time.Duration(intevalMs) * time.Millisecond)

	go func() {
		for range ticker.C {
			replicas := c.net.GetNodes()
			dStates := make(map[string]*counterDState)

			c.bufLock.Lock()
			for _, replica := range replicas {

				for replicaBuf, dStateInc := range c.bufferInc {
					if replica != replicaBuf && replica != c.replica {
						if dState, ok := dStates[replicaBuf]; ok {
							dState.DStateInc = dStateInc
						} else {
							dStates[replica] = &counterDState{replicaBuf, dStateInc, 0}
						}
					}
				}

				for replicaBuf, dStateDec := range c.bufferDec {
					if replica != replicaBuf && replica != c.replica {
						if dState, ok := dStates[replicaBuf]; ok {
							dState.DStateDec = dStateDec
						} else {
							dStates[replica] = &counterDState{replicaBuf, 0, dStateDec}
						}
					}
				}
			}

			for replica, payload := range dStates {
				msg := network.Message{
					Id:      c.id,
					Payload: *payload,
				}
				c.net.Send(msg, replica)
			}

			c.bufferInc = make(map[string]int64)
			c.bufferDec = make(map[string]int64)
			c.bufLock.Unlock()
		}

	}()
}

func (c *counter) Id() uint32 {
	return c.id
}

func (c *counter) Inc() {
	delta := c.stateInc[c.replica] + 1
	dState := counterDState{c.replica, delta, 0}
	c.store(dState)
}

func (c *counter) Dec() {
	delta := c.stateDec[c.replica] + 1
	dState := counterDState{c.replica, 0, delta}
	c.store(dState)
}

func (c *counter) Read() int64 {
	var total int64 = 0
	for _, inc := range c.stateInc {
		total += inc
	}
	for _, dec := range c.stateDec {
		total -= dec
	}

	return total
}
