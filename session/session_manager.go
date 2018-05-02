package session

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/ServiceComb/go-chassis/client/rest"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/lager"

	"context"
	cache "github.com/patrickmn/go-cache"
)

// ErrResponseNil used for to represent the error response, when it is nil
var ErrResponseNil = errors.New("Can not Set session, resp is nil")

// SessionCache session cache variable
var SessionCache *cache.Cache

func init() {
	SessionCache = initCache()
	cookieMap = make(map[string]string)
}
func initCache() *cache.Cache {
	var value *cache.Cache

	value = cache.New(3e+10, time.Second*30)
	return value
}

var cookieMap map[string]string

// getLBCookie gets cookie from local map
func getLBCookie(key string) string {
	return cookieMap[key]
}

// setLBCookie sets cookie to local map
func setLBCookie(key, value string) {
	cookieMap[key] = value
}

// GetContextMetadata gets data from context
func GetContextMetadata(ctx context.Context, key string) string {
	md, ok := ctx.Value(common.ContextValueKey{}).(map[string]string)
	if ok {
		for k, v := range md {
			if k == key {
				return v
			}
		}
	}
	return ""
}

// SetContextMetadata sets data to context
func SetContextMetadata(ctx context.Context, key string, value string) context.Context {
	md, ok := ctx.Value(common.ContextValueKey{}).(map[string]string)
	if !ok {
		md = make(map[string]string)
	}

	if md[key] == value {
		return ctx
	}

	md[key] = value
	return context.WithValue(ctx, common.ContextValueKey{}, md)
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

// CheckForSessionIDFromContext check session id
func CheckForSessionIDFromContext(ctx context.Context, ep string, autoTimeout int) context.Context {

	timeValue := time.Duration(autoTimeout) * time.Second

	sessionIDStr := GetContextMetadata(ctx, common.LBSessionID)
	if sessionIDStr != "" {
		cookieKey := strings.Split(string(sessionIDStr), "=")
		if len(cookieKey) > 1 {
			sessionIDStr = cookieKey[1]
		}
	}

	ClearExpired()
	var sessBool bool
	if sessionIDStr != "" {
		_, sessBool = SessionCache.Get(sessionIDStr)
	}

	if sessionIDStr != "" && sessBool {
		cookie := common.LBSessionID + "=" + sessionIDStr
		setLBCookie(common.LBSessionID, cookie)
		Save(sessionIDStr, ep, timeValue)
		return ctx
	}

	sessionIDValue := generateCookieSessionID()
	cookie := common.LBSessionID + "=" + sessionIDValue
	setLBCookie(common.LBSessionID, cookie)
	Save(sessionIDValue, ep, timeValue)
	return SetContextMetadata(ctx, common.LBSessionID, cookie)
}

//Temporary responsewriter for SetCookie
type cookieResponseWriter http.Header

// Header implements ResponseWriter Header interface
func (c cookieResponseWriter) Header() http.Header {
	return http.Header(c)
}

//Write is a dummy function
func (c cookieResponseWriter) Write([]byte) (int, error) {
	panic("ERROR")
}

//WriteHeader is a dummy function
func (c cookieResponseWriter) WriteHeader(int) {
	panic("ERROR")
}

//setCookie appends cookie with already present cookie with ';' in between
func setCookie(resp *http.Response, value string) {
	Resp := rest.Response{Resp: resp}

	newCookie := common.LBSessionID + "=" + value
	oldCookie := string(Resp.GetCookie(common.LBSessionID))

	if oldCookie != "" {
		//If cookie is already set, append it with ';'
		newCookie = newCookie + ";" + oldCookie
	}

	c1 := http.Cookie{Name: common.LBSessionID, Value: newCookie}

	w := cookieResponseWriter(resp.Header)
	http.SetCookie(w, &c1)
}

// CheckForSessionID check session id
func CheckForSessionID(ep string, autoTimeout int, resp *http.Response, req *http.Request) {
	if resp == nil {
		lager.Logger.Warnf("", ErrResponseNil)
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
		setCookie(resp, sessionIDStr)
		Save(sessionIDStr, ep, timeValue)
	} else {
		sessionIDValue := generateCookieSessionID()
		setCookie(resp, sessionIDValue)
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

// DeletingKeySuccessiveFailure deleting key successes and failures
func DeletingKeySuccessiveFailure(resp *http.Response) {
	SessionCache.DeleteExpired()
	if resp == nil {
		valueChassisLb := getLBCookie(common.LBSessionID)
		if string(valueChassisLb) != "" {
			cookieKey := strings.Split(string(valueChassisLb), "=")
			if len(cookieKey) > 1 {
				Delete(cookieKey[1])
				setLBCookie(common.LBSessionID, "")
			}
		}
		return
	}

	valueChassisLb := GetSessionFromResp(common.LBSessionID, resp)
	if string(valueChassisLb) != "" {
		cookieKey := strings.Split(string(valueChassisLb), "=")
		if len(cookieKey) > 1 {
			Delete(cookieKey[1])
		}
	}
}

// GetSessionCookie getting session cookie
func GetSessionCookie(ctx context.Context, resp *http.Response) string {
	if ctx != nil {
		return GetContextMetadata(ctx, common.LBSessionID)
	}

	if resp == nil {
		lager.Logger.Warnf("", ErrResponseNil)
		return ""
	}

	valueChassisLb := GetSessionFromResp(common.LBSessionID, resp)
	if string(valueChassisLb) != "" {
		return string(valueChassisLb)
	}

	return ""
}
