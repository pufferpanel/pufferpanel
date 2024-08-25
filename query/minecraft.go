package query

import (
	"fmt"
	"github.com/dreamscached/minequery/v2"
)

type MinecraftResponse struct {
	NumPlayers int      `json:"numPlayers"`
	MaxPlayers int      `json:"maxPlayers"`
	Version    string   `json:"version"`
	Players    []string `json:"players"`
}

func Minecraft(ip string, port int) (MinecraftResponse, error) {
	if port == 0 {
		return MinecraftResponse{}, fmt.Errorf("port is required")
	}
	if ip == "" || ip == "0.0.0.0" {
		ip = "127.0.0.1"
	}

	res, err := minequery.Ping17(ip, port)
	if err != nil {
		return MinecraftResponse{}, err
	}

	var players []string
	for _, v := range res.SamplePlayers {
		players = append(players, v.Nickname)
	}

	return MinecraftResponse{
		NumPlayers: res.OnlinePlayers,
		MaxPlayers: res.MaxPlayers,
		Version:    res.VersionName,
		Players:    players,
	}, nil
}
