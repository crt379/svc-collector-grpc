package server

import (
	"time"

	"github.com/crt379/svc-collector-grpc/internal/types"

	pbappapi "github.com/crt379/svc-collector-grpc-proto/appapi"
	pbapp "github.com/crt379/svc-collector-grpc-proto/application"
	pbappproc "github.com/crt379/svc-collector-grpc-proto/appproc"
	pbappsvc "github.com/crt379/svc-collector-grpc-proto/appsvc"
	pbprocessor "github.com/crt379/svc-collector-grpc-proto/processor"
	pbservice "github.com/crt379/svc-collector-grpc-proto/service"
	pbsvcapi "github.com/crt379/svc-collector-grpc-proto/svcapi"
	pbsvcapieg "github.com/crt379/svc-collector-grpc-proto/svcapieg"
	pbtenant "github.com/crt379/svc-collector-grpc-proto/tenant"

	"github.com/golang/protobuf/ptypes/timestamp"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func TimeToPtime(t types.Time) *timestamp.Timestamp {
	return timestamppb.New(time.Time(t))
}

type TPBMeta[T any, PT any] interface {
	*T
	ToPbMeta() (PT, error)
}

func Metas2Pbmeta[T any, PT any, TP TPBMeta[T, PT]](metas *[]T) (l int32, pbmeta []*PT, err error) {
	l = int32(len(*metas))
	pbmeta = make([]*PT, l)
	for i, m := range *metas {
		var (
			pt PT
			mp TP = &m
		)

		pt, err = mp.ToPbMeta()
		if err != nil {
			return
		}

		pbmeta[i] = &pt
	}

	return l, pbmeta, err
}

type TenantMeta struct {
	Uuid       int        `json:"uuid" db:"uuid"`
	Name       string     `json:"name" db:"name"`
	Describe   string     `json:"describe" db:"describe"`
	CreateTime types.Time `json:"create_time" db:"create_time"`
	UpdateTime types.Time `json:"update_time" db:"update_time"`
}

func (m *TenantMeta) ToPbMeta() (pbtenant.TenantMeta, error) {
	return pbtenant.TenantMeta{
		Uuid:       int32(m.Uuid),
		Name:       m.Name,
		Describe:   m.Describe,
		CreateTime: TimeToPtime(m.CreateTime),
		UpdateTime: TimeToPtime(m.UpdateTime),
	}, nil
}

func (m *TenantMeta) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *TenantMeta) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

type ServiceMeta struct {
	Uuid       int        `json:"uuid" db:"uuid"`
	Name       string     `json:"name" db:"name"`
	Describe   string     `json:"describe" db:"describe"`
	CreateTime types.Time `json:"create_time" db:"create_time"`
	UpdateTime types.Time `json:"update_time" db:"update_time"`
	TenantId   int        `json:"tenant_id" db:"tenant_id"`
}

func (m *ServiceMeta) ToPbMeta() (pbservice.ServiceMeta, error) {
	return pbservice.ServiceMeta{
		Uuid:       int32(m.Uuid),
		Name:       m.Name,
		Describe:   m.Describe,
		CreateTime: TimeToPtime(m.CreateTime),
		UpdateTime: TimeToPtime(m.UpdateTime),
		TenantId:   int32(m.TenantId),
	}, nil
}

func (m *ServiceMeta) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *ServiceMeta) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

type SvcapiMeta struct {
	Uuid       int        `json:"uuid" db:"uuid"`
	Path       string     `json:"path" db:"path"`
	Method     string     `json:"method" db:"method"`
	Describe   string     `json:"describe" db:"describe"`
	CreateTime types.Time `json:"create_time" db:"create_time"`
	UpdateTime types.Time `json:"update_time" db:"update_time"`
	ServiceId  int        `json:"service_id" db:"sid"`
	TenantId   int        `json:"tenant_id" db:"tenant_id"`
}

func (m *SvcapiMeta) ToPbMeta() (pbsvcapi.SvcapiMeta, error) {
	return pbsvcapi.SvcapiMeta{
		Uuid:       int32(m.Uuid),
		Path:       m.Path,
		Method:     m.Method,
		Describe:   m.Describe,
		CreateTime: TimeToPtime(m.CreateTime),
		UpdateTime: TimeToPtime(m.UpdateTime),
		TenantId:   int32(m.TenantId),
		ServiceId:  int32(m.ServiceId),
	}, nil
}

type SvcapiegMeta struct {
	Uuid       int        `json:"uuid" db:"uuid"`
	Data       any        `json:"data" db:"data"`
	CreateTime types.Time `json:"create_time" db:"create_time"`
	UpdateTime types.Time `json:"update_time" db:"update_time"`
	SvcapiId   int        `json:"svcapi_id" db:"aid"`
	ServiceId  int        `json:"service_id" db:"-"`
	TenantId   int        `json:"tenant_id" db:"tenant_id"`
	JdataId    int        `json:"-" db:"jid"`
}

func (m *SvcapiegMeta) DataToMap() error {
	var (
		data   []byte
		grades map[string]interface{}
	)

	switch v := m.Data.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return nil
	}

	m.Data = grades

	return json.Unmarshal(data, &m.Data)
}

func (m *SvcapiegMeta) GetData() (string, error) {
	var (
		jd  []byte
		err error
	)

	err = m.DataToMap()
	if err != nil {
		return "", err
	}

	jd, err = json.Marshal(m.Data)

	return string(jd), err
}

