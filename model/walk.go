package model

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"log"
	"os"
	"path/filepath"
	"time"
)

var AppMainWindow *OssMainWindow

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
	return m
}

func (m *EnvModel) ItemCount() int {
	return len(m.Items)
}

func (m *EnvModel) Value(index int) interface{} {
	return m.Items[index].Value
}

type OssMainWindow struct {
	*walk.MainWindow
	App             *walk.Application
	NotifyIcon      *walk.NotifyIcon
	OssSecretDialog *OssSecretDialog
	UploadListBox   *walk.ListBox
	DisplayListBox  *walk.ListBox
	SelectedFile    *walk.LineEdit
	EnvModel        *EnvModel
	PushButton      *walk.PushButton
	ProgressBar     *walk.ProgressBar
}

type OssSecretDialog struct {
	*walk.Dialog
	Endpoint *walk.LineEdit //用户名
	Bucket   *walk.LineEdit //密码
	AK       *walk.LineEdit
	SK       *walk.LineEdit
}

func (mw *OssMainWindow) saveDiaLogInfo() {
	if mw.OssSecretDialog.Endpoint.Text() == "" || mw.OssSecretDialog.Bucket.Text() == "" || mw.OssSecretDialog.AK.Text() == "" || mw.OssSecretDialog.SK.Text() == "" {
		walk.MsgBox(mw, "错误", "请填写完整信息", walk.MsgBoxIconError)
		return
	}
	endpointConfig := GetGlobalEndpointConfig()
	endpointConfig.Endpoint = mw.OssSecretDialog.Endpoint.Text()
	endpointConfig.Bucket = mw.OssSecretDialog.Bucket.Text()
	endpointConfig.AccessKey = mw.OssSecretDialog.AK.Text()
	endpointConfig.SecretKey = mw.OssSecretDialog.SK.Text()
	mw.OssSecretDialog.Close(0)

}

func (mw *OssMainWindow) SelectFile() {

	dlg := new(walk.FileDialog)
	dlg.Title = "选择文件"
	if _, err := dlg.ShowOpen(mw); err != nil {
		panic(err)
	}
	path := dlg.FilePath

	mw.SelectedFile.SetText(path)
	endpointConfig := GetGlobalEndpointConfig()
	_, items, _ := endpointConfig.CreateOssClient()
	model := &EnvModel{Items: items}
	mw.DisplayListBox.SetModel(model)

}

func (mw *OssMainWindow) UploadFile() {

	endpointConfig := GetGlobalEndpointConfig()
	if mw.SelectedFile.Text() == "" {
		walk.MsgBox(mw, "错误", "请选择文件", walk.MsgBoxIconError)
		return
	}
	if endpointConfig.Endpoint == "" || endpointConfig.Bucket == "" || endpointConfig.AccessKey == "" || endpointConfig.SecretKey == "" {
		walk.MsgBox(mw, "错误", "配置文件有误,请检查!", walk.MsgBoxIconError)
		return
	}
	//等待上传
	client, items, err := endpointConfig.CreateOssClient()
	if err != nil {
		walk.MsgBox(mw, "失败", "初始化配置失败,请检查配置!", walk.MsgBoxIconError)
		return
	}
	bucket, err := client.Bucket(endpointConfig.Bucket)
	if err != nil {
		walk.MsgBox(mw, "失败", "获取存储桶失败!", walk.MsgBoxIconError)
		return
	}
	mw.PushButton.SetEnabled(false)
	go func() {

		base := filepath.Base(mw.SelectedFile.Text())
		items = append(items, EnvItem{Name: "", Value: fmt.Sprintf("上传文件名:\r\n%s", base)})
		model := &EnvModel{Items: items}
		mw.DisplayListBox.SetModel(model)
		err = bucket.PutObjectFromFile(base, mw.SelectedFile.Text(), oss.Progress(&OssProgressListener{ProgressBar: mw.ProgressBar}))
		if err != nil {
			walk.MsgBox(mw, "失败", "上传失败!", walk.MsgBoxIconError)
			return
		} else {
			walk.MsgBox(mw, "成功", "上传成功!", walk.MsgBoxOK)

		}
		defer mw.PushButton.SetEnabled(true)
		mw.ProgressBar.SetValue(0)
	}()

}

func (mw *OssMainWindow) RemoveStyle(style int32) {
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
		endpointConfig := GetGlobalEndpointConfig()
		mw.App.Settings().Put("endpoint", endpointConfig.Endpoint)
		mw.App.Settings().Put("bucket", endpointConfig.Bucket)
		mw.App.Settings().Put("accessKey", endpointConfig.AccessKey)
		mw.App.Settings().Put("secretKey", endpointConfig.SecretKey)
		mw.App.Settings().Save()
		os.Exit(0)
	})
	if err := ni.ContextMenu().Actions().Add(exitAction); err != nil {
		log.Fatal(err)
	}
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
					PushButton{Text: "保存",
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
	endpointConfig := GetGlobalEndpointConfig()
	secretDialog.Endpoint.SetText(endpointConfig.Endpoint)
	secretDialog.Bucket.SetText(endpointConfig.Bucket)
	secretDialog.AK.SetText(endpointConfig.AccessKey)
	secretDialog.SK.SetText(endpointConfig.SecretKey)
	secretDialog.SetIcon(mw.Icon())
	secretDialog.Run()
}
