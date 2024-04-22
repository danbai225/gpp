package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/netip"
	"os"
	"strings"
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
	conf := &Config{PeerList: make([]*Peer, 0)}
	err = json.Unmarshal(file, &conf)
	return conf, err
}
func SaveConfig(config *Config) error {
	home, _ := os.UserHomeDir()
	path := "config.json"
	_, err := os.Stat(path)
	if err != nil {
		path = fmt.Sprintf("%s%c%s%c%s", home, os.PathSeparator, ".gpp", os.PathSeparator, "config.json")
	}
	file, _ := json.Marshal(config)
	return os.WriteFile(path, file, 0644)
}
func ParsePeer(token string) (error, *Peer) {
	tokenBytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return err, nil
	}
	token = string(tokenBytes)
	split := strings.Split(token, "@")
	protocol := strings.ReplaceAll(split[0], "gpp://", "")
	switch protocol {
	case "vless":
	default:
		return fmt.Errorf("unknown protocol: %s", protocol), nil
	}
	if len(split) != 2 {
		return fmt.Errorf("invalid token: %s", token), nil
	}
	split = strings.Split(split[1], "/")
	addr, err := netip.ParseAddrPort(split[0])
	if err != nil {
		return err, nil
	}
	if len(split) != 2 {
		return fmt.Errorf("invalid token: %s", token), nil
	}
	uuid := split[1]
	split = strings.Split(token, "#")
	name := fmt.Sprintf("%s:%d", addr.Addr().String(), addr.Port())
	if len(split) == 2 {
		name = split[1]
	}
	return nil, &Peer{
		Name:     name,
		Protocol: protocol,
		Port:     addr.Port(),
		Addr:     addr.Addr().String(),
		UUID:     uuid,
	}
}
