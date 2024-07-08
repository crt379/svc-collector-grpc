package appsvc

import (
	pb "github.com/crt379/svc-collector-grpc-proto/appsvc"
	"google.golang.org/grpc"
)

func RegisterServer(srv *grpc.Server) {
	pb.RegisterAppsvcServer(srv, &AppsvcImp{})
}
