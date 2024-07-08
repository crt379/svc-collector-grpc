package svcapieg

import (
	pb "github.com/crt379/svc-collector-grpc-proto/svcapieg"
	"google.golang.org/grpc"
)

func RegisterServer(srv *grpc.Server) {
	pb.RegisterSvcapiegServer(srv, &SvcapiegImp{})
}
