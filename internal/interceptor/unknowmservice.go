package interceptor

import (
	"github.com/crt379/svc-collector-grpc/internal/ctxvalue"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnknownServiceHandler(srv any, ss grpc.ServerStream) error {
	ctx := ss.Context()
	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("UnknownService")

	return status.New(codes.Unimplemented, "unknown service").Err()
}
