package interceptor

import (
	"context"

	"github.com/crt379/svc-collector-grpc/internal/ctxvalue"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryRecovery(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ any, err error) {
	return FWithUnaryRecovery()(ctx, req, info, handler)
}

func FWithUnaryRecovery(fs ...func()) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ any, err error) {
		defer func() {
			logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
			if r := recover(); r != nil {
				for _, f := range fs {
					f()
				}
				err = status.New(codes.Internal, "服务内部错误").Err()
				logger.Error("panic", zap.Any("error", r))
			}
		}()

		return handler(ctx, req)
	}
}

func StreamRecovery(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	return FWithStreamRecovery()(srv, ss, info, handler)
}

func FWithStreamRecovery(fs ...func()) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			logger, _ := ctxvalue.LoggerContext{}.GetValue(ss.Context())
			if r := recover(); r != nil {
				for _, f := range fs {
					f()
				}
				err = status.New(codes.Internal, "服务内部错误").Err()
				logger.Error("panic", zap.Any("error", r))
			}
		}()

		return handler(srv, ss)
	}
}
