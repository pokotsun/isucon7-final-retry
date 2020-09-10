package main

var (
	M_ITEM_ARRAY []mItem
	M_ITEM_DICT  map[int]mItem
)

func InitMItem() error {
	var items []mItem
	err := db.Select(&items, "SELECT * FROM m_item")
	if err != nil {
		return err
	}
	M_ITEM_ARRAY = items
	M_ITEM_DICT = make(map[int]mItem)
	for _, v := range items {
		M_ITEM_DICT[v.ItemID] = v
	}
	return nil
}
