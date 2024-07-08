package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"syscall"

	"github.com/crt379/svc-collector-grpc/internal/config"
	"github.com/crt379/svc-collector-grpc/internal/interceptor"
	"github.com/crt379/svc-collector-grpc/internal/logging"
	"github.com/crt379/svc-collector-grpc/internal/server/appapi"
	"github.com/crt379/svc-collector-grpc/internal/server/application"
	"github.com/crt379/svc-collector-grpc/internal/server/appproc"
	"github.com/crt379/svc-collector-grpc/internal/server/appsvc"
	"github.com/crt379/svc-collector-grpc/internal/server/processor"
	"github.com/crt379/svc-collector-grpc/internal/server/register"
	"github.com/crt379/svc-collector-grpc/internal/server/service"
	"github.com/crt379/svc-collector-grpc/internal/server/svcapi"
	"github.com/crt379/svc-collector-grpc/internal/server/svcapieg"
	"github.com/crt379/svc-collector-grpc/internal/server/tenant"

	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func main() {
	defer logging.LoggerSync()

	logger := zap.L()
	logger.Info("starting")

	zlogger := logging.NewZapLogger(logger)
	grpclog.SetLoggerV2(zlogger)

	panicsTotal := promauto.NewCounter(interceptor.PanicsTotal)

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptor.UnaryMeta,
			interceptor.UnaryTrace,
			interceptor.UnaryTraceSpanLog,
			interceptor.UnaryLatency,
			interceptor.UnaryReqRepLog,
			interceptor.WithUnaryPrometheus(),
			interceptor.FWithUnaryRecovery(panicsTotal.Inc),
		),
		grpc.ChainStreamInterceptor(
			interceptor.StreamMeta,
			interceptor.StreamTrace,
			interceptor.StreamTraceSpanLog,
			interceptor.StreamLatency,
			interceptor.StreamHandlerLog,
			interceptor.WithStreamPrometheus(),
			interceptor.FWithStreamRecovery(panicsTotal.Inc),
		),
		grpc.UnknownServiceHandler(interceptor.UnknownServiceHandler),
	}

	srv := grpc.NewServer(opts...)
	tenant.RegisterServer(srv)
	service.RegisterServer(srv)
	svcapi.RegisterServer(srv)
	svcapieg.RegisterServer(srv)
	application.RegisterServer(srv)
	appsvc.RegisterServer(srv)
	register.RegisterServer(srv)
	processor.RegisterServer(srv)
	appapi.RegisterServer(srv)
	appproc.RegisterServer(srv)

	g := &run.Group{}

	g.Add(func() error {
		listenAddr := fmt.Sprintf("%s:%s", config.AppConfig.Listen.Host, config.AppConfig.Listen.Port)
		logger.Info("starting listen addr: " + listenAddr)

		lis, err := net.Listen("tcp", listenAddr)
		if err != nil {
			return err
		}
		return srv.Serve(lis)
	}, func(error) {
		srv.GracefulStop()
		srv.Stop()
	})

	httpaddr := fmt.Sprintf("%s:%s", config.AppConfig.Prometheus.Host, config.AppConfig.Prometheus.Port)
	httpSrv := &http.Server{Addr: httpaddr}
	g.Add(func() error {
		m := http.NewServeMux()
		m.Handle("/metrics", promhttp.Handler())
		httpSrv.Handler = m
		logger.Info("starting HTTP server addr: " + httpaddr)

		return httpSrv.ListenAndServe()
	}, func(error) {
		if err := httpSrv.Close(); err != nil {
			logger.Info("failed to stop server", zap.String("error", err.Error()))
		}
	})

	c := make(chan struct{})
	g.Add(func() error {
		logger.Info("InternalRegister")
		err := register.InternalRegister(
			context.Background(),
			config.AppConfig.Register.Name,
			config.AppConfig.Addr,
		)
		if err == nil {
			<-c
		}
		return err
	}, func(error) {
		close(c)
	})

	g.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))

	if err := g.Run(); err != nil {
		logger.Info("g run error", zap.Any("error", err))
		os.Exit(1)
	}
}
