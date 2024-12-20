package spaceship_aar

import (
	"encoding/json"
	"log"

	"github.com/SuzukiHonoka/spaceship/api"
	"github.com/SuzukiHonoka/spaceship/pkg/config"
	"github.com/SuzukiHonoka/spaceship/pkg/config/client"
	"github.com/SuzukiHonoka/spaceship/pkg/config/server"
	"github.com/SuzukiHonoka/spaceship/pkg/dns"
	"github.com/SuzukiHonoka/spaceship/pkg/logger"
)

// LauncherWrapper wraps the api.Launcher and provides some proxy methods.
type LauncherWrapper struct {
	*api.Launcher
}

// NewLauncher creates a new *LauncherWrapper.
func NewLauncher() *LauncherWrapper {
	return &LauncherWrapper{
		api.NewLauncher(),
	}
}

// Launch parses the particular cfg string to the client-oriented configuration and launch the inner launcher.
// Note that this is not native configuration format. if using native, go for LaunchFromString instead.
func (l *LauncherWrapper) Launch(s string) bool {
	var cfg Config
	if err := json.Unmarshal([]byte(s), &cfg); err != nil {
		log.Printf("launch: unmarshal cfg failed, err=%s", err)
		return false
	}
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
			ServerAddr:      cfg.ServerAddr,
			Host:            cfg.Host,
			UUID:            cfg.Uuid,
			ListenSocks:     cfg.ListenSocks,
			ListenSocksUnix: cfg.ListenSocksUnix,
			ListenHttp:      cfg.ListenHttp,
			Mux:             uint8(cfg.Mux),
			EnableTLS:       cfg.Tls,
		},
		Server: server.Server{
			Path:   cfg.Path,
			Buffer: uint16(cfg.Buffer),
			IPv6:   cfg.IPv6,
		},
	}
	return l.Launcher.Launch(c)
}

// LaunchFromFile reads the native configuration string from file and passes to the inner launcher.
func (l *LauncherWrapper) LaunchFromFile(path string) bool {
	return l.Launcher.LaunchFromFile(path)
}

// LaunchFromString passes the native configuration string to the inner launcher.
func (l *LauncherWrapper) LaunchFromString(c string) bool {
	return l.Launcher.LaunchFromString(c)
}

// Stop releases the blocking from inner api
func (l *LauncherWrapper) Stop() {
	l.Launcher.Stop()
}
