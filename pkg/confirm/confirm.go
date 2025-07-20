package confirm

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/umichan0621/steam/pkg/auth"
	"github.com/umichan0621/steam/pkg/common"
)

type ConfirmationResponse struct {
	Success       bool            `json:"success"`
	Confirmations []*Confirmation `json:"conf"`
}

type Confirmation struct {
	ID           string   `json:"id"`
	Type         uint8    `json:"type"`
	Creator      string   `json:"creator_id"`
	Nonce        string   `json:"nonce"`
	CreationTime uint64   `json:"creation"`
	TypeName     string   `json:"type_name"`
	Cancel       string   `json:"cancel"`
	Accept       string   `json:"accept"`
	Icon         string   `json:"icon"`
	Multi        bool     `json:"multi"`
	Headline     string   `json:"headline"`
	Summary      []string `json:"summary"`
}

func GetConfirmations(auth *auth.Core) ([]*Confirmation, error) {
	identitySecret := auth.IdentitySecret()
	if identitySecret == "" {
		return nil, fmt.Errorf("empty identity secret")
	}
	current := time.Now().Unix()

	key, err := generateConfirmationCode(identitySecret, "conf", current)
	if err != nil {
		return nil, err
	}
	params := url.Values{
		"p":   {auth.DeviceID()},
		"a":   {auth.SteamID()},
		"k":   {key},
		"t":   {strconv.FormatInt(current, 10)},
		"m":   {"android"},
		"tag": {"conf"},
	}

	getUrl := fmt.Sprintf("%s/mobileconf/getlist?%s", common.URI_STEAM_COMMUNITY, params.Encode())
	httpRes, err := auth.HttpClient().Get(getUrl)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	data, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}
	res := ConfirmationResponse{}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	if !res.Success {
		return nil, fmt.Errorf("fail to get ConfirmationResponse: %s", string(data))
	}
	return res.Confirmations, nil
}

func AnswerConfirmation(auth *auth.Core, confirmation *Confirmation, answer string) error {
	identitySecret := auth.IdentitySecret()
	if identitySecret == "" {
		return fmt.Errorf("empty identity secret")
	}
	current := time.Now().Unix()

	key, err := generateConfirmationCode(identitySecret, answer, current)
	if err != nil {
		return err
	}
	params := url.Values{
		"p":   {auth.DeviceID()},
		"a":   {auth.SteamID()},
		"k":   {key},
		"t":   {strconv.FormatInt(current, 10)},
		"m":   {"android"},
		"tag": {answer},
		"op":  {answer},
		"cid": {confirmation.ID},
		"ck":  {confirmation.Nonce},
	}
	reqUrl := fmt.Sprintf("%s/mobileconf/ajaxop?%s", common.URI_STEAM_COMMUNITY, params.Encode())
	httpReq, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return err
	}
	httpReq.Header.Set("X-Requested-With", "XMLHttpRequest")
	httpRes, err := auth.HttpClient().Do(httpReq)
	if err != nil {
		return err
	}
	if httpRes != nil {
		defer httpRes.Body.Close()
	}
	data, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return err
	}
	type Response struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	res := Response{}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return err
	}

	if !res.Success {
		return errors.New(res.Message)
	}
	return nil
}

func generateConfirmationCode(identitySecret, tag string, current int64) (string, error) {
	data, err := base64.StdEncoding.DecodeString(identitySecret)
	if err != nil {
		return "", err
	}

	ful := make([]byte, 8+len(tag))
	binary.BigEndian.PutUint32(ful[4:], uint32(current))
	copy(ful[8:], tag)

	hash := hmac.New(sha1.New, data)
	_, err = hash.Write(ful)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(hash.Sum(nil)), nil
}
