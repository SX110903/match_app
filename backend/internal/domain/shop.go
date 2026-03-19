package domain

import "time"

type ShopTransaction struct {
	ID        string    `db:"id"         json:"id"`
	UserID    string    `db:"user_id"    json:"user_id"`
	ItemType  string    `db:"item_type"  json:"item_type"`
	ItemValue int       `db:"item_value" json:"item_value"`
	Cost      int       `db:"cost"       json:"cost"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
