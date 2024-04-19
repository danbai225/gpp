package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/danbai225/gpp/core"
	"github.com/getlantern/elevate"
	box "github.com/sagernet/sing-box"
	_ "image/png"
	"log"
	"os"
	"runtime"
	"sync"
)

func init() {
	switch runtime.GOOS {
	case "windows":
		_ = os.Setenv("FYNE_FONT", "C:\\windows\\Fonts\\simfang.ttf")
	case "darwin":
		_ = os.Setenv("FYNE_FONT", "/System/Library/Fonts/Supplemental/Arial Unicode.ttf")
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
	// 检查是否有管理员权限
	// 如果没有管理员权限则重新启动程序
	// 如果有则继续运行
	if len(os.Args) == 1 {
		cmd := elevate.Command(os.Args[0], "sudo")
		// 开始运行
		_ = cmd.Run()
		// 结束
		os.Exit(0)
	}
	home, _ := os.UserHomeDir()
	config := core.Config{}
	path := "config.json"
	bytes, err := os.ReadFile(path)
	if err != nil {
		path = fmt.Sprintf("%s%c%s%c%s", home, os.PathSeparator, ".gpp", os.PathSeparator, "config.json")
		bytes, err = os.ReadFile(path)
		if err != nil {
			log.Println("读取配置文件失败:", err)
			return
		}
	}
	_ = json.Unmarshal(bytes, &config)
	mainApp := App{
		App:    app.New(),
		logo:   fyne.NewStaticResource("logo.png", logo),
		config: config,
	}
	mainApp.SetIcon(mainApp.logo)
	if desk, ok := mainApp.App.(desktop.App); ok {
		m := fyne.NewMenu("GPP", fyne.NewMenuItem("显示", func() { mainApp.window.Show() }))
		desk.SetSystemTrayMenu(m)
		desk.SetSystemTrayIcon(mainApp.logo)
	}
	mainApp.btn = widget.NewButton("加速", mainApp.Switch)
	mainApp.window = mainApp.NewWindow("GPP加速器")
	mainApp.window.SetCloseIntercept(func() { mainApp.window.Hide() })
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
	log.Println("退出程序")
}
func (a *App) Switch() {
	a.lock.Lock()
	defer a.lock.Unlock()
	log.Println("切换加速状态:", !a.run)
	if a.run {
		err := a.b.Close()
		if err != nil {
			log.Println(err)
		}
		a.btn.SetText("加速")
		a.run = false
	} else {
		client, err := core.Client(a.config)
		if err != nil {
			log.Println(err)
			return
		}
		a.b = client
		a.btn.SetText("加速中")
		a.tip.Show()
		defer a.tip.Hide()
		err = a.b.Start()
		if err != nil {
			a.btn.SetText("加速失败,重新加速")
			log.Println(err)
			return
		}
		a.run = true
		a.btn.SetText("停止")
	}
	log.Println("加速状态:", a.run)
}
