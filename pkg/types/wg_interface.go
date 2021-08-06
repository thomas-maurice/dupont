package types

type WireguardInterface struct {
	Interface `yaml:",inline"`

	Port  int             `yaml:"port"`
	Key   Key             `yaml:"key"`
	Peers []WireguardPeer `yaml:"peers"`
}

type Key struct {
	PrivateKey string `yaml:"privateKey"`
	PublicKey  string `yaml:"publicKey"`
}

type Endpoint struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

type WireguardPeer struct {
	Key         Key       `yaml:"key"`
	Name        string    `yaml:"name"`
	Description string    `yaml:"description"`
	Endpoint    *Endpoint `yaml:"endpoint"`
	AllowedIPs  []string  `yaml:"allowedIPs"`
	KeepAlive   int       `yaml:"keepAlive"`
}
