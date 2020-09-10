package main

import (
	"database/sql"
	"errors"
	"sync"
	"syscall"
	"time"
)

var GoRouineFuncMap = make(map[string]bool)
var mutex sync.Mutex

func GoFuncGetStatus(roomName string) {
	mutex.Lock()
	if _, ok := GoRouineFuncMap[roomName]; ok {
		mutex.Unlock()
		return
	}
	GoRouineFuncMap[roomName] = true
	mutex.Unlock()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		<-ticker.C
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
	}
}
