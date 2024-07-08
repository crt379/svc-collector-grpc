package application

import (
	"context"
	"fmt"
	"time"

	pb "github.com/crt379/svc-collector-grpc-proto/application"
	"github.com/crt379/svc-collector-grpc/internal/ctxvalue"
	"github.com/crt379/svc-collector-grpc/internal/server"
	svrtenant "github.com/crt379/svc-collector-grpc/internal/server/tenant"
	"github.com/crt379/svc-collector-grpc/internal/storage"
	"github.com/crt379/svc-collector-grpc/internal/types"
	"github.com/crt379/svc-collector-grpc/internal/util"
)

type ApplicationImp struct {
	pb.UnimplementedApplicationServer
}

func (imp *ApplicationImp) pre(ctx context.Context) (tenant server.TenantMeta, err error) {
	tenant, err = svrtenant.CheckByMeta(ctx)
	if err != nil {
		return
	}

	return
}

func (imp *ApplicationImp) Create(ctx context.Context, req *pb.CreateRequest) (resp *pb.CreateReply, err error) {
	var (
		tenant server.TenantMeta
		apps   []server.ApplicationMeta
	)
	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("ApplicationImp Create")

	tenant, err = imp.pre(ctx)
	if err != nil {
		return resp, err
	}

	resp = new(pb.CreateReply)

	if req.Name == "" {
		return server.ParamterResp(&CResp{resp}, "name 不能为空")
	}
	if util.StrPunctIllegal(req.Name, '-') {
		return server.ParamterResp(&CResp{resp}, "name 不能含有非'-'的字符")
	}

	dao := ApplicationPgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	apps, err = dao.Select(&server.ApplicationMeta{Name: req.Name, TenantId: tenant.Uuid})
	if err != nil {
		return server.SqlErrResp(&CResp{resp}, err)
	}
	if len(apps) > 0 {
		return server.AlreadyExistsResp(&CResp{resp}, fmt.Sprintf("name 为 %s 的 application 已经存在", req.Name))
	}

	app := server.ApplicationMeta{
		Name:       req.Name,
		Describe:   req.Describe,
		CreateTime: types.Time(time.Now()),
		TenantId:   tenant.Uuid,
	}
	app.UpdateTime = app.CreateTime

	app.Uuid, err = dao.Insert(&app)
	if err != nil {
		return server.SqlErrResp(&CResp{resp}, err)
	}

	pbmeta, _ := app.ToPbMeta()
	resp.Application = &pbmeta

	return server.OkResp(&CResp{resp})
}

func (imp *ApplicationImp) Get(ctx context.Context, req *pb.GetRequest) (resp *pb.GetReply, err error) {
	var (
		tenant server.TenantMeta
		app    server.ApplicationMeta
		apps   []server.ApplicationMeta
		total  int
	)
	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("ApplicationImp Get")

	tenant, err = imp.pre(ctx)
	if err != nil {
		return resp, err
	}

	resp = new(pb.GetReply)
	dao := ApplicationPgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	app.Uuid = int(req.Uuid)
	app.Name = req.Name
	app.TenantId = tenant.Uuid

	total, err = dao.Count(&app)
	if err != nil {
		return server.SqlErrResp(&GResp{resp}, err)
	}

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

		apps, err = dao.Select(&app, limitoption)
		if err != nil {
			return server.SqlErrResp(&GResp{resp}, err)
		}
	}

	resp.Count, resp.Applications, _ = server.Metas2Pbmeta[server.ApplicationMeta, pb.ApplicationMete](&apps)
	resp.Total = int32(total)

	return server.OkResp(&GResp{resp})
}

func (imp *ApplicationImp) Delete(ctx context.Context, req *pb.DeleteRequest) (resp *pb.DeleteReply, err error) {
	var (
		tenant server.TenantMeta
		app    server.ApplicationMeta
		total  int
	)
	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("ApplicationImp Delete")

	tenant, err = imp.pre(ctx)
	if err != nil {
		return resp, err
	}

	resp = new(pb.DeleteReply)

	if req.Uuid <= 0 {
		return server.NotFoundResp(&DResp{resp}, fmt.Sprintf("application: %d 不存在", req.Uuid))
	}

	dao := ApplicationPgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	app.Uuid = int(req.Uuid)
	app.TenantId = tenant.Uuid

	total, err = dao.Count(&app)
	if err != nil {
		return server.SqlErrResp(&DResp{resp}, err)
	}

	if total == 0 {
		return server.NotFoundResp(&DResp{resp}, fmt.Sprintf("application: %d 不存在", req.Uuid))
	}

	err = dao.Delete(&app)
	if err != nil {
		return server.SqlErrResp(&DResp{resp}, err)
	}

	return server.OkResp(&DResp{resp})
}

func (imp *ApplicationImp) Update(ctx context.Context, req *pb.UpdateRequest) (resp *pb.UpdateReply, err error) {
	var (
		tenant server.TenantMeta
		app    server.ApplicationMeta
		apps   []server.ApplicationMeta
	)

	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("ApplicationImp Update")

	tenant, err = imp.pre(ctx)
	if err != nil {
		return resp, err
	}

	resp = new(pb.UpdateReply)

	if req.Uuid <= 0 {
		return server.NotFoundResp(&UResp{resp}, fmt.Sprintf("application: %d 不存在", req.Uuid))
	}

	dao := ApplicationPgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	app.Uuid = int(req.Uuid)
	app.TenantId = tenant.Uuid

	apps, err = dao.Select(&app)
	if err != nil {
		return server.SqlErrResp(&UResp{resp}, err)
	}
	if len(apps) == 0 {
		return server.NotFoundResp(&UResp{resp}, fmt.Sprintf("application: %d 不存在", req.Uuid))
	}

	var newapp server.ApplicationMeta
	app = apps[0]
	newapp.Name = req.Name
	newapp.Describe = req.Describe

	if !util.UpdateValueSame(&newapp, &app) {
		return server.ParamterResp(&UResp{resp}, "没有需要修改的内容")
	}

	app.UpdateTime = types.Time(time.Now())
	_, err = dao.Update(&app)
	if err != nil {
		return server.SqlErrResp(&UResp{resp}, err)
	}

	pbmeta, _ := app.ToPbMeta()
	resp.Application = &pbmeta

	return server.OkResp(&UResp{resp})
}
