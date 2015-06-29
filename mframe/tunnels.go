package mframe
import (
    "net"
    "crypto/rsa"
    "log"
//    "crypto/tls"
//    "crypto/rand"
//    "fmt"
//    "time"
//    "crypto/x509"
//    "crypto/x509/pkix"
//    "math/big"
////    "crypto/ecdsa"
//    "encoding/pem"
    "fmt"
    "strconv"
    "time"
)

type Tunnel interface{
    Send(message interface{}) error
}

const UDP_TYPE = "udp"
const SOURCE_PORT = 12345
const LISTEN_PORT = 12346

type DefaultTunnel struct {
    ip string
    publicKey rsa.PublicKey
    privateKey rsa.PrivateKey
    hostKey rsa.PublicKey
}

func NewUdpAddress(ip string , port int) *net.UDPAddr{
    return &net.UDPAddr{IP: net.ParseIP(ip), Port: port}
}

func e(err error) {
    if err  != nil {
        panic(err)
    }
}

func OpenTunnel(targetHost string) *net.UDPConn{
//    addrs, err := net.LookupHost(host)
//    if err != nil {
//        panic(err)
//    }

//    udpConn, _ := net.ListenUDP(UDP_TYPE, NewUdpAddress("localhost", 0))
//    localport := udpConn.LocalAddr().(*net.UDPAddr).Port;




//    conn, err := net.DialUDP(UDP_TYPE, NewUdpAddress("localhost", localport), NewUdpAddress(targetHost, LISTEN_PORT))
    conn, err := net.DialUDP(UDP_TYPE, nil, NewUdpAddress(targetHost, LISTEN_PORT))
    e(err)
//    fmt.Printf("%v \n\n", conn.LocalAddr().(*net.UDPAddr).Port)
    return conn
}

func doStuff(conn *net.UDPConn) {
    i := 0
    for {
        msg := strconv.Itoa(i)
        i++
        buf := []byte(msg)
        _,err := conn.Write(buf)
        if err != nil {
            fmt.Println(msg, err)
        }
        time.Sleep(time.Second * 1)
    }
}


func handleConnection(c net.Conn) {

    log.Printf("Client(TLS) %v connected via secure channel.", c.RemoteAddr())

    // stuff to do... like read data from client, process it, write back to client
    // see what you can do with (c net.Conn) at
    // http://golang.org/pkg/net/#Conn

    // buffer := make([]byte, 4096)

    //for {
    //		n, err := c.Read(buffer)
    //		if err != nil || n == 0 {
    //			c.Close()
    //			break
    //		}
    //		n, err = c.Write(buffer[0:n])
    //		if err != nil {
    //			c.Close()
    //			break
    //		}
    //	}
    log.Printf("Connection from %v closed.", c.RemoteAddr())
}


//func makeKey() {
//
//    serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
//    serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
//    if err != nil {
//        log.Fatalf("failed to generate serial number: %s", err)
//    }
//
//    template := x509.Certificate{
//        SerialNumber: serialNumber,
//        Subject: pkix.Name{
//            Organization: []string{"Acme Co"},
//        },
//        NotBefore: time.Now(),
//        NotAfter:  time.Now() + time.Hour * 24,
//
//        KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
//        ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
//        BasicConstraintsValid: true,
//    }
//
//    //set the host for this cert
//    template.DNSNames = append(template.DNSNames, "localhost")
//    //Set the cert to be self signed
//    template.IsCA = true
//    template.KeyUsage |= x509.KeyUsageCertSign
//
//    priv, err := rsa.GenerateKey(rand.Reader, 4096)
//    if err != nil {
//        panic(err)
//    }
//
//    derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, priv.PublicKey, priv)
//    if err != nil {
//        log.Fatalf("Failed to create certificate: %s", err)
//    }
//    pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
//
//}


func main() {

    readData := func (conn *net.UDPConn){
        buf := make([]byte, 1024)

        for {
            n,addr,err := conn.ReadFromUDP(buf)
            fmt.Println("Received ",string(buf[0:n]), " from ",addr)

            if err != nil {
                fmt.Println("Error: ",err)
            }
        }
    }

    /* Now listen at selected port */
    serverConn, err := net.ListenUDP(UDP_TYPE, NewUdpAddress("localhost", LISTEN_PORT))
    e(err)
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