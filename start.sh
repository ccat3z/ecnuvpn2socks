#! /bin/bash

faketty () {
    script -qfec "$(printf "%q " "$@")"
}

if [ -z "$USERNAME" ] || [ -z "$PASSWORD" ]; then
    echo "\$USERNAME or \$PASSWORD is empty" >&2
    exit 1
fi

function wait_vpn {
    while [ ! -f /etc/resolv.conf.array ]; do
        echo "waiting vpn..." >&2
        sleep 2s
    done

    cat /etc/resolv.conf.array > /etc/resolv.conf
}

function reachable {
    http_code="$(curl --connect-timeout 15 -m 15 -o /dev/null -w "%{http_code}" "$1" 2> /dev/null)"

    if [ "$http_code" = 204 ]; then
        true
    else
        false
    fi
}

echo "start socks5 server" >&2
socks5&

cd /usr/local/share/motionpro/ || exit 1
while :; do
    echo "start vpn" >&2

    faketty ./vpn_cmdline -h "$VPN_HOST" -o "$VPN_PORT" -u "$USERNAME" -p "$PASSWORD" &
    VPN_PID=$!

    wait_vpn
    echo "vpn is ready!" >&2

    while reachable http://g.cn/generate_204; do
        sleep 30s
    done

    echo "killing vpn" >&2
    ./vpn_cmdline -s
    wait "${VPN_PID}"
    echo "vpn stopped" >&2
done
