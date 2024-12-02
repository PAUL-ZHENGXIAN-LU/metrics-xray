package vo

import "metrics-xray/calc"

///data type base on uri
type QueryAggreagetForm struct {
	Tenant    string
	AppId     string
	Namespace string

	//contision
	FilterTagType  string
	FilterTagVaule string
	Freq           string

	Ts      int64
	BuReqId string
}

type QueryBuRequestForm struct {
	AppId     string
	Namespace string
	BuReqId   string
}

type TransactionAggreagtionVo struct {
	AppId          string
	Namespace      string
	Ts             int64
	FilterTagType  string
	FilterTagVaule string
	Freq           string
	Datas          map[string]*TransactionVo
}

func NewTransactionAggreagteVo(appId string, ns string, filterTagType string, filterTag string) *TransactionAggreagtionVo {
	ret := TransactionAggreagtionVo{
		AppId:          appId,
		Namespace:      ns,
		FilterTagType:  filterTagType,
		FilterTagVaule: filterTag,
		Datas:          make(map[string]*TransactionVo),
	}

	return &ret
}

func Conver2TransactionAggreagteVo(query *QueryAggreagetForm, aggregation *calc.TransactionAggregate) *TransactionAggreagtionVo {
	ret := NewTransactionAggreagteVo(query.AppId, query.Namespace, query.FilterTagType, query.FilterTagVaule)
	mapTran := aggregation.GetTrans()
	for name, tran := range *mapTran {
		vo := Conver2TransactionVo(tran)
		ret.Datas[name] = vo
	}
	return ret
}

type TransactionVo struct {
	Name      string
	Count     int64
	Failed    int64
	TotalTime int64
	MaxTime   int64
	MinTime   int64
	AvgTime   float32
	Percent99 int64
	Percent95 int64
	Percent90 int64
	Percent80 int64
}

func Conver2TransactionVo(obj *calc.TransactionStatistics) *TransactionVo {
	ret := TransactionVo{
		Name:      obj.Name,
		Count:     obj.Count,
		Failed:    obj.Failed,
		TotalTime: obj.TotalTime,
		MaxTime:   obj.MaxTime,
		MinTime:   obj.MinTime,
		AvgTime:   float32(obj.AvgTime) / 1000.0,
	}

	return &ret
}

type CounterAggreagtionVo struct {
	AppId          string
	Namespace      string
	Ts             int64
	FilterTagType  string
	FilterTagVaule string
	Freq           string
	Datas          map[string]*CounterVo
}

func NewCounterAggreagtionVo(appId string, ns string, filterTagType string, filterTag string) *CounterAggreagtionVo {
	ret := CounterAggreagtionVo{
		AppId:          appId,
		Namespace:      ns,
		FilterTagType:  filterTagType,
		FilterTagVaule: filterTag,
		Datas:          make(map[string]*CounterVo),
	}

	return &ret
}

type CounterVo struct {
	Name      string
	Count     int64
	Failed    int64
	TotalTime int64
	MaxTime   int64
	MinTime   int64
	AvgTime   float32
	Percent99 int64
	Percent95 int64
	Percent90 int64
	Percent80 int64
}

func Conver2CounterAggreagteVo(query *QueryAggreagetForm, aggregation *calc.CounterAggregate) *CounterAggreagtionVo {
	ret := NewCounterAggreagtionVo(query.AppId, query.Namespace, query.FilterTagType, query.FilterTagVaule)
	m := aggregation.GetCounters()
	for name, obj := range *m {
		vo := Conver2CounterVo(obj)
		ret.Datas[name] = vo
	}
	return ret
}

func Conver2CounterVo(obj *calc.CounterStatistics) *CounterVo {
	ret := CounterVo{
		Name:   obj.Name,
		Count:  obj.Count,
		Failed: obj.Failed,
	}
	return &ret
}

////////////////////////////////////////////////////////////
//
// business request event
//
///////////////////////////////////////////////////////////

type SoeGroupRequestVo struct {
	AppId     string
	Namespace string
	//Ts            int64
	BuRequestId string
	Datas       map[string]*SoeGroupVo
}

type SoeGroupVo struct {
	Name       string
	EType      string
	Completed  bool
	Count      int32
	Failed     int64
	StartTime  int64
	EndTime    int64
	DurationMs int32
	List       []*SoeVo
}

type SoeVo struct {
	Step   string
	Ts     int64
	Status int32
	Info   string
}
