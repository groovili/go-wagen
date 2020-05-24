install:
	go get -u github.com/gobuffalo/packr/v2/packr2

build:
	${GOBIN}/packr2 build

clean:
	${GOBIN}/packr2 clean