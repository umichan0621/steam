package inventory

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"

	"github.com/umichan0621/steam/pkg/auth"
	"github.com/umichan0621/steam/pkg/common"
)

func AllItems(auth *auth.Core, language, appID, contextID, startAssetID string, count uint64, items *[]InventoryItem) (hasMore bool, lastAssetID string, err error) {
	params := url.Values{
		"l":     {language},
		"count": {strconv.FormatUint(count, 10)},
	}
	if startAssetID != "" {
		params.Set("start_assetid", startAssetID)
	}

	url := fmt.Sprintf("http://steamcommunity.com/inventory/%s/%s/%s?%s", auth.SteamID(), appID, contextID, params.Encode())
	res, err := auth.HttpClient().Get(url)
	if err != nil {
		return false, "", err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return false, "", err
	}
	resp := Response{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return false, "", err
	}
	fmt.Println(string(data))
	descriptions := make(map[string]int)
	for i, desc := range resp.Descriptions {
		key := fmt.Sprintf("%d_%d", desc.ClassID, desc.InstanceID)
		descriptions[key] = i
	}

	for _, asset := range resp.Assets {
		var desc *common.EconItemDesc
		key := fmt.Sprintf("%d_%d", asset.ClassID, asset.InstanceID)
		if d, ok := descriptions[key]; ok {
			desc = resp.Descriptions[d]
		}

		item := InventoryItem{
			AppID:      asset.AppID,
			ContextID:  asset.ContextID,
			AssetID:    asset.AssetID,
			ClassID:    asset.ClassID,
			InstanceID: asset.InstanceID,
			Amount:     asset.Amount,
			Desc:       desc,
		}
		*items = append(*items, item)
	}
	hasMore = resp.HasMore != 0
	return hasMore, resp.LastAssetID, nil
}
