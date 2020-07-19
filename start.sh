#! /bin/bash

set -e

faketty () {
    script -qfec "$(printf "%q " "$@")"
}

log () {
    echo -e "\033[1m$(date "+%b %d %H:%M:%S") $*\033[0m" >&2
}

if [ -z "$USERNAME" ] || [ -z "$PASSWORD" ]; then
    echo "\$USERNAME or \$PASSWORD is empty" >&2
    exit 1
fi

function wait_vpn {
    while [ ! -f /etc/resolv.conf.array ]; do
        log "waiting vpn..."
        sleep 2s
    done

    log "update resolve.conf"
    cat /etc/resolv.conf.array > /etc/resolv.conf

    log "remove array temp files"
    rm -v /etc/*.array
}

function reachable {
    max_attempt=3
    attempt=1
    while [ "$attempt" -le "$max_attempt" ]; do
        log "try to connect $1 ($attempt/$max_attempt)"
        http_code="$(curl --connect-timeout 10 -m 10 -o /dev/null -w "%{http_code}" "$1" 2> /dev/null)"

        if [ "$http_code" = 204 ]; then
            log "success"
            return 0
        fi

        attempt=$((attempt + 1))
    done
    log "failed"
    return 1
}

log "start socks5 server"
socks5 &> /var/log/socks5.log &

cd /usr/local/share/motionpro/ || exit 1
while :; do
    log "start vpn"

    (
        wait_vpn

        log "vpn is ready! start health checker"
        while reachable http://g.cn/generate_204; do
            sleep 30s
        done

        log "cannot connect to internet. killing vpn process"
        ./vpn_cmdline -s
    )&
    DAEMON_PID=$!

    faketty ./vpn_cmdline -h "$VPN_HOST" -o "$VPN_PORT" -u "$USERNAME" -p "$PASSWORD"
    kill "$DAEMON_PID" &> /dev/null || true
    log "vpn stopped"
done
