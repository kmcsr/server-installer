#!/bin/sh

cd $(dirname $0)

echo
go version || exit $?
echo

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

if [ "$ASYNC_BUILD" = 0 ] || [ "$ASYNC_BUILD" = false ]; then
	ASYNC_BUILD=
else
	echo ">>> Making logs dir 'logs.tmp'"
	mkdir -p 'logs.tmp'
	[ -f .flag.error ] && rm .flag.error
fi

for arch in "${arch32[@]}" "${arch64[@]}"; do
	if [ -n "$ASYNC_BUILD" ]; then
		logf="$(pwd)/logs.tmp/linux-${arch}.log"
		{
			if ! GOOS=linux GOARCH=$arch _build; then
				touch .flag.error
				exit 1
			fi
			echo "[+++] Done" >&2
		} 1>"$logf" &
	else
		GOOS=linux GOARCH=$arch _build || exit $?
	fi
done

for arch in "${arch32[@]}" "${arch64[@]}"; do
	if [ -n "$ASYNC_BUILD" ]; then
		logf="$(pwd)/logs.tmp/windows-${arch}.log"
		{
			if ! GOOS=windows GOARCH=$arch _build; then
				touch .flag.error
				exit 1
			fi
			echo "[+++] Done" >&2
		} 1>"$logf" &
	else
		GOOS=windows GOARCH=$arch _build || exit $?
	fi
done

for arch in "${arch64[@]}"; do
	if [ -n "$ASYNC_BUILD" ]; then
		logf="$(pwd)/logs.tmp/darwin-${arch}.log"
		{
			if ! GOOS=darwin GOARCH=$arch _build; then
				touch .flag.error
				exit 1
			fi
			echo "[+++] Done" >&2
		} 1>"$logf" &
	else
		GOOS=darwin GOARCH=$arch _build || exit $?
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
