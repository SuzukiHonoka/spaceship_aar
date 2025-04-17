package spaceship_aar

import (
	"errors"
	"fmt"
	"net"
	"net/netip"
	"os/exec"
	"sync"
	"time"

	"github.com/docker/go-units"
	"github.com/google/shlex"
	"gvisor.dev/gvisor/pkg/tcpip"
	"gvisor.dev/gvisor/pkg/tcpip/stack"

	"github.com/xjasonlyu/tun2socks/v2/core"
	"github.com/xjasonlyu/tun2socks/v2/core/device"
	"github.com/xjasonlyu/tun2socks/v2/core/option"
	"github.com/xjasonlyu/tun2socks/v2/dialer"
	"github.com/xjasonlyu/tun2socks/v2/log"
	"github.com/xjasonlyu/tun2socks/v2/proxy"
	"github.com/xjasonlyu/tun2socks/v2/restapi"
	"github.com/xjasonlyu/tun2socks/v2/tunnel"
)

type Engine struct {
	_engineMu sync.Mutex

	// _defaultKey holds the default key for the engine.
	_defaultKey *EngineKey

	// _defaultProxy holds the default proxy for the engine.
	_defaultProxy proxy.Proxy

	// _defaultDevice holds the default device for the engine.
	_defaultDevice device.Device

	// _defaultStack holds the default stack for the engine.
	_defaultStack *stack.Stack
}

// Start starts the default engine up.
func (e *Engine) Start() error {
	if err := e.start(); err != nil {
		return fmt.Errorf("[ENGINE] failed to start: %v", err)
	}
	return nil
}

// Stop shuts the default engine down.
func (e *Engine) Stop() error {
	if err := e.stop(); err != nil {
		return fmt.Errorf("[ENGINE] failed to stop: %v", err)
	}
	return nil
}

// Insert loads *EngineKey to the default engine.
func (e *Engine) Insert(k *EngineKey) {
	e._engineMu.Lock()
	e._defaultKey = k
	e._engineMu.Unlock()
}

func (e *Engine) start() error {
	e._engineMu.Lock()
	defer e._engineMu.Unlock()

	if e._defaultKey == nil {
		return errors.New("empty key")
	}

	for _, f := range []func(*EngineKey) error{
		e.general,
		e.restAPI,
		e.netstack,
	} {
		if err := f(e._defaultKey); err != nil {
			return err
		}
	}
	return nil
}

func (e *Engine) stop() (err error) {
	e._engineMu.Lock()
	if e._defaultDevice != nil {
		e._defaultDevice.Close()
	}
	if e._defaultStack != nil {
		e._defaultStack.Close()
		e._defaultStack.Wait()
	}
	e._engineMu.Unlock()
	return nil
}

func (e *Engine) execCommand(cmd string) error {
	parts, err := shlex.Split(cmd)
	if err != nil {
		return err
	}
	if len(parts) == 0 {
		return errors.New("empty command")
	}
	_, err = exec.Command(parts[0], parts[1:]...).Output()
	return err
}

func (e *Engine) general(k *EngineKey) error {
	level, err := log.ParseLevel(k.LogLevel)
	if err != nil {
		return err
	}
	log.SetLogger(log.Must(log.NewLeveled(level)))

	if k.Interface != "" {
		iface, err := net.InterfaceByName(k.Interface)
		if err != nil {
			return err
		}
		dialer.DefaultDialer.InterfaceName.Store(iface.Name)
		dialer.DefaultDialer.InterfaceIndex.Store(int32(iface.Index))
		log.Infof("[DIALER] bind to interface: %s", k.Interface)
	}

	if k.Mark != 0 {
		dialer.DefaultDialer.RoutingMark.Store(int32(k.Mark))
		log.Infof("[DIALER] set fwmark: %#x", k.Mark)
	}

	if k.UDPTimeout > 0 {
		if k.UDPTimeout < time.Second {
			return errors.New("invalid udp timeout value")
		}
		tunnel.T().SetUDPTimeout(k.UDPTimeout)
	}
	return nil
}

func (e *Engine) restAPI(k *EngineKey) error {
	if k.RestAPI != "" {
		u, err := parseRestAPI(k.RestAPI)
		if err != nil {
			return err
		}
		host, token := u.Host, u.User.String()

		restapi.SetStatsFunc(func() tcpip.Stats {
			e._engineMu.Lock()
			defer e._engineMu.Unlock()

			// default stack is not initialized.
			if e._defaultStack == nil {
				return tcpip.Stats{}
			}
			return e._defaultStack.Stats()
		})

		go func() {
			if err := restapi.Start(host, token); err != nil {
				log.Errorf("[RESTAPI] failed to start: %v", err)
			}
		}()
		log.Infof("[RESTAPI] serve at: %s", u)
	}
	return nil
}

func (e *Engine) netstack(k *EngineKey) (err error) {
	if k.Proxy == "" {
		return errors.New("empty proxy")
	}
	if k.Device == "" {
		return errors.New("empty device")
	}

	if k.TUNPreUp != "" {
		log.Infof("[TUN] pre-execute command: `%s`", k.TUNPreUp)
		if preUpErr := e.execCommand(k.TUNPreUp); preUpErr != nil {
			log.Errorf("[TUN] failed to pre-execute: %s: %v", k.TUNPreUp, preUpErr)
		}
	}

	defer func() {
		if k.TUNPostUp == "" || err != nil {
			return
		}
		log.Infof("[TUN] post-execute command: `%s`", k.TUNPostUp)
		if postUpErr := e.execCommand(k.TUNPostUp); postUpErr != nil {
			log.Errorf("[TUN] failed to post-execute: %s: %v", k.TUNPostUp, postUpErr)
		}
	}()

	if e._defaultProxy, err = parseProxy(k.Proxy); err != nil {
		return
	}
	tunnel.T().SetDialer(e._defaultProxy)

	if e._defaultDevice, err = parseDevice(k.Device, uint32(k.MTU)); err != nil {
		return
	}

	var multicastGroups []netip.Addr
	if multicastGroups, err = parseMulticastGroups(k.MulticastGroups); err != nil {
		return err
	}

	var opts []option.Option
	if k.TCPModerateReceiveBuffer {
		opts = append(opts, option.WithTCPModerateReceiveBuffer(true))
	}

	if k.TCPSendBufferSize != "" {
		size, err := units.RAMInBytes(k.TCPSendBufferSize)
		if err != nil {
			return err
		}
		opts = append(opts, option.WithTCPSendBufferSize(int(size)))
	}

	if k.TCPReceiveBufferSize != "" {
		size, err := units.RAMInBytes(k.TCPReceiveBufferSize)
		if err != nil {
			return err
		}
		opts = append(opts, option.WithTCPReceiveBufferSize(int(size)))
	}

	if e._defaultStack, err = core.CreateStack(&core.Config{
		LinkEndpoint:     e._defaultDevice,
		TransportHandler: tunnel.T(),
		MulticastGroups:  multicastGroups,
		Options:          opts,
	}); err != nil {
		return
	}

	log.Infof(
		"[STACK] %s://%s <-> %s://%s",
		e._defaultDevice.Type(), e._defaultDevice.Name(),
		e._defaultProxy.Proto(), e._defaultProxy.Addr(),
	)
	return nil
}
