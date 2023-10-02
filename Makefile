build:
	go build -o gvs main.go

vet:
	go vet -tests=false ./...

lint:
	golangci-lint run

test:
	go test `go list ./... | grep -v internal`

test-expire-cache:
	go clean -testcache

test-no-cache: test-expire-cache test