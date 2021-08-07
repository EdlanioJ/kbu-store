package handler

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	successRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_stores_success_incoming_messages_total",
		Help: "The total number of success incoming success HTTP requests",
	})
	errorRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_stores_error_incoming_message_total",
		Help: "The total number of error incoming success HTTP requests",
	})
	createRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_stores_create_incoming_requests_total",
		Help: "The total number of incoming create product HTTP requests",
	})
	getRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_stores_get_incoming_requests_total",
		Help: "The total number of incoming get by id store HTTP requests",
	})
	indexRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_stores_index_incoming_requests_total",
		Help: "The total number of incoming index store HTTP requests",
	})
	ativateRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_stores_activate_incoming_requests_total",
		Help: "The total number of incoming activate store HTTP requests",
	})
	blockRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_stores_block_incoming_requests_total",
		Help: "The total number of incoming block store HTTP requests",
	})
	disableRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_stores_disable_incoming_requests_total",
		Help: "The total number of incoming disable store HTTP requests",
	})
	deleteRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_stores_delete_incoming_requests_total",
		Help: "The total number of incoming delete store HTTP requests",
	})
	updateRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_stores_update_incoming_requests_total",
		Help: "The total number of incoming update store HTTP requests",
	})
)
