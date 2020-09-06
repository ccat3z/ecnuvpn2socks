# shnuvpn2socks

Convert SHNU VPN to socks5 server

## Introdouction

Create tun/tap device in container and run a socks5 server.

### Problems of Motion Pro Client (Standalone version) for Linux

* Change /etc/hosts & /etc/resolv.conf without notification
* Broke main route table
* Silent failure if net-tools is missing
* Use 1.1.1.1 as routing address which used as DNS by Clouldflare

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