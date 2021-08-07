# dupont

Creates VXLAN tunnels over wireguard

## Why ?

Wireguard does not allow you to route arbitrary traffic through a tunnel, let's say I have this setup

* Network A: 10.0.0.0/24
* Network B: 10.1.0.0/24
* Tunnel A<>B 10.3.0.0/24

You cannot make A and B communicate without having to NAT traffic, hence masking the original IPs.

To get around this you can create an overlay network on top of the wireguard tunnel.

## How ?
Look at the config files in the `examples` directory

Then compile the binary
```bash
$ make
$ ./bin/dupont -config examples/host-1.yaml
```

## Example config

You can write the configurations both in yaml and HCL, the HCL being the more readable
one, as follows:

```hcl
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
```

Which produces something like that:
```
$ ip address
[...]
40: wg-0: <POINTOPOINT,NOARP,UP,LOWER_UP> mtu 1420 qdisc noqueue state UNKNOWN group default 
    link/none 
    inet 192.168.69.2/24 brd 192.168.69.255 scope global wg-0
       valid_lft forever preferred_lft forever
41: br-vx-0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1350 qdisc noqueue state UP group default 
    link/ether 8a:a5:6a:ec:81:e5 brd ff:ff:ff:ff:ff:ff
    inet 192.168.70.2/24 brd 192.168.70.255 scope global br-vx-0
       valid_lft forever preferred_lft forever
    inet6 fe80::88a5:6aff:feec:81e5/64 scope link 
       valid_lft forever preferred_lft forever
42: vx-0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1350 qdisc noqueue master br-vx-0 state UNKNOWN group default 
    link/ether 8a:a5:6a:ec:81:e5 brd ff:ff:ff:ff:ff:ff
    inet6 fe80::88a5:6aff:feec:81e5/64 scope link 
       valid_lft forever preferred_lft forever
```