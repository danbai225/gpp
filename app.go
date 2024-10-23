package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/cloverstd/tcping/ping"
	"github.com/danbai225/gpp/backend/client"
	"github.com/danbai225/gpp/backend/config"
	"github.com/danbai225/gpp/backend/data"
	"github.com/danbai225/gpp/systray"
	box "github.com/sagernet/sing-box"
	netutils "github.com/shirou/gopsutil/v3/net"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
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
	lock     sync.Mutex
}

// NewApp creates a new App application struct
func NewApp() *App {
	conf := config.Config{}
	app := App{
		conf: &conf,
	}
	return &app
}
func (a *App) systemTray() {
	systray.SetIcon(logo) // read the icon from a file
	show := systray.AddMenuItem("显示窗口", "显示窗口")
	systray.AddSeparator()
	exit := systray.AddMenuItem("退出加速器", "退出加速器")
	show.Click(func() { runtime.WindowShow(a.ctx) })
	exit.Click(func() {
		a.Stop()
		runtime.Quit(a.ctx)
		systray.Quit()
		time.Sleep(time.Second)
		os.Exit(0)
	})
	systray.SetOnClick(func(menu systray.IMenu) { runtime.WindowShow(a.ctx) })
	go func() {
		listener, err := net.Listen("tcp", "127.0.0.1:54713")
		if err != nil {
			_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
				Type:    runtime.ErrorDialog,
				Title:   "监听错误",
				Message: fmt.Sprintln("Error listening0:", err),
			})
		}
		var conn net.Conn
		for {
			conn, err = listener.Accept()
			if err != nil {
				_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
					Type:    runtime.ErrorDialog,
					Title:   "监听错误",
					Message: fmt.Sprintln("Error listening1:", err),
				})
				continue
			}
			// 读取指令
			buffer := make([]byte, 1024)
			n, err := conn.Read(buffer)
			if err != nil {
				_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
					Type:    runtime.ErrorDialog,
					Title:   "监听错误",
					Message: fmt.Sprintln("Error read:", err),
				})
				continue
			}
			command := string(buffer[:n])
			// 如果收到显示窗口的命令，则显示窗口
			if command == "SHOW_WINDOW" {
				// 展示窗口的代码
				runtime.WindowShow(a.ctx)
			}
			_ = conn.Close()
		}
	}()
}

func (a *App) testPing() {
	a.PingAll()
	tick := time.Tick(time.Second * 5)
	for range tick {
		a.PingAll()
	}
}
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	go systray.Run(a.systemTray, func() {})
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
		rdata, err2 := httpGet("https://ghp.ci/https://github.com/SagerNet/sing-geoip/releases/latest/download/geoip.db.sha256sum")
		if err2 == nil {
			sum256 := sha256.Sum256(file)
			if fmt.Sprintf("%x", sum256) != string(rdata) {
				rdata, err2 = httpGet("https://ghp.ci/https://github.com/SagerNet/sing-geoip/releases/latest/download/geoip.db")
				if err2 == nil {
					_ = os.WriteFile(geoPath, rdata, 0o644)
				}
			}
		}
	} else {
		rdata, err2 := httpGet("https://ghp.ci/https://github.com/SagerNet/sing-geoip/releases/latest/download/geoip.db")
		if err2 == nil {
			_ = os.WriteFile(geoPath, rdata, 0o644)
		}
	}
	geoPath = fmt.Sprintf("%s%c%s%c%s", home, os.PathSeparator, ".gpp", os.PathSeparator, "geosite.db")
	file, err = os.ReadFile(geoPath)
	if err == nil {
		rdata, err2 := httpGet("https://ghp.ci/https://github.com/SagerNet/sing-geosite/releases/latest/download/geosite.db.sha256sum")
		if err2 == nil {
			sum256 := sha256.Sum256(file)
			if fmt.Sprintf("%x", sum256) != string(rdata) {
				rdata, err2 = httpGet("https://ghp.ci/https://github.com/SagerNet/sing-geosite/releases/latest/download/geosite.db")
				if err2 == nil {
					_ = os.WriteFile(geoPath, rdata, 0o644)
				}
			}
		}
	} else {
		rdata, err2 := httpGet("https://ghp.ci/https://github.com/SagerNet/sing-geosite/releases/latest/download/geosite.db")
		if err2 == nil {
			_ = os.WriteFile(geoPath, rdata, 0o644)
		}
	}
}
func (a *App) PingAll() {
	a.lock.Lock()
	if a.box != nil {
		a.lock.Unlock()
		return
	}
	a.lock.Unlock()
	group := sync.WaitGroup{}
	for i := range a.conf.PeerList {
		if a.conf.PeerList[i].Protocol == "direct" {
			continue
		}
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
	a.lock.Lock()
	defer a.lock.Unlock()
	status := data.Status{
		Running:  a.box != nil,
		GamePeer: a.gamePeer,
		HttpPeer: a.httpPeer,
	}

	counters, _ := netutils.IOCounters(true)
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
	if strings.HasPrefix(token, "http") {
		_, err := http.Get(token)
		if err != nil {
			_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
				Type:    runtime.ErrorDialog,
				Title:   "订阅错误",
				Message: err.Error(),
			})
			return err.Error()
		}
		a.conf.SubAddr = token
	} else {
		err, peer := config.ParsePeer(token)
		if err != nil {
			_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
				Type:    runtime.ErrorDialog,
				Title:   "导入错误",
				Message: err.Error(),
			})
			return err.Error()
		}
		for _, p := range a.conf.PeerList {
			if p.Name == peer.Name {
				_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
					Type:    runtime.ErrorDialog,
					Title:   "导入错误",
					Message: fmt.Sprintf("节点 %s 已存在", peer.Name),
				})
				return fmt.Sprintf("peer %s already exists", peer.Name)
			}
		}
		a.conf.PeerList = append(a.conf.PeerList, peer)
	}
	err := config.SaveConfig(a.conf)
	if err != nil {
		_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "导入错误",
			Message: err.Error(),
		})
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
	a.lock.Lock()
	defer a.lock.Unlock()
	if a.box != nil {
		return "running"
	}
	var err error
	a.box, err = client.Client(a.gamePeer, a.httpPeer, a.conf.ProxyRule, a.conf.DirectRule, a.conf.ProxyDNS, a.conf.LocalDNS)
	if err != nil {
		_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "加速失败",
			Message: err.Error(),
		})
		a.box = nil
		return err.Error()
	}
	err = a.box.Start()
	if err != nil {
		_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "加速失败",
			Message: err.Error(),
		})
		a.box = nil
		return err.Error()
	}
	return "ok"
}

// Stop 停止加速
func (a *App) Stop() string {
	a.lock.Lock()
	defer a.lock.Unlock()
	if a.box == nil {
		return "not running"
	}
	err := a.box.Close()
	if err != nil {
		_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "停止失败",
			Message: err.Error(),
		})
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
		Counter:  1,
		Interval: time.Millisecond * 200,
		Timeout:  time.Second * 3,
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
	defer func() { _ = resp.Body.Close() }()
	return io.ReadAll(resp.Body)
}
