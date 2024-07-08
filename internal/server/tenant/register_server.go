package tenant

import (
	pb "github.com/crt379/svc-collector-grpc-proto/tenant"
	"google.golang.org/grpc"
)

func RegisterServer(srv *grpc.Server) {
	pb.RegisterTenantServer(srv, &TenantImp{})
}
