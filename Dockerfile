FROM golang:1.16

WORKDIR /go/src
ENV PATH="/go/bin:${PATH}"
ENV GO111MODULE=on
ENV CGO_ENABLED=1

RUN apt-get update && \
  apt-get install build-essential protobuf-compiler -y && \
  go get github.com/spf13/cobra/cobra && \
  go get -u github.com/swaggo/swag/cmd/swag && \
  go get google.golang.org/grpc/cmd/protoc-gen-go-grpc && \
  go get github.com/vektra/mockery/v2/.../ && \
  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest


CMD ["tail", "-f", "/dev/null"]
