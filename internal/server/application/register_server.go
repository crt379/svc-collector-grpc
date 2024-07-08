package application

import (
	pb "github.com/crt379/svc-collector-grpc-proto/application"
	"google.golang.org/grpc"
)

func RegisterServer(srv *grpc.Server) {
	pb.RegisterApplicationServer(srv, &ApplicationImp{})
}
