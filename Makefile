build:
	go build -o gvs cmd/gvs/main.go

run:
	go run cmd/gvs/main.go $(FLAGS)

vet:
	go vet -tests=false ./...

lint:
	golangci-lint run

test:
	go test `go list ./... | grep -v internal`

test-file:
	go test -v $(file)

test-expire-cache:
	go clean -testcache

test-no-cache: test-expire-cache test

docs:
	godoc -http=:6060