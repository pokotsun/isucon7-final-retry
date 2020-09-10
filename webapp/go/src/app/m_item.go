package main

import (
	"fmt"
	"math/big"
)

var (
	// M_ITEM_ARRAY []mItem
	M_ITEM_DICT map[int]mItem

	POWER_DICT = make(map[string]*big.Int)
)

func InitMItem() error {
	var items []mItem
	err := db.Select(&items, "SELECT * FROM m_item")
	if err != nil {
		return err
	}
	// M_ITEM_ARRAY = items
	M_ITEM_DICT = make(map[int]mItem)
	for _, v := range items {
		M_ITEM_DICT[v.ItemID] = v
	}
	return nil
}

func FetchMItem(itemID int) mItem {
	item, _ := M_ITEM_DICT[itemID]
	return item
}

func (item *mItem) BuildCacheKeyByCount(count int) string {
	return fmt.Sprintf("%d-%d", item.ItemID, count)
}
