IMAGE=github-exporter
IMAGE_VERSION=2.1.0
GOVERSION=1.26

.PHONY: dockerhub clean

dockerhub:
	buildah build --build-arg GOVERSION=${GOVERSION} --build-arg IMAGE_VERSION=${IMAGE_VERSION} -t ${IMAGE}:${IMAGE_VERSION} .
	#buildah login docker.io
	buildah push ${IMAGE}:${IMAGE_VERSION} gounix/${IMAGE}:${IMAGE_VERSION}

clean:
	buildah rmi -f `buildah images --json ${IMAGE}:${IMAGE_VERSION} |jq ".[] | .id" | sed 's/"//g'`

