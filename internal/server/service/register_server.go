package service

import (
	pb "github.com/crt379/svc-collector-grpc-proto/service"
	"google.golang.org/grpc"
)

func RegisterServer(srv *grpc.Server) {
	pb.RegisterServiceServer(srv, &ServiceImp{})
}
