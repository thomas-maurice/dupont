package types

import "fmt"

type VXLANInterface struct {
	Name       string `yaml:"name" hcl:",label"`
	Address    string `yaml:"address" hcl:"address"`
	Parent     string `yaml:"parent" hcl:"parent"`
	VNI        int    `yaml:"vni" hcl:"vni"`
	Port       int    `yaml:"port" hcl:"port,optional"`
	Neighbours []struct {
		// These are the neighbours inside of the vxlan overlay
		Address string `yaml:"address" hcl:"address"`
	} `yaml:"neighbours" hcl:"neighbour,block"`
}

func (vx *VXLANInterface) BridgeName() string {
	return fmt.Sprintf("br-%s", vx.Name)
}
