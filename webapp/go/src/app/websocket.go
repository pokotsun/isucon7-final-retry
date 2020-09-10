package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WS struct {
	ID       int
	RoomName string
	Conn     *websocket.Conn
	Mux      sync.Mutex
}

var (
	connMap = make(map[string]map[int]WS)
)

func AddConn(ws WS) {
	if _, ok := connMap[ws.RoomName]; !ok {
		connMap[ws.RoomName] = make(map[int]WS)
	}
	connMap[ws.RoomName][ws.ID] = ws
}

func (ws WS) WriteJSON(v interface{}) error {
	ws.Mux.Lock()
	defer ws.Mux.Unlock()
	return ws.Conn.WriteJSON(v)
}

func (ws WS) Close() {
	ws.Conn.Close()
	delete(connMap[ws.RoomName], ws.ID)
}
