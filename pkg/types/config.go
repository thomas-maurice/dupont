package types

type Config struct {
	EnsureSysctl bool            `yaml:"ensureSysctl" hcl:"ensureSysctl,optional"`
	Interfaces   InterfacesBlock `yaml:"interfaces" hcl:"interfaces,block"`
}

type InterfacesBlock struct {
	Wireguard []WireguardInterface `yaml:"wireguard" hcl:"wireguard,block"`
	VXLAN     []VXLANInterface     `yaml:"vxlan" hcl:"vxlan,block"`
}

func (cfg *Config) CheckConfig() []error {
	return nil
}
