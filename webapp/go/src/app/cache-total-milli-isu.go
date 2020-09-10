package main

import (
	"encoding/json"
	"fmt"
)

func getKeyCurrentStatus(roomName string) string {
	return fmt.Sprintf("CUR_TOTAL_MILISU-rname-%s", roomName)
}

func setCurrentStatusToCache(roomName string, item CurrentStatus) {
	key := getKeyCurrentStatus(roomName)
	v, err := json.Marshal(item)
	if err != nil {
		logger.Errorf("Json Marshal Err On Set CurTotalMillIsu: %s", err)
	}
	err = cacheClient.SingleSet(key, v)
	if err != nil {
		logger.Errorf("Failed to Set CurTotalMillIsu to cache Err: %s", err)
	}
}

func getCurrentStatusFromCache(roomName string) (CurrentStatus, error) {
	key := getKeyCurrentStatus(roomName)
	bytes, err := cacheClient.SingleGet(key)
	if err != nil {
		logger.Errorf("Failed to Get Cache CurTotalMillIsu: %s", err)
	}
	item := CurrentStatus{}
	err = json.Unmarshal(bytes, &item)
	return item, err
}
