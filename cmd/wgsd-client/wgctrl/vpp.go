package wgctrl

import (
	"errors"
	"time"

	"github.com/jwhited/wgsd/vpplink"
	"github.com/jwhited/wgsd/vpplink/types"
	log "github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type VppWgCtrl struct {
	vpp *vpplink.VppLink
}

func NewVpp(socket string) (ctrl *VppWgCtrl, err error) {
	// Get an API connection, with a few retries to accomodate VPP startup time
	for i := 0; i < 10; i++ {
		vpp, err := vpplink.NewVppLink(
			socket,
			log.WithFields(log.Fields{"component": "vpp"}),
		)
		if err != nil {
			log.Printf("Try [%d/10] %v", i, err)
			err = nil
			time.Sleep(2 * time.Second)
		} else {
			return &VppWgCtrl{
				vpp: vpp,
			}, nil
		}
	}
	return nil, errors.New("Cannot connect to VPP after 10 tries")
}

func (c *VppWgCtrl) Device(name string) (dev *WgDevice, err error) {
	swIfIndex, err := c.vpp.SearchInterfaceWithName(name)
	if err != nil {
		return nil, err
	}
	tunnel, err := c.vpp.FindWireguardTunnel(swIfIndex, true /* showPrivateKey */)
	if err != nil {
		return nil, err
	}
	peers, err := c.vpp.ListWireguardPeers(swIfIndex)
	if err != nil {
		return nil, err
	}
	pubKey, err := wgtypes.NewKey(tunnel.PublicKey)
	if err != nil {
		return nil, err
	}
	priKey, err := wgtypes.NewKey(tunnel.PrivateKey)
	if err != nil {
		return nil, err
	}
	dev = &WgDevice{
		Name:       name,
		SwIfIndex:  swIfIndex,
		PublicKey:  pubKey,
		PrivateKey: priKey,
		Peers:      make([]WgPeer, 0, len(peers)),
	}
	for _, peer := range peers {
		pubKey, err := wgtypes.NewKey(peer.PublicKey)
		if err != nil {
			return nil, err
		}
		dev.Peers = append(dev.Peers, WgPeer{
			Port:       int(peer.Port),
			Addr:       peer.Addr,
			AllowedIPs: peer.AllowedIps,
			PublicKey:  pubKey,
		})
	}
	return dev, nil
}

func (c *VppWgCtrl) AddPeer(dev *WgDevice, peer *WgPeer) (err error) {
	pubKey := make([]byte, 32)
	copy(pubKey, peer.PublicKey[:])
	_, err = c.vpp.AddWireguardPeer(&types.WireguardPeer{
		PublicKey:  pubKey,
		Port:       uint16(peer.Port),
		Addr:       peer.Addr,
		SwIfIndex:  dev.SwIfIndex,
		AllowedIps: peer.AllowedIPs,
	})
	return err
}
