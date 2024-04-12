package utils

import (
	"ossTool/view/desktop"
)

func RunWidowsDesktop(ossWindows *desktop.OssMainWindow) {
	ossWindows.Show()
}

func HideWidowsDesktop(ossWindows *desktop.OssMainWindow) {
	ossWindows.Hide()
}

func CloseWidowsDesktop(ossWindows *desktop.OssMainWindow) {
	ossWindows.Close()
}

func HideAndShowDesktop(ossWindows *desktop.OssMainWindow) {
	if visible := ossWindows.Visible(); visible {
		ossWindows.Hide()
	} else {
		ossWindows.Show()
	}
}
