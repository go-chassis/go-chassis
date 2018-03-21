package session

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/lager"

	cache "github.com/patrickmn/go-cache"
)

// ErrResponseNil used for to represent the error response, when it is nil
var ErrResponseNil = errors.New("Can not Set session, resp is nil")

// SessionCache session cache variable
var SessionCache *cache.Cache

func init() {
	SessionCache = initCache()
}
func initCache() *cache.Cache {
	var value *cache.Cache

	value = cache.New(3e+10, time.Second*30)
	return value
}

//GetSessionFromResp return session uuid in resp if there is
func GetSessionFromResp(cookieKey string, resp *http.Response) string {
	for _, c := range resp.Cookies() {
		if c.Name == cookieKey {
			return c.Value
		}
	}
	return ""
}

// CheckForSessionID check session id
func CheckForSessionID(ep string, autoTimeout int, resp *http.Response, req *http.Request) {
	if resp == nil {
		lager.Logger.Warn("", ErrResponseNil)
		return
	}

	timeValue := time.Duration(autoTimeout) * time.Second

	var sessionIDStr string

	if c, err := req.Cookie(common.LBSessionID); err == http.ErrNoCookie {
		sessionIDStr = ""
	} else {
		sessionIDStr = c.Value
	}

	ClearExpired()
	var sessBool bool
	if sessionIDStr != "" {
		_, sessBool = SessionCache.Get(sessionIDStr)
	}

	valueChassisLb := GetSessionFromResp(common.LBSessionID, resp)
	//if session is in resp, then just save it
	if string(valueChassisLb) != "" {
		Save(valueChassisLb, ep, timeValue)
	} else if sessionIDStr != "" && sessBool {
		c1 := new(http.Cookie)
		c1.Name = common.LBSessionID
		c1.Value = sessionIDStr
		setCookie(c1, resp)
		Save(sessionIDStr, ep, timeValue)
	} else {
		c1 := new(http.Cookie)
		c1.Name = common.LBSessionID
		sessionIDValue := generateCookieSessionID()
		c1.Value = sessionIDValue
		setCookie(c1, resp)
		Save(sessionIDValue, ep, timeValue)

	}

}

// generateCookieSessionID generate cookies for session id
func generateCookieSessionID() string {

	result := make([]byte, 16)

	rand.Seed(time.Now().UTC().UnixNano())
	tmp := rand.Int63()
	rand.Seed(tmp)
	for i := 0; i < 16; i++ {
		result[i] = byte(rand.Intn(16))
	}

	result[6] = (result[6] & 0xF) | (4 << 4)
	result[8] = (result[8] | 0x40) & 0x7F

	return fmt.Sprintf("%x-%x-%x-%x-%x", result[0:4], result[4:6], result[6:8], result[8:10], result[10:])

}

// setCookie set cookie
func setCookie(cookie *http.Cookie, resp *http.Response) {
	resp.Header.Add("Set-Cookie", cookie.String())
}

// DeletingKeySuccessiveFailure deleting key successes and failures
func DeletingKeySuccessiveFailure(resp *http.Response) {
	if resp == nil {
		lager.Logger.Warn("", ErrResponseNil)
		return
	}
	SessionCache.DeleteExpired()
	valueChassisLb := GetSessionFromResp(common.LBSessionID, resp)
	if string(valueChassisLb) != "" {
		cookieKey := strings.Split(string(valueChassisLb), "=")
		if len(cookieKey) > 1 {
			Delete(cookieKey[1])
		}
	}
}

// GetSessionCookie getting session cookie
func GetSessionCookie(resp *http.Response) string {
	if resp == nil {
		lager.Logger.Warn("", ErrResponseNil)
		return ""
	}

	valueChassisLb := GetSessionFromResp(common.LBSessionID, resp)
	if string(valueChassisLb) != "" {
		return string(valueChassisLb)
	}

	return ""
}
