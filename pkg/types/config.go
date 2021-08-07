package types

type Config struct {
	Interfaces struct {
		Wireguard []WireguardInterface `yaml:"wireguard" hcl:"wireguard,block"`
		VXLAN     []VXLANInterface     `yaml:"vxlan" hcl:"vxlan,block"`
	} `yaml:"interfaces" hcl:"interfaces,block"`
	EnsureSysctl bool `yaml:"ensureSysctl" hcl:"ensureSysctl,optional"`
}

func (cfg *Config) CheckConfig() []error {
	return nil
}
