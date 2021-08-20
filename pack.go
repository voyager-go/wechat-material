package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// PackageArticle
func Handler() {
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

// DoFetch 开始抓取素材列表信息
func DoFetch(mType string, offset, newsCount int) {
	lists := GetMaterialList(mType, offset, newsCount)
	obj := new(Article)
	articles := make([]*Article, newsCount)
	for idx := range lists.Item {
		for itemIdx := range lists.Item[idx].Content.NewsItem {
			// 解析出的单个图文消息
			// 持久化到数据库
			NewsItem := lists.Item[idx].Content.NewsItem
			newsArticle := new(Article)
			newsArticle.Title = NewsItem[itemIdx].Title
			newsArticle.Cover = UploadRemoteFile(NewsItem[itemIdx].ThumbURL)
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
	affected, err := obj.Create(canUseArticles)
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
		log.Fatalf("图文素材解析失败:%v \n", err)
	}
	document.Find("img").Each(func(i int, selection *goquery.Selection) {
		imgUrl, _ := selection.Attr("src")
		imgSrc := UploadRemoteFile(imgUrl)
		selection.SetAttr("src", imgSrc)
	})
	html, _ := document.Html()
	return html
}
