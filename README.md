# shnuvpn2socks

## Usage

```
docker run --cap-add=NET_ADMIN --device /dev/net/tun:/dev/net/tun --port 1080:1080 $image
```

## Enviroments

* VPN config
  * `USERNAME`
  * `PASSWORD`
* SOCKS5 config
  * `PROXY_USER`
  * `PROXY_PASSWORD`
  * `PROXY_PORT`
