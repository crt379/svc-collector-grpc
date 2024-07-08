package appapi

import (
	pb "github.com/crt379/svc-collector-grpc-proto/appapi"
	"google.golang.org/grpc"
)

func RegisterServer(srv *grpc.Server) {
	pb.RegisterAppapiServer(srv, &AppapiImp{})
}
