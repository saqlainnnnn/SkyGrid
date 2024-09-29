package p2p

import (
	"fmt"
	"net"
	"sync"
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

type TCPTransportOps struct {
	ListenAddr string
	HandshakeFunc HandshakeFunc
	Decoder Decoder
}

type TCPTransport struct {
	TCPTransportOps
	listener      net.Listener

	mu 		sync.RWMutex // read-write mutex, reading can be done
						// by multiple threads but exclusive writing
	peers 	map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOps) *TCPTransport {
	return &TCPTransport{
	// 	listenAddress: listenAddr,
	// 	shakeHands: NOPHandshakeFunc,
	// 
		TCPTransportOps: opts, 
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


func (t *TCPTransport) handleConn(conn net.Conn)  {
	peer := NewTCPPeer(conn, true)
	
	if err := t.HandshakeFunc(peer); err !=nil {
		conn.Close()
		fmt.Printf("TCP handshake error: %s\n", err)
		return	
	}
	
	// Read LOOP
	msg := &Message{}

	for {
		if err := t.Decoder.Decode(conn, msg); err != nil {
			fmt.Printf("TCP error: %s\n", err)
		}
		fmt.Printf("message: %v\n", msg)
	}
}
	


