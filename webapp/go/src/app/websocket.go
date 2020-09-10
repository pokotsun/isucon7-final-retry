package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocket struct {
	ID       int
	RoomName string

	Conn *websocket.Conn
	Mux  sync.Mutex
}

var (
	ConnMap = make(map[string]map[int]*WebSocket)
)

func BuildWebSocket(roomName string, conn *websocket.Conn) *WebSocket {
	return &WebSocket{
		ID:       autoIncrement.FetchID(),
		Conn:     conn,
		RoomName: roomName,
	}
}

func AppendConn(ws *WebSocket) {
	if _, ok := ConnMap[ws.RoomName]; !ok {
		ConnMap[ws.RoomName] = make(map[int]*WebSocket)
	}

	ConnMap[ws.RoomName][ws.ID] = ws
}

func (ws *WebSocket) WriteJSON(v interface{}) error {
	ws.Mux.Lock()
	defer ws.Mux.Unlock()
	return ws.Conn.WriteJSON(v)
}
