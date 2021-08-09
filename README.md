<h1 align="center">KBU Store Microservice</h1>

[![GitHub release](https://img.shields.io/github/tag/EdlanioJ/kbu-store.svg)](https://GitHub.com/EdlanioJ/kbu-store/releases/) 
[![Coverage Status](https://coveralls.io/repos/github/EdlanioJ/kbu-store/badge.svg?branch=main)](https://coveralls.io/github/EdlanioJ/kbu-store?branch=main)

</hr>

<br>
<p>This is a microservice example. I have been studying Golang for a while and I am doing this project to show what I have learned in process.</p>
<br/>

# Technologies
This Project was developed with the following technologies:
* [Golang v1.16](https://golang.org/) - Stable version of go
* [Fiber](https://docs.gofiber.io/) -  Web framework
* [gRPC](https://grpc.io/) - RPC framework
* [Kafka](https://github.com/segmentio/kafka-go) - Kafka library in Go
* [Docker](https://www.docker.com/) - Docker
* [Prometheus](https://prometheus.io/) - Prometheus
* [Grafana](https://grafana.com/) - Grafana
* [Jaeger](https://www.jaegertracing.io/) - Jaeger tracing
* [Gorm](https://gorm.io) - ORM library for Golang

<br/>

# Getting Started

## Prerequisits
Install:
* [Git](https://git-scm.com/)
* [Docker](https://www.docker.com/)


## Setup

```bash
git clone https://github.com/EdlanioJ/kbu-store.git

cd kbu-store

docker-compose up -d

docker-compose exec app bash

# Database migration

make migrate.up

 # Create env file

 make env
```

## Run
 ```bash
 docker-compose exec app bash

# Start http server

make start.http

# Start gRPC server

make start.grpc
  ```

<b>Swagger UI:</b>
- http://localhost:3333/api/v1/docs/index.html

<b>Jaeger UI:</b>
- http://localhost:16686

<b>Prometheus UI:</b>
- http://localhost:9090

<b>Grafana UI:</b>
- http://localhost:3000


<b>Kafka Control Center:</b>
- http://localhost:9021

<b>gRPC Service:</b>
```bash 
docker-compose exec app bash

evans -r repl
```