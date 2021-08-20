package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/google/uuid"
)

// SavePic 保存图片
func SavePic(url string) string {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	var fileName string
	if isExists, _ := PathExists("./images"); isExists {
		fileName = uuid.New().String() + ".jpeg"
		_ = ioutil.WriteFile("./images/"+fileName, body, 0755)
	}
	return "/images/" + fileName
}

// PathExists 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		// 创建文件夹
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
		} else {
			return true, nil
		}
	}
	return false, err
}
