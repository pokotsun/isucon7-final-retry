package main

import "github.com/gorilla/websocket"

type WebSocket struct {
	ID       int
	RoomName string

	Conn *websocket.Conn
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
	return ws.Conn.WriteJSON(v)
}
