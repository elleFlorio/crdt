package main

type gCounterDState struct {
	Replica string
	DState  uint64
}

type gSetDState struct {
	Replica string
	DState  []GSetElement
}

type counterDState struct {
	Replica   string
	DStateInc int64
	DStateDec int64
}
