package desktop

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"log"
	"os"
	"ossTool/config"
	"time"
)

type OssMainWindow struct {
	*walk.MainWindow
	App             *walk.Application
	NotifyIcon      *walk.NotifyIcon
	OssSecretDialog *OssSecretDialog
	EndpointConfig  *config.EndpointConfig
	UploadListBox   *walk.ListBox
	DisplayListBox  *walk.ListBox
	SelectedFile    *walk.LineEdit
	EnvModel        *EnvModel
}

type OssSecretDialog struct {
	*walk.Dialog
	Endpoint *walk.LineEdit //用户名
	Bucket   *walk.LineEdit //密码
	AK       *walk.LineEdit
	SK       *walk.LineEdit
}

type EnvItem struct {
	Name  string
	Value string
}

type EnvModel struct {
	walk.ListModelBase
	Items []EnvItem
}

func NewEnvModel() *EnvModel {
	m := &EnvModel{Items: make([]EnvItem, 0)}
	/*
		m.items[0] = EnvItem{"name", "value"}
		m.items[1] = EnvItem{"name", "value"}
		m.items[2] = EnvItem{"name", "value"}*/
	return m
}

func (m *EnvModel) ItemCount() int {
	return len(m.Items)
}

func (m *EnvModel) Value(index int) interface{} {
	return m.Items[index].Value
}

func InitMainWindow() *OssMainWindow {

	mw := new(OssMainWindow)
	var openAction, showAboutBoxAction *walk.Action
	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Visible:  false,
		Title:    "OssTool",
		OnDropFiles: func(files []string) {
			mw.SelectedFile.SetText(files[0])
		},
		Size: Size{
			Width:  500,
			Height: 600,
		},
		MenuItems: []MenuItem{
			Menu{
				Text: "&选项",
				Items: []MenuItem{
					Action{
						AssignTo: &openAction,
						Text:     "&密钥",
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
			Menu{
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
			},
		},
		Layout: VBox{},

		Children: []Widget{
			GroupBox{
				MaxSize: Size{Width: 500, Height: 500},
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
						Text:      "上传",
						OnClicked: mw.UploadFile, //上传
					},
				},
			},
			ListBox{ //记录框
				Enabled:  false,
				AssignTo: &mw.DisplayListBox,
				//OnCurrentIndexChanged: mw.lb_CurrentIndexChanged, //单击
				//OnItemActivated:       mw.lb_ItemActivated, //双击
			},
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
	mw.removeStyle(^win.WS_MINIMIZEBOX)
	// 去使能最大化按钮
	mw.removeStyle(^win.WS_MAXIMIZEBOX)
	// 去使能关闭按钮
	hMenu := win.GetSystemMenu(mw.Handle(), false)
	win.RemoveMenu(hMenu, win.SC_CLOSE, win.MF_BYCOMMAND)
	// 去使能调整大小
	mw.removeStyle(^win.WS_SIZEBOX)
	// 去使能位置移动
	//hMenu = win.GetSystemMenu(mw.Handle(), false)
	//win.RemoveMenu(hMenu, win.SC_MOVE, win.MF_BYCOMMAND)
	// 设置窗口居中

	//mw.MainWindow.SetX((int(win.GetSystemMetrics(0)) - mw.MainWindow.Width()/2) / 2)

	//mw.MainWindow.SetY((int(win.GetSystemMetrics(1)) - mw.MainWindow.Height()) / 2)
	mw.OssSecretDialog = &OssSecretDialog{}

	return mw
}

func (mw *OssMainWindow) NewNotifyIcon() {
	notifyIcon, err := walk.NewNotifyIcon(mw)
	if err != nil {
		panic(err)
	}
	mw.NotifyIcon = notifyIcon

	icon, err := walk.Resources.Icon("tray.ico")
	if err != nil {
		panic(err)
	}
	err = mw.SetIcon(icon)
	if err != nil {
		panic(err)
	}

	if err := mw.NotifyIcon.SetIcon(icon); err != nil {
		panic(err)
	}
	if err := mw.NotifyIcon.SetToolTip("点击使用更多功能"); err != nil {
		panic(err)
	}

	mw.NotifyIcon.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button != walk.LeftButton {
			return
		}
		//此处左键点击托盘复原主界面
		if visible := mw.Visible(); visible {
			mw.Hide()
		} else {
			mw.Show()

		}
	})
	addHelpMenu(mw, mw.NotifyIcon)
	addExitMenu(mw, mw.NotifyIcon)

	if err := mw.NotifyIcon.SetVisible(true); err != nil {
		log.Fatal(err)
	}

	// Now that the icon is visible, we can bring up an info balloon.
	/*	if err := mw.NotifyIcon.ShowInfo("右下角弹窗标题", "右下角弹窗内容哦"); err != nil {
		log.Fatal(err)
	}*/

	if err := mw.NotifyIcon.ShowCustom(
		"OssTool",
		"使用GO语言开发的数据桶上传插件",
		icon); err != nil {
		log.Fatal(err)
	}

}

