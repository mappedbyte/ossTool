package main

import (
	"ossTool/initialize"
	"ossTool/model"
)

func main() {
	initialize.InitDesktop()
	initialize.InitHotKey()
	initialize.InitEndpointConfig()
	initialize.InitOss()
	model.AppMainWindow.Run()

	endpointConfig := model.GetGlobalEndpointConfig()
	model.AppMainWindow.App.Settings().Put("endpoint", endpointConfig.Endpoint)
	model.AppMainWindow.App.Settings().Put("bucket", endpointConfig.Bucket)
	model.AppMainWindow.App.Settings().Put("accessKey", endpointConfig.AccessKey)
	model.AppMainWindow.App.Settings().Put("secretKey", endpointConfig.SecretKey)
	model.AppMainWindow.App.Settings().Save()
	//fmt.Println(".........")
}
