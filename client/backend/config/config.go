package config

type Peer struct {
	Protocol string `json:"protocol"`
	Port     uint16 `json:"port"`
	Addr     string `json:"addr"`
	UUID     string `json:"uuid"`
}
type Config struct {
	GamePeers *Peer `json:"game_peers"`
	HTTP      *Peer `json:"http"`
}
