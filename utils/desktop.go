package utils

import (
	"ossTool/model"
)

func RunWidowsDesktop(ossWindows *model.OssMainWindow) {
	ossWindows.Show()
}

func HideWidowsDesktop(ossWindows *model.OssMainWindow) {
	ossWindows.Hide()
}

func CloseWidowsDesktop(ossWindows *model.OssMainWindow) {
	ossWindows.Close()
}

func HideAndShowDesktop(ossWindows *model.OssMainWindow) {
	if visible := ossWindows.Visible(); visible {
		ossWindows.Hide()
	} else {
		ossWindows.Show()
	}
}
