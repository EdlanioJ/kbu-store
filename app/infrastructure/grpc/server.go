package grpc

import (
	"fmt"
	"net"
	"net/http"

	"github.com/EdlanioJ/kbu-store/app/domain"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/grpc/interceptors"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/grpc/pb"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/grpc/service"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	reg = prometheus.NewRegistry()

	grpcMetrics = grpc_prometheus.NewServerMetrics()
)

func init() {
	reg.MustRegister(grpcMetrics)
}

type grpcServer struct {
	Port         int
	MetricPort   int
	StoreUsecase domain.StoreUsecase
}

func NewGrpcServer() *grpcServer {
	return &grpcServer{}
}

func (s *grpcServer) Serve() {
	errorInterceptor := interceptors.NewErrorInterceptor()

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			errorInterceptor.Unary(),
			grpcMetrics.UnaryServerInterceptor(),
		),
	))
	reflection.Register(grpcServer)

	storeService := service.NewStoreServer(s.StoreUsecase)

	pb.RegisterStoreServiceServer(grpcServer, storeService)

	address := fmt.Sprintf("0.0.0.0:%d", s.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Error(err)
	}

	grpcMetrics.InitializeMetrics(grpcServer)
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	go func() {
		log.Infof("metric server started at port \u001b[92m%d\u001b[0m", s.MetricPort)
		if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", s.MetricPort), nil); err != nil {
			log.Fatal("Unable to start a http server.")
		}
	}()

	log.Infof("gRPC server started at port \u001b[92m%d\u001b[0m", s.Port)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}
