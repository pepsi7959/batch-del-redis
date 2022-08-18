all: build

format:
	go fmt

clean:
	rm ./batch-del-redis

build:
	go build -o batch-del-redis main.go redis_cluster.go

rebuild:
	rm -f ./batch-del-redis && go build -o batch-del-redis main.go redis_cluster.go
