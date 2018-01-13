.PHONY: server
server:
	cd cmd/wasgeit-server && go install 

.PHONY: chelper
chelper: 
	cd cmd/crawlerhelper && go install

