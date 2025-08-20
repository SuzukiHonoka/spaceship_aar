//go:build unix

package spaceship_aar

import (
	"net/url"

	"github.com/SuzukiHonoka/tun2socks/v2/core/device"
	"github.com/SuzukiHonoka/tun2socks/v2/core/device/tun"
)

func parseTUN(u *url.URL, mtu uint32) (device.Device, error) {
	return tun.Open(u.Host, mtu)
}
