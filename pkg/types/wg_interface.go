package types

type WireguardInterface struct {
	Name    string           `yaml:"name" hcl:",label"`
	Address string           `yaml:"address" hcl:"address"`
	Port    int              `yaml:"port" hcl:"port"`
	Key     Key              `yaml:"key" hcl:"key,block"`
	Peers   []*WireguardPeer `yaml:"peers" hcl:"peer,block"`
}

type Key struct {
	PrivateKey string `yaml:"privateKey" hcl:"privateKey,optional"`
	PublicKey  string `yaml:"publicKey" hcl:"publicKey,optional"`
}

type Endpoint struct {
	Address string `yaml:"address" hcl:"address"`
	Port    int    `yaml:"port" hcl:"port"`
}

type WireguardPeer struct {
	Key         Key       `yaml:"key" hcl:"key,block"`
	Name        string    `yaml:"name" hcl:",label"`
	Description string    `yaml:"description" hcl:"description,optional"`
	Endpoint    *Endpoint `yaml:"endpoint" hcl:"endpoint,block"`
	AllowedIPs  []string  `yaml:"allowedIPs" hcl:"allowedIPs"`
	KeepAlive   int       `yaml:"keepAlive" hcl:"keepAlive"`
}
