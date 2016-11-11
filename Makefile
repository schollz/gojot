SOURCEDIR=.

BINARY=sdees

VERSION=2.0.0
BUILD_TIME=`date +%FT%T%z`
BUILD=`git rev-parse HEAD`
BUILDSHORT = `git rev-parse --short HEAD`
OS=something
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME}"

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	go build -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME} -X main.OS=linux_amd64" -o ${BINARY}

.PHONY: update
update:
	go get -u -v github.com/codahale/chacha20
	go get -u -v gopkg.in/cheggaaa/pb.v1
	go get -u -v github.com/jcelliott/lumber
	go get -u -v github.com/mitchellh/go-homedir
	go get -u -v github.com/urfave/cli
	go get -u -v golang.org/x/crypto/ssh/terminal
	go get -u -v golang.org/x/crypto/openpgp/armor
	go get -u -v golang.org/x/crypto/openpgp
	go get -u -v github.com/kardianos/osext
	go get -u -v github.com/aktau/github-release
	go get -u -v github.com/jteeuwen/go-bindata/...
	go get -u -v github.com/olekukonko/tablewriter

.PHONY: test
test:
	cd src && $(GOPATH)/bin/go-bindata bin
	cd src && sed -i -- 's/package main/package sdees/g' bindata.go
	cd src && go test -v -cover

.PHONY: cloc
cloc:
	echo `grep -v "^$$" src/*.go | grep -v "//" | wc -l` lines of code
	echo `grep -v "^$$" src/*_test.go | grep -v "//" | wc -l` lines of testing code


.PHONY: install
install:
	sudo mv sdees /usr/local/bin/

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	rm -rf binaries
	rm -rf vim
	rm -rf tempsdees
	rm -rf src/gittest
	rm -rf src/test
	rm -rf src/gittest10

