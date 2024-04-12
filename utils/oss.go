package utils

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"ossTool/config"
	"ossTool/global"
	"ossTool/view/desktop"
	"strings"
)

func NewOssClient(endpointConfig config.EndpointConfig) {
	client, err := oss.New(endpointConfig.Endpoint, endpointConfig.AccessKey, endpointConfig.SecretKey)
	if err != nil {
		panic(err)
	}
	bucket, err := client.Bucket(endpointConfig.Bucket)
	if err != nil {
		panic(err)
	}
	buckets, err := client.ListBuckets()
	if err != nil {
		panic(err)
	}
	var bucketList []string
	for _, bucket := range buckets.Buckets {
		name := bucket.Name
		bucketList = append(bucketList, name)
	}

	bucketName := bucket.BucketName

	model := desktop.NewEnvModel()
	model.Items = append(model.Items, desktop.EnvItem{Value: "使用oss版本:" + oss.Version, Name: "version"})
	model.Items = append(model.Items, desktop.EnvItem{Value: "存储桶列表:" + strings.Join(bucketList, ", "), Name: "bucketList"})
	model.Items = append(model.Items, desktop.EnvItem{Value: "当前使用存储桶:" + bucketName, Name: "bucketList"})
	global.MainWindow.DisplayListBox.SetModel(model)
}