func (mw *OssMainWindow) InitAppSetting() {
	app := walk.App()
	// These specify the app data sub directory for the settings file.
	app.SetOrganizationName("MappedByte")
	app.SetProductName("OssTool Settings")

	// Settings file name.
	settings := walk.NewIniFileSettings("settings.ini")
	settings.SetExpireDuration(time.Hour * 24 * 30 * 3)
	if err := settings.Load(); err != nil {
		panic(err)
	}
	app.SetSettings(settings)
	mw.App = app
}

func (mw *OssMainWindow) OperationSecret() {
	secretDialog := mw.OssSecretDialog

	err := Dialog{
		AssignTo: &secretDialog.Dialog,
		Title:    "密钥",
		MaxSize:  Size{Width: 300, Height: 200},
		Size: Size{
			Width:  300,
			Height: 200,
		},
		Layout: VBox{},
		Name:   "Dialog",
		Children: []Widget{
			GroupBox{
				MinSize: Size{Width: 300, Height: 200},
				Layout:  Grid{Columns: 5},
				Children: []Widget{
					Label{
						//Font:   Font{Family: "Consolas", PointSize: 13},
						Text:   "Endpoint:",
						Row:    1,
						Column: 1,
					},

					LineEdit{
						AssignTo: &secretDialog.Endpoint,
						//Font: Font{Family: "Consolas", PointSize: 13},
						// MinSize:    Size{90, 30},不好用，和默认值冲突以默认值为准
						Row:         1,
						Column:      3,
						ToolTipText: "填写端点(目前仅支持和阿里云)",
						// ColumnSpan: 2,
					},

					Label{
						//Font:   Font{Family: "Consolas", PointSize: 13},
						Text:   "Bucket:",
						Row:    2,
						Column: 1,
					},

					LineEdit{
						AssignTo: &secretDialog.Bucket,
						//Font: Font{Family: "Consolas", PointSize: 13},
						// MinSize:    Size{90, 30},不好用，和默认值冲突以默认值为准
						Row:         2,
						Column:      3,
						ToolTipText: "存储桶名称",
						// ColumnSpan: 2,
					},

					Label{
						//Font:   Font{Family: "Consolas", PointSize: 13},
						Text:   "AK:",
						Row:    3,
						Column: 1,
					},
					LineEdit{
						AssignTo: &secretDialog.AK,
						//Font:         Font{Family: "Consolas", PointSize: 13},
						PasswordMode: true, //密码用密文
						Row:          3,
						Column:       3,
						ToolTipText:  "AK",
						// ColumnSpan: 2,
					},

					Label{
						//Font:   Font{Family: "Consolas", PointSize: 13},
						Text:   "SK:",
						Row:    4,
						Column: 1,
					},
					LineEdit{
						AssignTo: &secretDialog.SK,
						//Font:         Font{Family: "Consolas", PointSize: 13},
						PasswordMode: true, //密码用密文
						Row:          4,
						Column:       3,
						ToolTipText:  "SK",
						// ColumnSpan: 2,
					},
					PushButton{Text: "登录",
						Font:      Font{Family: "Consolas", PointSize: 11},
						OnClicked: mw.saveDiaLogInfo,
						Row:       5,
						Column:    3,
						// ColumnSpan: 3,
					},
				},
			}},
	}.Create(mw)
	currStyle := win.GetWindowLong(secretDialog.Handle(), win.GWL_STYLE)
	win.SetWindowLong(secretDialog.Handle(), win.GWL_STYLE, currStyle&^win.WS_SIZEBOX)
	if err != nil {
		panic(err)
	}
	//回填内容
	secretDialog.Endpoint.SetText(mw.EndpointConfig.Endpoint)
	secretDialog.Bucket.SetText(mw.EndpointConfig.Bucket)
	secretDialog.AK.SetText(mw.EndpointConfig.AccessKey)
	secretDialog.SK.SetText(mw.EndpointConfig.SecretKey)
	secretDialog.SetIcon(mw.Icon())
	secretDialog.Run()
}

