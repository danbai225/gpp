package main

import (
	"encoding/json"
	"fmt"
	"github.com/danbai225/gpp/core"
	box "github.com/sagernet/sing-box"
	"os"
)

func main() {
	path := "config.json"
	home, _ := os.UserHomeDir()
	if len(os.Args) < 2 {
		fmt.Println("Usage: gpp [server|client] [config.json]")
		return
	} else if len(os.Args) > 2 {
		path = os.Args[2]
	} else {
		_, err := os.Stat(path)
		if err != nil {
			path = fmt.Sprintf("%s%c%s%c%s", home, os.PathSeparator, ".gpp", os.PathSeparator, "config.json")
		}
	}
	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("read config err:", err)
	}
	config := core.Config{}
	_ = json.Unmarshal(bytes, &config)
	if os.Args[1] == "server" {
		err = core.Server(config)
	} else if os.Args[1] == "client" {
		var b *box.Box
		b, err = core.Client(config)
		if err == nil {
			err = b.Start()
		}
	}
	if err != nil {
		fmt.Println("run err:", err)
	} else {
		fmt.Println("启动成功！！！")
		select {}
	}
}
