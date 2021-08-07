package dupont

import (
	"fmt"
	"net"
	"strings"

	"github.com/thomas-maurice/dupont/pkg/types"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// from https://stackoverflow.com/questions/31191313/how-to-get-the-next-ip-address
func nextIP(ip net.IP, inc uint) net.IP {
	i := ip.To4()
	v := uint(i[0])<<24 + uint(i[1])<<16 + uint(i[2])<<8 + uint(i[3])
	v += inc
	v3 := byte(v & 0xFF)
	v2 := byte((v >> 8) & 0xFF)
	v1 := byte((v >> 16) & 0xFF)
	v0 := byte((v >> 24) & 0xFF)
	return net.IPv4(v0, v1, v2, v3)
}

func GenerateTopology(t *types.Topology) (map[string]*types.Config, error) {
	result := make(map[string]*types.Config)

	if _, err := t.Valid(); err != nil {
		return nil, err
	}

	if err := t.Fill(); err != nil {
		return nil, err
	}

	topologyID := t.ID()
	idx := 0
	// Fill in the hosts, with they ip addresses and wireguard keys
	for host := range t.Hosts {
		wgKey, err := wgtypes.GenerateKey()
		if err != nil {
			return nil, err
		}
		msk, _ := t.Network.OverlayCIDR.Mask.Size()

		cfg := &types.Config{
			EnsureSysctl: true,
			Interfaces: types.InterfacesBlock{
				Wireguard: []types.WireguardInterface{
					{
						Name: fmt.Sprintf("wg-%s", topologyID),
						Address: fmt.Sprintf(
							"%s/%d",
							nextIP(t.Network.WireguardAddress, uint(idx)).String(),
							32,
						),
						Port: t.Network.WireguardPort,
						Key: types.Key{
							PrivateKey: wgKey.String(),
							PublicKey:  wgKey.PublicKey().String(),
						},
						Peers: make([]*types.WireguardPeer, 0),
					},
				},
				VXLAN: []types.VXLANInterface{
					{
						Name: fmt.Sprintf("vx-%s", topologyID),
						Address: fmt.Sprintf(
							"%s/%d",
							nextIP(t.Network.OverlayAddress, uint(idx)).String(),
							msk,
						),
						VNI: t.Network.VNI,
					},
				},
			},
		}

		idx++
		result[host] = cfg
	}

	// Fill in all the rest
	for name, config := range result {
		for peerName, peerConfig := range result {
			if name == peerName {
				continue
			}

			peer := types.WireguardPeer{
				Name: peerName,
				Description: fmt.Sprintf(
					"%s: %s - %s",
					peerName,
					peerConfig.Interfaces.Wireguard[0].Address,
					topologyID,
				),
				Key: types.Key{
					PublicKey: peerConfig.Interfaces.Wireguard[0].Key.PublicKey,
				},
				Endpoint: &types.Endpoint{
					Address: t.Hosts[peerName],
					Port:    t.Network.WireguardPort,
				},
				KeepAlive:  5,
				AllowedIPs: []string{peerConfig.Interfaces.Wireguard[0].Address},
			}

			config.Interfaces.Wireguard[0].Peers = append(config.Interfaces.Wireguard[0].Peers, &peer)

			config.Interfaces.VXLAN[0].Neighbours = append(config.Interfaces.VXLAN[0].Neighbours, struct {
				Address string "yaml:\"address\" hcl:\"address\""
			}{Address: strings.Split(peerConfig.Interfaces.VXLAN[0].Address, "/")[0]})

			config.Interfaces.VXLAN[0].Parent = config.Interfaces.Wireguard[0].Name
		}
	}

	return result, nil
}
