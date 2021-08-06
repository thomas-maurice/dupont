package dupont

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/thomas-maurice/dupont/pkg/types"
	"github.com/vishvananda/netlink"
	"go.uber.org/zap"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

const (
	wireguardInterfaceMTU = 1420
)

// ensureWireguardInterface makes sure the interface exists and is of the correct type.
// if not the interface will be destroyed and re-created
func ensureWireguardInterface(log *zap.Logger, name string) error {
	link, _ := netlink.LinkByName(name)

	if link != nil {
		if link.Type() != "wireguard" {
			err := netlink.LinkDel(link)
			if err != nil {
				return err
			}
		}
	}

	err := netlink.LinkAdd(&netlink.Wireguard{LinkAttrs: netlink.LinkAttrs{Name: name}})
	if err != nil && !os.IsExist(err) {
		return err
	}

	link, _ = netlink.LinkByName(name)
	if link == nil {
		return fmt.Errorf("could not get a handle on %s", name)
	}
	if err := netlink.LinkSetMTU(link, wireguardInterfaceMTU); err != nil {
		return err
	}
	if err := netlink.LinkSetUp(link); err != nil {
		return err
	}

	return nil
}

func configureWireguardInterface(log *zap.Logger, iface *types.WireguardInterface) error {
	client, err := wgctrl.New()
	if err != nil {
		return err
	}

	key, err := wgtypes.ParseKey(iface.Key.PrivateKey)
	if err != nil {
		return err
	}

	var peers []wgtypes.PeerConfig
	for _, peer := range iface.Peers {
		keepaliveDuration := time.Duration(peer.KeepAlive) * time.Second

		var peerIPs []net.IPNet
		peerKey, err := wgtypes.ParseKey(peer.Key.PublicKey)
		if err != nil {
			return err
		}

		var udpEndpoint *net.UDPAddr
		if peer.Endpoint != nil {
			udpEndpoint = &net.UDPAddr{
				IP:   net.ParseIP(peer.Endpoint.Address),
				Port: int(peer.Endpoint.Port),
			}
		}

		for _, nw := range peer.AllowedIPs {
			_, peerNet, err := net.ParseCIDR(nw)
			if peerNet != nil {
				peerIPs = append(peerIPs, *peerNet)
			}
			if err != nil {
				return err
			}
		}

		peers = append(peers, wgtypes.PeerConfig{
			PublicKey:                   peerKey,
			PersistentKeepaliveInterval: &keepaliveDuration,
			ReplaceAllowedIPs:           true,
			AllowedIPs:                  peerIPs,
			Endpoint:                    udpEndpoint,
		})
	}

	return client.ConfigureDevice(iface.Name, wgtypes.Config{
		PrivateKey:   &key,
		ListenPort:   &iface.Port,
		ReplacePeers: true,
		Peers:        peers,
	})
}
