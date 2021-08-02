# ecnuvpn2socks

Access ECNU VPN via socks5

## Usage

``` sh
make docker-image
docker run \
  --cap-add=NET_ADMIN --device /dev/net/tun:/dev/net/tun \
  -p 127.0.0.1:8290:1080/tcp -p 127.0.0.1:8289:1080/udp \
  ecnuvpn2socks \
  -vpn-username name \
  -vpn-password pass \
  -socks-port 1080 \
  -ip 127.0.0.1
# Or use docker-compose. Don't forget to set the environment variables required by docker-compose.yml.
```

### Flags

```
-fail-threshold int
      Threshold of immediate failure time in second (default 60)
-ip string
      socks server ip (default "127.0.0.1")
-max-try int
      Maximum attempts of reconnect (default 3)
-socks-password string
      password of socks account
-socks-port int
      socks server port (default 1080)
-socks-username string
      username of socks account
-vpn-host string
      hostname of vpn (default "vpn-cn.ecnu.edu.cn")
-vpn-password string
      password of vpn
-vpn-username string
      username of vpn
```
