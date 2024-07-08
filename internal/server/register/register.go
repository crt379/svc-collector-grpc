package register

import (
	"context"
	"sync"

	"github.com/crt379/registerdiscovery"
	pb "github.com/crt379/svc-collector-grpc-proto/register"
	"github.com/crt379/svc-collector-grpc/internal/config"
	"github.com/crt379/svc-collector-grpc/internal/ctxvalue"
	"github.com/crt379/svc-collector-grpc/internal/server"
	"github.com/crt379/svc-collector-grpc/internal/storage"
)

var (
	mu    sync.Mutex
	sctxc = make(map[string]ctxc)
	saddr = make(map[string]string)
	addrs = make(map[string]*[]string)

	lcache = cache[*pb.RegisterInfo]{
		f: ToPbRegisterInfo,
	}
)

type ctxc struct {
	ctx    context.Context
	cancel context.CancelFunc
}

type cache[T any] struct {
	f func(ctx context.Context, service string, add string) T
}

func (c *cache[T]) add(ctx ctxc, service string, addr string) bool {
	if _, ok := sctxc[service]; ok {
		return false
	}

	sctxc[service] = ctx
	saddr[service] = addr

	slptr, ok := addrs[addr]
	if ok {
		*slptr = append(*slptr, service)
	} else {
		sl := []string{service}
		addrs[addr] = &sl
	}

	return true
}

func (c *cache[T]) del(service string) (metas []T) {
	if service == "" {
		for s, ctx := range sctxc {
			addr := saddr[s]
			metas = append(metas, c.f(ctx.ctx, s, addr))

			ctx.cancel()
		}

		sctxc = make(map[string]ctxc)
		saddr = make(map[string]string)
		addrs = make(map[string]*[]string)
	} else {
		ctx, ok := sctxc[service]
		if !ok {
			return
		}

		addr := saddr[service]
		metas = append(metas, c.f(ctx.ctx, service, addr))

		delete(sctxc, service)
		delete(saddr, service)
		var nsl []string
		sl := *(addrs[addr])
		for i, s := range sl {
			if s == service {
				nsl = sl[:i]
				if i+1 < len(sl) {
					nsl = append(nsl, sl[i+1:]...)
				}
				addrs[addr] = &nsl
				break
			}
		}

		ctx.cancel()
	}

	return
}

func (c *cache[T]) get() (metas []T) {
	for s, addr := range saddr {
		ctx := sctxc[s]
		metas = append(metas, c.f(ctx.ctx, s, addr))
	}

	return
}

func (c *cache[T]) getByService(service string) (metas []T) {
	metas = make([]T, 0, 1)
	ctx, ok := sctxc[service]
	if !ok {
		return
	}

	addr := saddr[service]
	metas = append(metas, c.f(ctx.ctx, service, addr))

	return
}

func (c *cache[T]) getByAddress(addr string) (metas []T) {
	slptr, ok := addrs[addr]
	if !ok || len(*slptr) == 0 {
		return
	}

	metas = make([]T, len(*slptr))
	for i, s := range *slptr {
		ctx := sctxc[s]
		metas[i] = c.f(ctx.ctx, s, addr)
	}

	return
}

func InternalRegister(ctx context.Context, service string, addr string) (err error) {
	mu.Lock()
	defer mu.Unlock()

	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	if !lcache.add(ctxc{ctx: ctx, cancel: cancel}, service, addr) {
		return
	}

	err = registerdiscovery.RegisterService(
		ctx,
		service,
		addr,
		storage.EtcdClient,
	)
	if err != nil {
		lcache.del(service)
		return err
	}

	return nil
}

func InternalUnregister[T any](service string, f func(ctx context.Context, service string, add string) T) (ts []T) {
	mu.Lock()
	defer mu.Unlock()

	cache := cache[T]{
		f: f,
	}

	return cache.del(service)
}

func InternalGetRegister[T any](service string, address string, f func(ctx context.Context, service string, add string) T) (ts []T) {
	mu.Lock()
	defer mu.Unlock()

	cache := cache[T]{
		f: f,
	}

	if service == "" && address == "" {
		ts = cache.get()
		return
	}

	if service != "" {
		ts = append(ts, cache.getByService(service)...)
	}
	if address != "" {
		ts = append(ts, cache.getByAddress(address)...)
	}

	return
}

func ToPbRegisterInfo(_ context.Context, service string, add string) *pb.RegisterInfo {
	return &pb.RegisterInfo{
		Service: service,
		Address: add,
	}
}

type RegisterImp struct {
	pb.UnimplementedRegisterServer
}

func (imp *RegisterImp) GetRegister(ctx context.Context, req *pb.GetRegisterRequest) (resp *pb.GetRegisterReply, err error) {
	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("RegisterImp GetRegister")

	resp = new(pb.GetRegisterReply)

	resp.Infos = InternalGetRegister(req.Service, req.Address, ToPbRegisterInfo)

	return server.OkResp(&GResp{resp})
}

func (imp *RegisterImp) Register(ctx context.Context, req *pb.RegisterRequest) (resp *pb.RegisterReply, err error) {
	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("RegisterImp Register")

	resp = new(pb.RegisterReply)

	if req.Verify != config.AppConfig.Addr {
		return server.ParamterResp(&CResp{resp}, "verify 校验不通过")
	}
	if req.Service == "" {
		return server.ParamterResp(&CResp{resp}, "service 不能为空")
	}

	InternalRegister(context.Background(), req.Service, config.AppConfig.Addr)

	resp.Info = &pb.RegisterInfo{
		Service: req.Service,
		Address: config.AppConfig.Addr,
	}

	return server.OkResp(&CResp{resp})
}

func (imp *RegisterImp) Unregister(ctx context.Context, req *pb.UnregisterRequest) (resp *pb.UnregisterReply, err error) {
	logger, _ := ctxvalue.LoggerContext{}.GetValue(ctx)
	logger.Info("RegisterImp Unregister")

	resp = new(pb.UnregisterReply)

	if req.Verify != config.AppConfig.Addr {
		return server.ParamterResp(&DResp{resp}, "verify 校验不通过")
	}

	resp.Infos = InternalUnregister(req.Service, ToPbRegisterInfo)

	return server.OkResp(&DResp{resp})
}
