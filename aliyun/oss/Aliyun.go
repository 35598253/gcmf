package oss

import (
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// Aliyun 管理对象
type Aliyun struct {
	AccessKeyId     string
	AccessKeySecret string
	Endpoint        string
	BucketName      string
	UpPath          string
}

func (a *Aliyun) upload(UpFile string) (string, error) {

	var upPath = a.UpPath + "/" + time.Now().Format("2006/01/02/")

	obName := upPath + "123.jpg"
	client, err := oss.New(a.Endpoint, a.AccessKeyId, a.AccessKeySecret)
	if err != nil {
		return "", err
	}
	// 获取存储空间。
	bucket, err := client.Bucket(a.BucketName)
	if err != nil {
		return "", err
	}
	err = bucket.PutObjectFromFile(obName, UpFile)
	if err != nil {
		return "", err
	}

	return obName, nil
}
