package main

import (
	"encoding/json"
	"fmt"
	"github.com/danbai225/gpp/core"
	box "github.com/sagernet/sing-box"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: gpp [server|client] [config.json]")
		return
	}
	bytes, err := os.ReadFile(os.Args[2])
	if err != nil {
		fmt.Println("read config err:", err)
	}
	config := core.Config{}
	_ = json.Unmarshal(bytes, &config)
	if os.Args[1] == "server" {
		err = core.Server(config)
	} else if os.Args[1] == "client" {
		var box *box.Box
		box, err = core.Client(config)
		if err == nil {
			err = box.Start()
		}
	}
	if err != nil {
		fmt.Println("run err:", err)
	}
}
