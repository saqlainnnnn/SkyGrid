package p2p

//Peer is an interface that represents the remote connection
type Peer interface {

}

//Transport is anything that handles the communication 
//between the nodes in the network, this can be of (TCP, UPD, WEBSOCKETS)
type Transport interface {
	ListenAndAccept() error
}