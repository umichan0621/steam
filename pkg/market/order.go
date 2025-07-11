package market

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/umichan0621/steam/pkg/auth"
	"github.com/umichan0621/steam/pkg/common"
	"golang.org/x/net/html"
)

type SteamOrder struct {
	AssetID        string `json:"id"`
	MarketName     string `json:"market_name"`
	MarketHashName string `json:"market_hash_name"`
	Commodity      uint64 `json:"commodity"`
	Seq            uint64
	Price          float64
	DateString     string
}

func HistoryOrder(auth *auth.Core, language, appID, contextID string, start, count uint64) ([]*SteamOrder, error) {
	params := url.Values{
		"l":     {language},
		"start": {strconv.FormatUint(start, 10)},
		"count": {strconv.FormatUint(count, 10)},
	}
	reqUrl := fmt.Sprintf("%s/market/myhistory?%s", common.URI_STEAM_COMMUNITY, params.Encode())
	res, err := auth.HttpClient().Get(reqUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fail to get market history, code: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	jsonData := string(data)
	success := gjson.Get(jsonData, "success").Bool()
	totalCount := gjson.Get(jsonData, "total_count").Uint()
	if !success {
		return nil, fmt.Errorf("fail to get market history")
	}
	assets := gjson.Get(jsonData, "assets")
	assetsList := assets.Get(appID).Get(contextID)

	ordersList := []*SteamOrder{}
	assetsList.ForEach(func(_, val gjson.Result) bool {
		order := &SteamOrder{}
		err := json.Unmarshal([]byte(val.String()), order)
		if err != nil {
			return false
		}
		ordersList = append(ordersList, order)
		return true
	})

	htmlData := gjson.Get(jsonData, "results_html").String()
	hoversData := gjson.Get(jsonData, "hovers").String()
	htmlNode, _ := html.Parse(strings.NewReader(htmlData))

	historyRow2AssetID := generateHistoryRow2AssetIDMap(hoversData)
	historyRow2Price := map[string]float64{}
	historyRow2DateString := map[string]string{}
	assetID2Price := map[string]float64{}
	assetID2DateString := map[string]string{}
	generateHistoryRow2PriceMap(htmlNode, &historyRow2Price, "empty")
	generateHistoryRow2DateStringMap(htmlNode, &historyRow2DateString, "empty")
	for historyRow, assetID := range historyRow2AssetID {
		price, ok := historyRow2Price[historyRow]
		if ok {
			assetID2Price[assetID] = price
		}
		dateString, ok := historyRow2DateString[historyRow]
		if ok {
			assetID2DateString[assetID] = dateString
		}
	}
	for i, order := range ordersList {
		assetID := order.AssetID
		price, ok := assetID2Price[assetID]
		if ok {
			order.Price = price
		}
		dateString, ok := assetID2DateString[assetID]
		if ok {
			order.DateString = dateString
		}
		order.Seq = totalCount - start - uint64(i)
	}
	return ordersList, nil
}

func generateHistoryRow2AssetIDMap(hovers string) map[string]string {
	res := map[string]string{}
	hoversList := strings.Split(hovers, ";")
	for _, tmp := range hoversList {
		elements := strings.Split(tmp, ",")
		if len(elements) <= 4 {
			continue
		}
		historyRow := elements[1]
		asstIDStr := elements[4]
		if strings.Contains(historyRow, "_image") {
			continue
		}
		historyRow = strings.ReplaceAll(historyRow, "_name", "")
		historyRow = strings.ReplaceAll(historyRow, "'", "")
		historyRow = strings.ReplaceAll(historyRow, " ", "")
		asstIDStr = strings.ReplaceAll(asstIDStr, "'", "")
		asstIDStr = strings.ReplaceAll(asstIDStr, " ", "")
		res[historyRow] = asstIDStr
	}
	return res
}

func generateHistoryRow2PriceMap(n *html.Node, priceMap *map[string]float64, historyRow string) {
	if n.Type == html.ElementNode {
		class := ""
		id := ""
		for _, attr := range n.Attr {
			switch attr.Key {
			case "class":
				class = attr.Val
			case "id":
				id = attr.Val
			}
			switch class {
			case "market_listing_row market_recent_listing_row":
				historyRow = id
			case "market_listing_price":
				priceStr := strings.ReplaceAll(n.FirstChild.Data, "\t", "")
				index := strings.Index(priceStr, " ")
				if index >= 0 {
					priceStr = priceStr[index+1:]
					price, err := strconv.ParseFloat(priceStr, 64)
					if err == nil {
						(*priceMap)[historyRow] = price
					} else {
						(*priceMap)[historyRow] = -0.1
					}
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		generateHistoryRow2PriceMap(c, priceMap, historyRow)
	}
}

func generateHistoryRow2DateStringMap(n *html.Node, dateStringMap *map[string]string, historyRow string) {
	if n.Type == html.ElementNode {
		class := ""
		id := ""
		for _, attr := range n.Attr {
			switch attr.Key {
			case "class":
				class = attr.Val
			case "id":
				id = attr.Val
			}
			switch class {
			case "market_listing_row market_recent_listing_row":
				historyRow = id
			case "market_listing_listed_date_combined":
				dateString := n.FirstChild.Data
				dateString = strings.ReplaceAll(dateString, "\t", "")
				dateString = strings.ReplaceAll(dateString, "\n", "")
				(*dateStringMap)[historyRow] = dateString
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		generateHistoryRow2DateStringMap(c, dateStringMap, historyRow)
	}
}
