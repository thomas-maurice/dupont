# Make sure we enable ip forward and co
ensureSysctl = true

# Our interfaces definitions
interfaces {
  # Wireguard interfaces definitions
  wireguard "wg-0" {
    # First interface definition
    address = "192.168.69.1/32"
    port    = 6969
    key {
      privateKey = "4CQWNQylWDWoZGgWDj58skAQuC84v1JXBKKqLTwcb3c="
      # Note that specifying the public key here is a matter
      # of convenience, you would not have that (prolly) on
      # an actual deployment
      publicKey = "bScGfgslFnmIEcuAdU8PQla6OtE29VntPOd3rOb5phs="
    }
    peer "wg-0" {
      description = "Laptop"
      key {
        publicKey = "NYNj4shJcxucrhgNTwRg1sshlCT9cGKvClWEsycm/28="
      }
      allowedIPs = [
        "192.168.69.2/32",
      ]
      endpoint {
        address = "10.99.1.200"
        port    = 6969
      }
      keepAlive = 5
    }
  }
  
  vxlan "vx-0" {
    address = "192.168.70.1/24"
    vni     = 60
    parent  = "wg-0"
    neighbour {
      address = "192.168.70.2"
    }
  }
}