package initialize

import (
	"ossTool/global"
	"ossTool/utils"
)

func InitOss() {

	utils.NewOssClient(*global.MainWindow.EndpointConfig)

}
