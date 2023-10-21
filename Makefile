TEST_PATH = $(shell go list ./... | grep -v internal)

build:
	go build -o gvs cmd/gvs/main.go

run:
	go run cmd/gvs/main.go $(FLAGS)

format:
	gofmt -d -s .

vet:
	go vet -tests=false ./...

lint:
	golangci-lint run

test:
	go test $(TEST_PATH)

test-file:
	go test -v $(FILE)

test-coverage:
	go test -cover $(TEST_PATH)

test-coverage-list:
	go test -v -coverpkg=./... -coverprofile=profile.cov $(TEST_PATH) && go tool cover -func profile.cov

docs:
	godoc -http=:6060