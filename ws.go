package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)


var upgrader = websocket.Upgrader{
	HandshakeTimeout: 0,
	ReadBufferSize:   0,
	WriteBufferSize:  0,
	WriteBufferPool:  nil,
	Subprotocols:     []string{},
	Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
	},
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	WS_PING *websocket.PreparedMessage
)

func init() {
	WS_PING, _ = websocket.NewPreparedMessage(1, []byte(`{"e": "P"}`))
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	WSMgr.Add(SnowNode.Generate().String(), conn)
}

type WSMgrT struct {
	Lock        *sync.RWMutex
	Connections map[string]*websocket.Conn
}

var WSMgr = &WSMgrT{
	Lock:        &sync.RWMutex{},
	Connections: map[string]*websocket.Conn{},
}

func (mgr *WSMgrT) Add(id string, conn *websocket.Conn) {
	mgr.Lock.Lock()
	defer mgr.Lock.Unlock()
	mgr.Connections[id] = conn
	b, _ := json.Marshal(GetStats())
	conn.WriteMessage(websocket.TextMessage, b)
}

func (mgr *WSMgrT) Remove(id string) {
	mgr.Lock.Lock()
	defer mgr.Lock.Unlock()
	conn := mgr.Connections[id]
	conn.WriteControl(websocket.CloseMessage, []byte{}, time.Time{})
	conn.Close()
	delete(mgr.Connections, id)
}

func (mgr *WSMgrT) SendStats() {
	enc, _ := json.Marshal(GetStats())
	prepared, _ := websocket.NewPreparedMessage(websocket.TextMessage, enc)

	mgr.Lock.Lock()
	defer mgr.Lock.Unlock()

	for _, c := range mgr.Connections {
		c.WritePreparedMessage(prepared)
	}
}

func (mgr *WSMgrT) PingLoop() {
	for {
		time.Sleep(30 * time.Second)
		mgr.Ping()
	}
}

func (mgr *WSMgrT) Ping() {
	mgr.Lock.Lock()
	defer mgr.Lock.Unlock()
	fmt.Println("PING")

	wg := &sync.WaitGroup{}

	for id, c := range mgr.Connections {
		wg.Add(1)
		go func (id string, c *websocket.Conn) {	
			c.WritePreparedMessage(WS_PING)		
			c.SetReadDeadline(time.Now().Add(5 * time.Second))
			_, p, err := c.ReadMessage()

			if err != nil || len(p) == 0 || p[0] != 'P' {
				go mgr.Remove(id)
			} else {
				c.SetReadDeadline(time.Time{})
			}
			wg.Done()
		}(id, c)
	}
	wg.Wait()
}

func init() {
	go WSMgr.PingLoop()
}
