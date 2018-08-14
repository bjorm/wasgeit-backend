.PHONY: server crawler chelper container-server container-crawler

server:
	go install -i github.com/bjorm/wasgeit/cmd/wasgeit-server

crawler:
	go install -i github.com/bjorm/wasgeit/cmd/wasgeit-crawler

chelper:
	go install -i github.com/bjorm/wasgeit/cmd/crawlerhelper

container-server:
	sudo docker build --compress --build-arg MAKE_TARGET=server -t wasgeit/server .

container-crawler:
	sudo docker build --compress --build-arg MAKE_TARGET=crawler -t wasgeit/crawler .
