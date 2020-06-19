FROM serjs/go-socks5-proxy as build-socks5

FROM debian:9

# debian cn mirror
RUN sed -i "2s/^/#/" /etc/apt/sources.list && \
    sed -i "s/deb.debian.org/ftp2.cn.debian.org/g" /etc/apt/sources.list

# install deps
RUN apt-get update && \
    apt-get install -y net-tools curl zip && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists

# install motionpro
ENV MOTIONPRO_STANDALONE_URL https://support.arraynetworks.net/prx/001/http/supportportal.arraynetworks.net/downloads/pkg_9_4_0_253/vpn_cmdline_linux64_v1.0.54.zip
RUN curl -Lo /tmp/motionpro.zip "$MOTIONPRO_STANDALONE_URL" && \
    mkdir -p /usr/local/share/motionpro && \
    cd /usr/local/share/motionpro && \
    unzip /tmp/motionpro.zip && \
    rm /tmp/motionpro.zip

# install dumb-init
ENV DUMB_INIT_URL https://github.com/Yelp/dumb-init/releases/download/v1.2.2/dumb-init_1.2.2_amd64.deb
RUN curl -Lo /tmp/dumb-init.deb "$DUMB_INIT_URL" && \
    dpkg -i /tmp/dumb-init.deb && \
    rm /tmp/dumb-init.deb

# install socks5 server
COPY --from=build-socks5 /socks5 /usr/local/bin/socks5

ENV VPN_HOST vpn.shnu.edu.cn
ENV VPN_PORT 443

COPY start.sh /start.sh
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/start.sh"]
