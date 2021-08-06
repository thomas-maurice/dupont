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

```yaml
# Make sure we enable ip forward and co
ensureSysctl: true

# Our interfaces definitions
interfaces:
  wireguard:
    - name: wg-0
      address: 192.168.69.2/24
      port: 6969
      key:
        privateKey: yEkDVgYvtyssPeK1cxCKYoZjNt65GL7NqBNceuOQQlY=
        # Note that specifying the public key here is a matter
        # of convenience, you would not have that (prolly) on
        # an actual deployment
        publicKey: NYNj4shJcxucrhgNTwRg1sshlCT9cGKvClWEsycm/28=
      peers:
        - name: "wg-0"
          description: "Desktop"
          key:
            publicKey: bScGfgslFnmIEcuAdU8PQla6OtE29VntPOd3rOb5phs=
          allowedIPs:
            - 192.168.69.1/24
          endpoint:
            address: 10.99.1.230
            port: 6969
          keepAlive: 5
  vxlan:
    - name: vx-0
      address: 192.168.70.2/24
      vni: 60
      parent: wg-0
      neighbours:
        - address: 192.168.70.1
```