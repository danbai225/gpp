package main

import (
	"client/backend/client"
	"client/backend/config"
	"client/backend/data"
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/cloverstd/tcping/ping"
	"github.com/energye/systray"
	box "github.com/sagernet/sing-box"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

// App struct
type App struct {
	ctx         context.Context
	conf        *config.Config
	gamePeer    *config.Peer
	httpPeer    *config.Peer
	box         *box.Box
	processList []string
}

// NewApp creates a new App application struct
func NewApp() *App {
	conf := config.Config{}
	app := App{
		conf:        &conf,
		processList: make([]string, 0),
	}
	return &app
}
func (a *App) systemTray() {
	systray.SetIcon(logo) // read the icon from a file
	show := systray.AddMenuItem("显示窗口", "显示窗口")
	systray.AddSeparator()
	exit := systray.AddMenuItem("退出加速器", "退出加速器")
	show.Click(func() { runtime.WindowShow(a.ctx) })
	exit.Click(func() { os.Exit(0) })
	systray.SetOnClick(func(menu systray.IMenu) { runtime.WindowShow(a.ctx) })
}

func (a *App) testPing() {
	a.PingAll()
	tick := time.Tick(time.Second * 60)
	for range tick {
		a.PingAll()
	}
}
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	systray.Run(a.systemTray, func() {})
	loadConfig, err := config.LoadConfig()
	if err != nil {
		_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.WarningDialog,
			Title:   "配置加载错误",
			Message: err.Error(),
		})
	} else {
		a.conf = loadConfig
	}
	if len(a.conf.PeerList) > 0 {
		a.gamePeer = a.conf.PeerList[0]
		a.httpPeer = a.conf.PeerList[0]
	}
	go a.testPing()
	home, _ := os.UserHomeDir()
	geoPath := fmt.Sprintf("%s%c%s%c%s", home, os.PathSeparator, ".gpp", os.PathSeparator, "geoip.db")
	file, err := os.ReadFile(geoPath)
	if err == nil {
		rdata, err2 := httpGet("https://mirror.ghproxy.com/https://github.com/SagerNet/sing-geoip/releases/latest/download/geoip.db.sha256sum")
		if err2 == nil {
			sum256 := sha256.Sum256(file)
			fmt.Println(fmt.Sprintf("%x", sum256), string(rdata))
			if fmt.Sprintf("%x", sum256) != string(rdata) {
				rdata, err2 = httpGet("https://mirror.ghproxy.com/https://github.com/SagerNet/sing-geoip/releases/latest/download/geoip.db")
				if err2 == nil {
					_ = os.WriteFile(geoPath, rdata, 0644)
				}
			}
		}
	} else {
		rdata, err2 := httpGet("https://mirror.ghproxy.com/https://github.com/SagerNet/sing-geoip/releases/latest/download/geoip.db")
		if err2 == nil {
			_ = os.WriteFile(geoPath, rdata, 0644)
		}
	}
	geoPath = fmt.Sprintf("%s%c%s%c%s", home, os.PathSeparator, ".gpp", os.PathSeparator, "geosite.db")
	file, err = os.ReadFile(geoPath)
	if err == nil {
		rdata, err2 := httpGet("https://mirror.ghproxy.com/https://github.com/SagerNet/sing-geosite/releases/latest/download/geosite.db.sha256sum")
		if err2 == nil {
			sum256 := sha256.Sum256(file)
			fmt.Println(fmt.Sprintf("%x", sum256), string(rdata))
			if fmt.Sprintf("%x", sum256) != string(rdata) {
				rdata, err2 = httpGet("https://mirror.ghproxy.com/https://github.com/SagerNet/sing-geosite/releases/latest/download/geosite.db")
				if err2 == nil {
					_ = os.WriteFile(geoPath, rdata, 0644)
				}
			}
		}
	} else {
		rdata, err2 := httpGet("https://mirror.ghproxy.com/https://github.com/SagerNet/sing-geosite/releases/latest/download/geosite.db")
		if err2 == nil {
			_ = os.WriteFile(geoPath, rdata, 0644)
		}
	}
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

func (a *App) Status() *data.Status {
	status := data.Status{
		Running:  a.box != nil,
		GamePeer: a.gamePeer,
		HttpPeer: a.httpPeer,
	}

	counters, _ := net.IOCounters(true)
	for _, counter := range counters {
		if counter.Name == "utun225" {
			status.Up = counter.BytesSent
			status.Down = counter.BytesRecv
		}
	}
	return &status
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
	a.box, err = client.Client(a.gamePeer, a.httpPeer, a.processList)
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
		Interval: time.Millisecond * 200,
		Timeout:  time.Second,
	})
	start := tcPing.Start()
	<-start
	result := tcPing.Result()
	return uint(result.Avg().Milliseconds())
}
func httpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
