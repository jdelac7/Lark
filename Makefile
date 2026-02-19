.PHONY: build server cli run-server run-cli clean web-install web-dev web-build

build: server cli

server:
	go build -o server/server ./server

cli:
	go build -o cli/cli ./cli

run-server: server
	./server/server

run-cli: cli
	./cli/cli

clean:
	rm -f server/server cli/cli

deps:
	go mod tidy

web-install:
	cd web && npm install

web-dev:
	cd web && npm run dev

web-build:
	cd web && npm run build
