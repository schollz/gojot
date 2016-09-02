SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=sdees

VERSION=1.1.3
BUILD_TIME=`date +%FT%T%z`
BUILD=`git rev-parse HEAD`

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME}"

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	go get github.com/jcelliott/lumber
	go get github.com/pkg/sftp
	go get github.com/speps/go-hashids
	go get github.com/mitchellh/go-homedir
	go get github.com/urfave/cli
	go get github.com/zyedidia/micro/...
	go get github.com/zyedidia/micro/cmd/micro
	go get -u github.com/jteeuwen/go-bindata/...
	rm -rf bin
	mkdir bin
	cp $(GOPATH)/bin/micro bin/
	rm -rf bin
	go-bindata bin && go build ${LDFLAGS} -o ${BINARY} ${SOURCES}

.PHONY: install
install:
	sudo mv sdees /usr/local/bin/

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	rm -rf binaries
	rm -rf tempsdees

.PHONY: binaries
binaries:
	rm -rf bin
	rm -rf binaries
	rm -rf micro
	mkdir binaries
	mkdir bin
	git clone https://github.com/zyedidia/micro.git && cd micro/cmd/micro && env GOOS=linux GOARCH=amd64 go build -o ../../../bin/micro
	go-bindata bin
	env GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o binaries/sdees
	zip -j binaries/sdees_linux_amd64.zip binaries/sdees
	rm binaries/sdees
	# env GOOS=linux GOARCH=arm go build ${LDFLAGS} -o binaries/sdees
	# zip -j binaries/sdees_linux_arm.zip binaries/sdees
	# rm binaries/sdees
	# env GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o binaries/sdees
	# zip -j binaries/sdees_linux_arm64.zip binaries/sdees
	# rm binaries/sdees
	# wget ftp://ftp.vim.org/pub/vim/pc/vim74w32.zip
	# unzip vim74w32.zip
	# mv vim/vim74/vim.exe ./binaries/
	# rm -rf vim*
	# env GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o binaries/sdees.exe
	# zip -j binaries/sdees_windows_amd64.zip binaries/sdees.exe binaries/vim.exe
	# rm -rf binaries/vim.exe
	# rm -rf ./vim/
	# rm binaries/sdees.exe
