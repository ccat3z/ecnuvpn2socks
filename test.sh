#! /bin/sh

set -ex

export IMAGE_NAME=shnuvpn2socks-test
make -e docker-image
exec docker run -ti --rm --name $IMAGE_NAME --cap-add=NET_ADMIN --device /dev/net/tun:/dev/net/tun -p 127.0.0.1:8289:8289/tcp -p 127.0.0.1:8289:8289/udp $IMAGE_NAME "$@"