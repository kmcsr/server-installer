#!/bin/sh

cd $(dirname $0)

[ -n "$PUBLIC_PREFIX" ] || PUBLIC_PREFIX=craftmine/server-installer

platforms=(linux/amd64 linux/arm64 linux/arm64/v8)

export DOCKER_BUILDKIT=1 # for docker build cache

function build(){
	tag=$1
	platform=$2
	fulltag="${PUBLIC_PREFIX}:${tag}"
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
	build v1 $platform || exit $?
done
