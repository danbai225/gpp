package core

type Peer struct {
	Port uint16 `json:"port"`
	Addr string `json:"addr"`
	UUID string `json:"uuid"`
}
type Config struct {
	Peer
	Default Peer `json:"default"`
	HTTP    Peer `json:"http"`
}
