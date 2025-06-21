package inventory

import "github.com/umichan0621/steam/pkg/common"

type Asset struct {
	AppID      uint32 `json:"appid"`
	ContextID  uint64 `json:"contextid,string"`
	AssetID    string `json:"assetid"`
	ClassID    uint64 `json:"classid,string"`
	InstanceID uint64 `json:"instanceid,string"`
	Amount     uint64 `json:"amount,string"`
}

type Response struct {
	Assets              []Asset                `json:"assets"`
	Descriptions        []*common.EconItemDesc `json:"descriptions"`
	Success             int                    `json:"success"`
	HasMore             int                    `json:"more_items"`
	LastAssetID         string                 `json:"last_assetid"`
	TotalInventoryCount int                    `json:"total_inventory_count"`
	ErrorMsg            string                 `json:"error"`
}

type InventoryItem struct {
	AppID      uint32               `json:"appid"`
	ContextID  uint64               `json:"contextid"`
	AssetID    string               `json:"id,string,omitempty"`
	ClassID    uint64               `json:"classid,string,omitempty"`
	InstanceID uint64               `json:"instanceid,string,omitempty"`
	Amount     uint64               `json:"amount,string"`
	Desc       *common.EconItemDesc `json:"-"` /* May be nil  */
}
