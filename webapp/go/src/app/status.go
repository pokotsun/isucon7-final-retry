package main

import (
	"errors"
	"syscall"
	"time"
)

func GoFuncGetStatus() {
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		<-ticker.C

		for roomName, conns := range connMap {
			go func() {
				status, err := getStatus(roomName)
				if err != nil {
					logger.Infow("GoFuncGetStatus", "err", err)
					return
				}
				for _, ws := range conns {
					err = ws.WriteJSON(status)
					if err != nil && errors.Is(err, syscall.EPIPE) {
						logger.Error("GoFuncGetStatus WriteJSON", "err", err)
						ws.Close()
					}
				}
			}()
		}
	}
}
