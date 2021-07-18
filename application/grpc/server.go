package grpc

import (
	"fmt"
	"net"
	"time"

	"github.com/EdlanioJ/kbu-store/application/factory"
	"github.com/EdlanioJ/kbu-store/application/grpc/middleware"
	"github.com/EdlanioJ/kbu-store/application/grpc/pb"
	"github.com/EdlanioJ/kbu-store/application/grpc/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

func StartServer(database *gorm.DB, tc time.Duration, port int) {
	errorInterceptor := middleware.NewErrorInterceptor()
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(errorInterceptor.Unary()))
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

	log.Infof("gRPC server started at port %d", port)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}
