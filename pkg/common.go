package dupont

import (
	"fmt"
	"net"
	"os"

	"github.com/vishvananda/netlink"
	"go.uber.org/zap"
)

func configureInterfaceRoute(name string, route *net.IPNet) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return err
	}

	return netlink.RouteReplace(&netlink.Route{
		LinkIndex: link.Attrs().Index,
		Dst:       route,
		Scope:     netlink.SCOPE_LINK,
	})
}

// ensureBridge makes sure the bridge exists and is of the correct type.
// if not the bridge will be destroyed and re-created
func ensureBridge(log *zap.Logger, name string) error {
	link, _ := netlink.LinkByName(name)

	if link != nil {
		if link.Type() != "bridge" {
			err := netlink.LinkDel(link)
			if err != nil {
				return err
			}
		}
	}

	err := netlink.LinkAdd(&netlink.Bridge{LinkAttrs: netlink.LinkAttrs{Name: name}})
	if err != nil && !os.IsExist(err) {
		return err
	}

	link, err = netlink.LinkByName(name)
	if err != nil {
		return err
	}
	if link == nil {
		return fmt.Errorf("could not get a handle on %s", name)
	}
	if err := netlink.LinkSetUp(link); err != nil {
		return err
	}

	return nil
}
