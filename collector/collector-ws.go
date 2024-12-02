package collector

import (
	"errors"
	"fmt"
	"metrics-xray/calc"
	"metrics-xray/parser"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	WS_IDLE_TIMEOUT int64 = 120
	WS_IDLE_PING    int64 = 30
)

var upgrader = websocket.Upgrader{
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ClientSession struct {
	conn *websocket.Conn
	sax  *parser.LogSax

	instanceId string
	appId      string
	ip         string
	port       int32
	//status
	recvLastTime int64

	messageQueue chan []byte
	mu           sync.Mutex
	//had been closed
	closed bool
}

func NewClient(app string, instanceId string, conn *websocket.Conn) *ClientSession {
	return &ClientSession{
		conn:         conn,
		appId:        app,
		instanceId:   instanceId,
		recvLastTime: 0,
		messageQueue: make(chan []byte, 100),
		sax:          parser.NewLogSax(calc.G_queueClient),
		closed:       false,
	}
}

func (pThis *ClientSession) ReadPump() {
	defer func() {

	}()
	pThis.readRun()
}
func (pThis *ClientSession) readRun() {
	defer func() {
		G_manager.removeClient(pThis.appId, pThis.instanceId)
		pThis.Close(errors.New("recv err"))
	}()
	pThis.recvLastTime = time.Now().Unix()

	for {
		mt, message, err := pThis.conn.ReadMessage()
		if err != nil {
			logger.Info("read:", err)
			break
		}
		pThis.recvLastTime = time.Now().Unix()

		if mt == websocket.TextMessage || mt == websocket.PingMessage {
			pThis.sax.LoadData(string(message))
		}
	}
}

func (pThis *ClientSession) Send(returnMessage []byte) {
	defer func() {

	}()
	pThis.mu.Lock()
	err := pThis.unexceptionSend(returnMessage)
	pThis.mu.Unlock()
	if err != nil {
		//logger.Errorf("client.conn.WriteMessage error %s", err.Error())
		pThis.Close(errors.New("send exption"))
	}
}

func (pThis *ClientSession) unexceptionSend(returnMessage []byte) (ret error) {
	ret = errors.New("send unknow exception")
	defer func() {

	}()
	ret = pThis.conn.WriteMessage(websocket.TextMessage, returnMessage)
	return ret
}

func (pThis *ClientSession) SendHeadbeat() error {
	defer func() {

	}()
	pThis.mu.Lock()
	err := pThis.unexceptionSend([]byte(""))
	pThis.mu.Unlock()
	if err != nil {
		//logger.Errorf("client.conn.WriteMessage error %s", err.Error())
		pThis.Close(errors.New("sendHeadbeat err"))
		return err
	}
	return nil
}

func (pThis *ClientSession) Close(err error) {
	defer func() {

	}()
	if pThis.conn != nil && !pThis.closed {
		pThis.mu.Lock()
		pThis.closed = true
		pThis.mu.Unlock()
		pThis.conn.Close()
		if err != nil {
			logger.Info("on close, the err:" + err.Error())
		} else {
			logger.Info("on close succ")
		}
	}
}

func onTick1s(pThis *ClientSession) {
	//
	t := time.Now().Unix()
	if pThis.recvLastTime == 0 {
		pThis.recvLastTime = t
	} else {
		if t-pThis.recvLastTime > WS_IDLE_TIMEOUT {
			pThis.Close(errors.New("the recv idle timeout"))
		} else if t-pThis.recvLastTime > WS_IDLE_PING {
			pThis.SendHeadbeat()
		}
	}

}

func wsHandle(c *gin.Context) {
	fmt.Println("wsHandle onConnect, ip=" + c.RemoteIP())
	logger.Debug("wsHandle onConnect, ip=", c.RemoteIP())
	//upgrade the http connect to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error("upgrade:", err)
		return
	}
	//defer conn.Close()
	instanceId := c.Query("instance") //  get app
	appId := c.Query("app")           //  get app
	//sign := c.Query("token") //  get app
	if instanceId == "" || appId == "" {
		logger.Error("instanceId is empty:")
		_ = conn.Close()
		return
	}
	app := G_manager.getApp(appId)
	if app == nil {
		logger.Error("instanceId is empty:")
		_ = conn.Close()
		return
	}

	client := NewClient(appId, instanceId, conn)
	app.addClient(client)
	go client.ReadPump()
}
