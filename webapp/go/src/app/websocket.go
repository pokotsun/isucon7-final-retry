package main

import (
	"errors"
	"sync"
	"syscall"

	"github.com/gorilla/websocket"
)

type WS struct {
	ID       int
	RoomName string
	Conn     *websocket.Conn
	Mux      sync.Mutex

	CloseChannel chan int8
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
	err := ws.Conn.WriteJSON(v)

	if err != nil && errors.Is(err, syscall.EPIPE) {
		logger.Errorw("WriteJSON ", "err", err)
		logger.Info("WriteJSON and Close")
		ws.Close()

		return nil
	}

	return err
}

func (ws WS) Close() {
	ws.Conn.Close()
	delete(connMap[ws.RoomName], ws.ID)
	ws.CloseChannel <- 1
}
