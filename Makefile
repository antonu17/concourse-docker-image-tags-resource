# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build

all: clean build

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o assets/check github.com/antonu17/concourse-docker-image-tags-resource/cmd/check
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o assets/in github.com/antonu17/concourse-docker-image-tags-resource/cmd/in
clean:
	rm -f assets/check
	rm -f assets/in
