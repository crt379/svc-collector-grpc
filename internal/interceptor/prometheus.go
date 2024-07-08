package interceptor

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/crt379/svc-collector-grpc/internal/ctxvalue"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
)

var (
	once sync.Once

	Labels = []string{"grpc_service", "grpc_method"}

	PanicsTotal = prometheus.CounterOpts{
		Name: "grpc_req_panics_recovered_total",
		Help: "Total number of gRPC requests recovered from internal panic.",
	}
	StartedCounter = prometheus.CounterOpts{
		Name: "grpc_server_started_total",
		Help: "Total number of RPCs started on the server.",
	}
	HandledCounter = prometheus.CounterOpts{
		Name: "grpc_server_handled_total",
		Help: "Total number of RPCs completed on the server, regardless of success or failure.",
	}
	StreamMsgReceived = prometheus.CounterOpts{
		Name: "grpc_server_msg_received_total",
		Help: "Total number of RPC stream messages received on the server.",
	}
	StreamMsgSent = prometheus.CounterOpts{
		Name: "grpc_server_msg_sent_total",
		Help: "Total number of gRPC stream messages sent by the server.",
	}
	HandledHistogram = prometheus.HistogramOpts{
		Name:    "grpc_server_handling_seconds",
		Help:    "Histogram of response latency (seconds) of gRPC that had been application-level handled by the server.",
		Buckets: []float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120},
	}
)

var (
	pc pCounterBox
)

type pCounterBox struct {
	StartedCounterC, HandledCounterC, StreamMsgReceivedC, StreamMsgSentC *prometheus.CounterVec
	HandledHistogramC                                                    *prometheus.HistogramVec
}

func _init() {
	pc.StartedCounterC = promauto.NewCounterVec(StartedCounter, Labels)
	pc.HandledCounterC = promauto.NewCounterVec(HandledCounter, Labels)
	pc.StreamMsgReceivedC = promauto.NewCounterVec(StreamMsgReceived, Labels)
	pc.StreamMsgSentC = promauto.NewCounterVec(StreamMsgSent, Labels)
	pc.HandledHistogramC = promauto.NewHistogramVec(HandledHistogram, Labels)
}

func ef(ctx context.Context) prometheus.Labels {
	traceid, ok := ctxvalue.TraceContext{}.GetValue(ctx)
	if ok {
		return prometheus.Labels{"traceID": *traceid}
	}
	return nil
}

func PregWithUnaryPrometheus(metrics *grpcprom.ServerMetrics) grpc.UnaryServerInterceptor {
	return metrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(ef))
}

func PregWithStreamPrometheus(metrics *grpcprom.ServerMetrics) grpc.StreamServerInterceptor {
	return metrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(ef))
}

func splitFullMethodName(fullMethod string) (string, string) {
	fullMethod = strings.TrimPrefix(fullMethod, "/") // remove leading slash
	if i := strings.Index(fullMethod, "/"); i >= 0 {
		return fullMethod[:i], fullMethod[i+1:]
	}
	return "unknown", "unknown"
}

func incrementWithExemplar(c *prometheus.CounterVec, lable prometheus.Labels, lvals ...string) {
	c.WithLabelValues(lvals...).(prometheus.ExemplarAdder).AddWithExemplar(1, lable)
}

func observeWithExemplar(h *prometheus.HistogramVec, lable prometheus.Labels, value float64, lvals ...string) {
	h.WithLabelValues(lvals...).(prometheus.ExemplarObserver).ObserveWithExemplar(value, lable)
}

func WithUnaryPrometheus() grpc.UnaryServerInterceptor {
	once.Do(_init)
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		start := time.Now()

		svc, method := splitFullMethodName(info.FullMethod)
		lable := ef(ctx)

		incrementWithExemplar(pc.StartedCounterC, lable, svc, method)
		incrementWithExemplar(pc.StreamMsgReceivedC, lable, svc, method)
		resp, err = handler(ctx, req)
		incrementWithExemplar(pc.StreamMsgSentC, lable, svc, method)
		incrementWithExemplar(pc.HandledCounterC, lable, svc, method)
		observeWithExemplar(pc.HandledHistogramC, lable, time.Since(start).Seconds(), svc, method)

		return resp, err
	}
}

func WithStreamPrometheus() grpc.StreamServerInterceptor {
	once.Do(_init)
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		start := time.Now()
		svc, method := splitFullMethodName(info.FullMethod)
		lable := ef(ss.Context())

		incrementWithExemplar(pc.StartedCounterC, lable, svc, method)
		err = handler(srv, ss)
		incrementWithExemplar(pc.HandledCounterC, lable, svc, method)
		observeWithExemplar(pc.HandledHistogramC, lable, time.Since(start).Seconds(), svc, method)

		return err
	}
}
