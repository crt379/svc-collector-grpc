package svcapieg

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	pb "github.com/crt379/svc-collector-grpc-proto/svcapieg"
	"github.com/crt379/svc-collector-grpc/internal/ctxvalue"
	"github.com/crt379/svc-collector-grpc/internal/server"
	svrjdata "github.com/crt379/svc-collector-grpc/internal/server/jdata"
	svrsvc "github.com/crt379/svc-collector-grpc/internal/server/service"
	svrsvcapi "github.com/crt379/svc-collector-grpc/internal/server/svcapi"
	svrtenant "github.com/crt379/svc-collector-grpc/internal/server/tenant"
	"github.com/crt379/svc-collector-grpc/internal/storage"
	"github.com/crt379/svc-collector-grpc/internal/types"
	"go.uber.org/zap"
)

type SvcapiegImp struct {
	pb.UnimplementedSvcapiegServer
}

func (imp *SvcapiegImp) pre(ctx context.Context, sid, aid int) (tenant server.TenantMeta, service server.ServiceMeta, svcapi server.SvcapiMeta, err error) {
	tenant, err = svrtenant.CheckByMeta(ctx)
	if err != nil {
		return
	}

	service, err = svrsvc.CheckByMeta(ctx, tenant.Uuid, sid)
	if err != nil {
		return
	}

	svcapi, err = svrsvcapi.CheckByMeta(ctx, service.Uuid, aid)
	if err != nil {
		return
	}

	return
}

func (imp *SvcapiegImp) Create(ctx context.Context, req *pb.CreateRequest) (resp *pb.CreateReply, err error) {
	var (
		svcapi server.SvcapiMeta
		eg     server.SvcapiegMeta
		egs    []server.SvcapiegMeta
		pbmeta pb.SvcapiegMeta
	)
	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("SvcapiegImp Create")

	_, _, svcapi, err = imp.pre(ctx, int(req.ServiceId), int(req.SvcapiId))
	if err != nil {
		return resp, err
	}

	resp = new(pb.CreateReply)
	dao := SvcapiegPgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	logger.Debug("req data", zap.Any("body", req.Data))

	var bodyjson any
	err = json.Unmarshal([]byte(req.Data), &bodyjson)
	if err != nil {
		return server.InternalResp(&CResp{resp}, err)
	}
	logger.Debug("req data json", zap.Any("body", bodyjson))

	var bodyjsonbytes []byte
	bodyjsonbytes, err = json.Marshal(bodyjson)
	if err != nil {
		return server.InternalResp(&CResp{resp}, err)
	}

	var (
		jdata    server.Jdata
		jdatas   []server.Jdata
		hashtype string = "md5"
	)
	jdata_dao := svrjdata.JdataPgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	datajsonmd5 := md5.Sum(bodyjsonbytes)
	datajsonmd5str := hex.EncodeToString(datajsonmd5[:])
	logger.Info("req data json md5", zap.String(hashtype, datajsonmd5str))

	jdatas, err = jdata_dao.Select(&server.Jdata{HashType: hashtype, HashValue: datajsonmd5str})
	if err != nil {
		return server.SqlErrResp(&CResp{resp}, err)
	}

	if len(jdatas) == 0 {
		jdata.Data = bodyjson
		jdata.CreateTime = types.Time(time.Now())
		jdata.UpdateTime = jdata.CreateTime
		jdata.HashType = hashtype
		jdata.HashValue = datajsonmd5str
		jdata.Uuid, err = jdata_dao.Insert(&jdata)
		if err != nil {
			return server.SqlErrResp(&CResp{resp}, err)
		}
	} else {
		jdata = jdatas[0]
		egs, err = dao.Select(&server.SvcapiegMeta{JdataId: jdata.Uuid, SvcapiId: svcapi.Uuid})
		if err != nil {
			return server.SqlErrResp(&CResp{resp}, err)
		}

		if len(egs) > 0 {
			return server.AlreadyExistsResp(&CResp{resp}, "已有相同的 svcapieg")
		}
	}

	eg.Data = bodyjson
	eg.CreateTime = types.Time(time.Now())
	eg.UpdateTime = eg.CreateTime
	eg.SvcapiId = svcapi.Uuid
	eg.TenantId = svcapi.TenantId
	eg.ServiceId = svcapi.ServiceId
	eg.JdataId = jdata.Uuid

	eg.Uuid, err = dao.Insert(&eg)
	if err != nil {
		return server.SqlErrResp(&CResp{resp}, err)
	}

	pbmeta, err = eg.ToPbMeta()
	if err != nil {
		return server.InternalResp(&CResp{resp}, err)
	}

	resp.Svcapieg = &pbmeta

	return server.OkResp(&CResp{resp})
}

func (imp *SvcapiegImp) Get(ctx context.Context, req *pb.GetRequest) (resp *pb.GetReply, err error) {
	var (
		svcapi server.SvcapiMeta
		eg     server.SvcapiegMeta
		egs    []server.SvcapiegMeta
		total  int
	)

	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("SvcapiegImp Get")

	_, _, svcapi, err = imp.pre(ctx, int(req.ServiceId), int(req.SvcapiId))
	if err != nil {
		return resp, err
	}

	resp = new(pb.GetReply)
	dao := SvcapiegPgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	eg.Uuid = int(req.Uuid)
	eg.SvcapiId = int(svcapi.Uuid)

	total, err = dao.Count(&eg)
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

		egs, err = dao.SelectAndJdata(&eg, limitoption)
		if err != nil {
			return server.SqlErrResp(&GResp{resp}, err)
		}
	}

	resp.Count, resp.Svcapiegs, err = server.Metas2Pbmeta[server.SvcapiegMeta, pb.SvcapiegMeta](&egs)
	if err != nil {
		server.SqlErrResp(&GResp{resp}, err)
	}
	resp.Total = int32(total)

	return server.OkResp(&GResp{resp})
}

