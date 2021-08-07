package types

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
)

const (
	TopologyIDLength = 6
)

type Topology struct {
	Name    string `yaml:"name" hcl:"name"`
	Network struct {
		Wireguard        string     `yaml:"wireguard" hcl:"wireguard"`
		WireguardCIDR    *net.IPNet `yaml:"-"`
		WireguardAddress net.IP     `yaml:"-"`
		Overlay          string     `yaml:"overlay" hcl:"overlay"`
		OverlayCIDR      *net.IPNet `yaml:"-"`
		OverlayAddress   net.IP     `yaml:"-"`
		VNI              int        `yaml:"vni" hcl:"vni"`
		WireguardPort    int        `yaml:"wireguardPort" hcl:"wireguardPort"`
	} `yaml:"network" hcl:"network,block"`
	Hosts map[string]string `yaml:"hosts" hcl:"hosts"`
}

func (t *Topology) Valid() (bool, error) {
	if t.Name == "" {
		return false, fmt.Errorf("topology name cant be empty")
	}

	if t.Network.VNI == 0 {
		return false, fmt.Errorf("VNI cannot be 0")
	}

	if t.Network.WireguardPort == 0 {
		return false, fmt.Errorf("Wireguard port cannot be 0")
	}

	return true, nil
}

func (t *Topology) Fill() error {
	t.Network.WireguardAddress = net.ParseIP(strings.Split(t.Network.Wireguard, "/")[0])
	t.Network.OverlayAddress = net.ParseIP(strings.Split(t.Network.Overlay, "/")[0])

	var err error
	_, t.Network.WireguardCIDR, err = net.ParseCIDR(t.Network.Wireguard)
	if err != nil {
		return err
	}
	_, t.Network.OverlayCIDR, err = net.ParseCIDR(t.Network.Overlay)
	if err != nil {
		return err
	}

	return nil
}

func (t *Topology) ID() string {
	return hex.EncodeToString(sha1.New().Sum([]byte(t.Name)))[:TopologyIDLength]
}
