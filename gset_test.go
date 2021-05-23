package main

import (
	"hash/fnv"
	"sort"
	"testing"
	"time"

	"github.com/elleFlorio/crdt/network"
)

func getId(value string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(value))
	return h.Sum32()
}

func TestGSet(t *testing.T) {
	var state1 []string = make([]string, 0, 3)
	var state2 []string = make([]string, 0, 3)

	node1 := network.CreateLocalNode()
	node2 := network.CreateLocalNode()

	gset1 := NewGSet(node1, 10)
	gset2 := ConnectGSet(gset1.Id(), node2, 10)

	gset1.Add(GSetElement{Id: getId("a"), Value: "a"})
	gset1.Add(GSetElement{Id: getId("b"), Value: "b"})

	time.Sleep(time.Duration(20) * time.Millisecond)

	for _, e := range gset1.Read() {
		state1 = append(state1, e.Value.(string))
	}

	for _, e := range gset2.Read() {
		state2 = append(state2, e.Value.(string))
	}

	if len(state1) != len(state2) {
		t.Fatalf("States are not equals: State1 %v - State2 %v", state1, state2)
	}

	sort.Strings(state1)
	sort.Strings(state2)

	for i := range state1 {
		if state1[i] != state2[i] {
			t.Fatalf("States are not equals: State1 %v - State2 %v", state1, state2)
		}
	}

	gset2.Add(GSetElement{Id: getId("c"), Value: "c"})

	time.Sleep(time.Duration(20) * time.Millisecond)

	for _, e := range gset1.Read() {
		state1 = append(state1, e.Value.(string))
	}

	for _, e := range gset2.Read() {
		state2 = append(state2, e.Value.(string))
	}

	if len(state1) != len(state2) {
		t.Fatalf("States are not equals: State1 %v - State2 %v", state1, state2)
	}

	sort.Strings(state1)
	sort.Strings(state2)

	for i := range state1 {
		if state1[i] != state2[i] {
			t.Fatalf("States are not equals: State1 %v - State2 %v", state1, state2)
		}
	}

	gset1.Add(GSetElement{Id: getId("a"), Value: "a"})
	gset2.Add(GSetElement{Id: getId("c"), Value: "c"})

	time.Sleep(time.Duration(20) * time.Millisecond)

	for _, e := range gset1.Read() {
		state1 = append(state1, e.Value.(string))
	}

	for _, e := range gset2.Read() {
		state2 = append(state2, e.Value.(string))
	}

	if len(state1) != len(state2) {
		t.Fatalf("States are not equals: State1 %v - State2 %v", state1, state2)
	}

	sort.Strings(state1)
	sort.Strings(state2)

	for i := range state1 {
		if state1[i] != state2[i] {
			t.Fatalf("States are not equals: State1 %v - State2 %v", state1, state2)
		}
	}
}
