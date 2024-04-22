package main

import (
	"client/backend/client"
	"client/backend/config"
	"context"
	box "github.com/sagernet/sing-box"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx      context.Context
	conf     *config.Config
	gamePeer *config.Peer
	httpPeer *config.Peer
	box      *box.Box
}

// NewApp creates a new App application struct
func NewApp() *App {
	conf := config.Config{}
	return &App{
		conf: &conf,
	}
}

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
}

// Start 启动加速
func (a *App) Start() string {
	if a.box != nil {
		return "running"
	}
	var err error
	a.box, err = client.Client(a.gamePeer, a.httpPeer)
	if err != nil {
		return err.Error()
	}
	err = a.box.Start()
	if err != nil {
		return err.Error()
	}
	return "ok"
}

// Stop 停止加速
func (a *App) Stop() string {
	if a.box == nil {
		return "not running"
	}
	err := a.box.Close()
	if err != nil {
		return err.Error()
	}
	a.box = nil
	return "ok"
}
