package initialize

import (
	"github.com/MakeNowJust/hotkey"
	"ossTool/global"
	"ossTool/utils"
)

func InitHotKey() {

	manager := hotkey.New()
	_, _ = manager.Register(hotkey.Ctrl, 'U', func() {
		utils.HideAndShowDesktop(global.MainWindow)
	})

}
