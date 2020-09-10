package main

import (
	"time"
)

func GoFuncGetStatus() {
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		<-ticker.C

		for roomName, conns := range connMap {
			logger.Infow("GoFuncGetStatus", "conns", conns)
			go func() {
				status, err := getStatus(roomName)
				if err != nil {
					logger.Infow("GoFuncGetStatus", "err", err)
					return
				}
				for _, ws := range conns {
					logger.Infow("GoFuncGetStatus", "ws", ws, "ws.Conn", ws.Conn)
					err = ws.WriteJSON(status)
					if err != nil {
						logger.Info(err)
						return
					}
				}
			}()
		}
	}
}