.PHONY: release
release:
	echo "Deleting old release"
	git tag -d ${VERSION};
	git push origin :${VERSION};
	github-release delete \
			--user schollz \
			--repo sdees \
			--tag ${VERSION}
	echo "Moving tag"
	git tag --force latest ${BUILD}
	git push --force --tags
	echo "Creating new release"
	github-release release \
	    --user schollz \
	    --repo sdees \
	    --tag ${VERSION} \
	    --name "${VERSION}" \
	    --description "This is a standalone latest of sdees."
	echo "Uploading Windows 32"
	rm -rf micro*
	wget https://github.com/zyedidia/micro/releases/download/v1.1.2/micro-1.1.2-win32.zip
	unzip micro*.zip
	mv micro*/micro.exe ./src/bin
	cd src && $(GOPATH)/bin/go-bindata ./bin
	rm -rf ./src/bin/micro*
	cd src && sed -i -- 's/package main/package sdees/g' bindata.go
	env GOOS=windows GOARCH=386 go build -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME} -X main.OS=win32" -o sdees.exe
	zip -j sdees-${VERSION}-win32.zip sdees.exe README.md LICENSE
	github-release upload \
			--user schollz \
			--repo sdees \
			--tag ${VERSION} \
			--name "sdees-${VERSION}-win32.zip" \
			--file sdees-${VERSION}-win32.zip
	rm sdees.exe
	cd src && git reset --hard HEAD
	rm -f *.zip
	rm -rf micro*
	echo "Uploading Windows 64"
	rm -rf micro*
	wget https://github.com/zyedidia/micro/releases/download/v1.1.2/micro-1.1.2-win64.zip
	unzip micro*.zip
	mv micro*/micro.exe ./src/bin
	wget ftp://ftp.vim.org/pub/vim/pc/vim80-069w32.zip
	unzip vim*zip
	mv vim/vim80/vim.exe ./src/bin/
	cd src && $(GOPATH)/bin/go-bindata ./bin
	rm -rf ./src/bin/micro*
	cd src && sed -i -- 's/package main/package sdees/g' bindata.go
	env GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME} -X main.OS=win64" -o sdees.exe
	zip -j sdees-${VERSION}-win64.zip sdees.exe README.md LICENSE
	github-release upload \
			--user schollz \
			--repo sdees \
			--tag ${VERSION} \
			--name "sdees-${VERSION}-win64.zip" \
			--file sdees-${VERSION}-win64.zip
	rm sdees.exe
	cd src && git reset --hard HEAD
	rm -f *.zip
	rm -rf micro*
	rm -rf vim80*zip
	echo "Uploading Linux 64"
	wget https://github.com/zyedidia/micro/releases/download/v1.1.2/micro-1.1.2-linux64.tar.gz
	tar -xvzf micro*.tar.gz
	mv micro*/micro ./src/bin
	cd src && $(GOPATH)/bin/go-bindata ./bin
	rm -rf ./src/bin/micro*
	cd src && sed -i -- 's/package main/package sdees/g' bindata.go
	env GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME} -X main.OS=linux64" -o sdees
	tar -czvf sdees-${VERSION}-linux64.tar.gz sdees README.md LICENSE
	github-release upload \
	    --user schollz \
	    --repo sdees \
	    --tag ${VERSION} \
	    --name "sdees-${VERSION}-linux64.tar.gz" \
	    --file sdees-${VERSION}-linux64.tar.gz
	rm sdees
	rm -f *.tar.gz
	rm -rf micro*
	echo "Uploading Linux 32"
	wget https://github.com/zyedidia/micro/releases/download/v1.1.2/micro-1.1.2-linux32.tar.gz
	tar -xvzf micro*.tar.gz
	mv micro*/micro ./src/bin
	cd src && $(GOPATH)/bin/go-bindata ./bin
	rm -rf ./src/bin/micro*
	cd src && sed -i -- 's/package main/package sdees/g' bindata.go
	env GOOS=linux GOARCH=386 go build -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME} -X main.OS=linux32" -o sdees
	tar -czvf sdees-${VERSION}-linux32.tar.gz sdees README.md LICENSE
	github-release upload \
	    --user schollz \
	    --repo sdees \
	    --tag ${VERSION} \
	    --name "sdees-${VERSION}-linux32.tar.gz" \
	    --file sdees-${VERSION}-linux32.tar.gz
	rm sdees
	rm -f *.tar.gz
	rm -rf micro*
	echo "---------- Uploading Linux ARM ------------"
	wget https://github.com/zyedidia/micro/releases/download/v1.1.2/micro-1.1.2-linux-arm.tar.gz
	tar -xvzf micro*.tar.gz
	mv micro*/micro ./src/bin
	cd src && $(GOPATH)/bin/go-bindata ./bin
	rm -rf ./src/bin/micro*
	cd src && sed -i -- 's/package main/package sdees/g' bindata.go
	env GOOS=linux GOARCH=arm go build -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME} -X main.OS=linux-arm" -o sdees
	tar -czvf sdees-${VERSION}-linux-arm.tar.gz sdees README.md LICENSE
	github-release upload \
	    --user schollz \
	    --repo sdees \
	    --tag ${VERSION} \
	    --name "sdees-${VERSION}-linux-arm.tar.gz" \
	    --file sdees-${VERSION}-linux-arm.tar.gz
	rm sdees
	rm -f *.tar.gz
	rm -rf micro*
	echo "---------- Uploading OSX ------------"
	wget https://github.com/zyedidia/micro/releases/download/v1.1.2/micro-1.1.2-osx.tar.gz
	tar -xvzf micro*.tar.gz
	mv micro*/micro ./src/bin
	cd src && $(GOPATH)/bin/go-bindata ./bin
	rm -rf ./src/bin/micro*
	cd src && sed -i -- 's/package main/package sdees/g' bindata.go
	env GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME} -X main.OS=osx" -o sdees
	tar -czvf sdees-${VERSION}-osx.tar.gz sdees README.md LICENSE
	github-release upload \
	    --user schollz \
	    --repo sdees \
	    --tag ${VERSION} \
	    --name "sdees-${VERSION}-osx.tar.gz" \
	    --file sdees-${VERSION}-osx.tar.gz
	rm sdees
	rm -f *.tar.gz
	rm -rf micro*
	echo "---------- Uploading freebsd32 ------------"
	wget https://github.com/zyedidia/micro/releases/download/v1.1.2/micro-1.1.2-freebsd32.tar.gz
	tar -xvzf micro*.tar.gz
	mv micro*/micro ./src/bin
	cd src && $(GOPATH)/bin/go-bindata ./bin
	rm -rf ./src/bin/micro*
	cd src && sed -i -- 's/package main/package sdees/g' bindata.go
	env GOOS=freebsd GOARCH=386 go build -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME} -X main.OS=freebsd" -o sdees
	tar -czvf sdees-${VERSION}-freebsd32.tar.gz sdees README.md LICENSE
	github-release upload \
	    --user schollz \
	    --repo sdees \
	    --tag ${VERSION} \
	    --name "sdees-${VERSION}-freebsd32.tar.gz" \
	    --file sdees-${VERSION}-freebsd32.tar.gz
	rm sdees
	rm -f *.tar.gz
	rm -rf micro*
	echo "---------- Uploading freebsd64 ------------"
	wget https://github.com/zyedidia/micro/releases/download/v1.1.2/micro-1.1.2-freebsd64.tar.gz
	tar -xvzf micro*.tar.gz
	mv micro*/micro ./src/bin
	cd src && $(GOPATH)/bin/go-bindata ./bin
	rm -rf ./src/bin/micro*
	cd src && sed -i -- 's/package main/package sdees/g' bindata.go
	env GOOS=freebsd GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME} -X main.OS=freebsd" -o sdees
	tar -czvf sdees-${VERSION}-freebsd64.tar.gz sdees README.md LICENSE
	github-release upload \
	    --user schollz \
	    --repo sdees \
	    --tag ${VERSION} \
	    --name "sdees-${VERSION}-freebsd64.tar.gz" \
	    --file sdees-${VERSION}-freebsd64.tar.gz
	rm sdees
	rm -f *.tar.gz
	rm -rf micro*
	echo "---------- Uploading openbsd32 ------------"
	wget https://github.com/zyedidia/micro/releases/download/v1.1.2/micro-1.1.2-openbsd32.tar.gz
	tar -xvzf micro*.tar.gz
	mv micro*/micro ./src/bin
	cd src && $(GOPATH)/bin/go-bindata ./bin
	rm -rf ./src/bin/micro*
	cd src && sed -i -- 's/package main/package sdees/g' bindata.go
	env GOOS=openbsd GOARCH=386 go build -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME} -X main.OS=openbsd" -o sdees
	tar -czvf sdees-${VERSION}-openbsd32.tar.gz sdees README.md LICENSE
	github-release upload \
	    --user schollz \
	    --repo sdees \
	    --tag ${VERSION} \
	    --name "sdees-${VERSION}-openbsd32.tar.gz" \
	    --file sdees-${VERSION}-openbsd32.tar.gz
	rm sdees
	rm -f *.tar.gz
	rm -rf micro*
	echo "---------- Uploading openbsd64 ------------"
	wget https://github.com/zyedidia/micro/releases/download/v1.1.2/micro-1.1.2-openbsd64.tar.gz
	tar -xvzf micro*.tar.gz
	mv micro*/micro ./src/bin
	cd src && $(GOPATH)/bin/go-bindata ./bin
	rm -rf ./src/bin/micro*
	cd src && sed -i -- 's/package main/package sdees/g' bindata.go
	env GOOS=openbsd GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME} -X main.OS=openbsd" -o sdees
	tar -czvf sdees-${VERSION}-openbsd64.tar.gz sdees README.md LICENSE
	github-release upload \
	    --user schollz \
	    --repo sdees \
	    --tag ${VERSION} \
	    --name "sdees-${VERSION}-openbsd64.tar.gz" \
	    --file sdees-${VERSION}-openbsd64.tar.gz
	rm sdees
	rm -f *.tar.gz
	rm -rf micro*
	git reset --hard HEAD
	git pull
