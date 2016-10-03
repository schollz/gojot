SOURCEDIR=.

BINARY=sdees

VERSION=1.3.0
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
	go get github.com/jteeuwen/go-bindata/...
	rm -rf bin
	mkdir bin
	$(GOPATH)/bin/go-bindata bin
	go build ${LDFLAGS} -o ${BINARY}

.PHONY: install
install:
	sudo mv sdees /usr/local/bin/

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	rm -rf binaries
	rm -rf tempsdees

.PHONY: windows
windows:
	rm -rf bin
	mkdir bin
	wget ftp://ftp.vim.org/pub/vim/pc/vim80w32.zip
	unzip vim80w32.zip
	mv vim/vim80/vim.exe ./bin/
	rm -rf vim*
	rm -rf bindata.go
	$(GOPATH)/bin/go-bindata bin
	env GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o binaries/sdees.exe


.PHONY: release
release:
	go get github.com/kardianos/osext
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
	    --description "This is a standalone latest of sdees 1.X."
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
	echo "Uploading Windows 32 latest, bundled with VIM"
	rm -rf vim*
	wget ftp://ftp.vim.org/pub/vim/pc/vim80w32.zip
	unzip vim80w32.zip
	mv vim/vim80/vim.exe ./src/bin/
	cd src && $(GOPATH)/bin/go-bindata ./bin
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
	rm -rf ./src/bin/vim.exe
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
	rm -rf bin
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
	rm -rf bin
	mkdir bin
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
	rm -rf ./bin/
	rm -rf bindata.go
	rm binaries/sdees.exe
