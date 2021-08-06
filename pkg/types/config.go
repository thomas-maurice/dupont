package types

type Config struct {
	Interfaces struct {
		Wireguard []WireguardInterface `yaml:"wireguard"`
		VXLAN     []VXLANInterface     `yaml:"vxlan"`
	} `yaml:"interfaces"`
	EnsureSysctl    bool `yaml:"ensureSysctl"`
}

func (cfg *Config) CheckConfig() []error {
	return nil
}
