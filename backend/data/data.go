package data

import "client/backend/config"

type Status struct {
	Running  bool         `json:"running"`
	GamePeer *config.Peer `json:"game_peer"`
	HttpPeer *config.Peer `json:"http_peer"`
}