func (imp *SvcapiegImp) Delete(ctx context.Context, req *pb.DeleteRequest) (resp *pb.DeleteReply, err error) {
	var (
		svcapi server.SvcapiMeta
		eg     server.SvcapiegMeta
		total  int
	)

	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("SvcapiegImp Delete")

	_, _, svcapi, err = imp.pre(ctx, int(req.ServiceId), int(req.SvcapiId))
	if err != nil {
		return resp, err
	}

	resp = new(pb.DeleteReply)

	if req.Uuid <= 0 {
		return server.NotFoundResp(&DResp{resp}, fmt.Sprintf("svcapieg: %d 不存在", req.Uuid))
	}

	dao := SvcapiegPgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	eg.Uuid = int(req.Uuid)
	eg.SvcapiId = int(svcapi.Uuid)

	total, err = dao.Count(&eg)
	if err != nil {
		return server.SqlErrResp(&DResp{resp}, err)
	}
	if total == 0 {
		return server.NotFoundResp(&DResp{resp}, fmt.Sprintf("svcapieg: %d 不存在", req.Uuid))
	}

	err = dao.Delete(&eg)
	if err != nil {
		return server.SqlErrResp(&DResp{resp}, err)
	}

	return server.OkResp(&DResp{resp})
}

func (imp *SvcapiegImp) Update(ctx context.Context, req *pb.UpdateRequest) (resp *pb.UpdateReply, err error) {
	var (
		svcapi server.SvcapiMeta
		eg     server.SvcapiegMeta
		egs    []server.SvcapiegMeta
	)

	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("SvcapiegImp Delete")

	_, _, svcapi, err = imp.pre(ctx, int(req.ServiceId), int(req.SvcapiId))
	if err != nil {
		return resp, err
	}

	resp = new(pb.UpdateReply)

	if req.Uuid <= 0 {
		return server.NotFoundResp(&UResp{resp}, fmt.Sprintf("svcapieg: %d 不存在", req.Uuid))
	}

	dao := SvcapiegPgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	eg.Uuid = int(req.Uuid)
	eg.SvcapiId = int(svcapi.Uuid)

	egs, err = dao.Select(&eg)
	if err != nil {
		return server.SqlErrResp(&UResp{resp}, err)
	}
	if len(egs) == 0 {
		return server.NotFoundResp(&UResp{resp}, fmt.Sprintf("svcapieg: %d 不存在", req.Uuid))
	}
	eg = egs[0]

	logger.Debug("req data", zap.Any("body", req.Data))

	var bodyjson any
	err = json.Unmarshal([]byte(req.Data), &bodyjson)
	if err != nil {
		return server.InternalResp(&UResp{resp}, err)
	}
	logger.Debug("req data json", zap.Any("body", bodyjson))

	var bodyjsonbytes []byte
	bodyjsonbytes, err = json.Marshal(bodyjson)
	if err != nil {
		return server.InternalResp(&UResp{resp}, err)
	}

	var (
		jdata    server.Jdata
		jdatas   []server.Jdata
		hashtype string = "md5"
	)
	jdata_dao := svrjdata.JdataPgDao{
		W:      storage.WriteDB,
		R:      storage.ReadDB,
		Logger: logger,
	}

	datajsonmd5 := md5.Sum(bodyjsonbytes)
	datajsonmd5str := hex.EncodeToString(datajsonmd5[:])
	logger.Info("req data json md5", zap.String(hashtype, datajsonmd5str))

	jdatas, err = jdata_dao.Select(&server.Jdata{HashType: hashtype, HashValue: datajsonmd5str})
	if err != nil {
		return server.SqlErrResp(&UResp{resp}, err)
	}

	var (
		pbmeta pb.SvcapiegMeta
		total  int
	)

	if len(jdatas) == 0 {
		// todo: 如果原 jdata 没有关联 eg 了则可以 1.修改原jdata的value 2.删除原jdata，关联新jdata
		jdata.Data = bodyjson
		jdata.CreateTime = types.Time(time.Now())
		jdata.UpdateTime = jdata.CreateTime
		jdata.HashType = hashtype
		jdata.HashValue = datajsonmd5str
		jdata.Uuid, err = jdata_dao.Insert(&jdata)
		if err != nil {
			return server.SqlErrResp(&UResp{resp}, err)
		}
	} else {
		jdata = jdatas[0]
		if jdata.Uuid == eg.JdataId {
			eg.Data = bodyjson
			pbmeta, err = eg.ToPbMeta()
			if err != nil {
				return server.InternalResp(&UResp{resp}, err)
			}
			resp.Svcapieg = &pbmeta

			return server.OkResp(&UResp{resp})
		}

		total, err = dao.Count(&server.SvcapiegMeta{SvcapiId: svcapi.Uuid, JdataId: jdata.Uuid})
		if err != nil {
			return server.SqlErrResp(&UResp{resp}, err)
		}
		if total > 0 {
			return server.AlreadyExistsResp(&UResp{resp}, "已有相同的 svcapieg")
		}
	}

	eg.Data = bodyjson
	eg.UpdateTime = types.Time(time.Now())
	eg.JdataId = jdata.Uuid

	_, err = dao.Update(&eg)
	if err != nil {
		return server.SqlErrResp(&UResp{resp}, err)
	}

	pbmeta, err = eg.ToPbMeta()
	if err != nil {
		return server.InternalResp(&UResp{resp}, err)
	}
	resp.Svcapieg = &pbmeta

	return server.OkResp(&UResp{resp})
}
