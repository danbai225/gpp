package main

import (
	"client/backend/client"
	"client/backend/config"
	"context"
	"fmt"
	"github.com/cloverstd/tcping/ping"
	box "github.com/sagernet/sing-box"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"sync"
	"time"
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
	app := App{
		conf: &conf,
	}
	go app.testPing()
	return &app
}
func (a *App) testPing() {
	a.PingAll()
	tick := time.Tick(time.Second * 30)
	for range tick {
		a.PingAll()
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
func (a *App) PingAll() {
	group := sync.WaitGroup{}
	for i := range a.conf.PeerList {
		group.Add(1)
		peer := a.conf.PeerList[i]
		go func() {
			defer group.Done()
			peer.Ping = pingPort(peer.Addr, peer.Port)
		}()
	}
	group.Wait()
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
	for _, p := range a.conf.PeerList {
		if p.Name == peer.Name {
			return fmt.Sprintf("peer %s already exists", peer.Name)
		}
	}
	a.conf.PeerList = append(a.conf.PeerList, peer)
	err = config.SaveConfig(a.conf)
	if err != nil {
		return err.Error()
	}
	return "ok"
}
func (a *App) Del(Name string) string {
	for i, peer := range a.conf.PeerList {
		if peer.Name == Name {
			a.conf.PeerList = append(a.conf.PeerList[:i], a.conf.PeerList[i+1:]...)
			break
		}
	}
	err := config.SaveConfig(a.conf)
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
func pingPort(host string, port uint16) uint {
	tcPing := ping.NewTCPing()
	tcPing.SetTarget(&ping.Target{
		Host:     host,
		Port:     int(port),
		Counter:  3,
		Interval: 1,
		Timeout:  time.Second,
	})
	start := tcPing.Start()
	<-start
	result := tcPing.Result()
	return uint(result.Avg().Milliseconds())
}
