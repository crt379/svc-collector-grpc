package tenant

import (
	"context"
	"fmt"

	"github.com/crt379/svc-collector-grpc/internal/ctxvalue"
	"github.com/crt379/svc-collector-grpc/internal/server"
	"github.com/crt379/svc-collector-grpc/internal/storage"
	"go.uber.org/zap"
)

func CheckByMeta(ctx context.Context) (tenant server.TenantMeta, err error) {
	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Debug("Tenant CheckByMeta")

	md, ok := ctxvalue.GrpcMetaContext{}.GetValue(ctx)
	if !ok {
		return tenant, server.InternalErr("grpc mete not in context")
	}

	getvalue := md.Get("x-access-tenant")
	if len(getvalue) == 0 {
		return tenant, server.InvalidArgumentErr("grpc mete not found x-access-tenant")
	}

	dao := TenantPgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	cache := TenantCache{
		R: storage.ReadRedis,
		W: storage.WriteRedis,
	}

	tenantname := getvalue[0]
	tenant, err = cache.ZScoreGet(tenantname)
	if err == nil {
		return tenant, err
	}
	logger.Debug("ZScoreGet err", zap.String("error", err.Error()))

	var tenants []server.TenantMeta
	tenants, err = dao.Select(&server.TenantMeta{Name: tenantname})
	if err != nil {
		return tenant, server.InternalErr(err.Error())
	}

	if len(tenants) == 0 {
		return tenant, server.NotFoundErr(fmt.Sprintf("tenant: %s not found", tenantname))
	}

	tenant = tenants[0]

	err = cache.ZAddSet(&tenants)
	if err != nil {
		logger.Debug("ZAddSet err", zap.String("error", err.Error()))
	}

	return tenant, nil
}
