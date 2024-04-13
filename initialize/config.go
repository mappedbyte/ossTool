package initialize

import "ossTool/model"

func InitEndpointConfig() {
	endpointConfig := model.LoadOssConfigFromSettings(model.AppMainWindow.App.Settings())
	model.AppEndpointConfig = &endpointConfig
}
