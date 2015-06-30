package udp
import (
    "net"
    "fmt"
    "time"
    "strconv"
    "runtime"
)

const UDP_TYPE = "udp4"

type Tunnel interface {
    Send(bytes []byte)
}

type DefaultTunnel struct {
    conn *net.UDPConn
    targetAddr *net.UDPAddr
}


func (t *DefaultTunnel) Send(bytes []byte){
    _, err:= t.conn.WriteToUDP(bytes, t.targetAddr)
    E(err)
}

func NewUdpAddress(ip string, port int) *net.UDPAddr {
    return &net.UDPAddr{IP: net.ParseIP(ip), Port: port}
}

func Listen(port int, notifyConnection func(*net.UDPConn, *net.UDPAddr)) {
    handOffNewConn := func(remoteAddr *net.UDPAddr) {
        //Open a new listen port
        persistentConn, err := net.ListenUDP(UDP_TYPE, nil)
        E(err)
        //notify new connection
        go notifyConnection(persistentConn, remoteAddr)
        runtime.Gosched()
        //notify client of new port
        _, err = persistentConn.WriteToUDP(nil, remoteAddr)
        E(err)
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
    _,add,err := conn.ReadFromUDP(nil)
    E(err)
    fmt.Printf("established new connection to %v \n\n", add)
    return &DefaultTunnel{conn, add}
}

func OpenTunnel(targetHost string) *net.UDPConn {
    //    addrs, err := net.LookupHost(host)
    //    if err != nil {
    //        panic(err)
    //    }

    //    udpConn, _ := net.ListenUDP(UDP_TYPE, NewUdpAddress("localhost", 0))
    //    localport := udpConn.LocalAddr().(*net.UDPAddr).Port;


    //    conn, err := net.DialUDP(UDP_TYPE, NewUdpAddress("localhost", localport), NewUdpAddress(targetHost, LISTEN_PORT))
    conn, err := net.DialUDP(UDP_TYPE, nil, NewUdpAddress(targetHost, 0000))
    E(err)
    //    fmt.Printf("%v \n\n", conn.LocalAddr().(*net.UDPAddr).Port)
    return conn
}

func doStuff(conn *net.UDPConn) {
    i := 0
    for {
        msg := strconv.Itoa(i)
        i++
        buf := []byte(msg)
        _, err := conn.Write(buf)
        if err != nil {
            fmt.Println(msg, err)
        }
        time.Sleep(time.Second * 1)
    }
}

func main() {

    readData := func(conn *net.UDPConn) {
        buf := make([]byte, 1024)

        for {
            n, addr, err := conn.ReadFromUDP(buf)
            fmt.Println("Received ", string(buf[0:n]), " from ", addr)

            if err != nil {
                fmt.Println("Error: ", err)
            }
        }
    }

    /* Now listen at selected port */
    serverConn, err := net.ListenUDP(UDP_TYPE, NewUdpAddress("localhost", 5))
    E(err)
    go readData(serverConn)
    defer serverConn.Close()


    //    time.Sleep(time.Minute * 1)


    clientConn := OpenTunnel("localhost")
    defer clientConn.Close()
    doStuff(clientConn)





    //    cert, err := tls.LoadX509KeyPair("server.pem", "server.key")
    //
    //    if err != nil {
    //        log.Fatal(err)
    //    }
    //
    //    config := tls.Config{Certificates: []tls.Certificate{cert}, ClientAuth: tls.RequireAnyClientCert
    //    }
    //    config.Rand = rand.Reader
    //
    //    ln, err := tls.Listen("tcp", ":6600", &config)
    //    if err != nil {
    //        log.Fatal(err)
    //    }
    //
    //    fmt.Println("Server(TLS) up and listening on port 6600")
    //
    //    for {
    //        conn, err := ln.Accept()
    //        if err != nil {
    //            log.Println(err)
    //            continue
    //        }
    //        go handleConnection(conn)
    //    }
}

func E(err error) {
    if err  != nil {
        panic(err)
    }
}
