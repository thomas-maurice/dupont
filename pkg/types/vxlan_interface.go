package types

import "fmt"

type VXLANInterface struct {
	Interface `yaml:",inline"`

	Parent     string `yaml:"parent"`
	VNI        int    `yaml:"vni"`
	Port       int    `yaml:"port"`
	Neighbours []struct {
		// These are the neighbours inside of the vxlan overlay
		Address string `yaml:"address"`
	} `yaml:"neighbours"`
}

func (vx *VXLANInterface) BridgeName() string {
	return fmt.Sprintf("br-%s", vx.Name)
}
