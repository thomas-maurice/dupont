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
				return fmt.Errorf("could not delete wireguard interface: %w", err)
			}
		}
	}

	err := netlink.LinkAdd(&netlink.Wireguard{LinkAttrs: netlink.LinkAttrs{Name: name}})
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("could not create wireguard interface: %w", err)
	}

	link, err = netlink.LinkByName(name)
	if err != nil {
		return fmt.Errorf("could not get a handle on %s: %w", name, err)
	}
	if link == nil {
		return fmt.Errorf("could not get a handle on %s", name)
	}
	if err := netlink.LinkSetMTU(link, wireguardInterfaceMTU); err != nil {
		return fmt.Errorf("could not set interface MTU: %w", err)
	}
	if err := netlink.LinkSetUp(link); err != nil {
		return fmt.Errorf("could not set interface up: %w", err)
	}

	return nil
}

func configureWireguardInterface(log *zap.Logger, iface *types.WireguardInterface) error {
	client, err := wgctrl.New()
	if err != nil {
		return fmt.Errorf("could not create wireguard client: %w", err)

	}

	key, err := wgtypes.ParseKey(iface.Key.PrivateKey)
	if err != nil {
		return fmt.Errorf("could not parse wireguard private key on %s: %w", iface.Name, err)
	}

	link, err := netlink.LinkByName(iface.Name)
	if err != nil {
		return fmt.Errorf("could not configure interface %s: %w", iface.Name, err)
	}
	if link == nil {
		return fmt.Errorf("could not get a handle on %s", iface.Name)
	}

	routes := make([]*netlink.Route, 0)

	var peers []wgtypes.PeerConfig
	for _, peer := range iface.Peers {
		keepaliveDuration := time.Duration(peer.KeepAlive) * time.Second

		var peerIPs []net.IPNet
		peerKey, err := wgtypes.ParseKey(peer.Key.PublicKey)
		if err != nil {
			return fmt.Errorf("could not parse private key of peer %s of %s: %w", peer.Name, iface.Name, err)
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
			if err != nil {
				return fmt.Errorf("could not parse peer allowed IP: %w", err)
			}
			if peerNet != nil {
				peerIPs = append(peerIPs, *peerNet)
			}

			routes = append(routes, &netlink.Route{
				LinkIndex: link.Attrs().Index,
				Dst:       peerNet,
				Scope:     netlink.SCOPE_LINK,
			})
		}

		peers = append(peers, wgtypes.PeerConfig{
			PublicKey:                   peerKey,
			PersistentKeepaliveInterval: &keepaliveDuration,
			ReplaceAllowedIPs:           true,
			AllowedIPs:                  peerIPs,
			Endpoint:                    udpEndpoint,
		})
	}

	for _, route := range routes {
		err := netlink.RouteReplace(route)
		if err != nil {
			return fmt.Errorf("could not apply route: %w", err)
		}
	}

	return client.ConfigureDevice(iface.Name, wgtypes.Config{
		PrivateKey:   &key,
		ListenPort:   &iface.Port,
		ReplacePeers: true,
		Peers:        peers,
	})
}
