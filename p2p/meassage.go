package p2p

import "net"

//Meassage hold any arbitary data ie being sent over //
//the each transport between two nodes in the network
type Message struct {
	From net.Addr
	Payload []byte
}