package udp
import (
    "testing"
    "net"
    "fmt"
//    "time"
    "runtime"
    "time"
    "sync"
)



func Test_initial(t *testing.T) {

    c, err := net.ListenUDP("udp", nil)
    E(err)



    fmt.Printf("%v \n\n", c.LocalAddr())


}

func Test_listening(t *testing.T) {
    listenOnNewConnections := func(conn *net.UDPConn, remoteAddr *net.UDPAddr){
        fmt.Printf("tracing remoteAddr %v \n\n", remoteAddr)
        fmt.Printf("tracing localAddr %v \n\n", conn.LocalAddr())
        byteslice := make([]byte, 1)
        _, add, err := conn.ReadFromUDP(byteslice)
        E(err)
        fmt.Printf("%v %v %v\n\n", add, remoteAddr, byteslice)
    }

    //start the "server" listening
    targetPort := 10001
    go Listen(targetPort, listenOnNewConnections)
    runtime.Gosched()

    readresponseAddr := func(conn *net.UDPConn){
        byteslice := make([]byte, 1)
        _,add,err := conn.ReadFromUDP(byteslice)
        E(err)

        fmt.Printf("recieved new target address from server %v \n\n", add)
        conn.WriteToUDP([]byte("yo"), add)
    }
    conn, err := net.ListenUDP(UDP_TYPE, nil)
    fmt.Printf("client listening on %v \n\n", conn.LocalAddr())
    E(err)
    go readresponseAddr(conn)
    runtime.Gosched()
    fmt.Printf("calling from %v \n\n", conn.LocalAddr())
    conn.WriteToUDP([]byte("yo"), NewUdpAddress("localhost", targetPort))

    time.Sleep(time.Second)
}


func Test_listening3(t *testing.T) {
    listenOnNewConnections := func(conn *net.UDPConn, remoteAddr *net.UDPAddr){
        for {
            b := make([]byte, 2^16)
            conn.ReadFromUDP(b)
            fmt.Printf(string(b))
        }
    }

    //start the "server" listening
    targetPort := 10001
    go Listen(targetPort, listenOnNewConnections)
    runtime.Gosched()
    tunnel := OpenConnection("localhost", targetPort)

    var wg sync.WaitGroup

    wg.Add(1)

    go func() {
        for i := 0; i < 10; i++ {
            tunnel.Send([]byte("yo\n"))
        }
        wg.Done()
    }()

    OpenConnection("localhost", targetPort)
    OpenConnection("localhost", targetPort)

    wg.Wait()
}

func Test_listening2(t *testing.T) {
    con2, err := net.ListenUDP(UDP_TYPE, NewUdpAddress("localhost", 10003))
    E(err)
    fmt.Printf("con2 local %v \n\n", con2.LocalAddr())
    fmt.Printf("con2 remot %v \n\n", con2.RemoteAddr())


    con2, err = net.ListenUDP(UDP_TYPE, NewUdpAddress("localhost", 10004))
    E(err)
    fmt.Printf("con2 local %v \n\n", con2.LocalAddr())
    fmt.Printf("con2 remot %v \n\n", con2.RemoteAddr())



    conn, err := net.DialUDP(UDP_TYPE, NewUdpAddress("localhost", 10009), NewUdpAddress("localhost", 10006))
    E(err)
    fmt.Printf("conn local %v \n\n", conn.LocalAddr())
    fmt.Printf("conn remot %v \n\n", conn.RemoteAddr())


    time.Sleep(time.Minute)
}
