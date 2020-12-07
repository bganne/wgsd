package wgctrl

import (
	"errors"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"net"
)

type WgPeer struct {
	Port       int
	Addr       net.IP
	AllowedIPs []net.IPNet
	PublicKey  wgtypes.Key
}

type WgDevice struct {
	Name       string
	SwIfIndex  uint32
	Peers      []WgPeer
	PublicKey  wgtypes.Key
	PrivateKey wgtypes.Key
}

type WgClient interface {
	Device(name string) (dev *WgDevice, err error)
	AddPeer(dev *WgDevice, peer *WgPeer) (err error)
}

func New(dataplane string, socket string) (WgClient, error) {
	switch dataplane {
	case "linux":
		client, err := NewLinux()
		return client, err
	case "vpp":
		client, err := NewVpp(socket)
		return client, err
	default:
		return nil, errors.New("Unknown dataplane " + dataplane)
	}
}
