package grpc

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/EdlanioJ/kbu-store/application/factory"
	"github.com/EdlanioJ/kbu-store/application/grpc/middleware"
	"github.com/EdlanioJ/kbu-store/application/grpc/pb"
	"github.com/EdlanioJ/kbu-store/application/grpc/service"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

var (
	reg = prometheus.NewRegistry()

	grpcMetrics = grpc_prometheus.NewServerMetrics()
)

func init() {
	reg.MustRegister(grpcMetrics)
}

func StartServer(database *gorm.DB, tc time.Duration, port int) {
	errorInterceptor := middleware.NewErrorInterceptor()

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			errorInterceptor.Unary(),
			grpcMetrics.UnaryServerInterceptor(),
		),
	))
	reflection.Register(grpcServer)

	categoryUsecase := factory.CategoryUsecase(database, tc)
	storeUsecase := factory.StoreUsecase(database, tc)

	categoryGrpcServce := service.NewCategotyServer(categoryUsecase)
	storeService := service.NewStoreServer(storeUsecase)

	pb.RegisterCategoryServiceServer(grpcServer, categoryGrpcServce)
	pb.RegisterStoreServiceServer(grpcServer, storeService)

	address := fmt.Sprintf("0.0.0.0:%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Error(err)
	}

	grpcMetrics.InitializeMetrics(grpcServer)
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", 3330), nil); err != nil {
			log.Fatal("Unable to start a http server.")
		}
	}()

	log.Infof("gRPC server started at port %d", port)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}
