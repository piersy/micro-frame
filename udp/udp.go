package udp
import (
	"net"
	"fmt"
	"runtime"
)

const UDP_TYPE = "udp4"

type Tunnel interface {
	Send(bytes []byte)
	Read() []byte
}

type udpTunnel struct {
	conn       *net.UDPConn
	targetAddr *net.UDPAddr
}


func (t *udpTunnel) Send(bytes []byte) {
	_, err := t.conn.WriteToUDP(bytes, t.targetAddr)
	E(err)
}

func (t *udpTunnel) Read() []byte {
	b := make([]byte, 2^16)
	_, _, err := t.conn.ReadFromUDP(b)
	E(err)
	return b
}

func NewUdpAddress(ip string, port int) *net.UDPAddr {
	return &net.UDPAddr{IP: net.ParseIP(ip), Port: port}
}

func Listen(port int, notifyConnection func(t Tunnel)) {
	handOffNewConn := func(remoteAddr *net.UDPAddr) {
		//Open a new listen port
		persistentConn, err := net.ListenUDP(UDP_TYPE, nil)
		E(err)
		//notify new connection
		tunnel := &udpTunnel{persistentConn, remoteAddr}
		go notifyConnection(tunnel)
		runtime.Gosched()
		//notify client of new port
		tunnel.Send(nil)
	}
	listenConn, err := net.ListenUDP(UDP_TYPE, NewUdpAddress("localhost", port))
	fmt.Printf("server listening on %v \n\n", listenConn.LocalAddr())
	E(err)
	for {
		//Lsiten for incoming signals
		_, returnAddr, err := listenConn.ReadFromUDP(nil)
		E(err)
		fmt.Printf("server received incoming from %v\n\n", returnAddr)
		//Create new connection
		go handOffNewConn(returnAddr)
	}
}

func OpenConnection(targetHost string, targetPort int) Tunnel {
	//Listen on a port
	conn, err := net.ListenUDP(UDP_TYPE, nil)
	E(err)
	//InitiateNewConnection
	conn.WriteToUDP(nil, NewUdpAddress(targetHost, targetPort))
	//Read return adress
	_, add, err := conn.ReadFromUDP(nil)
	E(err)
	fmt.Printf("established new connection to %v \n\n", add)
	return &udpTunnel{conn, add}
}

func E(err error) {
	if err  != nil {
		panic(err)
	}
}
