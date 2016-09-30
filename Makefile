SOURCEDIR=.

BINARY=sdees

VERSION=1.9.98
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
	go build -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME} -X main.OS=linux_amd64" -o ${BINARY}

.PHONY: test
test:
	go get -u github.com/schollz/sdees/src
	go get github.com/jcelliott/lumber
	go get github.com/mitchellh/go-homedir
	go get github.com/urfave/cli
	go get github.com/jteeuwen/go-bindata/...
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
	rm -rf vim*
	wget ftp://ftp.vim.org/pub/vim/pc/vim80w32.zip
	unzip vim80w32.zip
	mv vim/vim80/vim.exe ./src/bin/
	cd src && $(GOPATH)/bin/go-bindata ./bin
	cd src && sed -i -- 's/package main/package sdees/g' bindata.go
	env GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o sdees-vim.exe
	cd src && git reset --hard HEAD
	rm -rf ./src/bin/vim.exe

.PHONY: latest
latest:
	go get github.com/aktau/github-release
	echo "Deleting old release"
	github-release delete \
	    --user schollz \
	    --repo sdees \
	    --tag latest
	echo "Moving tag"
	git tag --force latest ${BUILD}
	git push --force --tags
	echo "Creating new release"
	github-release release \
	    --user schollz \
	    --repo sdees \
	    --tag latest \
	    --name "Latest" \
	    --description "This is a standalone latest of sdees." \
	    --pre-release
	echo "Uploading Windows 64 latest"
	env GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME} -X main.OS=windows_amd64_novim" -o sdees.exe
	zip -j sdees_windows_amd64_novim.zip sdees.exe
	github-release upload \
	    --user schollz \
	    --repo sdees \
	    --tag latest \
	    --name "sdees_windows_amd64_novim.zip" \
	    --file sdees_windows_amd64_novim.zip
	rm sdees.exe
	rm -f *.zip
	echo "Uploading Windows 64 latest, bundled with VIM"
	rm -rf vim*
	wget ftp://ftp.vim.org/pub/vim/pc/vim80w32.zip
	unzip vim80w32.zip
	mv vim/vim80/vim.exe ./src/bin/
	cd src && $(GOPATH)/bin/go-bindata ./bin
	cd src && sed -i -- 's/package main/package sdees/g' bindata.go
	env GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME} -X main.OS=windows_amd64" -o sdees.exe
	zip -j sdees_windows_amd64.zip sdees.exe
	github-release upload \
			--user schollz \
			--repo sdees \
			--tag latest \
			--name "sdees_windows_amd64.zip" \
			--file sdees_windows_amd64.zip
	rm sdees.exe
	cd src && git reset --hard HEAD
	rm -rf ./src/bin/vim.exe
	rm -f *.zip
	echo "Uploading Linux Amd64"
	env GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME} -X main.OS=linux_amd64" -o sdees
	zip -j sdees_linux_amd64.zip sdees
	github-release upload \
	    --user schollz \
	    --repo sdees \
	    --tag latest \
	    --name "sdees_linux_amd64.zip" \
	    --file sdees_linux_amd64.zip
	rm sdees
	rm -f *.zip
	echo "Uploading Linux Arm"
	env GOOS=linux GOARCH=arm go build ${LDFLAGS} -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME} -X main.OS=linux_arm" -o sdees
	zip -j sdees_linux_arm.zip sdees
	github-release upload \
	    --user schollz \
	    --repo sdees \
	    --tag latest \
	    --name "sdees_linux_arm.zip" \
	    --file sdees_linux_arm.zip
	rm sdees
	rm -f *.zip
	echo "Uploading Linux Arm64"
	env GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME} -X main.OS=linux_arm64" -o sdees
	zip -j sdees_linux_arm64.zip sdees
	github-release upload \
	    --user schollz \
	    --repo sdees \
	    --tag latest \
	    --name "sdees_linux_arm64.zip" \
	    --file sdees_linux_arm64.zip
	rm sdees
	rm -f *.zip
	echo "Uploading OSX"
	env GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME} -X main.OS=osx" -o sdees
	zip -j sdees_osx.zip sdees
	github-release upload \
	    --user schollz \
	    --repo sdees \
	    --tag latest \
	    --name "sdees_osx.zip" \
	    --file sdees_osx.zip
	rm sdees
	rm -f *.zip


