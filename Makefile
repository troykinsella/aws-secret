PACKAGE=github.com/troykinsella/aws-secret
BINARY=aws-secret

VERSION=1.0.0

LDFLAGS=-ldflags "-X main.AppVersion=${VERSION}"

build:
	go build ${LDFLAGS} ${COMMAND}

dev-deps:
	go get -u github.com/golang/dep/cmd/dep

deps:
	dep ensure

install:
	go install ${LDFLAGS}

dist: deps
	GOOS=linux  GOARCH=amd64 go build ${LDFLAGS} -o aws-secret_linux_amd64 ${PACKAGE}
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o aws-secret_darwin_amd64 ${PACKAGE}

clean:
	rm ${BINARY} || true
	rm ${BINARY}_* || true

.PHONY: build dev-deps deps install dist clean
