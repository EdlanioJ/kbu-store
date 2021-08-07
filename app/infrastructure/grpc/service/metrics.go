package service

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	successMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "stores_success_incoming_grpc_messages_total",
		Help: "The total number of success incoming success gRPC messages",
	})
	errorMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "stores_error_incoming_grpc_message_total",
		Help: "The total number of error incoming success gRPC messages",
	})
	createMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "stores_create_incoming_grpc_requests_total",
		Help: "The total number of incoming create store gRPC messages",
	})
	getMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "stores_get_incoming_grpc_requests_total",
		Help: "The total number of incoming get by id store gRPC messages",
	})
	listMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "stores_list_incoming_grpc_requests_total",
		Help: "The total number of incoming list stores gRPC messages",
	})
	updateMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "stores_update_incoming_grpc_requests_total",
		Help: "The total number of incoming update store gRPC messages",
	})
	deleteMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "stores_delete_incoming_grpc_requests_total",
		Help: "The total number of incoming delete store gRPC messages",
	})
	activateMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "stores_activate_incoming_grpc_requests_total",
		Help: "The total number of incoming activate store gRPC messages",
	})
	blockMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "stores_block_incoming_grpc_requests_total",
		Help: "The total number of incoming block store gRPC messages",
	})
	disableMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "stores_disable_incoming_grpc_requests_total",
		Help: "The total number of incoming disable store gRPC messages",
	})
)
