appname = home

.PHONY: lint
lint:
	golangci-lint run -v

.PHONY: build
build:
	go build -o $(appname) ./cmd/server

.PHONY: clean
clean:
	rm -f $(appname) coverage.out coverage.html
	go mod tidy
	go mod vendor

.PHONY: rebuild
rebuild: clean build

.PHONY: deploy
deploy:
	gcloud app deploy

.PHONY: rerun
rerun: rebuild
	./$(appname)

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
