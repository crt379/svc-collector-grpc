package appapi

import (
	"context"

	pb "github.com/crt379/svc-collector-grpc-proto/appapi"
	"github.com/crt379/svc-collector-grpc/internal/ctxvalue"
	"github.com/crt379/svc-collector-grpc/internal/server"
	svrtenant "github.com/crt379/svc-collector-grpc/internal/server/tenant"
	"github.com/crt379/svc-collector-grpc/internal/storage"
)

type AppapiImp struct {
	pb.UnimplementedAppapiServer
}

func (imp *AppapiImp) pre(ctx context.Context) (tenant server.TenantMeta, err error) {
	tenant, err = svrtenant.CheckByMeta(ctx)
	if err != nil {
		return
	}

	return
}

func (imp *AppapiImp) Get(ctx context.Context, req *pb.GetRequest) (resp *pb.GetReply, err error) {
	var (
		tenant server.TenantMeta
	)
	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("AppapiImp Get")

	tenant, err = imp.pre(ctx)
	if err != nil {
		return resp, err
	}

	resp = new(pb.GetReply)
	dao := AppapiPgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}
	appapi := server.AppapiMeta{
		Appid:    int(req.Appid),
		Appsvcid: int(req.Appsvcid),
		Svcid:    int(req.Svcid),
		Svcname:  req.Svcname,
		TenantId: tenant.Uuid,
	}

	var appapis []server.AppapiMeta
	appapis, err = dao.Select(&appapi)
	if err != nil {
		return server.SqlErrResp(&GResp{resp}, err)
	}

	resp.Count, resp.Appapis, _ = server.Metas2Pbmeta[server.AppapiMeta, pb.AppapiMeta](&appapis)
	resp.Total = resp.Count

	return server.OkResp(&GResp{resp})
}
