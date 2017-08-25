# Make a release with
# make -j4 release

VERSION=$(shell git describe)
LDFLAGS=-ldflags "-s -w -X main.version=${VERSION}"

.PHONY: build
build:
	go-bindata -debug src/bundle/vim.exe
	go build ${LDFLAGS} -o dist/gojot
	rm -rf bindata.go

.PHONY: linuxarm
linuxarm:
	go-bindata -debug src/bundle/vim.exe
	env GOOS=linux GOARCH=arm go build ${LDFLAGS} -o dist/gojot_linux_arm
	rm -rf bindata.go
	# cd dist && upx --brute gojot_linux_arm

.PHONY: linux64
linux64:
	go-bindata -debug src/bundle/vim.exe
	env GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o dist/gojot_linux_amd64
	cd dist && upx --brute gojot_linux_amd64
	rm -rf bindata.go

.PHONY: windows
windows:
	go-bindata src/bundle/vim.exe
	env GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o dist/gojot_windows_amd64.exe
	rm -rf bindata.go
	# cd dist && upx --brute gojot_windows_amd64.exe

.PHONY: osx
osx:
	go-bindata -debug src/bundle/vim.exe
	env GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o dist/gojot_osx_amd64
	rm -rf bindata.go
	# cd dist && upx --brute gojot_osx_amd64

.PHONY: release
release: osx windows linux64 linuxarm