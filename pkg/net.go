package dupont

import (
	"fmt"
	"net"
	"strings"
	"syscall"

	"github.com/thomas-maurice/dupont/pkg/types"
	"github.com/vishvananda/netlink"
	"go.uber.org/zap"
)

func ApplyConfiguration(log *zap.Logger, cfg *types.Config) error {
	wgIfaces := make(map[string]*types.WireguardInterface)

	if cfg.EnsureSysctl {
		err := ensureSysctl(log)
		if err != nil {
			return err
		}
	}

	for _, wgIface := range cfg.Interfaces.Wireguard {
		wgIfaces[wgIface.Name] = &wgIface
		err := ensureWireguardInterface(log, wgIface.Name)
		if err != nil {
			return err
		}
		err = ensureIPAddress(log, wgIface.Name, wgIface.Address)
		if err != nil {
			return err
		}
		err = configureWireguardInterface(log, &wgIface)
		if err != nil {
			return err
		}
	}

	for _, vxlanIface := range cfg.Interfaces.VXLAN {
		err := ensureVXLANInterface(log, &vxlanIface, wgIfaces[vxlanIface.Parent])
		if err != nil {
			return err
		}
		err = ensureIPAddress(log, vxlanIface.BridgeName(), vxlanIface.Address)
		if err != nil {
			return err
		}
	}

	return nil
}

func TearDownConfiguration(log *zap.Logger, cfg *types.Config) error {
	for _, vxlanIface := range cfg.Interfaces.VXLAN {
		link, err := netlink.LinkByName(vxlanIface.Name)
		if strings.Contains(err.Error(), "Link not found") {
			continue
		} else if err != nil {
			return fmt.Errorf("could not get link %s: %w", vxlanIface.Name, err)
		}
		err = netlink.LinkDel(link)
		if err != nil {
			return fmt.Errorf("could not delete link %s: %w", vxlanIface.Name, err)
		}
	}

	for _, vxlanIface := range cfg.Interfaces.VXLAN {
		link, err := netlink.LinkByName(vxlanIface.BridgeName())
		if strings.Contains(err.Error(), "Link not found") {
			continue
		} else if err != nil {
			return fmt.Errorf("could not get bridge %s: %w", vxlanIface.BridgeName(), err)
		}
		err = netlink.LinkDel(link)
		if err != nil {
			return fmt.Errorf("could not delete bridge %s: %w", vxlanIface.BridgeName(), err)
		}
	}

	for _, wgIface := range cfg.Interfaces.Wireguard {
		link, err := netlink.LinkByName(wgIface.Name)
		if strings.Contains(err.Error(), "Link not found") {
			continue
		} else if err != nil {
			return fmt.Errorf("could not get wg %s: %w", wgIface.Name, err)
		}
		err = netlink.LinkDel(link)
		if err != nil {
			return fmt.Errorf("could not delete wg %s: %w", wgIface.Name, err)
		}
	}

	return nil
}

func ensureIPAddress(log *zap.Logger, name string, addr string) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return err
	}

	ipaddr, address, err := net.ParseCIDR(addr)
	if err != nil {
		return err
	}

	address.IP = ipaddr

	addrs, err := netlink.AddrList(link, syscall.AF_INET)
	if err != nil {
		return err
	}
	for _, addr := range addrs {
		if !addr.IP.Equal(address.IP) || addr.Mask.String() != address.Mask.String() {
			err = netlink.AddrDel(link, &addr)
			if err != nil {
				return err
			}
		}
	}

	err = netlink.AddrReplace(link, &netlink.Addr{
		IPNet: address,
	})

	if err != nil {
		return err
	}

	return nil
}
