package main

import (
	"encoding/json"
	"fmt"
	"github.com/danbai225/gpp/core"
	"os"
)

func init() {
	Server()
}

//export Server
func Server() {
	bytes, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Println("read config err:", err)
	}
	config := core.Config{}
	_ = json.Unmarshal(bytes, &config)
	_ = core.Server(config)
}
func main() {}
