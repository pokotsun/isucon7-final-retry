package main

import (
	"time"
)

func GoFuncGetStatus() {
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		<-ticker.C

		for roomName, conns := range connMap {
			if len(connMap) == 0 {
				continue
			}
			go func() {
				status, err := getStatus(roomName)
				if err != nil {
					logger.Infow("GoFuncGetStatus", "err", err)
				}
				for _, ws := range conns {
					err = ws.WriteJSON(status)
					if err != nil {
						logger.Infow("GoFuncGetStatus", "err", err)
					}
				}
			}()
		}
	}
}
