build:
	go build -o gvs main.go

vet:
	go vet -tests=false ./...

test:
	go test ./...