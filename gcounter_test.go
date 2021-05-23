package main

import (
	"testing"
	"time"

	"github.com/elleFlorio/crdt/network"
)

func TestGCounter(t *testing.T) {
	var state1 uint64 = 0
	var state2 uint64 = 0

	node1 := network.CreateLocalNode()
	node2 := network.CreateLocalNode()

	gcounter1 := NewGCounter(node1, 10)
	gcounter2 := ConnectGCounter(gcounter1.Id(), node2, 10)

	gcounter1.Inc()
	gcounter1.Inc()

	time.Sleep(time.Duration(20) * time.Millisecond)

	state1 = gcounter1.Read()
	state2 = gcounter2.Read()

	if state1 != state2 {
		t.Fatalf("States are not equals: State1 %d - State2 %d", state1, state2)
	}

	gcounter2.Inc()

	time.Sleep(time.Duration(20) * time.Millisecond)

	state1 = gcounter1.Read()
	state2 = gcounter2.Read()

	if state1 != state2 {
		t.Fatalf("States are not equals: State1 %d - State2 %d", state1, state2)
	}
}
