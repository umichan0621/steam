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

// Success while Code == 1
type BuyOrderResponse struct {
	Code    int    `json:"success"`
	Msg     string `json:"message"`
	OrderID uint64 `json:"buy_orderid,string"`
}

func (core *Core) CreateBuyOrder(auth *auth.Core, appID string, paymentPrice float64, quantity uint64, currencyID, hashName string) (*BuyOrderResponse, error) {
	reqUrl := fmt.Sprintf("%s/market/createbuyorder/", common.URI_STEAM_COMMUNITY)
	reqHeader := http.Header{}
	referer := strings.ReplaceAll(hashName, " ", "%20")
	referer = strings.ReplaceAll(referer, "#", "%23")
	referer = fmt.Sprintf("%s/market/listings/%s/%s", common.URI_STEAM_COMMUNITY, appID, referer)
	reqHeader.Add("Referer", referer)
	reqHeader.Add("Content-Type", "application/x-www-form-urlencoded")
	reqBody := url.Values{
		"appid":            {appID},
		"currency":         {currencyID},
		"market_hash_name": {hashName},
		"price_total":      {strconv.FormatUint(uint64(paymentPrice*100), 10)},
		"quantity":         {strconv.FormatUint(quantity, 10)},
		"sessionid":        {auth.SessionID()},
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
	response := &BuyOrderResponse{}
	err = json.Unmarshal(data, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (core *Core) CancelBuyOrder(auth *auth.Core, orderID uint64) error {
	reqUrl := fmt.Sprintf("%s/market/cancelbuyorder/", common.URI_STEAM_COMMUNITY)
	reqHeader := http.Header{}
	reqHeader.Add("Referer", fmt.Sprintf("%s/market", common.URI_STEAM_COMMUNITY))
	reqHeader.Add("Content-Type", "application/x-www-form-urlencoded")
	reqBody := url.Values{
		"sessionid":   {auth.SessionID()},
		"buy_orderid": {strconv.FormatUint(orderID, 10)},
	}

	req, err := http.NewRequest(http.MethodPost, reqUrl, strings.NewReader(reqBody.Encode()))
	if err != nil {
		return err
	}

	req.Header = reqHeader

	res, err := auth.HttpClient().Do(req)
	if res != nil {
		res.Body.Close()
	}
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("cannot cancel %d: %d", orderID, res.StatusCode)
	}
	return nil
}
