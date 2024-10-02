package p2p

import "net"

//RPC hold any arbitary data ie being sent over //
//the each transport between two nodes in the network
type RPC struct {
	From net.Addr
	Payload []byte
}