.PHONY: release
release:
	go get github.com/aktau/github-release
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
	echo "Uploading Windows 64 latest"
	env GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME} -X main.OS=windows_amd64_novim" -o sdees.exe
	zip -j sdees_windows_amd64_novim.zip sdees.exe
	github-release upload \
	    --user schollz \
	    --repo sdees \
	    --tag ${VERSION} \
	    --name "sdees_windows_amd64_novim.zip" \
	    --file sdees_windows_amd64_novim.zip
	rm sdees.exe
	rm -f *.zip
	echo "Uploading Windows 64 latest, bundled with VIM"
	rm -rf vim*
	wget ftp://ftp.vim.org/pub/vim/pc/vim80w32.zip
	unzip vim80w32.zip
	mv vim/vim80/vim.exe ./src/bin/
	cd src && $(GOPATH)/bin/go-bindata ./bin
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
	rm -rf ./src/bin/vim.exe
	rm -f *.zip
	echo "Uploading Linux Amd64"
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
	echo "Uploading Linux Arm"
	env GOOS=linux GOARCH=arm go build ${LDFLAGS} -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME} -X main.OS=linux_arm" -o sdees
	zip -j sdees_linux_arm.zip sdees
	github-release upload \
	    --user schollz \
	    --repo sdees \
	    --tag ${VERSION} \
	    --name "sdees_linux_arm.zip" \
	    --file sdees_linux_arm.zip
	rm sdees
	rm -f *.zip
	echo "Uploading Linux Arm64"
	env GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME} -X main.OS=linux_arm64" -o sdees
	zip -j sdees_linux_arm64.zip sdees
	github-release upload \
	    --user schollz \
	    --repo sdees \
	    --tag ${VERSION} \
	    --name "sdees_linux_arm64.zip" \
	    --file sdees_linux_arm64.zip
	rm sdees
	rm -f *.zip
	echo "Uploading OSX"
	env GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME} -X main.OS=osx" -o sdees
	zip -j sdees_osx.zip sdees
	github-release upload \
	    --user schollz \
	    --repo sdees \
	    --tag ${VERSION} \
	    --name "sdees_osx.zip" \
	    --file sdees_osx.zip
	rm sdees
	rm -f *.zip

.PHONY: binaries
binaries:
	go get github.com/jteeuwen/go-bindata/...
	rm -rf binaries
	mkdir binaries
	mkdir bin
	$(GOPATH)/bin/go-bindata bin
	## OS X
	env GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o binaries/sdees
	zip -j binaries/sdees_osx_amd64.zip binaries/sdees
	rm binaries/sdees

	## LINUX
	env GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o binaries/sdees
	zip -j binaries/sdees_linux_amd64.zip binaries/sdees
	rm binaries/sdees
	env GOOS=linux GOARCH=arm go build ${LDFLAGS} -o binaries/sdees
	zip -j binaries/sdees_linux_arm.zip binaries/sdees
	rm binaries/sdees
	env GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o binaries/sdees
	zip -j binaries/sdees_linux_arm64.zip binaries/sdees
	rm binaries/sdees
	## WINDOWS
	wget ftp://ftp.vim.org/pub/vim/pc/vim80w32.zip
	unzip vim80w32.zip
	mv vim/vim80/vim.exe ./bin/
	rm -rf vim*
	rm -rf bindata.go
	$(GOPATH)/bin/go-bindata bin
	env GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o binaries/sdees.exe
	zip -j binaries/sdees_windows_amd64.zip binaries/sdees.exe
	rm -rf binaries/vim.exe
	rm -rf ./vim/
	rm binaries/sdees.exe
