package calc

import (
	"metrics-xray/model"
)

type ITransactionStatistics interface {
	ToTransaction(*model.MetricsRecordBase) *model.TransactionV1
	Add(*TransactionStatistics)
}

type TransactionStatistics struct {
	Name      string
	Namespace string
	Count     int64
	TotalTime int64
	MaxTime   int64
	MinTime   int64
	AvgTime   int64
	Failed    int64
	Bucks     []int64

	ITransactionStatistics
}

func NewTransactionStatistics(rec *model.TransactionV1) *TransactionStatistics {
	tas := rec.Tas
	p := TransactionStatistics{
		Name:      rec.Name,
		Namespace: rec.Namespace,
		Count:     tas[model.TAG_TAS_COUNT],
		TotalTime: tas[model.TAG_TAS_TOTAL_TIME],
		MaxTime:   tas[model.TAG_TAS_MAX_TIME],
		MinTime:   tas[model.TAG_TAS_MIN_TIME],
		AvgTime:   tas[model.TAG_TAS_AVG_TIME],
		Failed:    tas[model.TAG_TAS_FAILED],
		Bucks:     rec.Bucks,
	}
	p.calcAvgTime()
	return &p
}

func (pThis *TransactionStatistics) ToTransaction(base *model.MetricsRecordBase) *model.TransactionV1 {
	rec := model.NewTransactionV1()
	rec.MetricsRecordBase = *base
	pThis.fillTas(&rec.Tas)
	rec.Bucks = pThis.Bucks
	return rec
}
func (pThis *TransactionStatistics) fillTas(tas *map[string]int64) {
	(*tas)[model.TAG_TAS_COUNT] = pThis.Count
	(*tas)[model.TAG_TAS_TOTAL_TIME] = pThis.TotalTime
	(*tas)[model.TAG_TAS_MAX_TIME] = pThis.MaxTime
	(*tas)[model.TAG_TAS_MIN_TIME] = pThis.MaxTime
	(*tas)[model.TAG_TAS_AVG_TIME] = pThis.AvgTime
	(*tas)[model.TAG_TAS_FAILED] = pThis.Failed
}

func (pThis *TransactionStatistics) Add(p *TransactionStatistics) {
	pThis.Name = p.Name
	pThis.Namespace = p.Namespace

	pThis.Count += p.Count
	pThis.Failed += p.Failed
	pThis.TotalTime += p.TotalTime
	if p.MaxTime > pThis.MaxTime {
		pThis.MaxTime = p.MaxTime
	}
	if p.MinTime < pThis.MinTime {
		pThis.MinTime = p.MinTime
	}
	pThis.calcAvgTime()
}

func (pThis *TransactionStatistics) calcAvgTime() {
	if pThis.Count > 0 {
		pThis.AvgTime = pThis.TotalTime * 1000 / pThis.Count
	} else {
		pThis.AvgTime = 0
	}
}

func (pThis *TransactionStatistics) AddRec(p *model.TransactionV1) {
	pThis.Add(NewTransactionStatistics(p))
}

type CounterStatistics struct {
	Name      string
	Namespace string
	Count     int64
	Failed    int64
}

func NewCounterStatistics(rec *model.CounterV1) *CounterStatistics {

	p := CounterStatistics{
		Name:      rec.Name,
		Namespace: rec.Namespace,
		Count:     rec.Count,
		Failed:    rec.Failed,
	}
	return &p
}

func (pThis *CounterStatistics) ToCounter(base *model.MetricsRecordBase) *model.CounterV1 {
	rec := model.NewCounterV1()
	rec.MetricsRecordBase = *base
	rec.Count = pThis.Count
	rec.Failed = pThis.Failed
	return rec
}

func (pThis *CounterStatistics) Add(p *CounterStatistics) {
	pThis.Name = p.Name
	pThis.Namespace = p.Namespace
	pThis.Count += p.Count
	pThis.Failed += p.Failed
}

func (pThis *CounterStatistics) AddRec(p *model.CounterV1) {
	pThis.Count += p.Count
	pThis.Failed += p.Failed
}

/**
/////////////////////////////////////////////////////////////////////////////////////////////////
//
// aggregate, grouby two way:
	1. app, namespace,
	2. node, idc, version
//
///////////////////////////////////////////////////////////////////////////////////////////
**/

type TransactionAggregate struct {
	tagType   string
	tagValue  string
	namespace string
	tranMap   map[string]*TransactionStatistics
}

func NewTransactionAggregate(tagType string, tagVaule string, namespace string) *TransactionAggregate {
	p := TransactionAggregate{
		tagType:   tagType,
		tagValue:  tagVaule,
		namespace: namespace,
		tranMap:   make(map[string]*TransactionStatistics),
	}
	return &p
}

func (pThis *TransactionAggregate) AddTrans(p *TransactionStatistics) *TransactionStatistics {
	//lock
	if t, ok := pThis.tranMap[p.Name]; ok {
		t.Add(p)
		return t
	} else {
		pThis.tranMap[p.Name] = p
		return t
	}
}

func (pThis *TransactionAggregate) AddNodeTrans(node *TransactionAggregate) {
	//lock
	for _, tran := range node.tranMap {
		pThis.AddTrans(tran)
	}
}
func (pThis *TransactionAggregate) GetTrans() *(map[string]*TransactionStatistics) {
	return &pThis.tranMap
}

func (pThis *TransactionAggregate) Clear() {
	//lock
	pThis.tranMap = make(map[string]*TransactionStatistics)
}

type CounterAggregate struct {
	tagType    string
	tagValue   string
	namespace  string
	counterMap map[string]*CounterStatistics
}

func NewCounterAggregate(tagType string, tagVaule string, namespace string) *CounterAggregate {
	p := CounterAggregate{
		tagType:    tagType,
		tagValue:   tagVaule,
		namespace:  namespace,
		counterMap: make(map[string]*CounterStatistics),
	}
	return &p
}

func (pThis *CounterAggregate) AddRec(p *CounterStatistics) *CounterStatistics {
	//lock
	if t, ok := pThis.counterMap[p.Name]; ok {
		t.Add(p)
		return t
	} else {
		pThis.counterMap[p.Name] = p
		return t
	}
}

func (pThis *CounterAggregate) AddNode(node *CounterAggregate) {

	for _, rec := range node.counterMap {
		pThis.AddRec(rec)
	}
}

func CalcPercentLine(bucks *[]int64) (*model.PercentLine, error) {
	p := model.PercentLine{}
	return &p, nil
}

func (pThis *CounterAggregate) GetCounters() *(map[string]*CounterStatistics) {
	return &pThis.counterMap
}

func (pThis *CounterAggregate) Clear() {
	//lock
	pThis.counterMap = make(map[string]*CounterStatistics)
}
