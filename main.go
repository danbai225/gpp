package main

import (
	_ "embed"
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	logs "github.com/danbai225/go-logs"
	"github.com/danbai225/gpp/core"
	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/constant/goos"
	_ "image/png"
	"os"
	"sync"
)

func init() {
	switch {
	case goos.IsWindows == 1:
		_ = os.Setenv("FYNE_FONT", "C:\\windows\\Fonts\\simfang.ttf")
	}
}

//go:embed logo.png
var logo []byte

type App struct {
	fyne.App
	btn    *widget.Button
	run    bool
	b      *box.Box
	window fyne.Window
	lock   sync.Mutex
	tip    *widget.Label
	logo   fyne.Resource
	config core.Config
}

func main() {
	config := core.Config{}
	var bytes []byte
	if len(os.Args) < 2 {
		bytes, _ = os.ReadFile("config.json")
	} else {
		bytes, _ = os.ReadFile(os.Args[1])
	}
	_ = json.Unmarshal(bytes, &config)
	mainApp := App{
		App:    app.New(),
		logo:   fyne.NewStaticResource("logo.png", logo),
		config: config,
	}
	mainApp.SetIcon(mainApp.logo)
	if desk, ok := mainApp.App.(desktop.App); ok {
		m := fyne.NewMenu("GPP")
		desk.SetSystemTrayMenu(m)
		desk.SetSystemTrayIcon(mainApp.logo)
	}
	mainApp.btn = widget.NewButton("加速", mainApp.Switch)
	mainApp.window = mainApp.NewWindow("GPP加速器")
	h0 := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), widget.NewLabel("欢迎使用GPP加速器，本加速器内测免费使用。"), layout.NewSpacer())
	mainApp.tip = widget.NewLabel("第一次使用加速器需要下载资源加速过程需要1-2分钟。")
	h1 := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), mainApp.tip, layout.NewSpacer())
	mainApp.tip.Hide()
	h2 := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), mainApp.btn, layout.NewSpacer())
	mainApp.window.SetContent(container.New(layout.NewVBoxLayout(), h0, h1, h2))
	mainApp.window.Resize(fyne.NewSize(360, 200))
	mainApp.window.SetFixedSize(true)
	mainApp.window.ShowAndRun()
	if mainApp.b != nil {
		_ = mainApp.b.Close()
	}
	logs.Info("退出程序")
}
func (a *App) Switch() {
	a.lock.Lock()
	defer a.lock.Unlock()
	logs.Info("切换加速状态:", !a.run)
	if a.run {
		err := a.b.Close()
		if err != nil {
			logs.Err(err)
		}
		a.btn.SetText("加速")
		a.run = false
	} else {
		client, err := core.Client(a.config)
		if err != nil {
			logs.Err(err)
			return
		}
		a.b = client
		a.btn.SetText("加速中")
		a.tip.Show()
		defer a.tip.Hide()
		err = a.b.Start()
		if err != nil {
			a.btn.SetText("加速失败,重新加速")
			logs.Err(err)
			return
		}
		a.run = true
		a.btn.SetText("停止")
	}
	logs.Info("加速状态:", a.run)
}
