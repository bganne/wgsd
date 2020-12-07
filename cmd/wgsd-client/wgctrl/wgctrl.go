package wgctrl

import (
	"net"

	linuxwgctrl "golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type LinuxWgCtrl struct {
	linux *linuxwgctrl.Client
}

func NewLinux() (client *LinuxWgCtrl, err error) {
	linux, err := linuxwgctrl.New()
	if err != nil {
		return nil, err
	}
	return &LinuxWgCtrl{
		linux: linux,
	}, nil
}

func (c *LinuxWgCtrl) Device(name string) (dev *WgDevice, err error) {
	wgDevice, err := c.linux.Device(name)
	if err != nil {
		return nil, err
	}
	peers := make([]WgPeer, len(wgDevice.Peers))
	for _, wpeer := range wgDevice.Peers {
		peers = append(peers, WgPeer{
			Port:       wpeer.Endpoint.Port,
			Addr:       wpeer.Endpoint.IP,
			AllowedIPs: wpeer.AllowedIPs,
			PublicKey:  wpeer.PublicKey,
		})
	}
	return &WgDevice{
		Name:       name,
		PublicKey:  wgDevice.PublicKey,
		PrivateKey: wgDevice.PrivateKey,
		Peers:      peers,
	}, nil
}

func (c *LinuxWgCtrl) AddPeer(dev *WgDevice, peer *WgPeer) (err error) {
	wgDevice, err := c.linux.Device(dev.Name)
	if err != nil {
		return err
	}
	peerConfig := wgtypes.PeerConfig{
		PublicKey:  peer.PublicKey,
		UpdateOnly: false,
		Endpoint: &net.UDPAddr{
			IP:   peer.Addr,
			Port: peer.Port,
		},
		ReplaceAllowedIPs: true,
		AllowedIPs:        peer.AllowedIPs,
	}
	deviceConfig := wgtypes.Config{
		PrivateKey:   &wgDevice.PrivateKey,
		ReplacePeers: false,
		Peers:        []wgtypes.PeerConfig{peerConfig},
	}
	if wgDevice.FirewallMark > 0 {
		deviceConfig.FirewallMark = &wgDevice.FirewallMark
	}
	err = c.linux.ConfigureDevice(dev.Name, deviceConfig)
	return err
}
