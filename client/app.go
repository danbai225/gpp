package main

import (
	"client/backend/config"
	"context"
	"fmt"
	"github.com/getlantern/elevate"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"os/exec"
)

var home, _ = os.UserHomeDir()
var boxPath = fmt.Sprintf("%s%c%s%c%s", home, os.PathSeparator, ".gpp", os.PathSeparator, "box.exe")

// App struct
type App struct {
	ctx      context.Context
	conf     *config.Config
	gamePeer *config.Peer
	httpPeer *config.Peer
	boxCmd   *exec.Cmd
}

// NewApp creates a new App application struct
func NewApp() *App {
	conf := config.Config{}
	return &App{
		conf: &conf,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	loadConfig, err := config.LoadConfig()
	if err != nil {
		_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.WarningDialog,
			Title:   "配置加载错误",
			Message: err.Error(),
		})
	}
	a.conf = loadConfig
	a.gamePeer = a.conf.PeerList[0]
	a.httpPeer = a.conf.PeerList[0]
	a.Start()
}

// Start 启动加速
func (a *App) Start() string {
	if a.boxCmd != nil && a.boxCmd.ProcessState != nil {
		return "running"
	}
	a.boxCmd = elevate.Command(boxPath)
	err := a.boxCmd.Start()
	if err != nil {
		return err.Error()
	}
	go func() {
		_ = a.boxCmd.Wait()
	}()
	return ""
}

// Stop 停止加速
func (a *App) Stop() string {
	if a.boxCmd == nil {
		return "not running"
	}
	err := a.boxCmd.Process.Kill()
	if err != nil {
		return err.Error()
	}
	return ""
}

func DownloadAssets() string {
	return ""
}
