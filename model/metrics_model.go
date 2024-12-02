package model

import (
	"metrics-xray/model/mtype"
)

const (
	TAG_TAS_COUNT      string = "count"
	TAG_TAS_TOTAL_TIME string = "totalTime"
	TAG_TAS_MAX_TIME   string = "maxTime"
	TAG_TAS_MIN_TIME   string = "minTime"
	TAG_TAS_AVG_TIME   string = "avgTime"
	TAG_TAS_FAILED     string = "failed"
)

type AppDef struct {
	AppId string
	SysId string
}

type KVItem struct {
	Key   string
	Value string
}

type AppStoreInfo struct {
	AppId string
	SysId string
	//key is the type=ta/ca/sa
	Namespaces map[string](*[]string)
	//key is type= idc/ver/org/cver
	FilterTags map[string](*[]string)
	Nodes      *[]string
}

func (pThis *AppStoreInfo) GetMapGroupValues(m *(map[string](*[]string)), key string) *[]string {
	if r, ok := (*m)[key]; ok {
		return r
	}
	var ret []string = make([]string, 0)
	return &ret
}

type MetricsRecordBase struct {
	Ts         int64
	TypeFormat mtype.FType
	TypePeriod string
	App        string
	Namespace  string
	Name       string
	Tags       *map[string]string
	ErrLink    string
	BuReqId    string
	//Spp       string
}

type TransactionV1 struct {
	MetricsRecordBase
	Tas   map[string]int64
	Bucks []int64
}

type PercentLine struct {
	Percent99 int64
	Percent95 int64
	Percent90 int64
	Percent80 int64
}

func NewTransactionV1() *TransactionV1 {
	p := &TransactionV1{
		//MetricsRecordBase: NewMetricsRecordBase(),
		Tas:   make(map[string]int64, 10),
		Bucks: make([]int64, 0),
	}
	return p
}

type TransactionItemV1 struct {
	MetricsRecordBase
	UUId     string
	StarTime int64
	EndTime  int64
	Duration int64
	Status   int32
	Info     string
}

func NewTransactionItemV1(UUId string) *TransactionItemV1 {
	p := &TransactionItemV1{
		UUId: UUId,
	}
	return p
}

type CounterV1 struct {
	MetricsRecordBase
	Count  int64
	Failed int64
}

func NewCounterV1() *CounterV1 {
	p := &CounterV1{}
	return p
}

type SOEV1 struct {
	MetricsRecordBase
	UUId       string
	TraceGroup string
	TypeEvent  string
	RpcStep    int32
	Status     int32
	ExTags     map[string]string
	Info       string
	RpcReqId   string
}

func NewSOEV1() *SOEV1 {
	p := &SOEV1{}
	return p
}

type GroupSOE struct {
	MetricsRecordBase
	RpcReqId    string
	TraceGroup  string
	TypeEvent   string
	BaseSoeList []*BaseSOE
}

type BaseSOE struct {
	UUId     string
	Ts       int64
	RpcStep  int32
	Status   int32
	ExTags   map[string]string
	Info     string
	RpcReqId string
}

type IRecordQueueClient interface {
	PostTransaction(rec *TransactionV1) int
	PostSOE(rec *SOEV1) int
	PostCounter(rec *CounterV1) int
}

type TransactionListener interface {
	OnTransaction(rec *TransactionV1) int
}
type CounterListener interface {
	OnCounter(rec *CounterV1) int
}
type SoeListener interface {
	OnSOE(rec *SOEV1) int
	OnReqStart(rec *SOEV1) int
	OnReqEnd(rec *SOEV1) int
}
type SoeGroupListener interface {
	OnReqEventComplete(rec *SOEV1) int
	OnEventStepComplete(soes []SOEV1) int
}

type TransactionGroup struct {
	Name string
	List []*TransactionV1
}

type CounterGroup struct {
	Name string
	List []*CounterV1
}
