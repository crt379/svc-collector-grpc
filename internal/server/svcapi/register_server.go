package svcapi

import (
	pb "github.com/crt379/svc-collector-grpc-proto/svcapi"
	"google.golang.org/grpc"
)

func RegisterServer(srv *grpc.Server) {
	pb.RegisterSvcapiServer(srv, &SvcapiImp{})
}
