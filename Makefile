.PHONY: server
server:
	go install -i github.com/bjorm/wasgeit/cmd/wasgeit-server

.PHONY: chelper
chelper:
	go install -i github.com/bjorm/wasgeit/cmd/crawler-helper