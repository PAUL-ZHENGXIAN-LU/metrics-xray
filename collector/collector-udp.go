package collector

import (
	"fmt"
	"metrics-xray/calc"
	"metrics-xray/parser"
	"sync"
	"time"

	"net"
)

const (
	UDP_DEFAULT_PORT int = 6090
)

var (
	G_collectorUdpPort int = UDP_DEFAULT_PORT
)

var udpService = UdpService{}

type UdpService struct {
	conn *net.UDPConn
	sax  *parser.LogSax

	ip   string
	port int32
	//status
	recvLastTime int64

	messageQueue chan []byte
	mu           sync.Mutex
	//had been closed
	closed bool
}

func NewUdpService(conn *net.UDPConn) *UdpService {
	return &UdpService{
		conn:         conn,
		recvLastTime: 0,
		messageQueue: make(chan []byte, 100),
		sax:          parser.NewLogSax(calc.G_queueClient),
		closed:       false,
	}
}

func (pThis *UdpService) ReadPump() {
	defer func() {
		fmt.Println("the udp is close now...")
		defer pThis.conn.Close()

	}()
	pThis.readRun()
}
func (pThis *UdpService) readRun() {
	defer func() {
		//G_manager.removeClient(pThis.appId, pThis.instanceId)
		//pThis.Close(errors.New("recv err"))
	}()
	pThis.recvLastTime = time.Now().Unix()

	for {
		var data [1024]byte
		n, addr, err := pThis.conn.ReadFromUDP(data[:])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("received message: %s from %s\n", data[:n], addr)
		pThis.recvLastTime = time.Now().Unix()

		//if mt == websocket.TextMessage || mt == websocket.PingMessage {
		pThis.sax.LoadData(string(data[:n]))
		//get the ip, and get the app
		//}

	}

}

func udpInit() {
	//init

	//addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:"+strconv.Itoa(G_collectorUdpPort))
	addr, err := net.ResolveUDPAddr("udp", ":6090")
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := NewUdpService(conn)

	go client.ReadPump()

}
