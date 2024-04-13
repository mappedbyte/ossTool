package model

import (
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/lxn/walk"
	"strings"
)

var AppEndpointConfig *EndpointConfig

type EndpointConfig struct {
	Endpoint  string `json:"endpoint"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	Bucket    string `json:"bucket"`
}

func LoadOssConfigFromSettings(settings walk.Settings) EndpointConfig {
	var endpointConfig = EndpointConfig{}
	if endpoint, exists := settings.Get("endpoint"); exists {
		endpointConfig.Endpoint = endpoint
	}
	if bucket, exists := settings.Get("bucket"); exists {
		endpointConfig.Bucket = bucket
	}
	if accessKey, exists := settings.Get("accessKey"); exists {
		endpointConfig.AccessKey = accessKey
	}
	if secretKey, exists := settings.Get("secretKey"); exists {
		endpointConfig.SecretKey = secretKey
	}
	return endpointConfig
}

func GetGlobalEndpointConfig() *EndpointConfig {
	return AppEndpointConfig
}

func (endpointConfig *EndpointConfig) CreateOssClient() (*oss.Client, []EnvItem, error) {
	items := make([]EnvItem, 0)

	if endpointConfig.Endpoint == "" || endpointConfig.Bucket == "" || endpointConfig.AccessKey == "" || endpointConfig.SecretKey == "" {
		items = append(items, EnvItem{Value: "系统提示: 配置文件未填写或有误,请检查!", Name: "bucketList"})
		return nil, items, errors.New("系统提示: 配置文件未填写或有误,请检查")
	}
	client, err := oss.New(endpointConfig.Endpoint, endpointConfig.AccessKey, endpointConfig.SecretKey)
	if err != nil {
		panic(err)
	}
	bucket, err := client.Bucket(endpointConfig.Bucket)
	if err != nil {
		//panic(err)
		return nil, items, err
	}
	buckets, err := client.ListBuckets()
	if err != nil {
		//panic(err)
		return nil, nil, err
	}
	var bucketList []string
	for _, bucket := range buckets.Buckets {
		name := bucket.Name
		bucketList = append(bucketList, name)
	}
	bucketName := bucket.BucketName
	items = append(items, EnvItem{Value: "使用oss版本:" + oss.Version, Name: "version"})
	items = append(items, EnvItem{Value: "存储桶列表:" + strings.Join(bucketList, ", "), Name: "bucketList"})
	items = append(items, EnvItem{Value: "当前使用存储桶:" + bucketName, Name: "bucketList"})

	return client, items, nil
}

// 定义进度条监听器。

type OssProgressListener struct {
	//Message chan string
	ProgressBar *walk.ProgressBar
}

// 定义进度变更事件处理函数。

func (listener *OssProgressListener) ProgressChanged(event *oss.ProgressEvent) {
	switch event.EventType {
	case oss.TransferStartedEvent:
		fmt.Printf("Transfer Started, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	/*	listener.Message <- fmt.Sprintf("Transfer Started, ConsumedBytes: %d, TotalBytes %d.\n",
		event.ConsumedBytes, event.TotalBytes)*/
	case oss.TransferDataEvent:
		fmt.Printf("\rTransfer Data, ConsumedBytes: %d, TotalBytes %d, %d%%.",
			event.ConsumedBytes, event.TotalBytes, event.ConsumedBytes*100/event.TotalBytes)
		/*listener.Message <- fmt.Sprintf("\rTransfer Data, ConsumedBytes: %d, TotalBytes %d, %d%%.",
		event.ConsumedBytes, event.TotalBytes, event.ConsumedBytes*100/event.TotalBytes)*/
		listener.ProgressBar.SetValue(int(event.ConsumedBytes * 100 / event.TotalBytes))
	case oss.TransferCompletedEvent:
		fmt.Printf("\nTransfer Completed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
		/*listener.Message <- fmt.Sprintf("\nTransfer Completed, ConsumedBytes: %d, TotalBytes %d.\n",
		event.ConsumedBytes, event.TotalBytes)*/
	case oss.TransferFailedEvent:
		fmt.Printf("\nTransfer Failed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
		/*listener.Message <- fmt.Sprintf("\nTransfer Failed, ConsumedBytes: %d, TotalBytes %d.\n",
		event.ConsumedBytes, event.TotalBytes)*/
	default:
	}
}
