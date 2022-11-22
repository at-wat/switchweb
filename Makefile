IMAGE := ghcr.io/at-wat/switchweb

.PHONY: docker-build
docker-build:
	docker build -t $(IMAGE) .
