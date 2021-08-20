package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// MaterialList 素材列表
type MaterialList struct {
	Item []struct {
		Content struct {
			CreateTime int64              `json:"create_time"`
			NewsItem   []MaterialNewsItem `json:"news_item"`
			UpdateTime int64              `json:"update_time"`
		} `json:"content"`
		MediaID    string `json:"media_id"`
		UpdateTime int64  `json:"update_time"`
	} `json:"item"`
	ItemCount  int64 `json:"item_count"`
	TotalCount int64 `json:"total_count"`
}

// MaterialNewsItem 单个素材
type MaterialNewsItem struct {
	Author             string `json:"author"`
	Content            string `json:"content"`
	ContentSourceURL   string `json:"content_source_url"`
	Digest             string `json:"digest"`
	NeedOpenComment    int64  `json:"need_open_comment"`
	OnlyFansCanComment int64  `json:"only_fans_can_comment"`
	ShowCoverPic       int64  `json:"show_cover_pic"`
	ThumbMediaID       string `json:"thumb_media_id"`
	ThumbURL           string `json:"thumb_url"`
	Title              string `json:"title"`
	URL                string `json:"url"`
}

// MaterialCount 素材总数
type MaterialCount struct {
	ImageCount int `json:"image_count"`
	NewsCount  int `json:"news_count"`
	VideoCount int `json:"video_count"`
	VoiceCount int `json:"voice_count"`
}

// ParseTemplate 解析网页模板
func ParseTemplate(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("/Users/artist/go/src/wxmedia/ueditor/content.html")
	if err != nil {
		fmt.Fprintf(w, "parse html errors: %s", err.Error())
		return
	}
	Html := template.HTML("")
	t.Execute(w, Html)
}

// GetMaterialCount 获取素材数量
func GetMaterialCount() *MaterialCount {
	url := "https://api.weixin.qq.com/cgi-bin/material/get_materialcount?access_token=%s"
	wholeUrl := fmt.Sprintf(url, GetAccessToken())
	resp, err := http.Get(wholeUrl)
	defer resp.Body.Close()
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	mc := new(MaterialCount)
	json.Unmarshal(body, &mc)
	return mc
}

// GetMaterialList 获取素材列表
func GetMaterialList(mType string, offset, pageSize int) *MaterialList {
	url := "https://api.weixin.qq.com/cgi-bin/material/batchget_material?access_token=%s"
	wholeUrl := fmt.Sprintf(url, GetAccessToken())
	data := make(map[string]interface{})
	data["type"] = mType
	data["offset"] = offset
	data["count"] = pageSize
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := http.Post(wholeUrl, "application/json", strings.NewReader(string(bytes)))
	defer resp.Body.Close()
	if err != nil {
		log.Fatalln(err)
	}
	reader, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	lists := new(MaterialList)
	json.Unmarshal(reader, &lists)
	return lists
}
