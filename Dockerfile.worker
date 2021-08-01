FROM golang:1.16 AS build

WORKDIR /go/src
COPY . .
ENV PATH="/go/bin:${PATH}"

RUN make env
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app/

COPY --from=build ./go/src/kbu-store .
COPY --from=build ./go/src/app.env .

CMD ["./kbu-store", "http"]
