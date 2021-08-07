name = "example topology"

network {
    wireguard     = "10.80.0.1/24"
    overlay       = "10.80.1.1/24"
    vni           = 42
    wireguardPort = 6060
}

hosts = {
    pi1 = "19.99.1.60"
    pi2 = "19.99.1.61"
    pi3 = "19.99.1.62"
    pi4 = "19.99.1.63"
}