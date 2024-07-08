package service

import (
	"context"
	"fmt"
	"time"

	pb "github.com/crt379/svc-collector-grpc-proto/service"
	"github.com/crt379/svc-collector-grpc/internal/ctxvalue"
	"github.com/crt379/svc-collector-grpc/internal/server"
	svrtenant "github.com/crt379/svc-collector-grpc/internal/server/tenant"
	"github.com/crt379/svc-collector-grpc/internal/storage"
	"github.com/crt379/svc-collector-grpc/internal/types"
	"github.com/crt379/svc-collector-grpc/internal/util"
)

type ServiceImp struct {
	pb.UnimplementedServiceServer
}

func (imp *ServiceImp) pre(ctx context.Context) (tenant server.TenantMeta, err error) {
	tenant, err = svrtenant.CheckByMeta(ctx)
	if err != nil {
		return
	}

	return
}

func (imp *ServiceImp) Create(ctx context.Context, req *pb.CreateRequest) (resp *pb.CreateReply, err error) {
	var (
		tenant   server.TenantMeta
		services []server.ServiceMeta
	)
	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("ServiceImp Create")

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

	dao := ServicePgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	services, err = dao.Select(&server.ServiceMeta{Name: req.Name, TenantId: tenant.Uuid})
	if err != nil {
		return server.SqlErrResp(&CResp{resp}, err)
	}
	if len(services) > 0 {
		return server.AlreadyExistsResp(&CResp{resp}, fmt.Sprintf("name 为 %s 的 tenant 已经存在", req.Name))
	}

	service := server.ServiceMeta{
		Name:       req.Name,
		Describe:   req.Describe,
		CreateTime: types.Time(time.Now()),
		TenantId:   tenant.Uuid,
	}
	service.UpdateTime = service.CreateTime

	service.Uuid, err = dao.Insert(&service)
	if err != nil {
		return server.SqlErrResp(&CResp{resp}, err)
	}

	pbmeta, _ := service.ToPbMeta()
	resp.Service = &pbmeta

	return server.OkResp(&CResp{resp})
}

func (imp *ServiceImp) Get(ctx context.Context, req *pb.GetRequest) (resp *pb.GetReply, err error) {
	var (
		tenant   server.TenantMeta
		service  server.ServiceMeta
		services []server.ServiceMeta
		total    int
	)

	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("ServiceImp Get")

	tenant, err = imp.pre(ctx)
	if err != nil {
		return resp, err
	}

	resp = new(pb.GetReply)
	dao := ServicePgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	service.Uuid = int(req.Uuid)
	service.Name = req.Name
	service.TenantId = tenant.Uuid

	total, err = dao.Count(&service)
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

		services, err = dao.Select(&service, limitoption)
		if err != nil {
			return server.SqlErrResp(&GResp{resp}, err)
		}
	}

	resp.Count, resp.Services, _ = server.Metas2Pbmeta[server.ServiceMeta, pb.ServiceMeta](&services)
	resp.Total = int32(total)

	return server.OkResp(&GResp{resp})
}

func (imp *ServiceImp) Delete(ctx context.Context, req *pb.DeleteRequest) (resp *pb.DeleteReply, err error) {
	var (
		tenant  server.TenantMeta
		service server.ServiceMeta
		total   int
	)

	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("ServiceImp Delete")

	tenant, err = imp.pre(ctx)
	if err != nil {
		return resp, err
	}

	resp = new(pb.DeleteReply)

	if req.Uuid <= 0 {
		return server.NotFoundResp(&DResp{resp}, fmt.Sprintf("service: %d 不存在", req.Uuid))
	}

	dao := ServicePgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	service.Uuid = int(req.Uuid)
	service.TenantId = tenant.Uuid

	total, err = dao.Count(&service)
	if err != nil {
		return server.SqlErrResp(&DResp{resp}, err)
	}

	if total == 0 {
		return server.NotFoundResp(&DResp{resp}, fmt.Sprintf("service: %d 不存在", req.Uuid))
	}

	err = dao.Delete(&service)
	if err != nil {
		return server.SqlErrResp(&DResp{resp}, err)
	}

	return server.OkResp(&DResp{resp})
}

func (imp *ServiceImp) Update(ctx context.Context, req *pb.UpdateRequest) (resp *pb.UpdateReply, err error) {
	var (
		tenant   server.TenantMeta
		service  server.ServiceMeta
		services []server.ServiceMeta
	)

	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("ServiceImp Update")

	tenant, err = imp.pre(ctx)
	if err != nil {
		return resp, err
	}

	resp = new(pb.UpdateReply)

	if req.Uuid <= 0 {
		return server.NotFoundResp(&UResp{resp}, fmt.Sprintf("tenant: %d 不存在", req.Uuid))
	}

	dao := ServicePgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	service.Uuid = int(req.Uuid)
	service.TenantId = tenant.Uuid

	services, err = dao.Select(&service)
	if err != nil {
		return server.SqlErrResp(&UResp{resp}, err)
	}

	if len(services) == 0 {
		return server.NotFoundResp(&UResp{resp}, fmt.Sprintf("tenant: %d 不存在", req.Uuid))
	}

	var newservice server.ServiceMeta
	service = services[0]
	newservice.Name = req.Name
	newservice.Describe = req.Describe

	if !util.UpdateValueSame(&newservice, &service) {
		return server.ParamterResp(&UResp{resp}, "没有需要修改的内容")
	}

	service.UpdateTime = types.Time(time.Now())
	_, err = dao.Update(&service)
	if err != nil {
		return server.SqlErrResp(&UResp{resp}, err)
	}

	pbmeta, _ := service.ToPbMeta()
	resp.Service = &pbmeta

	return server.OkResp(&UResp{resp})
}
