# movie-microservice

A WIP go microservice for movie metadata and ratings.

## Code Gen

```sh
# install protocol buffer compiler. Homebrew Sample:
brew install protobuf

# install code generator to generate go code from .proto files
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# export PATH
export PATH="$PATH:$(go env GOPATH)/bin"

# run from the root of the project
protoc -I=./src/contracts --go_out=./src movie.proto
```
