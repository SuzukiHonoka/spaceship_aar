package spaceship_aar

import (
	"encoding/json"
	"fmt"
	"github.com/SuzukiHonoka/spaceship/v2/api"
	"github.com/SuzukiHonoka/spaceship/v2/pkg/config"
	"github.com/SuzukiHonoka/spaceship/v2/pkg/config/client"
	"github.com/SuzukiHonoka/spaceship/v2/pkg/config/server"
	"github.com/SuzukiHonoka/spaceship/v2/pkg/dns"
	"github.com/SuzukiHonoka/spaceship/v2/pkg/logger"
)

// TotalResultWrapper is the result of the total bytes sent and received.
type TotalResultWrapper struct {
	r api.TotalResult
}

// BytesSent 8 bytes unsigned-integer bytes representation, little endian
func (r TotalResultWrapper) BytesSent() []byte {
	return r.r.BytesSent
}

// BytesReceived 8 bytes unsigned-integer bytes representation, little endian
func (r TotalResultWrapper) BytesReceived() []byte {
	return r.r.BytesReceived
}

func (r TotalResultWrapper) String() string {
	return r.r.String()
}

// SpeedResultWrapper is the result of the speed of transfer in seconds.
type SpeedResultWrapper struct {
	r api.SpeedResult
}

// BytesSent is the total bytes sent in the speed result.
func (r SpeedResultWrapper) BytesSent() float64 {
	return r.r.BytesSent
}

// BytesReceived is the total bytes received in the speed result.
func (r SpeedResultWrapper) BytesReceived() float64 {
	return r.r.BytesReceived
}

// String returns the string representation of the speed result.
func (r SpeedResultWrapper) String() string {
	return r.r.String()
}

// LauncherWrapper wraps the api.Launcher and provides some proxy methods.
type LauncherWrapper struct {
	l *api.Launcher
}

// NewLauncher creates a new *LauncherWrapper.
func NewLauncher() *LauncherWrapper {
	return &LauncherWrapper{
		l: api.NewLauncher(),
	}
}

// Launch parses the particular cfg string to the client-oriented configuration and launch the inner launcher.
// Note that this is not native configuration format. if using native, go for LaunchFromString instead.
func (l *LauncherWrapper) Launch(s string) error {
	var cfg Config
	if err := json.Unmarshal([]byte(s), &cfg); err != nil {
		return fmt.Errorf("launch: unmarshal cfg failed, err=%w", err)
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
		Client: &client.Client{
			ServerAddr:      cfg.ServerAddr,
			Host:            cfg.Host,
			UUID:            cfg.Uuid,
			ListenSocks:     cfg.ListenSocks,
			ListenSocksUnix: cfg.ListenSocksUnix,
			ListenHttp:      cfg.ListenHttp,
			BasicAuth:       cfg.BasicAuth,
			Mux:             uint8(cfg.Mux),
			EnableTLS:       cfg.Tls,
		},
		Server: &server.Server{
			Path:   cfg.Path,
			Buffer: uint16(cfg.Buffer),
			IPv6:   cfg.IPv6,
		},
	}
	return l.l.Launch(c)
}

// LaunchFromFile reads the native configuration string from file and passes to the inner launcher.
func (l *LauncherWrapper) LaunchFromFile(path string) error {
	return l.l.LaunchFromFile(path)
}

// LaunchFromString passes the native configuration string to the inner launcher.
func (l *LauncherWrapper) LaunchFromString(c string) error {
	return l.l.LaunchFromString(c)
}

// Stop releases the blocking from inner api
func (l *LauncherWrapper) Stop() {
	l.l.Stop()
}

// Speed returns the speed result from inner api
func (l *LauncherWrapper) Speed() *SpeedResultWrapper {
	return &SpeedResultWrapper{l.l.Speed()}
}

// Total returns the total result from inner api
func (l *LauncherWrapper) Total() *TotalResultWrapper {
	return &TotalResultWrapper{l.l.Total()}
}
