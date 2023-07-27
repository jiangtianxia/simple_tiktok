NAME=simple_tiktok
PKG=simple_tiktok
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
VERSION=git-$(subst /,-,$(BRANCH))-$(shell date +%Y%m%d%H)-$(shell git describe --always --dirty)
IMAGE_TAG=$(VERSION)
IMAGE_REPO=*
LOCAL_REPO=*
OutDir=proto


linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
		go build -v -a -tags netgo -installsuffix netgo -installsuffix cgo -ldflags '-w -s' -ldflags "-X main.Version=$(VERSION)" \
		-o ./build/linux/simple_tiktok $(PKG)

darwin:
	GOOS=darwin GOARCH=amd64 \
		go build -a -tags netgo -installsuffix netgo -ldflags "-X main.Version=$(VERSION)" \
		-o ./build/darwin/simple_tiktok $(PKG)

m1_osx:
	GOOS=darwin GOARCH=arm64 \
		go build -a -tags netgo -installsuffix netgo -ldflags "-X main.Version=$(VERSION)" \
		-o ./build/darwin/simple_tiktok $(PKG)

dist:
	env GOOS=linux GOARCH=amd64 go build -tags=jsoniter -v -o ./build/linux/simple_tiktok

build: linux docs

push: linux docs
	docker build -t ${IMAGE_REPO}/simple_tiktok:${IMAGE_TAG} .
# 	docker push ${IMAGE_REPO}/simple_tiktok:${IMAGE_TAG}

push_aarch64:
	docker build -f Dockerfile.aarch64 -t ${IMAGE_REPO}/aarch64/simple_tiktok:${IMAGE_TAG} .
#	docker push ${IMAGE_REPO}/aarch64/simple_tiktok:${IMAGE_TAG}

push_dev: linux docs
	docker build -t ${LOCAL_REPO}/simple_tiktok:${IMAGE_TAG} .
#   docker push ${LOCAL_REPO}/simple_tiktok:${IMAGE_TAG}

docs:
	swag fmt
	swag init --pd -g ./swagger.go -o ./apidocs/

gen:
	mkdir -p ./internal/${OutDir}
	protoc --go_out=./internal/${OutDir} --go-grpc_out=./internal/${OutDir} ./proto/*.proto

.PHONY: linux darwin m1_osx dist push push_aarch64 push_dev docs gen build