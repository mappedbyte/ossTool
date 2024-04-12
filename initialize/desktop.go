package initialize

import (
	"ossTool/global"
	"ossTool/view/desktop"
)

func InitDesktop() {
	global.MainWindow = desktop.InitMainWindow()
	global.MainWindow.InitAppSetting()
	global.MainWindow.Show()
	global.MainWindow.NewNotifyIcon()
}
