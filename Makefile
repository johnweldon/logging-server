IMAGE=docker.jw4.us/logsrv

ifeq ($(REVISION),)
	REVISION=$(shell git describe --dirty --first-parent --always --tags)
endif

all: image

clean:
	-rm ./logsrv
	go clean .

image:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags netgo -ldflags="-s -w" -o logsrv .
	docker build -t $(IMAGE):latest -t $(IMAGE):$(REVISION) .

push: clean image
	docker push $(IMAGE):$(REVISION)
	docker push $(IMAGE):latest

