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
	Name        string    `yaml:"name" hcl:",label"`
	Description string    `yaml:"description" hcl:"description,optional"`
	KeepAlive   int       `yaml:"keepAlive" hcl:"keepAlive"`
	Key         Key       `yaml:"key" hcl:"key,block"`
	Endpoint    *Endpoint `yaml:"endpoint" hcl:"endpoint,block"`
	AllowedIPs  []string  `yaml:"allowedIPs" hcl:"allowedIPs"`
}
