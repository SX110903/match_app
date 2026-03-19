package config

type ShopItem struct {
	ItemType  string `json:"item_type"`
	ItemValue int    `json:"item_value"`
	Cost      int    `json:"cost"`
	Name      string `json:"name"`
	Benefits  string `json:"benefits"`
}

var VIPItems = []ShopItem{
	{ItemType: "vip_upgrade", ItemValue: 1, Cost: 500, Name: "VIP Bronce", Benefits: "+10% más candidatos"},
	{ItemType: "vip_upgrade", ItemValue: 2, Cost: 1500, Name: "VIP Plata", Benefits: "+25% + super likes"},
	{ItemType: "vip_upgrade", ItemValue: 3, Cost: 3000, Name: "VIP Oro", Benefits: "+50% + boost diario"},
	{ItemType: "vip_upgrade", ItemValue: 4, Cost: 6000, Name: "VIP Platino", Benefits: "+75% + visible primero"},
	{ItemType: "vip_upgrade", ItemValue: 5, Cost: 10000, Name: "VIP Diamante", Benefits: "Sin límites + badge"},
}
