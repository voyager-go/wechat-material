package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func PackageArticle() {
	mType := ArticleTypeNews.String()
	count := GetMaterialCount()
	pageSize := 20
	if count.NewsCount <= pageSize {
		DoFetch(mType, 1, count.NewsCount)
	} else {
		for i := 1; i < count.NewsCount/pageSize; i++ {
			offset := (i - 1) * pageSize
			DoFetch(mType, offset, pageSize)
		}
	}
}

func DoFetch(mType string, offset, newsCount int) {
	lists := GetMaterialList(mType, offset, newsCount)
	obj := new(Article)
	articles := make([]*Article, 1)
	for idx := range lists.Item {
		for itemIdx := range lists.Item[idx].Content.NewsItem {
			// 解析出的单个图文消息
			// 持久化到数据库
			NewsItem := lists.Item[idx].Content.NewsItem
			newsArticle := new(Article)
			newsArticle.Title = NewsItem[itemIdx].Title
			newsArticle.Type = uint(ArticleTypeNews)
			newsArticle.Author = NewsItem[itemIdx].Author
			newsArticle.Digest = NewsItem[itemIdx].Digest
			newsArticle.Content = ParseItem(NewsItem[itemIdx].Content)
			newsArticle.UpdateTime = time.Unix(lists.Item[idx].Content.UpdateTime, 0).Format("2006-01-02 15:04:05")
			newsArticle.CreateTime = time.Unix(lists.Item[idx].Content.CreateTime, 0).Format("2006-01-02 15:04:05")
			articles = append(articles, newsArticle)
		}

	}
	// 处理切片中的 <nil>
	canUseArticles := make([]*Article, 0, len(articles))
	for _, item := range articles {
		if item != nil {
			canUseArticles = append(canUseArticles, item)
		}
	}
	affected, err := obj.Create(articles)
	if err != nil {
		log.Fatalf("持久化素材信息失败: %v \n", err)
	}
	fmt.Println(affected)
}

// ParseItem 解析单个图文
func ParseItem(article string) string {
	r, _ := regexp.Compile("data-src")
	article = r.ReplaceAllString(article, "src")
	document, err := goquery.NewDocumentFromReader(strings.NewReader(article))
	if err != nil {
		log.Fatalln(err)
	}
	document.Find("img").Each(func(i int, selection *goquery.Selection) {
		imgUrl, _ := selection.Attr("src")
		resp, err := http.Get(imgUrl)
		defer resp.Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		data := []byte(imgUrl)
		md5Sum := md5.Sum(data)
		fileName := fmt.Sprintf("%x", md5Sum) + ".jpg"
		imgSrc := Upload(fileName, bytes.NewReader(body))
		selection.SetAttr("src", imgSrc)
	})
	html, _ := document.Html()
	return html
}
