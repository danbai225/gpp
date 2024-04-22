package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Peer struct {
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	Port     uint16 `json:"port"`
	Addr     string `json:"addr"`
	UUID     string `json:"uuid"`
}
type Config struct {
	PeerList []*Peer `json:"peer_list"`
	GamePeer string  `json:"game_peer"`
	HTTPPeer string  `json:"http_peer"`
}

func LoadConfig() (*Config, error) {
	home, _ := os.UserHomeDir()
	path := "config.json"
	_, err := os.Stat(path)
	if err != nil {
		path = fmt.Sprintf("%s%c%s%c%s", home, os.PathSeparator, ".gpp", os.PathSeparator, "config.json")
	}
	file, _ := os.ReadFile(path)
	conf := &Config{}
	err = json.Unmarshal(file, &conf)
	return conf, err
}
