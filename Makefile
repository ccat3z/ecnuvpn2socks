IMAGE_NAME=shnuvpn2socks
DOCKER=docker
GO=go

.PHONY: controller
controller:
	CGO_ENABLED=0 $(GO) build -o controller

.PHONY: docker-image
docker-image: controller
	$(DOCKER) build -t $(IMAGE_NAME) .