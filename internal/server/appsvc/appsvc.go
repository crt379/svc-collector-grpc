package appsvc

import (
	"context"
	"fmt"
	"time"

	pb "github.com/crt379/svc-collector-grpc-proto/appsvc"
	"github.com/crt379/svc-collector-grpc/internal/ctxvalue"
	"github.com/crt379/svc-collector-grpc/internal/server"
	svrapp "github.com/crt379/svc-collector-grpc/internal/server/application"
	svrsvc "github.com/crt379/svc-collector-grpc/internal/server/service"
	svrtenant "github.com/crt379/svc-collector-grpc/internal/server/tenant"
	"github.com/crt379/svc-collector-grpc/internal/storage"
	"github.com/crt379/svc-collector-grpc/internal/types"
)

type AppsvcImp struct {
	pb.UnimplementedAppsvcServer
}

func (imp *AppsvcImp) pre(ctx context.Context, appid int) (tenant server.TenantMeta, app server.ApplicationMeta, err error) {
	tenant, err = svrtenant.CheckByMeta(ctx)
	if err != nil {
		return
	}

	app, err = svrapp.CheckByMeta(ctx, tenant.Uuid, appid)

	return
}

func (imp *AppsvcImp) Create(ctx context.Context, req *pb.CreateRequest) (resp *pb.CreateReply, err error) {
	var (
		tenant  server.TenantMeta
		app     server.ApplicationMeta
		svc     server.ServiceMeta
		appsvcs []server.AppsvcMeta
	)
	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("AppsvcImp Create")

	tenant, app, err = imp.pre(ctx, int(req.Appid))
	if err != nil {
		return resp, err
	}

	resp = new(pb.CreateReply)

	if req.Body == nil || req.Body.ServiceId < 1 {
		return server.ParamterResp(&CResp{resp}, "service_id 不能为空")
	}

	svc, err = svrsvc.CheckByMeta(ctx, tenant.Uuid, int(req.Body.ServiceId))
	if err != nil {
		return resp, err
	}

	dao := AppsvcPgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	appsvcs, err = dao.Select(&server.AppsvcMeta{AppId: app.Uuid, SvcId: svc.Uuid})
	if err != nil {
		return server.SqlErrResp(&CResp{resp}, err)
	}
	if len(appsvcs) > 0 {
		return server.AlreadyExistsResp(&CResp{resp}, fmt.Sprintf("service %d 已经关联", req.Body.ServiceId))
	}

	appsvc := server.AppsvcMeta{
		AppId:      app.Uuid,
		SvcId:      svc.Uuid,
		CreateTime: types.Time(time.Now()),
		Service:    svc,
	}
	appsvc.UpdateTime = appsvc.CreateTime

	appsvc.Uuid, err = dao.Insert(&appsvc)
	if err != nil {
		return server.SqlErrResp(&CResp{resp}, err)
	}

	pbmeta, _ := appsvc.ToPbMeta()
	resp.Appsvc = &pbmeta

	return server.OkResp(&CResp{resp})
}

func (imp *AppsvcImp) Get(ctx context.Context, req *pb.GetRequest) (resp *pb.GetReply, err error) {
	var (
		app server.ApplicationMeta
	)
	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("AppsvcImp Create")

	_, app, err = imp.pre(ctx, int(req.Appid))
	if err != nil {
		return resp, err
	}

	resp = new(pb.GetReply)
	dao := AppsvcPgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	var total int
	appsvc := server.AppsvcMeta{
		Uuid:    int(req.Uuid),
		AppId:   int(app.Uuid),
		SvcId:   int(req.Svcid),
		SvcName: req.Svcname,
	}

	total, err = dao.Count(&appsvc)
	if err != nil {
		return server.SqlErrResp(&GResp{resp}, err)
	}

	var appsvcs []server.AppsvcMeta
	resp.Page = 0
	resp.Limit = 100
	if total > 0 {
		if req.Page > 0 {
			resp.Page = req.Page
		}
		if req.Limit > 0 {
			resp.Limit = req.Limit
		}
		var limitoption *server.LimitOption
		if int(resp.Limit) < total {
			limitoption = server.NewLimitOption(int(resp.Page), int(resp.Limit), "row_number")
		}

		appsvcs, err = dao.SelectAndService(&appsvc, limitoption)
		if err != nil {
			return server.SqlErrResp(&GResp{resp}, err)
		}
	}

	resp.Count, resp.Appsvcs, _ = server.Metas2Pbmeta[server.AppsvcMeta, pb.AppsvcMeta](&appsvcs)
	resp.Total = int32(total)

	return server.OkResp(&GResp{resp})
}

func (imp *AppsvcImp) Delete(ctx context.Context, req *pb.DeleteRequest) (resp *pb.DeleteReply, err error) {
	var (
		app server.ApplicationMeta
	)
	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("AppsvcImp Create")

	_, app, err = imp.pre(ctx, int(req.Appid))
	if err != nil {
		return resp, err
	}

	resp = new(pb.DeleteReply)

	if req.Uuid <= 0 {
		return server.NotFoundResp(&DResp{resp}, fmt.Sprintf("appsvc: %d 不存在", req.Uuid))
	}

	appsvc := server.AppsvcMeta{
		Uuid:  int(req.Uuid),
		AppId: app.Uuid,
	}
	dao := AppsvcPgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	var total int
	total, err = dao.Count(&appsvc)
	if err != nil {
		return server.SqlErrResp(&DResp{resp}, err)
	}

	if total == 0 {
		return server.NotFoundResp(&DResp{resp}, fmt.Sprintf("appsvc: %d 不存在", req.Uuid))
	}

	err = dao.Delete(&appsvc)
	if err != nil {
		return server.SqlErrResp(&DResp{resp}, err)
	}

	return server.OkResp(&DResp{resp})
}
