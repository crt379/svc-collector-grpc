package service

import (
	"context"
	"fmt"

	"github.com/crt379/svc-collector-grpc/internal/ctxvalue"
	"github.com/crt379/svc-collector-grpc/internal/server"
	"github.com/crt379/svc-collector-grpc/internal/storage"
)

func CheckByMeta(ctx context.Context, tenantid, uuid int) (service server.ServiceMeta, err error) {
	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Debug("Service CheckByMeta")

	if uuid <= 0 {
		return service, server.NotFoundErr(fmt.Sprintf("service: %d not found", uuid))
	}

	dao := ServicePgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	var services []server.ServiceMeta
	services, err = dao.Select(&server.ServiceMeta{Uuid: uuid, TenantId: tenantid})
	if err != nil {
		return service, server.InternalErr(err.Error())
	}

	if len(services) == 0 {
		return service, server.NotFoundErr(fmt.Sprintf("service: %d not found", uuid))
	}

	service = services[0]

	return service, nil
}
