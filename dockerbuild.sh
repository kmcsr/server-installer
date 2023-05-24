#!/bin/sh

cd $(dirname $0)

[ -n "$PUBLIC_PREFIX" ] || PUBLIC_PREFIX=craftmine/server-installer

[ -n "$TAG" ] || { echo 'Error: You must give a env TAG' ; exit 1; }
# TAG=refs/tags/v1.2.0
TAG=$(basename ${TAG})

platforms=(linux/amd64 linux/arm64 linux/arm64/v8)

export DOCKER_BUILDKIT=1 # for docker build cache

function build(){
	platform=$1
	fulltag="${PUBLIC_PREFIX}:${TAG}"
	echo
	echo "==> building $fulltag for $platform"
	echo
	docker build --platform ${platform} \
	 --tag "$fulltag" \
	 . || return $?
	echo
	echo "==> tag $fulltag as ${PUBLIC_PREFIX}:latest"
	docker tag "$fulltag" "${PUBLIC_PREFIX}:latest" || return $?
	echo
	echo "==> pushing $fulltag"
	echo
	docker push "$fulltag" || return $?
	echo
	echo "==> pushing ${PUBLIC_PREFIX}:latest"
	echo
	docker push "${PUBLIC_PREFIX}:latest" || return $?
	return 0
}

for platform in "${platforms[@]}"; do
	build $platform || exit $?
done
