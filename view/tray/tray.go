package tray

import (
	"fmt"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"ossTool/global"
	"ossTool/utils"
)

func InitTray() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("OssTool")
	systray.SetTooltip("文件上传小插件")
	runWindow := systray.AddMenuItem("主界面", "Open the main window")
	systray.AddSeparator()
	mdCopy := systray.AddMenuItem("复制md格式", "copy for markdown")
	systray.AddSeparator()
	mURL := systray.AddMenuItem("文档地址", "visit the home page")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("退出程序", "Quit the whole app")
	go func() {
		for {
			select {
			case <-runWindow.ClickedCh:
				go func() {
					//运行主窗口
					utils.RunWidowsDesktop(global.MainWindow)
				}()
			case <-mdCopy.ClickedCh:
				if mdCopy.Checked() {
					mdCopy.Uncheck() //取消选中
					//mdFlag = false
					fmt.Println("取消选中")
				} else {
					mdCopy.Check() //选中
					//mdFlag = true
					fmt.Println("选中")
				}
			case <-mURL.ClickedCh:
				//open.Run("https://www.cnblogs.com/ludg/")
			case <-mQuit.ClickedCh:
				systray.Quit() //退出托盘
				utils.CloseWidowsDesktop(global.MainWindow)
				return
			}
		}
	}()
}
