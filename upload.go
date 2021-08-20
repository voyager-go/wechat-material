package main

import (
	"fmt"
	"io"
	"log"
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
	err    error
	bucket *oss.Bucket
)

// NewClient 新建OSS客户端
func NewClient() {
	EndPoint, AccessKeyID, AccessKeySecret, BucketName = GlobalCfg.OSSConfig.EndPoint, GlobalCfg.OSSConfig.AccessKeyID, GlobalCfg.OSSConfig.AccessKeySecret, GlobalCfg.OSSConfig.BucketName
	client, err = oss.New(EndPoint, AccessKeyID, AccessKeySecret)
	if err != nil {
		log.Fatalln(err)
	}
}

// CreateBucket 创建一个bucket
func CreateBucket() {
	NewClient()
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
