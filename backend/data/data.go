package data

import "github.com/danbai225/gpp/backend/config"

type Status struct {
	Running  bool         `json:"running"`
	GamePeer *config.Peer `json:"game_peer"`
	HttpPeer *config.Peer `json:"http_peer"`
	Up       uint64       `json:"up"`
	Down     uint64       `json:"down"`
}
