package network

import (
	"strconv"
	"sync"

	"github.com/google/uuid"
)

var (
	lnCounter int                  = 0
	lnNodes   map[string]localNode = make(map[string]localNode)
	lnMutex   sync.Mutex
)

type localNode struct {
	Id      uuid.UUID
	Address int
	ch      chan Message
}

func CreateLocalNode() *localNode {
	lnId := uuid.New()
	lnChan := make(chan Message, 10)

	lnMutex.Lock()
	ln := localNode{lnId, lnCounter, lnChan}
	lnNodes[strconv.Itoa(ln.Address)] = ln
	lnCounter += 1
	lnMutex.Unlock()

	return &ln
}

func (ln *localNode) GetLocalAddr() string {
	return strconv.Itoa(ln.Address)
}

func (ln *localNode) Connect(seed string) {
	// No need locally
}

func (ln *localNode) GetNodes() []string {
	addresses := make([]string, 0, len(lnNodes))
	for address := range lnNodes {
		addresses = append(addresses, address)
	}

	return addresses
}

func (ln *localNode) Listen(ch chan Message) {
	go func() {
		for msg := range ln.ch {
			ch <- msg
		}
	}()
}

func (ln *localNode) Send(msg Message, node string) {
	receiver := lnNodes[node]
	receiver.ch <- msg
}
