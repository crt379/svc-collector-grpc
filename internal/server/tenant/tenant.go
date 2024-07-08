package tenant

import (
	"context"
	"fmt"
	"time"

	pb "github.com/crt379/svc-collector-grpc-proto/tenant"
	"github.com/crt379/svc-collector-grpc/internal/ctxvalue"
	"github.com/crt379/svc-collector-grpc/internal/server"
	"github.com/crt379/svc-collector-grpc/internal/storage"
	"github.com/crt379/svc-collector-grpc/internal/types"
	"github.com/crt379/svc-collector-grpc/internal/util"
	"go.uber.org/zap"
)

type TenantImp struct {
	pb.UnimplementedTenantServer
}

func (imp *TenantImp) Create(ctx context.Context, req *pb.CreateRequest) (resp *pb.CreateReply, err error) {
	var (
		tenant  server.TenantMeta
		tenants []server.TenantMeta
		pbmeta  pb.TenantMeta
	)
	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("TenantImp Create")

	dao := TenantPgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	cache := TenantCache{
		R: storage.ReadRedis,
		W: storage.WriteRedis,
	}

	resp = new(pb.CreateReply)
	if req.Name == "" {
		return server.ParamterResp(&CResp{resp}, "name 不能为空")
	}
	if util.StrPunctIllegal(req.Name, '-') {
		return server.ParamterResp(&CResp{resp}, "name 不能含有非'-'的字符")
	}

	_, err = cache.ZRank(req.Name)
	if err == nil {
		return server.AlreadyExistsResp(&CResp{resp}, fmt.Sprintf("name 为 %s 的 tenant 已经存在", req.Name))
	}

	tenants, err = dao.Select(&server.TenantMeta{Name: req.Name})
	if err != nil {
		return server.SqlErrResp(&CResp{resp}, err)
	}

	if len(tenants) > 0 {
		err = cache.ZAddSet(&tenants)
		if err != nil {
			logger.Warn("ZAddSet err", zap.String("error", err.Error()))
		}
		return server.AlreadyExistsResp(&CResp{resp}, fmt.Sprintf("name 为 %s 的 tenant 已经存在", req.Name))
	}

	tenant.Name = req.Name
	tenant.Describe = req.Describe
	tenant.CreateTime = types.Time(time.Now())
	tenant.UpdateTime = tenant.CreateTime

	tenant.Uuid, err = dao.Insert(&tenant)
	if err != nil {
		return server.SqlErrResp(&CResp{resp}, err)
	}

	pbmeta, _ = tenant.ToPbMeta()
	resp.Tenant = &pbmeta

	return server.OkResp(&CResp{resp})
}

func (imp *TenantImp) Get(ctx context.Context, req *pb.GetRequest) (resp *pb.GetReply, err error) {
	var (
		tenant  server.TenantMeta
		rtenant server.TenantMeta
		tenants []server.TenantMeta
		total   int
		one     bool
	)

	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("TenantImp Get")

	dao := TenantPgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	cache := TenantCache{
		R: storage.ReadRedis,
		W: storage.WriteRedis,
	}

	resp = new(pb.GetReply)

	if req.Uuid != 0 {
		one = true
		tenant.Uuid = int(req.Uuid)
		rtenant, err = cache.Get(tenant.Uuid)
	}
	if req.Name != "" {
		one = true
		tenant.Name = req.Name
		rtenant, err = cache.ZScoreGet(tenant.Name)
	}

	if one && err == nil {
		pt, _ := rtenant.ToPbMeta()
		resp.Tenants = []*pb.TenantMeta{&pt}
		resp.Count, resp.Total = 1, 1
		return server.OkResp(&GResp{resp})
	}

	total, err = dao.Count(&tenant)
	if err != nil {
		return server.SqlErrResp(&GResp{resp}, err)
	}

	if total > 0 {
		tenants, err = dao.Select(&tenant)
		if err != nil {
			return server.SqlErrResp(&GResp{resp}, err)
		}
		if one {
			err = cache.ZAddSet(&tenants)
			if err != nil {
				logger.Warn("ZAddSet err", zap.String("error", err.Error()))
			}
		}
	}

	resp.Count, resp.Tenants, _ = server.Metas2Pbmeta[server.TenantMeta, pb.TenantMeta](&tenants)
	resp.Total = int32(total)

	return server.OkResp(&GResp{resp})
}

func (imp *TenantImp) Delete(ctx context.Context, req *pb.DeleteRequest) (resp *pb.DeleteReply, err error) {
	var (
		tenant server.TenantMeta
		total  int
	)

	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("TenantImp Delete")

	dao := TenantPgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	resp = new(pb.DeleteReply)

	if req.Uuid <= 0 {
		return server.NotFoundResp(&DResp{resp}, fmt.Sprintf("tenant: %d 不存在", req.Uuid))
	}

	tenant.Uuid = int(req.Uuid)
	total, err = dao.Count(&tenant)
	if err != nil {
		return server.SqlErrResp(&DResp{resp}, err)
	}

	if total == 0 {
		return server.NotFoundResp(&DResp{resp}, fmt.Sprintf("tenant: %d 不存在", req.Uuid))
	}

	err = dao.Delete(&tenant)
	if err != nil {
		return server.SqlErrResp(&DResp{resp}, err)
	}

	cache := TenantCache{
		R: storage.ReadRedis,
		W: storage.WriteRedis,
	}
	err = cache.ZRemDel(tenant.Name, tenant.Uuid)
	if err != nil {
		logger.Warn("ZRemDel err", zap.String("error", err.Error()))
	}

	return server.OkResp(&DResp{resp})
}

func (imp *TenantImp) Update(ctx context.Context, req *pb.UpdateRequest) (resp *pb.UpdateReply, err error) {
	var (
		tenant    server.TenantMeta
		newtenant server.TenantMeta
		tenants   []server.TenantMeta
	)

	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("TenantImp Update")

	dao := TenantPgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	resp = new(pb.UpdateReply)

	if req.Uuid <= 0 {
		return server.NotFoundResp(&UResp{resp}, fmt.Sprintf("tenant: %d 不存在", req.Uuid))
	}

	tenant.Uuid = int(req.Uuid)
	tenants, err = dao.Select(&tenant)
	if err != nil {
		return server.SqlErrResp(&UResp{resp}, err)
	}

	if len(tenants) == 0 {
		return server.NotFoundResp(&UResp{resp}, fmt.Sprintf("tenant: %d 不存在", req.Uuid))
	}

	tenant = tenants[0]
	newtenant.Name = req.Name
	newtenant.Describe = req.Describe

	if !util.UpdateValueSame(&newtenant, &tenant) {
		return server.ParamterResp(&UResp{resp}, "没有需要修改的内容")
	}

	tenant.UpdateTime = types.Time(time.Now())
	_, err = dao.Update(&tenant)
	if err != nil {
		return server.SqlErrResp(&UResp{resp}, err)
	}

	cache := TenantCache{
		R: storage.ReadRedis,
		W: storage.WriteRedis,
	}
	err = cache.ZRemDel(tenant.Name, tenant.Uuid)
	if err != nil {
		logger.Warn("ZRemDel err", zap.String("error", err.Error()))
	}

	pbmeta, _ := tenant.ToPbMeta()
	resp.Tenant = &pbmeta

	return server.OkResp(&UResp{resp})
}
