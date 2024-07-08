package processor

import (
	pb "github.com/crt379/svc-collector-grpc-proto/processor"
	"google.golang.org/grpc"
)

func RegisterServer(srv *grpc.Server) {
	pb.RegisterProcessorServer(srv, &ProcessorImp{})
}
