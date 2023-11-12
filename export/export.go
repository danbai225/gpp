package main

import (
	"github.com/danbai225/gpp/core"
)

func init() {
	Server()
}

//export Server
func Server() {
	//bytes, err := os.ReadFile("config.json")
	//if err != nil {
	//	fmt.Println("read config err:", err)
	//}
	config := core.Config{
		Port: 5123,
		Addr: "0.0.0.0",
		UUID: "badb17ef-eb22-4e03-9b17-efeb224e03e7",
	}
	//_ = json.Unmarshal(bytes, &config)
	core.Server(config)
}
func main() {
	config := core.Config{
		Port: 5123,
		Addr: "0.0.0.0",
		UUID: "badb17ef-eb22-4e03-9b17-efeb224e03e7",
	}
	//_ = json.Unmarshal(bytes, &config)
	core.Server(config)
	select {}
}
