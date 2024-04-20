package config

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
