package udp
import (
    "net"
    "fmt"
    "time"
    "strconv"
    "runtime"
)

const UDP_TYPE = "udp4"

func NewUdpAddress(ip string, port int) *net.UDPAddr {
    return &net.UDPAddr{IP: net.ParseIP(ip), Port: port}
}


//        readToChannel := func(conn *net.UDPConn, channel chan []byte){
//            //max size of udp packet
//            packet := make([]byte, 2**16)
//            for {
//                conn.ReadFromUDP(packet)
//                channel <- packet
//            }
//        }


func Listen(port int, newConnection func(*net.UDPConn, *net.UDPAddr)) {
    handOffNewConn := func(remoteAddr *net.UDPAddr) {
        persistentConn, err := net.ListenUDP(UDP_TYPE, nil)
        E(err)
        fmt.Printf("new server listen on %v \n\n", persistentConn.LocalAddr())
        fmt.Printf("notify of new connection \n")
        go newConnection(persistentConn, remoteAddr)
        runtime.Gosched()
        _, err = persistentConn.WriteToUDP([]byte("yo"), remoteAddr)
        E(err)
        fmt.Printf("server return address to %v \n\n", remoteAddr)
    }
    listenConn, err := net.ListenUDP(UDP_TYPE, NewUdpAddress("localhost", port))
    fmt.Printf("server listening on %v \n\n", listenConn.LocalAddr())
    E(err)
    readBytes := make([]byte, 1)
    for {
        _, returnAddr, err := listenConn.ReadFromUDP(readBytes)
        E(err)

        fmt.Printf("server received incoming from %v\n\n", returnAddr)
        go handOffNewConn(returnAddr)
    }
}

func OpenConnection(targetHost string, targetPort int) (*net.UDPConn, *net.UDPAddr) {
    conn, err := net.ListenUDP(UDP_TYPE, nil)
    E(err)
    conn.WriteToUDP([]byte("yo"), NewUdpAddress(targetHost, targetPort))
    _,add,err := conn.ReadFromUDP(nil)
    E(err)
    fmt.Printf("established nee connection to %v \n\n", add)
    return conn, add
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
