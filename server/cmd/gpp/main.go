package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/danbai225/gpp/core"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type ip struct {
	IP string `json:"ip"`
}

func main() {
	path := "config.json"
	home, _ := os.UserHomeDir()
	if len(os.Args) > 1 {
		path = os.Args[1]
	} else {
		_, err := os.Stat(path)
		if err != nil {
			path = fmt.Sprintf("%s%c%s%c%s", home, os.PathSeparator, ".gpp", os.PathSeparator, "config.json")
		}
	}
	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("read config err:", err, path)
		return
	}
	config := core.Peer{}
	_ = json.Unmarshal(bytes, &config)
	if config.Port == 0 {
		config.Port = 34555
	}
	if config.Addr == "" {
		config.Addr = "0.0.0.0"
	}
	if config.UUID == "" {
		config.UUID = uuid.New().String()
	}
	err = core.Server(config)
	if err != nil {
		fmt.Println("run err:", err)
	} else {
		fmt.Println("starting success！！！")
		ipStr := ""
		req, _ := http.NewRequest("GET", "https://api.ip.sb/jsonip", nil)
		req.Header.Set("User-Agent", "gpp")
		resp, err2 := http.DefaultClient.Do(req)
		if err2 != nil {
			fmt.Println("get ip err:", err2)
		} else {
			all, _ := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			i := ip{}
			_ = json.Unmarshal(all, &i)
			ipStr = i.IP
		}
		if ipStr != "" {
			fmt.Println("server ip:", ipStr)
			// 进行Base64编码
			encoded := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("gpp://%s@%s:%d/%s", config.Protocol, ipStr, config.Port, config.UUID)))
			fmt.Println("server token:", encoded)
		} else {
			fmt.Println("get ip err, please check your network")
		}
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		s := <-sigCh
		fmt.Printf("Received signal: %v\n", s)
		fmt.Println("Exiting...")
		os.Exit(0)
	}
}
