package core

type Config struct {
	Port uint16 `json:"port"`
	Addr string `json:"addr"`
	UUID string `json:"uuid"`
}
