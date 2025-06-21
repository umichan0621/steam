package market

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/umichan0621/steam/pkg/auth"
	"github.com/umichan0621/steam/pkg/common"
)

type MarketSellResponse struct {
	Success                    bool   `json:"success"`
	RequiresConfirmation       uint32 `json:"requires_confirmation"`
	MobileConfirmationRequired bool   `json:"needs_mobile_confirmation"`
	EmailConfirmationRequired  bool   `json:"needs_email_confirmation"`
	EmailDomain                string `json:"email_domain"`
}

func (core *Core) CreateSellOrder(auth *auth.Core, appID, contextID, assetID string, amount, paymentPrice uint64) (*MarketSellResponse, error) {
	reqUrl := fmt.Sprintf("%s/market/sellitem/", common.URI_STEAM_COMMUNITY)
	referUrl := fmt.Sprintf("%s/profiles/%s/inventory/", common.URI_STEAM_COMMUNITY, auth.SteamID())
	reqHeader := http.Header{}
	reqHeader.Add("Content-Type", "application/x-www-form-urlencoded")
	reqHeader.Add("Referer", referUrl)
	reqBody := url.Values{
		"appid":     {appID},
		"contextid": {contextID},
		"assetid":   {assetID},
		"sessionid": {auth.SessionID()},
		"amount":    {strconv.FormatUint(amount, 10)},
		"price":     {strconv.FormatUint(paymentPrice, 10)},
	}

	req, err := http.NewRequest(http.MethodPost, reqUrl, strings.NewReader(reqBody.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header = reqHeader

	res, err := auth.HttpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http error: %d, %s", res.StatusCode, string(data))
	}

	response := &MarketSellResponse{}
	err = json.Unmarshal(data, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (core *Core) CancelSellOrder() {}
