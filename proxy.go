package spaceship_aar

import (
	"encoding/json"
	"github.com/SuzukiHonoka/spaceship/api"
	"github.com/SuzukiHonoka/spaceship/pkg/config"
	"github.com/SuzukiHonoka/spaceship/pkg/config/client"
	"github.com/SuzukiHonoka/spaceship/pkg/config/server"
	"github.com/SuzukiHonoka/spaceship/pkg/dns"
	"github.com/SuzukiHonoka/spaceship/pkg/logger"
)

func Launch(s string) {
	var cfg Config
	_ = json.Unmarshal([]byte(s), &cfg)
	c := &config.MixedConfig{
		Role: config.RoleClient,
		DNS: &dns.DNS{
			Server: cfg.DNS,
		},
		CAs: []string{
			cfg.CA,
		},
		LogMode: logger.Mode(cfg.Log),
		Client: client.Client{
			ServerAddr:  cfg.ServerAddr,
			Host:        cfg.Host,
			UUID:        cfg.Uuid,
			ListenSocks: cfg.ListenSocks,
			ListenHttp:  cfg.ListenHttp,
			Mux:         uint8(cfg.Mux),
			EnableTLS:   cfg.Tls,
		},
		Server: server.Server{
			Path:   cfg.Path,
			Buffer: uint16(cfg.Buffer),
		},
	}
	api.Launch(c)
}

func LaunchFromFile(path string) {
	api.LaunchFromFile(path)
}

func LaunchFromString(c string) {
	api.LaunchFromString(c)
}

func Stop() {
	api.Stop()
}
