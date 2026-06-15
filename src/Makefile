IMAGE=github-exporter
IMAGE_VERSION=2.0.0
GOVERSION=1.26.4

.PHONY: dockerhub clean

dockerhub:
	buildah build --build-arg GOVERSION=${GOVERSION} -t ${IMAGE}:${IMAGE_VERSION} .
	#buildah login docker.io
	buildah push ${IMAGE}:${IMAGE_VERSION} gounix/${IMAGE}:${IMAGE_VERSION}

clean:
	buildah rmi -f `buildah images --json ${IMAGE}:${IMAGE_VERSION} |jq ".[] | .id" | sed 's/"//g'`

