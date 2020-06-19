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

```
docker run --cap-add=NET_ADMIN --device /dev/net/tun:/dev/net/tun --port 1080:1080 $image
```

### Enviroments

* VPN config
  * `USERNAME`
  * `PASSWORD`
* SOCKS5 config
  * `PROXY_USER`
  * `PROXY_PASSWORD`
  * `PROXY_PORT`
