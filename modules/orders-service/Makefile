BINARY_NAME=orders_service
BINARY_UNIX=$(BINARY_NAME)_unix
DOCKER_OWNER=antoineaube
DOCKER_IMAGE=orders-service

all: clean docker-build

deps:
	go get -t -v github.com/gin-gonic/gin
	go get -t -v github.com/streadway/amqp
	go get -t -v github.com/go-pg/pg

clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

build:
	go build -o $(BINARY_NAME) -v

build-linux: $(BINARY_UNIX)
$(BINARY_UNIX):
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BINARY_UNIX) -v

docker-build: $(BINARY_UNIX)
	docker build -t $(DOCKER_OWNER)/$(DOCKER_IMAGE) .
