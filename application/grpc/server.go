package grpc

import (
	"fmt"
	"net"
	"net/http"

	"github.com/EdlanioJ/kbu-store/application/grpc/middleware"
	"github.com/EdlanioJ/kbu-store/application/grpc/pb"
	"github.com/EdlanioJ/kbu-store/application/grpc/service"
	"github.com/EdlanioJ/kbu-store/domain"
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
	Port            int
	MetricPort      int
	StoreUsecase    domain.StoreUsecase
	CategoryUsecase domain.CategoryUsecase
}

func NewGrpcServer() *grpcServer {
	return &grpcServer{}
}

func (s *grpcServer) Serve() {
	errorInterceptor := middleware.NewErrorInterceptor()

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			errorInterceptor.Unary(),
			grpcMetrics.UnaryServerInterceptor(),
		),
	))
	reflection.Register(grpcServer)

	categoryGrpcServce := service.NewCategotyServer(s.CategoryUsecase)
	storeService := service.NewStoreServer(s.StoreUsecase)

	pb.RegisterCategoryServiceServer(grpcServer, categoryGrpcServce)
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
