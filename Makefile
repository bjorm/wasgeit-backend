BUILD_COMMIT=$(shell git rev-list HEAD --max-count=1)
BUILD_TIME=$(shell date --iso-8601=seconds)

LD_FLAGS=-ldflags "-X main.BuildCommit=$(BUILD_COMMIT) -X main.BuildTime=$(BUILD_TIME)"

.PHONY: server crawler chelper container-server container-crawler

server:
	go build $(LD_FLAGS) github.com/bjorm/wasgeit/cmd/wasgeit-server

crawler:
	go build $(LD_FLAGS) github.com/bjorm/wasgeit/cmd/wasgeit-crawler

chelper:
	go build $(LD_FLAGS) github.com/bjorm/wasgeit/cmd/crawlerhelper

container-server:
	sudo docker build --compress --build-arg MAKE_TARGET=server -t wasgeit/server .

container-crawler:
	sudo docker build --compress --build-arg MAKE_TARGET=crawler -t wasgeit/crawler .
