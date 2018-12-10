BINARY_NAME=orders_service
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build
build:
	cd src && go build -o ../$(BINARY_NAME) -v
clean:
	cd src && go clean
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
deps:
	go get github.com/gin-gonic/gin
	go get github.com/streadway/amqp
	go get github.com/go-pg/pg

# Cross compilation
build-linux: $(BINARY_UNIX)

$(BINARY_UNIX):
	cd src && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ../$(BINARY_UNIX) -v
docker-build: $(BINARY_UNIX)
	docker build -t orders-service .