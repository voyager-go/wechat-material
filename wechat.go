package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/bluele/gcache"
)

type WechatAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

var (
	gc             = gcache.New(20).LRU().Build()
	accessTokenKey = "accessToken"
)

// SetAccessToken 设置AccessToken
func SetAccessToken() string {
	urlString := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	wholeUrl := fmt.Sprintf(urlString, GlobalCfg.WeChatConfig.AppId, GlobalCfg.WeChatConfig.AppSecret)
	resp, err := http.Get(wholeUrl)
	defer resp.Body.Close()
	if err != nil {
		log.Fatalln(err)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	at := new(WechatAccessToken)
	err = json.Unmarshal(bytes, &at)
	if err != nil {
		log.Fatalln(err)
	}
	err = gc.SetWithExpire(accessTokenKey, at.AccessToken, time.Second*7000)
	if err != nil {
		log.Fatalln(err)
	}
	return at.AccessToken
}

// GetAccessToken 获取 access token
func GetAccessToken() string {
	value, err := gc.Get(accessTokenKey)
	if err != nil {
		value = SetAccessToken()
	}
	return value.(string)
}
