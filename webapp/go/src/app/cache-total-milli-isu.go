package main

import (
	"encoding/json"
	"fmt"
)

func getKeyCurrentTotalMillIsu(roomName string) string {
	return fmt.Sprintf("CUR_TOTAL_MILISU-rname-%s", roomName)
}

func setCurrentTotalMillIsu(roomName string, item CurrentTotalMillIsu) {
	key := getKeyCurrentTotalMillIsu(roomName)
	v, err := json.Marshal(item)
	if err != nil {
		logger.Errorf("Json Marshal Err On Set CurTotalMillIsu: %s", err)
	}
	err = cacheClient.SingleSet(key, v)
	if err != nil {
		logger.Errorf("Failed to Set CurTotalMillIsu to cache Err: %s", err)
	}
}

func getCurrentTotalMillIsu(roomName string) (CurrentTotalMillIsu, error) {
	key := getKeyCurrentTotalMillIsu(roomName)
	bytes, err := cacheClient.SingleGet(key)
	if err != nil {
		logger.Errorf("Failed to Get Cache CurTotalMillIsu: %s", err)
	}
	item := CurrentTotalMillIsu{}
	err = json.Unmarshal(bytes, &item)
	return item, err
}
