package spaceship_aar

import "encoding/json"

// since go-mobile only support basic types, some fields are minimized

type Config struct {
	Path        string `json:"path"`
	Host        string `json:"host"`
	ServerAddr  string `json:"server_addr"`
	Uuid        string `json:"uuid"`
	ListenSocks string `json:"listen_socks"`
	ListenHttp  string `json:"listen_http"`
	Tls         bool   `json:"tls"`
	Mux         int    `json:"mux"`
	Buffer      int    `json:"buffer"`
	DNS         string `json:"dns"`
	IPv6        bool   `json:"ipv6"`
	CA          string `json:"ca"`
	Log         string `json:"log"`
}

func (c *Config) ToJson() string {
	b, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return string(b)
}
