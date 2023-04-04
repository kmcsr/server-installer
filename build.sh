#!/bin/sh

cd $(dirname $0)

export CGO_ENABLED=0 # cross-compile without cgo

function _build(){
	f="minecraft_installer-${GOOS}-${GOARCH}"
	[ "$GOOS" == 'windows' ] && f="${f}.exe"
	echo "==> Building '$f'..."
	go build\
	 -trimpath -gcflags "-trimpath=${GOPATH}" -asmflags "-trimpath=${GOPATH}" -ldflags "-w -s"\
	 -o "./output/$f" "./cli"
	return $?
}

arch32=(386 arm)
arch64=(amd64 arm64)

for arch in "${arch32[@]}" "${arch64[@]}"; do
	GOOS=linux GOARCH=$arch _build || exit $?
done

for arch in "${arch32[@]}" "${arch64[@]}"; do
	GOOS=windows GOARCH=$arch _build || exit $?
done

for arch in "${arch64[@]}"; do
	GOOS=darwin GOARCH=$arch _build || exit $?
done

echo "==> Done"
