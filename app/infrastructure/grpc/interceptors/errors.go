package interceptors

import (
	"context"

	"github.com/EdlanioJ/kbu-store/app/domain"
	"github.com/asaskevich/govalidator"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type ErrorInterceptor struct {
}

func NewErrorInterceptor() *ErrorInterceptor {
	return &ErrorInterceptor{}
}
func (i *ErrorInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			return resp, appError(err)
		}

		return resp, nil
	}
}

func appError(err error) error {
	logrus.Error(err)
	if _, ok := err.(govalidator.Errors); ok {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	switch err {
	case domain.ErrNotFound,
		gorm.ErrRecordNotFound:
		return status.Error(codes.NotFound, domain.ErrNotFound.Error())
	case domain.ErrActived,
		domain.ErrBlocked,
		domain.ErrBadParam,
		domain.ErrInactived,
		domain.ErrIsPending:
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Internal, domain.ErrInternal.Error())
	}
}
