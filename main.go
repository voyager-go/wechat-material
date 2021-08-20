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

	"github.com/PuerkitoBio/goquery"
)

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

func main() {
	LoadConfig()
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	http.Handle("/ueditor/", http.StripPrefix("/ueditor/", http.FileServer(http.Dir("ueditor"))))
	http.HandleFunc("/index", ParseTemplate)
	err := http.ListenAndServe(":8899", nil)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}
