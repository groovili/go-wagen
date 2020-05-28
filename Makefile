install:
	go get -u github.com/gobuffalo/packr/v2/packr2

build:
	${GOBIN}/packr2 build

clean:
	${GOBIN}/packr2 clean

build-linux:
	GOOS=linux GOARCH=amd64 ${GOBIN}/packr2 build

build-win:
	GOOS=windows GOARCH=amd64 ${GOBIN}/packr2 build