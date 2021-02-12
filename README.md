# shnuvpn2socks

Convert SHNU VPN to socks5 server

## Introdouction

Create tun/tap device in container and run a socks5 server.

### Problems of Motion Pro Client (Standalone version) for Linux

* Modify /etc/hosts & /etc/resolv.conf without notification
* Break main ip route table
* Silent failure if net-tools is missing, which may not be included in recent
  linux distributions
* Use 1.1.1.1 as address as tun/tap interface which used by Clouldflare DNS

## Usage

``` sh
make docker-image
docker run \
  --cap-add=NET_ADMIN --device /dev/net/tun:/dev/net/tun \
  -p 127.0.0.1:8289:8289/tcp -p 127.0.0.1:8289:8289/udp \
  shnuvpn2socks \
  -vpn-username name \
  -vpn-password pass \
  -socks-port 1080 \
  -ip 127.0.0.1
```

### Flags

```
-fail-threshold int
      Threshold of immediate failure time in second (default 60)
-ip string
      socks server ip (default "127.0.0.1")
-max-try int
      Maximum attempts of reconnect (default 3)
-motionpro-lib string
      motionpro lib path (default "/usr/local/share/motionpro/")
-socks-password string
      password of socks account
-socks-port int
      socks server port (default 1080)
-socks-username string
      username of socks account
-vpn-host string
      hostname of vpn (default "vpn.shnu.edu.cn")
-vpn-password string
      password of vpn
-vpn-port string
      port of vpn (default "443")
-vpn-username string
      username of vpn
```