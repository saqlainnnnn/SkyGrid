package p2p

import (
	"fmt"
	"net"
)

//TCPPeer represents the remote node over a TCP established connection.
type TCPPeer struct {
	//conn is the underlying connection of the peer
	conn 		net.Conn
	//if we dial and retrieve a conn => outbound ==true
	//if we listen and retrieve a conn => outbound == false
	outbound 	bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn: conn,
		outbound: outbound,
	}
}

//close implements the peer interface
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPTransportOps struct {
	ListenAddr 		string
	HandshakeFunc 	HandshakeFunc
	Decoder 		Decoder
	OnPeer 			func(Peer) error
}

type TCPTransport struct {
	TCPTransportOps
	listener      net.Listener
	rpcch			chan RPC

}

//Consume implements the Transport interface, which will return read-only channel
//for reading the incoming messages recieved from another peer in the network
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

func NewTCPTransport(opts TCPTransportOps) *TCPTransport {
	return &TCPTransport{
	// 	listenAddress: listenAddr,
	// 	shakeHands: NOPHandshakeFunc,
	// 
		TCPTransportOps: opts,
		rpcch: make(chan RPC),
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAccpetLoop()

	return nil
}

func (t *TCPTransport) startAccpetLoop()  {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)	
		}

		go t.handleConn(conn)
	}
}

//handles lifecycle of a conn
func (t *TCPTransport) handleConn(conn net.Conn)  {
	
	var err error	
	
	defer func ()  {
		fmt.Printf("dropping peer connection: %+v\n", err)
		conn.Close()	
	}()
	
	peer := NewTCPPeer(conn, true)
	
	//checks for handshake to authenticate the peer, if not then it drops
	if err = t.HandshakeFunc(peer); err !=nil {
		return	
	}

	// We check if the user has provided OnPeer if not we drop the function
	if t.OnPeer != nil {
		if err = t.OnPeer(peer);  err != nil {
			return
		}

	}
	
	// Read LOOP
	rpc := &RPC{}

	for {
		if err := t.Decoder.Decode(conn, rpc); err != nil {
			fmt.Printf("TCP error: %s\n", err)
		}
		//adds address 
		rpc.From = conn.RemoteAddr()
		t.rpcch <- *rpc	
	
	}
}
	


