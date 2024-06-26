# movie-microservice

A sample go microservice.

## Stack:

- [Go](https://go.dev/)
- [pgx](https://github.com/jackc/pgx)
- [Consul](https://github.com/hashicorp/consul)
- [gRPC](https://github.com/grpc/grpc-go)

## Requirements

Consul (for service discovery):

```sh
docker run -d -p 8500:8500 -p 8600:8600/udp --name=dev-consul hashicorp/consul agent -server -ui -node=server-1 -bootstrap-expect=1 -client=0.0.0.0
```

## Code Gen

```sh
# install protocol buffer compiler. Homebrew Sample:
brew install protobuf

# install code generator to generate go code from .proto files
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# export PATH
export PATH="$PATH:$(go env GOPATH)/bin"

# run from the root of the project. generates movie.pb.go
protoc -I=./contracts --go_out=. movie.proto

# run from the root of the project. generates movie_grpc.pb.go
protoc -I=./contracts --go_out=. --go-grpc_out=. movie.proto
```

## Calling gRPC Servers

Install grpcurl CLI

```sh
# basically curl for gRPC servers
brew install grpcurl

# sample usage
grpcurl -plaintext -d '{"record_id":"1", "record_type": "movie"}' localhost:8082 RatingService/GetAggregatedRatings
```
