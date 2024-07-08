package svcapi

import (
	"context"
	"fmt"
	"time"

	pb "github.com/crt379/svc-collector-grpc-proto/svcapi"
	"github.com/crt379/svc-collector-grpc/internal/ctxvalue"
	"github.com/crt379/svc-collector-grpc/internal/server"
	svrsvc "github.com/crt379/svc-collector-grpc/internal/server/service"
	svrtenant "github.com/crt379/svc-collector-grpc/internal/server/tenant"
	"github.com/crt379/svc-collector-grpc/internal/storage"
	"github.com/crt379/svc-collector-grpc/internal/types"
	"github.com/crt379/svc-collector-grpc/internal/util"
)

type SvcapiImp struct {
	pb.UnimplementedSvcapiServer
}

func (imp *SvcapiImp) pre(ctx context.Context, sid int) (tenant server.TenantMeta, service server.ServiceMeta, err error) {
	tenant, err = svrtenant.CheckByMeta(ctx)
	if err != nil {
		return
	}

	service, err = svrsvc.CheckByMeta(ctx, tenant.Uuid, sid)
	if err != nil {
		return
	}

	return
}

func (imp *SvcapiImp) Create(ctx context.Context, req *pb.CreateRequest) (resp *pb.CreateReply, err error) {
	var (
		service server.ServiceMeta
		svcapis []server.SvcapiMeta
		pbmeta  pb.SvcapiMeta
	)
	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("SvcapiImp Create")

	_, service, err = imp.pre(ctx, int(req.ServiceId))
	if err != nil {
		return resp, err
	}

	resp = new(pb.CreateReply)

	if req.Path == "" {
		return server.ParamterResp(&CResp{resp}, "path 不能为空")
	}
	if req.Method == "" {
		return server.ParamterResp(&CResp{resp}, "method 不能为空")
	}

	dao := SvcapiPgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	svcapis, err = dao.Select(&server.SvcapiMeta{Path: req.Path, Method: req.Method, ServiceId: service.Uuid})
	if err != nil {
		return server.SqlErrResp(&CResp{resp}, err)
	}

	if len(svcapis) > 0 {
		return server.AlreadyExistsResp(&CResp{resp}, "service 已有相同 path 和 method 的 api")
	}

	var svcapi = server.SvcapiMeta{
		Path:       req.Path,
		Method:     req.Method,
		Describe:   req.Describe,
		CreateTime: types.Time(time.Now()),
		ServiceId:  service.Uuid,
		TenantId:   service.TenantId,
	}
	svcapi.UpdateTime = svcapi.CreateTime

	svcapi.Uuid, err = dao.Insert(&svcapi)
	if err != nil {
		return server.SqlErrResp(&CResp{resp}, err)
	}

	pbmeta, _ = svcapi.ToPbMeta()
	resp.Svcapi = &pbmeta

	return server.OkResp(&CResp{resp})
}

func (imp *SvcapiImp) Get(ctx context.Context, req *pb.GetRequest) (resp *pb.GetReply, err error) {
	var (
		service server.ServiceMeta
		svcapi  server.SvcapiMeta
		svcapis []server.SvcapiMeta
		total   int
	)

	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("SvcapiImp Get")

	_, service, err = imp.pre(ctx, int(req.ServiceId))
	if err != nil {
		return resp, err
	}

	resp = new(pb.GetReply)
	dao := SvcapiPgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	svcapi.Uuid = int(req.Uuid)
	svcapi.Path = req.Path
	svcapi.Method = req.Method
	svcapi.ServiceId = service.Uuid

	total, err = dao.Count(&svcapi)
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

		svcapis, err = dao.Select(&svcapi, limitoption)
		if err != nil {
			return server.SqlErrResp(&GResp{resp}, err)
		}
	}

	resp.Count, resp.Svcapis, _ = server.Metas2Pbmeta[server.SvcapiMeta, pb.SvcapiMeta](&svcapis)
	resp.Total = int32(total)

	return server.OkResp(&GResp{resp})
}

func (imp *SvcapiImp) Delete(ctx context.Context, req *pb.DeleteRequest) (resp *pb.DeleteReply, err error) {
	var (
		service server.ServiceMeta
		svcapi  server.SvcapiMeta
		total   int
	)

	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("SvcapiImp Delete")

	_, service, err = imp.pre(ctx, int(req.ServiceId))
	if err != nil {
		return resp, err
	}

	resp = new(pb.DeleteReply)

	if req.Uuid <= 0 {
		return server.NotFoundResp(&DResp{resp}, fmt.Sprintf("svcapi: %d 不存在", req.Uuid))
	}

	dao := SvcapiPgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	svcapi.Uuid = int(req.Uuid)
	svcapi.ServiceId = service.Uuid

	total, err = dao.Count(&svcapi)
	if err != nil {
		return server.SqlErrResp(&DResp{resp}, err)
	}
	if total == 0 {
		return server.NotFoundResp(&DResp{resp}, fmt.Sprintf("svcapi: %d 不存在", req.Uuid))
	}

	err = dao.Delete(&svcapi)
	if err != nil {
		return server.SqlErrResp(&DResp{resp}, err)
	}

	return server.OkResp(&DResp{resp})
}

func (imp *SvcapiImp) Update(ctx context.Context, req *pb.UpdateRequest) (resp *pb.UpdateReply, err error) {
	var (
		service   server.ServiceMeta
		svcapi    server.SvcapiMeta
		newsvcapi server.SvcapiMeta
		svcapis   []server.SvcapiMeta
	)

	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("SvcapiImp Update")

	_, service, err = imp.pre(ctx, int(req.ServiceId))
	if err != nil {
		return resp, err
	}

	resp = new(pb.UpdateReply)

	if req.Uuid <= 0 {
		return server.NotFoundResp(&UResp{resp}, fmt.Sprintf("svcapi: %d 不存在", req.Uuid))
	}

	dao := SvcapiPgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	svcapi.Uuid = int(req.Uuid)
	svcapi.ServiceId = service.Uuid

	svcapis, err = dao.Select(&svcapi)
	if err != nil {
		return server.SqlErrResp(&UResp{resp}, err)
	}

	if len(svcapis) == 0 {
		return server.NotFoundResp(&UResp{resp}, fmt.Sprintf("svcapi: %d 不存在", req.Uuid))
	}

	svcapi = svcapis[0]
	newsvcapi.Path = req.Path
	newsvcapi.Method = req.Method
	newsvcapi.Describe = req.Describe

	if !util.UpdateValueSame(&newsvcapi, &svcapi) {
		return server.ParamterResp(&UResp{resp}, "没有需要修改的内容")
	}

	svcapi.UpdateTime = types.Time(time.Now())
	_, err = dao.Update(&svcapi)
	if err != nil {
		return server.SqlErrResp(&UResp{resp}, err)
	}

	pbmeta, _ := svcapi.ToPbMeta()
	resp.Svcapi = &pbmeta

	return server.OkResp(&UResp{resp})
}
