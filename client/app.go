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

type Status struct {
	Running  bool
	GamePeer *config.Peer
	HttpPeer *config.Peer
}

func (a *App) Status() *Status {
	return &Status{
		Running:  a.box != nil,
		GamePeer: a.gamePeer,
		HttpPeer: a.httpPeer,
	}
}

func (a *App) List() []*config.Peer {
	return a.conf.PeerList
}
func (a *App) Add(token string) string {
	if a.conf.PeerList == nil {
		a.conf.PeerList = make([]*config.Peer, 0)
	}
	err, peer := config.ParsePeer(token)
	if err != nil {
		return err.Error()
	}
	a.conf.PeerList = append(a.conf.PeerList, peer)
	err = config.SaveConfig(a.conf)
	if err != nil {
		return err.Error()
	}
	return "ok"
}
func (a *App) SetPeer(game, http string) string {
	for _, peer := range a.conf.PeerList {
		if peer.Name == game {
			a.gamePeer = peer
			break
		}
	}
	for _, peer := range a.conf.PeerList {
		if peer.Name == http {
			a.httpPeer = peer
			break
		}
	}
	return "ok"
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
