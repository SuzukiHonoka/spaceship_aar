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

type LauncherWrapper struct {
	*api.Launcher
}

func NewLauncher() *LauncherWrapper {
	return &LauncherWrapper{
		api.NewLauncher(),
	}
}

func (l *LauncherWrapper) Launch(s string) bool {
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
	return l.Launcher.Launch(c)
}

func (l *LauncherWrapper) LaunchFromFile(path string) bool {
	return l.Launcher.LaunchFromFile(path)
}

func (l *LauncherWrapper) LaunchFromString(c string) bool {
	return l.Launcher.LaunchFromString(c)
}

func (l *LauncherWrapper) Stop() {
	l.Launcher.Stop()
}
