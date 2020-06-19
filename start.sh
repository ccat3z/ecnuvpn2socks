#! /bin/bash

faketty () {
    script -qfec "$(printf "%q " "$@")"
}

if [ -z "$USERNAME" ] || [ -z "$PASSWORD" ]; then
    echo "\$USERNAME or \$PASSWORD is empty" >&2
    exit 1
fi


(
while [ ! -f /etc/resolv.conf.array ]; do
    echo "waiting vpn..." >&2
    sleep 2s
done

echo "vpn is ready, start socks5 server" >&2
cat /etc/resolv.conf.array > /etc/resolv.conf
socks5
)&

cd /usr/local/share/motionpro/ || exit 1
faketty ./vpn_cmdline -h "$VPN_HOST" -o "$VPN_PORT" -u "$USERNAME" -p "$PASSWORD"
