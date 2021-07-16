.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/getVm getInstance/fetchId.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/registry registry/addvm.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose
