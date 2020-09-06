# FROM go:1.15 as build-controller

FROM debian:9

# debian cn mirror
RUN sed -i "2s/^/#/" /etc/apt/sources.list && \
    sed -i "s/deb.debian.org/ftp2.cn.debian.org/g" /etc/apt/sources.list

# install deps
RUN apt-get update && \
    apt-get install -y net-tools curl zip psmisc && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists

# install motionpro
ENV MOTIONPRO_STANDALONE_URL https://support.arraynetworks.net/prx/001/http/supportportal.arraynetworks.net/downloads/pkg_9_4_0_253/vpn_cmdline_linux64_v1.0.54.zip
RUN curl -Lo /tmp/motionpro.zip "$MOTIONPRO_STANDALONE_URL" && \
    mkdir -p /usr/local/share/motionpro && \
    cd /usr/local/share/motionpro && \
    unzip /tmp/motionpro.zip && \
    rm /tmp/motionpro.zip

ENV TZ Asia/Shanghai

# install fake sudo
COPY fake_sudo /usr/local/bin/sudo

# install controller
# COPY --from=build-controller /controller /usr/local/bin/controller
COPY /controller /controller

ENTRYPOINT ["/controller"]
