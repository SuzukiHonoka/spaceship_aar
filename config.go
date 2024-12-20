package spaceship_aar

import "encoding/json"

// Config is a specified client-oriented configuration.
// Note that since go-mobile only support very basic types, some fields are minimized or omitted
type Config struct {
	Path            string `json:"path"`
	Host            string `json:"host"`
	ServerAddr      string `json:"server_addr"`
	Uuid            string `json:"uuid"`
	ListenSocks     string `json:"listen_socks"`
	ListenSocksUnix string `json:"listen_socks_unix"`
	ListenHttp      string `json:"listen_http"`
	Tls             bool   `json:"tls"`
	Mux             int    `json:"mux"`
	Buffer          int    `json:"buffer"`
	DNS             string `json:"dns"`
	IPv6            bool   `json:"ipv6"`
	CA              string `json:"ca"`
	Log             string `json:"log"`
}

// ToJson converts Config to string in json format
func (c *Config) ToJson() string {
	b, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return string(b)
}
