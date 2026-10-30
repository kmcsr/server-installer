#!/bin/sh

cd $(dirname $0)

echo
go version || exit $?
echo

export CGO_ENABLED=0 # cross-compile without cgo

pkg_path=$(go list -m)
tag=$(git describe --tags --match v[0-9]* --abbrev=0 2>/dev/null || git log -1 --format="dev-%H")
echo "Tag version: ${tag}"

function _build(){
	f="minecraft_installer-${GOOS}-${GOARCH}"
	[ "$GOOS" == 'windows' ] && f="${f}.exe"
	echo "==> Building '$f'..."
	go build \
	 -trimpath -gcflags "-trimpath=${GOPATH}" -asmflags "-trimpath=${GOPATH}" \
	 -ldflags "-w -s -X '${pkg_path}.PkgVersion=${tag}'" \
	 -o "./output/$f" "./minecraft_installer"
	return $?
}

available_platforms=(
	darwin/amd64 darwin/arm64
	linux/386 linux/amd64 linux/arm linux/arm64
	windows/386 windows/amd64
)

if [ "$ASYNC_BUILD" = 0 ] || [ "$ASYNC_BUILD" = false ]; then
	ASYNC_BUILD=
else
	echo ">>> Making logs dir 'logs.tmp'"
	mkdir -p 'logs.tmp'
	[ -f .flag.error ] && rm .flag.error
fi


for p in "${available_platforms[@]}"; do
	os=${p%/*}
	arch=${p#*/}
	echo "Building $p..."
	if [ -n "$ASYNC_BUILD" ]; then
		logf="$(pwd)/logs.tmp/${os}-${arch}.log"
		{
			if ! GOOS=$os GOARCH=$arch _build; then
				touch .flag.error
				exit 1
			fi
			echo "[+++] Done" >&2
		} 1>"$logf" &
	else
		GOOS=$os GOARCH=$arch _build || exit $?
	fi
done

if [ -n "$ASYNC_BUILD" ]; then
	echo ">>> waiting works..."
	wait
	if [ -f .flag.error ]; then
		rm .flag.error
		echo "Some error have been happend, please check the console and the logs"
		if [ -n "$CI" ]; then # if it's in github action
			for log in `(ls logs.tmp/*.log)`; do
				echo
				echo "================START ${log}================"
				cat "$log"
				echo "================END ${log}================"
				echo
			done
		fi
		exit 1
	fi
fi

if [ -n "$ASYNC_BUILD" ] && ! [ -n "$KEEP_TMP_LOGS" ]; then
	echo ">>> Clearing logs at 'logs.tmp'"
	rm -rf logs.tmp
fi

echo "==> Done"
