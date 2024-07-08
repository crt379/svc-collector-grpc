package processor

import (
	"context"
	"fmt"
	"time"

	pb "github.com/crt379/svc-collector-grpc-proto/processor"
	"github.com/crt379/svc-collector-grpc/internal/ctxvalue"
	"github.com/crt379/svc-collector-grpc/internal/server"
	svrapp "github.com/crt379/svc-collector-grpc/internal/server/application"
	svrtenant "github.com/crt379/svc-collector-grpc/internal/server/tenant"
	"github.com/crt379/svc-collector-grpc/internal/storage"
	"github.com/crt379/svc-collector-grpc/internal/types"
	"github.com/crt379/svc-collector-grpc/internal/util"
)

type ProcessorImp struct {
	pb.UnimplementedProcessorServer
}

func (imp *ProcessorImp) pre(ctx context.Context, appid int) (tenant server.TenantMeta, app server.ApplicationMeta, err error) {
	tenant, err = svrtenant.CheckByMeta(ctx)
	if err != nil {
		return
	}

	app, err = svrapp.CheckByMeta(ctx, tenant.Uuid, appid)

	return
}

func (imp *ProcessorImp) Create(ctx context.Context, req *pb.CreateRequest) (resp *pb.CreateReply, err error) {
	var (
		tenant server.TenantMeta
		app    server.ApplicationMeta
	)
	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("ProcessorImp Create")

	tenant, app, err = imp.pre(ctx, int(req.Appid))
	if err != nil {
		return resp, err
	}

	resp = new(pb.CreateReply)

	if req.Addr == "" {
		return server.ParamterResp(&CResp{resp}, "addr 不能为空")
	}

	dao := ProcessorPgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	var (
		proc  server.ProcessorMeta
		procs []server.ProcessorMeta
	)
	proc.Addr = req.Addr
	proc.AppId = app.Uuid

	procs, err = dao.Select(&proc)
	if err != nil {
		return server.SqlErrResp(&CResp{resp}, err)
	}
	if len(procs) > 0 {
		return server.AlreadyExistsResp(&CResp{resp}, fmt.Sprintf("addr为 %s 的 processor 已经存在", req.Addr))
	}

	proc.Weight = int(req.Weight)
	proc.State = req.State
	proc.TanantId = tenant.Uuid
	if req.Weight == 0 {
		proc.Weight = 50
	}
	proc.CreateTime = types.Time(time.Now())
	proc.UpdateTime = proc.CreateTime

	proc.Uuid, err = dao.Insert(&proc)
	if err != nil {
		return server.SqlErrResp(&CResp{resp}, err)
	}

	pbmeta, _ := proc.ToPbMeta()
	resp.Processor = &pbmeta

	return server.OkResp(&CResp{resp})
}

func (imp *ProcessorImp) Get(ctx context.Context, req *pb.GetRequest) (resp *pb.GetReply, err error) {
	var (
		app server.ApplicationMeta
	)
	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("ProcessorImp Get")

	_, app, err = imp.pre(ctx, int(req.Appid))
	if err != nil {
		return resp, err
	}

	resp = new(pb.GetReply)
	dao := ProcessorPgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	var total int
	proc := server.ProcessorMeta{
		Uuid:   int(req.Uuid),
		Addr:   req.Addr,
		Weight: int(req.Weight),
		State:  req.State,
		AppId:  app.Uuid,
	}

	total, err = dao.Count(&proc)
	if err != nil {
		return server.SqlErrResp(&GResp{resp}, err)
	}

	var procs []server.ProcessorMeta
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

		procs, err = dao.Select(&proc, limitoption)
		if err != nil {
			return server.SqlErrResp(&GResp{resp}, err)
		}
	}

	resp.Count, resp.Processors, _ = server.Metas2Pbmeta[server.ProcessorMeta, pb.ProcessorMeta](&procs)
	resp.Total = int32(total)

	return server.OkResp(&GResp{resp})
}

func (imp *ProcessorImp) Delete(ctx context.Context, req *pb.DeleteRequest) (resp *pb.DeleteReply, err error) {
	var (
		app server.ApplicationMeta
	)
	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("ProcessorImp Delete")

	_, app, err = imp.pre(ctx, int(req.Appid))
	if err != nil {
		return resp, err
	}

	resp = new(pb.DeleteReply)

	if req.Uuid <= 0 {
		return server.NotFoundResp(&DResp{resp}, fmt.Sprintf("processor: %d 不存在", req.Uuid))
	}

	proc := server.ProcessorMeta{
		Uuid:  int(req.Uuid),
		AppId: app.Uuid,
	}
	dao := ProcessorPgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	var total int
	total, err = dao.Count(&proc)
	if err != nil {
		return server.SqlErrResp(&DResp{resp}, err)
	}
	if total == 0 {
		return server.NotFoundResp(&DResp{resp}, fmt.Sprintf("processor: %d 不存在", req.Uuid))
	}

	err = dao.Delete(&proc)
	if err != nil {
		return server.SqlErrResp(&DResp{resp}, err)
	}

	return server.OkResp(&DResp{resp})
}

func (imp *ProcessorImp) Update(ctx context.Context, req *pb.UpdateRequest) (resp *pb.UpdateReply, err error) {
	var (
		app server.ApplicationMeta
	)
	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("ProcessorImp Update")

	_, app, err = imp.pre(ctx, int(req.Appid))
	if err != nil {
		return resp, err
	}

	resp = new(pb.UpdateReply)

	if req.Uuid <= 0 {
		return server.NotFoundResp(&UResp{resp}, fmt.Sprintf("processor: %d 不存在", req.Uuid))
	}

	proc := server.ProcessorMeta{
		Uuid:  int(req.Uuid),
		AppId: app.Uuid,
	}
	dao := ProcessorPgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	var procs []server.ProcessorMeta
	procs, err = dao.Select(&proc)
	if err != nil {
		return server.SqlErrResp(&UResp{resp}, err)
	}
	if len(procs) == 0 {
		return server.NotFoundResp(&UResp{resp}, fmt.Sprintf("processor: %d 不存在", req.Uuid))
	}

	var newproc server.ProcessorMeta
	proc = procs[0]
	newproc.Addr = req.Addr
	newproc.Weight = int(req.Weight)
	newproc.State = req.State

	if !util.UpdateValueSame(&newproc, &proc) {
		return server.ParamterResp(&UResp{resp}, "没有需要修改的内容")
	}

	if req.Addr != "" {
		procs, err = dao.Select(&server.ProcessorMeta{Addr: req.Addr, AppId: app.Uuid})
		if err != nil {
			return server.SqlErrResp(&UResp{resp}, err)
		}
		if len(procs) > 0 {
			return server.NotFoundResp(&UResp{resp}, fmt.Sprintf("addr为 %s 的 processor 已经存在", req.Addr))
		}
	}

	proc.UpdateTime = types.Time(time.Now())
	_, err = dao.Update(&proc)
	if err != nil {
		return server.SqlErrResp(&UResp{resp}, err)
	}

	pbmeta, _ := proc.ToPbMeta()
	resp.Processor = &pbmeta

	return server.OkResp(&UResp{resp})
}
