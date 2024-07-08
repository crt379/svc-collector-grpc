package register

import (
	pb "github.com/crt379/svc-collector-grpc-proto/register"
	"google.golang.org/grpc"
)

func RegisterServer(srv *grpc.Server) {
	pb.RegisterRegisterServer(srv, &RegisterImp{})
}
