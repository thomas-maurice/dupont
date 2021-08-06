package dupont

import (
	"fmt"
	"net"
	"os"
	"syscall"

	"github.com/thomas-maurice/dupont/pkg/types"
	"github.com/vishvananda/netlink"
	"go.uber.org/zap"
)

const (
	vxlanInterfaceMTU = wireguardInterfaceMTU - 54
	nullMAC           = "00:00:00:00:00:00"
)

var (
	parsedNullMac net.HardwareAddr
)

func init() {
	var err error
	parsedNullMac, err = net.ParseMAC(nullMAC)
	if err != nil {
		panic(err)
	}
}

// ensureVXLANInterface makes sure the interface exists and is of the correct type.
// if not the interface will be destroyed and re-created
func ensureVXLANInterface(log *zap.Logger, iface *types.VXLANInterface, parentConfig *types.WireguardInterface) error {
	nl, err := netlink.NewHandle(netlink.FAMILY_ALL)
	if err != nil {
		return err
	}

	parent, err := netlink.LinkByName(iface.Parent)
	if err != nil {
		return err
	}

	parentAddresses, err := netlink.AddrList(parent, netlink.FAMILY_V4)
	if err != nil {
		return err
	}

	err = ensureBridge(log, iface.BridgeName())
	if err != nil && !os.IsExist(err) {
		return err
	}

	masterInterface, err := netlink.LinkByName("br-" + iface.Name)
	if err != nil {
		return err
	}

	if err := netlink.LinkSetUp(masterInterface); err != nil {
		return err
	}

	linkOpts := netlink.Vxlan{
		LinkAttrs: netlink.LinkAttrs{
			Name:        iface.Name,
			MasterIndex: masterInterface.Attrs().Index,
			MTU:         vxlanInterfaceMTU,
		},
		VxlanId:      iface.VNI,
		VtepDevIndex: parent.Attrs().Index,
		SrcAddr:      parentAddresses[0].IP,
		Port:         iface.Port,
		Learning:     false,
	}

	link, _ := netlink.LinkByName(iface.Name)

	if link != nil {
		if link.Type() != "vxlan" {
			err := netlink.LinkDel(link)
			if err != nil {
				return err
			}
		}
	}

	err = netlink.LinkAdd(&linkOpts)
	if err != nil && !os.IsExist(err) {
		return err
	}

	link, err = netlink.LinkByName(iface.Name)
	if err != nil {
		return err
	}
	if link == nil {
		return fmt.Errorf("could not get a handle on %s", iface.Name)
	}
	if err := netlink.LinkSetMTU(link, vxlanInterfaceMTU); err != nil {
		return err
	}
	if err := netlink.LinkSetUp(link); err != nil {
		return err
	}

	for _, wgPeer := range parentConfig.Peers {
		for _, addr := range wgPeer.AllowedIPs {
			peerIP, _, err := net.ParseCIDR(addr)
			if err != nil {
				return err
			}
			entry := netlink.Neigh{
				LinkIndex:    link.Attrs().Index,
				Family:       syscall.AF_BRIDGE,
				Flags:        netlink.NTF_SELF,
				State:        netlink.NUD_PERMANENT,
				IP:           peerIP,
				HardwareAddr: parsedNullMac,
			}

			err = nl.NeighAppend(&entry)
			if err != nil {
				return err
			}
		}
	}

	for _, neigh := range iface.Neighbours {
		entry := netlink.Neigh{
			LinkIndex:    link.Attrs().Index,
			State:        netlink.NUD_PERMANENT,
			Type:         syscall.RTN_UNICAST,
			IP:           net.ParseIP(neigh.Address),
			HardwareAddr: parsedNullMac,
		}

		err = nl.NeighAppend(&entry)
		if err != nil {
			return err
		}
	}

	return nil
}
