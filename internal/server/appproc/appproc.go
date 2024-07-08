package appproc

import (
	"context"

	pb "github.com/crt379/svc-collector-grpc-proto/appproc"
	"github.com/crt379/svc-collector-grpc/internal/ctxvalue"
	"github.com/crt379/svc-collector-grpc/internal/server"
	svrtenant "github.com/crt379/svc-collector-grpc/internal/server/tenant"
	"github.com/crt379/svc-collector-grpc/internal/storage"
)

type AppprocImp struct {
	pb.UnimplementedAppprocServer
}

func (imp *AppprocImp) pre(ctx context.Context) (tenant server.TenantMeta, err error) {
	tenant, err = svrtenant.CheckByMeta(ctx)
	if err != nil {
		return
	}

	return
}

func (imp *AppprocImp) Get(ctx context.Context, req *pb.GetRequest) (resp *pb.GetReply, err error) {
	var (
		tenant server.TenantMeta
	)
	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("AppprocImp Get")

	tenant, err = imp.pre(ctx)
	if err != nil {
		return resp, err
	}

	resp = new(pb.GetReply)
	dao := AppprocPgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	a3proc := server.AppprocMeta{
		Appid:    int(req.Appid),
		Appname:  req.Appname,
		State:    req.State,
		TenantId: tenant.Uuid,
	}
	if req.Weight != nil {
		w := int(req.Weight.Value)
		a3proc.Weight = &w
	}

	var a3procs []server.AppprocMeta
	a3procs, err = dao.Select(&a3proc)
	if err != nil {
		return server.SqlErrResp(&GResp{resp}, err)
	}

	resp.Count, resp.Appprocs, _ = server.Metas2Pbmeta[server.AppprocMeta, pb.AppprocMeta](&a3procs)
	resp.Total = resp.Count

	return server.OkResp(&GResp{resp})
}
