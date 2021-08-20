package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

var (
	GlobalCfg *Config
)

// WeChatConfig 微信公众号配置
type WeChatConfig struct {
	AppId     string `yaml:"AppId"`
	AppSecret string `yaml:"AppSecret"`
}

//  OSSConfig OSS上传配置
type OSSConfig struct {
	EndPoint        string `yaml:"EndPoint"`
	AccessKeyID     string `yaml:"AccessKeyID"`
	AccessKeySecret string `yaml:"AccessKeySecret"`
	BucketName      string `yaml:"BucketName"`
}

// DataBaseConfig 数据库配置
type DataBaseConfig struct {
	Driver string `yaml:"Driver"`
	Host   string `yaml:"Host"`
	Port   int    `yaml:"Port"`
	User   string `yaml:"User"`
	Pass   string `yaml:"Pass"`
	DBName string `yaml:"DBName"`
}

// Config 全局配置集合
type Config struct {
	WeChatConfig   WeChatConfig   `yaml:"WeChat"`
	OSSConfig      OSSConfig      `yaml:"OSS"`
	DataBaseConfig DataBaseConfig `yaml:"DataBase"`
}

// LoadConfig 加载配置文件
func LoadConfig() {
	yFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("ioutil.ReadFile 读取配置文件失败:%v \n", err)
	}
	err = yaml.Unmarshal(yFile, &GlobalCfg)
	if err != nil {
		log.Fatalf("yaml.Unmarshal 解析配置文件失败:%v \n", err)
	}
}