func (m *SvcapiegMeta) ToPbMeta() (pbsvcapieg.SvcapiegMeta, error) {
	data, err := m.GetData()

	return pbsvcapieg.SvcapiegMeta{
		Uuid:       int32(m.Uuid),
		Data:       data,
		CreateTime: TimeToPtime(m.CreateTime),
		UpdateTime: TimeToPtime(m.UpdateTime),
		SvcapiId:   int32(m.SvcapiId),
		// ServiceId:  int32(m.ServiceId),
		TenantId: int32(m.TenantId),
	}, err
}

type Jdata struct {
	Uuid       int        `db:"uuid"`
	Data       any        `db:"data"`
	CreateTime types.Time `db:"create_time"`
	UpdateTime types.Time `db:"update_time"`
	HashType   string     `db:"hash_type"`
	HashValue  string     `db:"hash_value"`
}

type ApplicationMeta struct {
	Uuid       int        `json:"uuid" db:"uuid"`
	Name       string     `json:"name" db:"name"`
	Describe   string     `json:"describe" db:"describe"`
	CreateTime types.Time `json:"create_time" db:"create_time"`
	UpdateTime types.Time `json:"update_time" db:"update_time"`
	TenantId   int        `json:"tenant_id" db:"tenant_id"`
}

func (m *ApplicationMeta) ToPbMeta() (pbapp.ApplicationMete, error) {
	return pbapp.ApplicationMete{
		Uuid:       int32(m.Uuid),
		Name:       m.Name,
		Describe:   m.Describe,
		CreateTime: TimeToPtime(m.CreateTime),
		UpdateTime: TimeToPtime(m.UpdateTime),
		TenantId:   int32(m.TenantId),
	}, nil
}

type AppsvcMeta struct {
	Uuid       int         `json:"uuid" db:"uuid"`
	AppId      int         `json:"appid" db:"aid"`
	SvcId      int         `json:"svcid" db:"sid"`
	CreateTime types.Time  `json:"create_time" db:"create_time"`
	UpdateTime types.Time  `json:"update_time" db:"update_time"`
	Service    ServiceMeta `json:"service" db:"-"`
	SvcName    string      `json:"-" db:"-"`
}

func (m *AppsvcMeta) ToPbMeta() (pbappsvc.AppsvcMeta, error) {
	svc, _ := m.Service.ToPbMeta()
	return pbappsvc.AppsvcMeta{
		Uuid:       int32(m.Uuid),
		Appid:      int32(m.AppId),
		Svcid:      int32(m.SvcId),
		CreateTime: TimeToPtime(m.CreateTime),
		UpdateTime: TimeToPtime(m.UpdateTime),
		Service:    &svc,
	}, nil
}

type ProcessorMeta struct {
	Uuid       int        `json:"uuid" db:"uuid"`
	Addr       string     `json:"addr" db:"addr"`
	Weight     int        `json:"weight" db:"weight"`
	State      string     `json:"state" db:"state"`
	CreateTime types.Time `json:"create_time" db:"create_time"`
	UpdateTime types.Time `json:"update_time" db:"update_time"`
	AppId      int        `json:"aid" db:"aid"`
	TanantId   int        `json:"tenant_id" db:"tenant_id"`
}

func (m *ProcessorMeta) ToPbMeta() (pbprocessor.ProcessorMeta, error) {
	return pbprocessor.ProcessorMeta{
		Uuid:       int32(m.Uuid),
		Addr:       m.Addr,
		Weight:     int32(m.Weight),
		State:      m.State,
		CreateTime: TimeToPtime(m.CreateTime),
		UpdateTime: TimeToPtime(m.UpdateTime),
		Appid:      int32(m.AppId),
		TenantId:   int32(m.TanantId),
	}, nil
}

type AAapi struct {
	Application ApplicationMeta `json:"application"`
	Service     ServiceMeta     `json:"service"`
	Svcapis     []SvcapiMeta    `json:"svcapis"`
}

type AppapiMeta struct {
	Appid    int    `json:"-"`
	Appname  string `json:"-"`
	Appsvcid int    `json:"-"`
	Svcid    int    `json:"-"`
	Svcname  string `json:"-"`
	TenantId int    `json:"-"`
	Appapi   AAapi  `json:"-"`
}

func (m *AppapiMeta) ToPbMeta() (pbappapi.AppapiMeta, error) {
	app, _ := m.Appapi.Application.ToPbMeta()
	svc, _ := m.Appapi.Service.ToPbMeta()
	_, apis, err := Metas2Pbmeta[SvcapiMeta, pbsvcapi.SvcapiMeta](&m.Appapi.Svcapis)

	return pbappapi.AppapiMeta{
		Application: &app,
		Service:     &svc,
		Svcapis:     apis,
	}, err
}

type A3p struct {
	Application ApplicationMeta `json:"application"`
	Processors  []ProcessorMeta `json:"processors"`
}

type AppprocMeta struct {
	Appid    int    `json:"-"`
	Appname  string `json:"-"`
	Weight   *int   `json:"-"`
	State    string `json:"-"`
	TenantId int    `json:"-"`
	A3p      A3p    `json:"-"`
}

func (m *AppprocMeta) ToPbMeta() (pbappproc.AppprocMeta, error) {
	app, _ := m.A3p.Application.ToPbMeta()
	_, procs, err := Metas2Pbmeta[ProcessorMeta, pbprocessor.ProcessorMeta](&m.A3p.Processors)

	return pbappproc.AppprocMeta{
		Application: &app,
		Processors:  procs,
	}, err
}
