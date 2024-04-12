package main

import (
	"ossTool/global"
	"ossTool/initialize"
)

func main() {
	initialize.InitDesktop()
	initialize.InitHotKey()
	initialize.InitEndpointConfig()
	initialize.InitOss()
	global.MainWindow.Run()
}
