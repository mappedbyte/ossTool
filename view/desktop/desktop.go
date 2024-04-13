package desktop

import (
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"ossTool/model"
)

func InitMainWindow() *model.OssMainWindow {

	mw := new(model.OssMainWindow)

	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Visible:  false,
		Title:    "OssTool",
		OnDropFiles: func(files []string) {
			mw.SelectedFile.SetText(files[0])
		},
		Size: Size{
			Width:  500,
			Height: 300,
		},
		MenuItems: []MenuItem{
			Menu{
				Text: "&选项",
				Items: []MenuItem{
					Action{

						Text: "&密钥",
						//Image:    ,
						Enabled: Bind("enabledCB.Checked"),
						Visible: Bind("!openHiddenCB.Checked"),
						//Shortcut:    Shortcut{Modifiers: walk.ModControl, Key: walk.KeyO},
						OnTriggered: mw.OperationSecret,
					},
					Separator{},
					Action{
						Text:        "&退出",
						OnTriggered: func() { mw.Close() },
					},
				},
			},
			/*		Menu{
						Text: "&View",
						Items: []MenuItem{
							Action{
								Text:    "Open / Special Enabled",
								Checked: Bind("enabledCB.Visible"),
							},
							Action{
								Text:    "Open Hidden",
								Checked: Bind("openHiddenCB.Visible"),
							},
						},
					},
					Menu{
						Text: "&Help",
						Items: []MenuItem{
							Action{
								AssignTo: &showAboutBoxAction,
								Text:     "About",
								//OnTriggered: mw.showAboutBoxAction_Triggered,
							},
						},
					},*/
		},
		Layout: VBox{},

		Children: []Widget{
			GroupBox{
				MaxSize: Size{Width: 300, Height: 300},
				Layout:  HBox{},
				Children: []Widget{
					PushButton{ //选择文件
						Text:      "打开文件",
						OnClicked: mw.SelectFile, //点击事件
					},
					Label{Text: "选中的文件 "},
					LineEdit{
						AssignTo: &mw.SelectedFile, //选中的文件
					},
					PushButton{ //上传
						AssignTo:  &mw.PushButton,
						Text:      "上传",
						OnClicked: mw.UploadFile, //上传
					},
				},
			},
			ListBox{ //记录框
				//Enabled:  false,
				MaxSize: Size{
					Width: 300, Height: 200,
				},
				AssignTo: &mw.DisplayListBox,
				//OnCurrentIndexChanged: mw.lb_CurrentIndexChanged, //单击
				//OnItemActivated:       mw.lb_ItemActivated, //双击
			},
			ProgressBar{AssignTo: &mw.ProgressBar},
			/*	Composite{
				Layout: Grid{Columns: 2, Spacing: 10},
				Children: []Widget{

					ListBox{ //记录框
						Enabled: false,
						//AssignTo: &mw.message,
						//OnCurrentIndexChanged: mw.lb_CurrentIndexChanged, //单击
						//OnItemActivated:       mw.lb_ItemActivated, //双击
					},

				},
			},*/
		},
	}.Create()); err != nil {
		panic(err)
	}
	// 去使能最小化按钮
	mw.RemoveStyle(^win.WS_MINIMIZEBOX)
	// 去使能最大化按钮
	mw.RemoveStyle(^win.WS_MAXIMIZEBOX)
	// 去使能关闭按钮
	//hMenu := win.GetSystemMenu(mw.Handle(), false)
	//win.RemoveMenu(hMenu, win.SC_CLOSE, win.MF_BYCOMMAND)
	// 去使能调整大小
	mw.RemoveStyle(^win.WS_SIZEBOX)
	// 去使能位置移动
	//hMenu = win.GetSystemMenu(mw.Handle(), false)
	//win.RemoveMenu(hMenu, win.SC_MOVE, win.MF_BYCOMMAND)
	// 设置窗口居中

	//mw.MainWindow.SetX((int(win.GetSystemMetrics(0)) - mw.MainWindow.Width()/2) / 2)

	//mw.MainWindow.SetY((int(win.GetSystemMetrics(1)) - mw.MainWindow.Height()) / 2)

	mw.OssSecretDialog = &model.OssSecretDialog{}

	return mw
}