func (mw *OssMainWindow) saveDiaLogInfo() {
	mw.EndpointConfig.Endpoint = mw.OssSecretDialog.Endpoint.Text()
	mw.EndpointConfig.Bucket = mw.OssSecretDialog.Bucket.Text()
	mw.EndpointConfig.AccessKey = mw.OssSecretDialog.AK.Text()
	mw.EndpointConfig.SecretKey = mw.OssSecretDialog.SK.Text()
	mw.OssSecretDialog.Close(0)

}

func (mw *OssMainWindow) SelectFile() {

	dlg := new(walk.FileDialog)
	dlg.Title = "选择文件"
	if _, err := dlg.ShowOpen(mw); err != nil {
		panic(err)
	}
	path := dlg.FilePath
	//fmt.Println(path)
	mw.SelectedFile.SetText(path)
}

func (mw *OssMainWindow) UploadFile() {
	if mw.SelectedFile.Text() == "" {
		walk.MsgBox(mw, "错误", "请选择文件", walk.MsgBoxIconError)
		return
	}

}

func (mw *OssMainWindow) removeStyle(style int32) {
	currStyle := win.GetWindowLong(mw.Handle(), win.GWL_STYLE)
	win.SetWindowLong(mw.Handle(), win.GWL_STYLE, currStyle&style)
}

func addHelpMenu(mw *OssMainWindow, ni *walk.NotifyIcon) {
	helpAction := walk.NewAction()
	if err := helpAction.SetText("帮助"); err != nil {
		log.Fatal(err)
	}
	helpAction.Triggered().Attach(func() {
		walk.MsgBox(mw, "帮助", "本软件使用GO语言进行开发,如有问题请联系作者!", walk.MsgBoxIconInformation)
	})
	if err := ni.ContextMenu().Actions().Add(helpAction); err != nil {
		log.Fatal(err)
	}
}

func addExitMenu(mw *OssMainWindow, ni *walk.NotifyIcon) {
	// We put an exit action into the context menu.
	exitAction := walk.NewAction()
	if err := exitAction.SetText("退出"); err != nil {
		log.Fatal(err)
	}
	exitAction.Triggered().Attach(func() {
		walk.App().Exit(0)
		ni.Dispose()
		mw.App.Settings().Put("endpoint", mw.EndpointConfig.Endpoint)
		mw.App.Settings().Put("bucket", mw.EndpointConfig.Bucket)
		mw.App.Settings().Put("accessKey", mw.EndpointConfig.AccessKey)
		mw.App.Settings().Put("secretKey", mw.EndpointConfig.SecretKey)
		mw.App.Settings().Save()
		os.Exit(0)
	})
	if err := ni.ContextMenu().Actions().Add(exitAction); err != nil {
		log.Fatal(err)
	}
}
