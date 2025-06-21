package trade

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/umichan0621/steam/pkg/auth"
	"github.com/umichan0621/steam/pkg/common"
)

type EconItem struct {
	AssetID    string `json:"assetid,omitempty"`
	InstanceID uint64 `json:"instanceid,string,omitempty"`
	ClassID    uint64 `json:"classid,string,omitempty"`
	AppID      uint32 `json:"appid"`
	ContextID  uint64 `json:"contextid,string"`
	Amount     uint32 `json:"amount,string"`
	Missing    bool   `json:"missing,omitempty"`
	EstUSD     uint32 `json:"est_usd,string"`
}

type TradeOffer struct {
	ID                 string      `json:"tradeofferid"`
	Partner            uint32      `json:"accountid_other"`
	ReceiptID          uint64      `json:"tradeid,string"`
	RecvItems          []*EconItem `json:"items_to_receive"`
	SendItems          []*EconItem `json:"items_to_give"`
	Message            string      `json:"message"`
	State              uint8       `json:"trade_offer_state"`
	ConfirmationMethod uint8       `json:"confirmation_method"`
	Created            int64       `json:"time_created"`
	Updated            int64       `json:"time_updated"`
	Expires            int64       `json:"expiration_time"`
	EscrowEndDate      int64       `json:"escrow_end_date"`
	RealTime           bool        `json:"from_real_time_trade"`
	IsOurOffer         bool        `json:"is_our_offer"`
}
type APIResponse struct {
	Inner *TradeOfferResponse `json:"response"`
}
type TradeOfferResponse struct {
	Offer          *TradeOffer            `json:"offer"`                 // GetTradeOffer
	SentOffers     []*TradeOffer          `json:"trade_offers_sent"`     // GetTradeOffers
	ReceivedOffers []*TradeOffer          `json:"trade_offers_received"` // GetTradeOffers
	Descriptions   []*common.EconItemDesc `json:"descriptions"`          // GetTradeOffers
}

func GetTradeOffers(auth *auth.Core, timeCutOff time.Time) (*TradeOfferResponse, error) {
	params := url.Values{
		"access_token":           {auth.AccessToken()},
		"get_sent_offers":        {"1"},
		"get_received_offers":    {"1"},
		"active_only":            {"1"},
		"get_descriptions":       {"1"},
		"language":               {"english"},
		"historical_only":        {"0"},
		"time_historical_cutoff": {strconv.FormatInt(timeCutOff.Unix(), 10)},
	}
	getUrl := fmt.Sprintf("%s/IEconService/GetTradeOffers/v1/?%s", common.URI_STEAM_API, params.Encode())
	httpRes, err := auth.HttpClient().Get(getUrl)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	res := APIResponse{}
	err = json.NewDecoder(httpRes.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	return res.Inner, nil
}

func GetTradeOffer(auth *auth.Core, offerID string) (*TradeOffer, error) {
	params := url.Values{
		"access_token": {auth.AccessToken()},
		"tradeofferid": {offerID},
	}
	getUrl := fmt.Sprintf("%s/IEconService/GetTradeOffer/v1/?%s", common.URI_STEAM_API, params.Encode())
	httpRes, err := auth.HttpClient().Get(getUrl)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	res := APIResponse{}
	err = json.NewDecoder(httpRes.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	return res.Inner.Offer, nil
}

func CancelTradeOffer(auth *auth.Core, offerID string) error {
	postUrl := fmt.Sprintf("%s/tradeoffer/%s/cancel", common.URI_STEAM_COMMUNITY, offerID)
	res, err := auth.HttpClient().PostForm(postUrl, url.Values{
		"sessionid": {auth.SessionID()},
	})
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func DeclineTradeOffer(auth *auth.Core, offerID string) error {
	postUrl := fmt.Sprintf("%s/tradeoffer/%s/decline", common.URI_STEAM_COMMUNITY, offerID)
	res, err := auth.HttpClient().PostForm(postUrl, url.Values{
		"sessionid": {auth.SessionID()},
	})
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
