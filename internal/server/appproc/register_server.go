package appproc

import (
	pb "github.com/crt379/svc-collector-grpc-proto/appproc"
	"google.golang.org/grpc"
)

func RegisterServer(srv *grpc.Server) {
	pb.RegisterAppprocServer(srv, &AppprocImp{})
}
