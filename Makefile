.PHONY: build clean

build:
	GOOS=linux go build -o main . && zip main.zip main && rm -f main


clean:
	@echo "  >  Cleaning build cache"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean
