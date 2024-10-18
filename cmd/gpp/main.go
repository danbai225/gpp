package main

import (
	"encoding/json"
	"fmt"
	"github.com/danbai225/gpp/server/core"
	"github.com/google/uuid"
	"os"
	"os/signal"
	"syscall"
)

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
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		s := <-sigCh
		fmt.Printf("Received signal: %v\n", s)
		fmt.Println("Exiting...")
		os.Exit(0)
	}
}
