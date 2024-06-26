package config

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Peer struct {
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	Port     uint16 `json:"port"`
	Addr     string `json:"addr"`
	UUID     string `json:"uuid"`
	Ping     uint   `json:"ping"`
}
type Config struct {
	PeerList []*Peer `json:"peer_list"`
	GamePeer string  `json:"game_peer"`
	HTTPPeer string  `json:"http_peer"`
}

func InitConfig() {
	home, _ := os.UserHomeDir()
	_path := "config.json"
	_, err := os.Stat(_path)
	if err != nil {
		_path = fmt.Sprintf("%s%c%s%c%s", home, os.PathSeparator, ".gpp", os.PathSeparator, "config.json")
	}
	_ = os.MkdirAll(filepath.Dir(_path), os.ModeDir|os.ModePerm)
	_, err = os.Stat(_path)
	if err != nil {
		file, _ := json.Marshal(Config{PeerList: make([]*Peer, 0)})
		err = os.WriteFile(_path, file, os.ModePerm)
	}
}
func LoadConfig() (*Config, error) {
	home, _ := os.UserHomeDir()
	_path := "config.json"
	_, err := os.Stat(_path)
	if err != nil {
		_path = fmt.Sprintf("%s%c%s%c%s", home, os.PathSeparator, ".gpp", os.PathSeparator, "config.json")
	}
	file, _ := os.ReadFile(_path)
	conf := &Config{PeerList: make([]*Peer, 0)}
	err = json.Unmarshal(file, &conf)
	var direct bool
	for _, peer := range conf.PeerList {
		if peer.Name == "直连" {
			direct = true
		}
	}
	if !direct {
		conf.PeerList = append(conf.PeerList, &Peer{Name: "直连", Protocol: "direct", Port: 0, Addr: "127.0.0.1", UUID: "", Ping: 0})
	}
	return conf, err
}
func SaveConfig(config *Config) error {
	home, _ := os.UserHomeDir()
	_path := "config.json"
	_, err := os.Stat(_path)
	if err != nil {
		_path = fmt.Sprintf("%s%c%s%c%s", home, os.PathSeparator, ".gpp", os.PathSeparator, "config.json")
	}
	file, _ := json.Marshal(config)
	return os.WriteFile(_path, file, os.ModePerm)
}
func ParsePeer(token string) (error, *Peer) {
	split := strings.Split(token, "#")
	name := ""
	if len(split) == 2 {
		token = split[0]
		name = split[1]
	}
	tokenBytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return err, nil
	}
	token = string(tokenBytes)
	split = strings.Split(token, "@")
	protocol := strings.ReplaceAll(split[0], "gpp://", "")
	switch protocol {
	case "vless", "shadowsocks", "socks":
	default:
		return fmt.Errorf("unknown protocol: %s", protocol), nil
	}
	if len(split) != 2 {
		return fmt.Errorf("invalid token: %s", token), nil
	}
	split = strings.Split(split[1], "/")
	addr := strings.Split(split[0], ":")
	if len(addr) != 2 {
		return errors.New("invalid addr: " + split[0]), nil
	}
	if len(split) != 2 {
		return fmt.Errorf("invalid token: %s", token), nil
	}
	uuid := split[1]
	if name == "" {
		name = fmt.Sprintf("%s:%s", addr[0], addr[1])
	}
	port, _ := strconv.ParseInt(addr[1], 10, 64)
	return nil, &Peer{
		Name:     name,
		Protocol: protocol,
		Port:     uint16(port),
		Addr:     addr[0],
		UUID:     uuid,
	}
}
