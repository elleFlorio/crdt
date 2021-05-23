package network

type Overlay interface {
	GetLocalAddr() string
	Connect(seed string)
	GetNodes() []string
	Listen(chan Message)
	Send(msg Message, node string)
}
