package main

import (
	"database/sql"
	"errors"
	"sync"
	"syscall"
	"time"
)

type GoRoutine struct {
	RoomName string
	Channel  chan int
}

var GoRouineFuncMap = make(map[string]GoRoutine)
var mutex sync.Mutex

func GoFuncGetStatus(roomName string) {
	mutex.Lock()
	if _, ok := GoRouineFuncMap[roomName]; ok {
		mutex.Unlock()
		return
	}
	r := GoRoutine{
		RoomName: roomName,
		Channel:  make(chan int),
	}

	GoRouineFuncMap[roomName] = r
	mutex.Unlock()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			status, err := getStatus(roomName)
			if err != nil {
				if err == sql.ErrNoRows {
					mutex.Lock()
					delete(GoRouineFuncMap, roomName)
					mutex.Unlock()
					return
				}
				logger.Infow("GoFuncGetStatus", "err", err)
			}
			for _, ws := range ConnMap[roomName] {
				err = ws.WriteJSON(status)
				if err != nil && errors.Is(err, syscall.EPIPE) {
					logger.Infow("GoFuncGetStatus", "err", err)
					delete(ConnMap[roomName], ws.ID)
				}
			}

		case <-r.Channel:
			return
		}
	}
}
