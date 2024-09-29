package p2p

//Meassage hold any arbitary data ie being sent over //
//the each transport between two nodes in the network
type Message struct {
	Payload []byte
}