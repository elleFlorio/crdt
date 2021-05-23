package main

import (
	"testing"
	"time"

	"github.com/elleFlorio/crdt/network"
)

func TestCounter(t *testing.T) {
	var state1 int64 = 0
	var state2 int64 = 0

	node1 := network.CreateLocalNode()
	node2 := network.CreateLocalNode()

	counter1 := NewCounter(node1, 10)
	counter2 := ConnectCounter(counter1.Id(), node2, 10)

	counter1.Inc()
	counter1.Inc()

	time.Sleep(time.Duration(20) * time.Millisecond)

	state1 = counter1.Read()
	state2 = counter2.Read()

	if state1 != 2 {
		t.Fatalf("Reference state is not correct: State %d - Expected %d", state1, 2)
	}

	if state1 != state2 {
		t.Fatalf("States are not equals: State1 %d - State2 %d", state1, state2)
	}

	counter2.Inc()

	time.Sleep(time.Duration(20) * time.Millisecond)

	state1 = counter1.Read()
	state2 = counter2.Read()

	if state2 != 3 {
		t.Fatalf("Reference state is not correct: State %d - Expected %d", state2, 3)
	}

	if state1 != state2 {
		t.Fatalf("States are not equals: State1 %d - State2 %d", state1, state2)
	}

	counter1.Dec()
	counter1.Dec()

	time.Sleep(time.Duration(20) * time.Millisecond)

	state1 = counter1.Read()
	state2 = counter2.Read()

	if state1 != 1 {
		t.Fatalf("Reference state is not correct: State %d - Expected %d", state1, 1)
	}

	if state1 != state2 {
		t.Fatalf("States are not equals: State1 %d - State2 %d", state1, state2)
	}

	counter2.Dec()

	time.Sleep(time.Duration(20) * time.Millisecond)

	state1 = counter1.Read()
	state2 = counter2.Read()

	if state2 != 0 {
		t.Fatalf("Reference state is not correct: State %d - Expected %d", state2, 0)
	}

	if state1 != state2 {
		t.Fatalf("States are not equals: State1 %d - State2 %d", state1, state2)
	}
}
