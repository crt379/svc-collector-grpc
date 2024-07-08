package interceptor

import (
	"context"
	"time"

	"github.com/crt379/svc-collector-grpc/internal/ctxvalue"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	max_resp_len = 1024
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type ServerStreamBox struct {
	grpc.ServerStream
	Ctx *context.Context
}

func (ss *ServerStreamBox) Context() context.Context {
	if ss.Ctx != nil {
		return *ss.Ctx
	}

	return ss.ServerStream.Context()
}

func UnaryLatency(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	start := time.Now()
	resp, err = handler(ctx, req)
	end := time.Now()
	latency := end.Sub(start)

	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("UnaryLatencyInterceptor", zap.String("latency", latency.String()))

	return resp, err
}

func StreamLatency(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	start := time.Now()

	ctx := ss.Context()
	ssb := ServerStreamBox{
		ServerStream: ss,
		Ctx:          &ctx,
	}
	err = handler(srv, &ssb)
	end := time.Now()
	latency := end.Sub(start)

	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("StreamLatencyInterceptor", zap.String("latency", latency.String()))

	return err
}

func UnaryMeta(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		ctx = ctxvalue.GrpcMetaContext{}.NewContext(ctx, &md)
	}

	return handler(ctx, req)
}

func StreamMeta(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	ctx := ss.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		ctx = ctxvalue.GrpcMetaContext{}.NewContext(ctx, &md)
	}

	ssb := ServerStreamBox{
		ServerStream: ss,
		Ctx:          &ctx,
	}

	return handler(srv, &ssb)
}

func UnaryTrace(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	var (
		getvalue []string
		traceid  string
	)

	md, ok := ctxvalue.GrpcMetaContext{}.GetValue(ctx)
	if ok {
		getvalue = md.Get("x-access-trace-id")
		if len(getvalue) > 0 {
			traceid = getvalue[0]
		} else {
			traceid = uuid.New().String()
		}
	} else {
		traceid = uuid.New().String()
	}

	ctx = ctxvalue.TraceContext{}.NewContext(ctx, &traceid)

	return handler(ctx, req)
}

func StreamTrace(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	var (
		getvalue []string
		traceid  string
	)

	ctx := ss.Context()
	md, ok := ctxvalue.GrpcMetaContext{}.GetValue(ctx)
	if ok {
		getvalue = md.Get("x-access-trace-id")
		if len(getvalue) > 0 {
			traceid = getvalue[0]
		} else {
			traceid = uuid.New().String()
		}
	} else {
		traceid = uuid.New().String()
	}

	ctx = ctxvalue.TraceContext{}.NewContext(ctx, &traceid)
	ssb := ServerStreamBox{
		ServerStream: ss,
		Ctx:          &ctx,
	}

	return handler(srv, &ssb)
}

func UnaryTraceSpanLog(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	traceid, _ := ctxvalue.TraceContext{}.GetValue(ctx)
	spanid := uuid.New().String()
	logger = logger.With(zap.String("trace_id", *traceid), zap.String("span_id", spanid))

	logger.Debug("TraceSpanLogInterceptor")

	ctx = ctxvalue.LoggerContext{}.NewContext(ctx, logger)

	return handler(ctx, req)
}

func StreamTraceSpanLog(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	ctx := ss.Context()
	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	traceid, _ := ctxvalue.TraceContext{}.GetValue(ctx)
	spanid := uuid.New().String()
	logger = logger.With(zap.String("trace_id", *traceid), zap.String("span_id", spanid))

	logger.Debug("StreamTraceSpanLogInterceptor")

	ctx = ctxvalue.LoggerContext{}.NewContext(ctx, logger)
	ssb := ServerStreamBox{
		ServerStream: ss,
		Ctx:          &ctx,
	}

	return handler(srv, &ssb)
}

func UnaryReqRepLog(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	var (
		buf  []byte
		jerr error
	)

	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Debug("UnaryReqRepLogInterceptor")

	buf, jerr = json.Marshal(req)
	if jerr != nil {
		logger.Info("req", zap.String("error", jerr.Error()))
	} else {
		logger.Info("req", zap.ByteString("data", buf))
	}

	resp, err = handler(ctx, req)
	if err != nil {
		if st, ok := status.FromError(err); !ok || st.Code() != codes.OK {
			logger.Error("handler", zap.String("error", err.Error()))
		} else {
			err = nil
		}
	}

	logger, _ = ctxvalue.LoggerContext{}.GetValue(ctx)

	buf, jerr = json.Marshal(resp)
	if jerr != nil {
		logger.Info("resp", zap.String("error", jerr.Error()))
	} else if err != nil || len(buf) < max_resp_len {
		logger.Info("resp", zap.ByteString("data", buf))
	}

	return resp, err
}

func StreamHandlerLog(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	ctx := ss.Context()
	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Debug("StreamHandlerLogInterceptor")

	err = handler(srv, ss)
	if err != nil {
		if st, ok := status.FromError(err); !ok || st.Code() != codes.OK {
			logger.Error("handler", zap.String("error", err.Error()), zap.String("FullMethod", info.FullMethod))
		} else {
			err = nil
		}
	}

	return err
}
