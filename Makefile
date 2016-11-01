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
	go get github.com/jcelliott/lumber
	go get github.com/mitchellh/go-homedir
	go get github.com/urfave/cli
	go get golang.org/x/crypto/ssh/terminal
	go get golang.org/x/crypto/openpgp/armor
	go get golang.org/x/crypto/openpgp
	go get github.com/kardianos/osext
	go get github.com/speps/go-hashids
	go build -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME} -X main.OS=linux_amd64" -o ${BINARY}

.PHONY: test
test:
	go get -u github.com/schollz/sdees/src
	go get github.com/jcelliott/lumber
	go get github.com/kardianos/osext
	go get github.com/mitchellh/go-homedir
	go get github.com/urfave/cli
	go get github.com/jteeuwen/go-bindata/...
	go get github.com/kardianos/osext
	go get github.com/speps/go-hashids
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

.PHONY: windows
windows:
	rm -rf micro*
	wget https://github.com/zyedidia/micro/releases/download/v1.1.2/micro-1.1.2-win64.zip
	unzip micro*.zip
	mv micro*/micro.exe ./src/bin
	cd src && C:/Go/work/bin/go-bindata.exe ./bin
	rm -rf ./src/bin/micro.exe
	cd src && sed -i -- 's/package main/package sdees/g' bindata.go
	env GOOS=windows GOARCH=386 go build -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME} -X main.OS=windows_amd64" -o sdees.exe
	zip -j sdees_windows_386.zip sdees.exe
	rm sdees.exe
	cd src && git reset --hard HEAD
	rm -f *.zip

.PHONY: release
release:
	go get github.com/kardianos/osext
	go get github.com/aktau/github-release
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
	echo "Uploading Windows 32 latest, bundled with VIM"
	rm -rf micro*
	wget https://github.com/zyedidia/micro/releases/download/v1.1.2/micro-1.1.2-win32.zip
	unzip micro*.zip
	mv micro*/micro.exe ./src/bin
	cd src && $(GOPATH)/bin/go-bindata ./bin
	rm -rf ./src/bin/micro*
	cd src && sed -i -- 's/package main/package sdees/g' bindata.go
	env GOOS=windows GOARCH=386 go build -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME} -X main.OS=windows_386" -o sdees.exe
	zip -j sdees_windows_386.zip sdees.exe
	github-release upload \
			--user schollz \
			--repo sdees \
			--tag ${VERSION} \
			--name "sdees_windows_386.zip" \
			--file sdees_windows_386.zip
	rm sdees.exe
	cd src && git reset --hard HEAD
	rm -f *.zip
	rm -rf micro*
	echo "Uploading Windows 64 latest, bundled with VIM"
	rm -rf micro*
	wget https://github.com/zyedidia/micro/releases/download/v1.1.2/micro-1.1.2-win64.zip
	unzip micro*.zip
	mv micro*/micro.exe ./src/bin
	cd src && $(GOPATH)/bin/go-bindata ./bin
	rm -rf ./src/bin/micro*
	cd src && sed -i -- 's/package main/package sdees/g' bindata.go
	env GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME} -X main.OS=windows_amd64" -o sdees.exe
	zip -j sdees_windows_amd64.zip sdees.exe
	github-release upload \
			--user schollz \
			--repo sdees \
			--tag ${VERSION} \
			--name "sdees_windows_amd64.zip" \
			--file sdees_windows_amd64.zip
	rm sdees.exe
	cd src && git reset --hard HEAD
	rm -f *.zip
	rm -rf micro*
	echo "Uploading Linux Amd64"
	wget https://github.com/zyedidia/micro/releases/download/v1.1.2/micro-1.1.2-linux64.tar.gz
	tar -xvzf micro*.tar.gz
	mv micro*/micro ./src/bin
	cd src && $(GOPATH)/bin/go-bindata ./bin
	rm -rf ./src/bin/micro*
	cd src && sed -i -- 's/package main/package sdees/g' bindata.go
	env GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME} -X main.OS=linux_amd64" -o sdees
	zip -j sdees_linux_amd64.zip sdees
	github-release upload \
	    --user schollz \
	    --repo sdees \
	    --tag ${VERSION} \
	    --name "sdees_linux_amd64.zip" \
	    --file sdees_linux_amd64.zip
	rm sdees
	rm -f *.zip
	rm -rf micro*
	# echo "Uploading Linux Arm"
	# env GOOS=linux GOARCH=arm go build ${LDFLAGS} -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME} -X main.OS=linux_arm" -o sdees
	# zip -j sdees_linux_arm.zip sdees
	# github-release upload \
	#     --user schollz \
	#     --repo sdees \
	#     --tag ${VERSION} \
	#     --name "sdees_linux_arm.zip" \
	#     --file sdees_linux_arm.zip
	# rm sdees
	# rm -f *.zip
	# echo "Uploading Linux Arm64"
	# env GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME} -X main.OS=linux_arm64" -o sdees
	# zip -j sdees_linux_arm64.zip sdees
	# github-release upload \
	#     --user schollz \
	#     --repo sdees \
	#     --tag ${VERSION} \
	#     --name "sdees_linux_arm64.zip" \
	#     --file sdees_linux_arm64.zip
	# rm sdees
	# rm -f *.zip
	# echo "Uploading OSX"
	# env GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME} -X main.OS=osx" -o sdees
	# zip -j sdees_osx.zip sdees
	# github-release upload \
	#     --user schollz \
	#     --repo sdees \
	#     --tag ${VERSION} \
	#     --name "sdees_osx.zip" \
	#     --file sdees_osx.zip
	# rm sdees
	# rm -f *.zip
