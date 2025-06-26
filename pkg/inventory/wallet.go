package inventory

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/umichan0621/steam/pkg/auth"
	"github.com/umichan0621/steam/pkg/common"
)

type WalletInfo struct {
	WalletCurrency       int32   `json:"wallet_currency"`
	WalletCountry        string  `json:"wallet_country"`
	WalletBalance        float32 `json:"wallet_balance,string"`
	WalletDelayedBalance float32 `json:"wallet_delayed_balance,string"`
	Success              int32   `json:"success"`
}

func WalletBalance(auth *auth.Core) (*WalletInfo, error) {
	reqUrl := fmt.Sprintf("%s/market/", common.URI_STEAM_COMMUNITY)
	res, err := auth.HttpClient().Get(reqUrl)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fail to get wallet balance, code: %d", res.StatusCode)
	}
	datax, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	data := string(datax)
	index := strings.Index(data, "g_rgWalletInfo")
	info := &WalletInfo{}
	info.Success = 0
	if index >= 0 {
		data = data[index:]
		start := strings.Index(data, "{")
		end := strings.Index(data, "}")
		if start >= 0 && end >= 0 {
			data = data[start : end+1]
			err := json.Unmarshal([]byte(data), info)
			if err != nil {
				return info, fmt.Errorf("fail to parse json: %s, data: %s", err.Error(), data)
			}
		}
	}
	info.WalletBalance /= 100.0
	info.WalletDelayedBalance /= 100.0
	return info, nil
}
