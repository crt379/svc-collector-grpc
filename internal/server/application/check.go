package application

import (
	"context"
	"fmt"

	"github.com/crt379/svc-collector-grpc/internal/ctxvalue"
	"github.com/crt379/svc-collector-grpc/internal/server"
	"github.com/crt379/svc-collector-grpc/internal/storage"
)

func CheckByMeta(ctx context.Context, tenantid, uuid int) (app server.ApplicationMeta, err error) {
	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Debug("Applicstion CheckByMeta")

	if uuid <= 0 {
		return app, server.NotFoundErr(fmt.Sprintf("application: %d not found", uuid))
	}

	dao := ApplicationPgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	var apps []server.ApplicationMeta
	apps, err = dao.Select(&server.ApplicationMeta{Uuid: uuid, TenantId: tenantid})
	if err != nil {
		return app, server.InternalErr(err.Error())
	}

	if len(apps) == 0 {
		return app, server.NotFoundErr(fmt.Sprintf("application: %d not found", uuid))
	}

	app = apps[0]

	return app, nil
}
