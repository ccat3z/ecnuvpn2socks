# FROM go:1.15 as build-controller

FROM debian:12

# install deps
RUN apt-get update && \
    apt-get install -y openconnect && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists

# install motionpro
ENV TZ Asia/Shanghai

# install controller
# COPY --from=build-controller /controller /usr/local/bin/controller
COPY /controller /controller

ENTRYPOINT ["/controller"]
