.PHONY: lint
lint:
	golangci-lint run -v

.PHONY: build
build:
	go build -o server ./cmd/server
	go build -o shortener ./cmd/shortener

.PHONY: clean
clean:
	rm -f home shortener coverage.out coverage.html
	go mod tidy
	go mod vendor

.PHONY: rebuild
rebuild: clean build

.PHONY: deploy
deploy:
	gcloud app deploy app.yaml shortener.yaml dispatch.yaml

.PHONY: server
server: rebuild
	./server

.PHONY: shortener
shortener: rebuild
	./shortener

.PHONY: test
test:
	go test ./...

.PHONY: coverage
coverage:
	go test -cover ./...

.PHONY: html
html:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

.PHONY: vendor
vendor:
	go mod vendor
