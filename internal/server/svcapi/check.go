package svcapi

import (
	"context"
	"fmt"

	"github.com/crt379/svc-collector-grpc/internal/ctxvalue"
	"github.com/crt379/svc-collector-grpc/internal/server"
	"github.com/crt379/svc-collector-grpc/internal/storage"
)

func CheckByMeta(ctx context.Context, svcid, uuid int) (svcapi server.SvcapiMeta, err error) {
	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Debug("Svcapi CheckByMeta")

	if uuid <= 0 {
		return svcapi, server.NotFoundErr(fmt.Sprintf("svcapi: %d not found", uuid))
	}

	dao := SvcapiPgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	var svcapis []server.SvcapiMeta
	svcapis, err = dao.Select(&server.SvcapiMeta{Uuid: uuid, ServiceId: svcid})
	if err != nil {
		return svcapi, server.InternalErr(err.Error())
	}

	if len(svcapis) == 0 {
		return svcapi, server.NotFoundErr(fmt.Sprintf("svcapi: %d not found", uuid))
	}

	svcapi = svcapis[0]

	return svcapi, nil
}
