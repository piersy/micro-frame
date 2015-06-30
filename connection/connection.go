package connection
import (
	"github.com/piersy/micro-frame/mframe"
	"github.com/piersy/micro-frame/udp"
)


type Connection interface {
	Health() int
	Send([]byte)
}


func OpenConnection(targetHost string, targetPort int) Connection {
	//Listen on a port
	data := udp.OpenTunnel(targetHost , targetPort)
	data.Write([]byte("data"))

	heartBeat := udp.OpenTunnel(targetHost , targetPort)
	heartBeat.Write([]byte("heartbeat"))

}