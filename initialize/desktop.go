package initialize

import (
	"ossTool/model"
	"ossTool/view/desktop"
)

func InitDesktop() {
	model.AppMainWindow = desktop.InitMainWindow()
	model.AppMainWindow.InitAppSetting()
	model.AppMainWindow.Show()
	model.AppMainWindow.NewNotifyIcon()
}
