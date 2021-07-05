package grpc

import (
	"fmt"
	"net"
	"time"

	"github.com/EdlanioJ/kbu-store/application/factory"
	"github.com/EdlanioJ/kbu-store/application/grpc/pb"
	"github.com/EdlanioJ/kbu-store/application/grpc/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

func StartServer(database *gorm.DB, tc time.Duration, port int) {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	categoryUsecase := factory.CategoryUsecase(database, tc)
	tagUsecase := factory.TagUsecase(database, tc)

	categoryGrpcServce := service.NewCategotyServer(categoryUsecase)
	tagGrpcService := service.NewTagServer(tagUsecase)

	pb.RegisterCategoryServiceServer(grpcServer, categoryGrpcServce)
	pb.RegisterTagServiceServer(grpcServer, tagGrpcService)

	address := fmt.Sprintf("0.0.0.0:%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Error(err)
	}

	log.Infof("gRPC server has been started on port %d", port)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}
