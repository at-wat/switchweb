IMAGE     := ghcr.io/at-wat/switchweb
IMAGE_TAG := latest

.PHONY: docker-build
docker-build:
	docker build -t $(IMAGE):$(IMAGE_TAG) .

.PHONY: show-image-full-name
show-image-full-name:
	@echo "$(IMAGE):$(IMAGE_TAG)"
