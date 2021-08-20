package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var (
	EndPoint        string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
)

var (
	client *oss.Client
	bucket *oss.Bucket
)

// NewClient 新建OSS客户端
func NewClient() {
	EndPoint, AccessKeyID, AccessKeySecret, BucketName = GlobalCfg.OSSConfig.EndPoint, GlobalCfg.OSSConfig.AccessKeyID, GlobalCfg.OSSConfig.AccessKeySecret, GlobalCfg.OSSConfig.BucketName
	var err error
	client, err = oss.New(EndPoint, AccessKeyID, AccessKeySecret)
	if err != nil {
		log.Fatalln(err)
	}
}

// CreateBucket 创建一个bucket
func CreateBucket() {
	NewClient()
	var err error
	err = client.CreateBucket(BucketName)
	if err != nil {
		log.Fatalln(err)
	}
}

// GetBucket 获取一个bucket
func GetBucket() {
	NewClient()
	ifExist, err := client.IsBucketExist(BucketName)
	if err != nil {
		log.Fatalln(err)
	}
	if ifExist {
		bucket, err = client.Bucket(BucketName)
	} else {
		CreateBucket()
	}
	if err != nil {
		log.Fatalln(err)
	}
}

// Upload 上传图片到OSS
func Upload(objectKey string, fileReader io.Reader) string {
	NewClient()
	GetBucket()
	var err error
	// 指定存储类型为标准存储，缺省也为标准存储。
	storageType := oss.ObjectStorageClass(oss.StorageStandard)

	// 指定访问权限为公共读，缺省为继承bucket的权限。
	objectAcl := oss.ObjectACL(oss.ACLPublicRead)

	err = bucket.PutObject(objectKey, fileReader, storageType, objectAcl)
	if err != nil {
		log.Fatalln(err)
	}
	exist, err := bucket.IsObjectExist(objectKey)
	if !exist || err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	return "https://" + BucketName + "." + EndPoint + "/" + objectKey
}

// UploadRemoteFile 获取远程图片并上传到OSS
func UploadRemoteFile(fileUrl string) string {
	resp, err := http.Get(fileUrl)
	defer resp.Body.Close()
	if err != nil {
		log.Fatalf("获取远程图片失败:%v \n", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("读取远程图片失败:%v \n", err)
	}
	data := []byte(fileUrl)
	md5Sum := md5.Sum(data)
	fileName := fmt.Sprintf("%x", md5Sum) + ".jpg"
	return Upload(fileName, bytes.NewReader(body))
}
