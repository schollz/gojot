# Make a release with
# make -j4 release

VERSION=$(shell git describe)
LDFLAGS=-ldflags "-s -w -X main.version=${VERSION}"

.PHONY: build
build:
	go build ${LDFLAGS} -o dist/gojot

.PHONY: linuxarm
linuxarm:
	env GOOS=linux GOARCH=arm go build ${LDFLAGS} -o dist/gojot_linux_arm
	# cd dist && upx --brute gojot_linux_arm

.PHONY: linux64
linux64:
	env GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o dist/gojot_linux_amd64
	cd dist && upx --brute gojot_linux_amd64

.PHONY: windows
windows:
	env GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o dist/gojot_windows_amd64.exe
	# cd dist && upx --brute gojot_windows_amd64.exe

.PHONY: osx
osx:
	env GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o dist/gojot_osx_amd64
	# cd dist && upx --brute gojot_osx_amd64

.PHONY: release
release: osx windows linux64 linuxarm