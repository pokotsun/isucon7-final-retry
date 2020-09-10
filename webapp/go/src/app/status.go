package main

import (
	"errors"
	"syscall"
	"time"
)

func GoFuncGetStatus(roomName string) {
	// ここもしかしたらバグ生むかも
	if len(ConnMap[roomName]) > 1 {
		return
	}
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		<-ticker.C
		status, err := getStatus(roomName)
		if err != nil {
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
