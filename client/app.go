package main

import (
	"client/backend/config"
	"client/backend/core"
	"context"
	"encoding/json"
	"fmt"
	box "github.com/sagernet/sing-box"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
)

// App struct
type App struct {
	ctx      context.Context
	conf     config.Config
	box      *box.Box
	gamePeer *config.Peer
	httpPeer *config.Peer
	isRun    bool
}

// NewApp creates a new App application struct
func NewApp() *App {
	conf := config.Config{}
	return &App{
		conf: conf,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	home, _ := os.UserHomeDir()
	path := "config.json"
	_, err := os.Stat(path)
	if err != nil {
		path = fmt.Sprintf("%s%c%s%c%s", home, os.PathSeparator, ".gpp", os.PathSeparator, "config.json")
	}
	runtime.LogDebugf(a.ctx, "config path: %s", path)
	file, _ := os.ReadFile(path)
	err = json.Unmarshal(file, &a.conf)
	if err != nil || len(a.conf.PeerList) == 0 {
		_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.WarningDialog,
			Title:   "配置加载错误",
			Message: "",
		})
	}
	a.gamePeer = a.conf.PeerList[0]
	a.httpPeer = a.conf.PeerList[0]
}

// Start 启动加速
func (a *App) Start() string {
	if a.isRun {
		return "running"
	}
	a.isRun = true
	var err error
	a.box, err = core.Client(a.gamePeer, a.httpPeer)
	if err != nil {
		return err.Error()
	}
	err = a.box.Start()
	if err != nil {
		return err.Error()
	}
	return ""
}

// Stop 停止加速
func (a *App) Stop() string {
	if !a.isRun {
		return "not running"
	}
	a.isRun = false
	err := a.box.Close()
	if err != nil {
		return err.Error()
	}
	return ""
}
