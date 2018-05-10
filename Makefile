.PHONY: server
server:
	go install -i github.com/bjorm/wasgeit/cmd/wasgeit-server

.PHONY: crawler
crawler:
	go install -i github.com/bjorm/wasgeit/cmd/wasgeit-crawler

.PHONY: chelper
chelper:
	go install -i github.com/bjorm/wasgeit/cmd/crawlerhelper