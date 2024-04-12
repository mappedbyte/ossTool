package initialize

import (
	"ossTool/config"
	"ossTool/global"
)

func InitEndpointConfig() {
	endpointConfig := config.EndpointConfig{}
	settings := global.MainWindow.App.Settings()
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
	global.MainWindow.EndpointConfig = &endpointConfig

}
