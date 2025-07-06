package auth

import (
	"crypto/md5"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type LoginInfo struct {
	UserName       string
	Password       string
	SharedSecret   string
	IdentitySecret string
}

type Core struct {
	httpClient *http.Client
	loginInfo  LoginInfo
	cookieData CookieData
	profileUrl string
	deviceID   string
}

func (core *Core) Init(info LoginInfo) {
	core.loginInfo = info
	core.httpClient = &http.Client{}
	core.profileUrl = ""
	sum := md5.Sum([]byte(info.UserName + info.Password))
	core.deviceID = fmt.Sprintf("android:%x-%x-%x-%x-%x",
		sum[:2], sum[2:4], sum[4:6], sum[6:8], sum[8:10])
}

func (core *Core) HttpClient() *http.Client { return core.httpClient }
func (core *Core) SteamID() string          { return core.cookieData.SteamID }
func (core *Core) SessionID() string        { return core.cookieData.SessionID }
func (core *Core) DeviceID() string         { return core.deviceID }
func (core *Core) IdentitySecret() string   { return core.loginInfo.IdentitySecret }

func (core *Core) AccessToken() string {
	temp := strings.Split(core.cookieData.SteamLoginSecure, "%7C%7C")
	if len(temp) >= 2 {
		return temp[1]
	}
	return ""
}

// timeout: millsecond, set only while timeout > 0;
// proxy: if proxyUrl == "", ignore
func (core *Core) SetHttpParam(timeout int, proxy string) error {
	transport := &http.Transport{}
	if proxy != "" {
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			return err
		}
		transport.Proxy = http.ProxyURL(proxyUrl)
	}
	if timeout > 0 {
		timeoutVal := time.Duration(timeout) * time.Millisecond
		dialer := net.Dialer{Timeout: timeoutVal}

		transport.DialContext = dialer.DialContext
		transport.TLSHandshakeTimeout = timeoutVal
		transport.ResponseHeaderTimeout = timeoutVal
		transport.ExpectContinueTimeout = timeoutVal
		core.httpClient.Timeout = timeoutVal
	}
	core.httpClient.Transport = transport
	return nil
}
