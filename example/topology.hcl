name = "example topology"

network {
    wireguard     = "10.80.0.1/24"
    overlay       = "10.80.1.1/24"
    vni           = 42
    wireguardPort = 1234
}

hosts = {
    pi1 = "10.99.1.60"
    pi2 = "10.99.1.61"
    pi3 = "10.99.1.62"
}
