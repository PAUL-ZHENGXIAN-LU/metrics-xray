package calc

import (
	"errors"
	"metrics-xray/model/mtype"
	"sync"
)

const (
	AGGREGATE_MAX_PERIED int   = 60
	CALC_STEP_SECOND     int64 = 60
	CALC_CLEAR_COUNT     int   = 30

	CACHE_MAX_CYCLE int = 181
	CACHE_MIN_CYCLE int = 120
)

var (
	l_aggregatePool AggregatePool = AggregatePool{
		appTimsseries: make(map[string]*AppAggregatePool),
	}
)

//////////////////////////////////////////////////////////////////////////////////////////////////
//
// cache pool
//
////////////////////////////////////////////////////////////////////////////////////////////

// /?? need support tenant

type AggregatePool struct {
	appTimsseries map[string](*AppAggregatePool)
}

func FindAppPool(appId string, env string) (*AppAggregatePool, bool) {
	if a, ok := l_aggregatePool.appTimsseries[appId]; ok {
		return a, true
	}
	return nil, false
}

func GetOrNewAppPool(appId string, env string) *AppAggregatePool {
	if a, ok := l_aggregatePool.appTimsseries[appId]; ok {
		return a
	}
	a := AppAggregatePool{
		app:           appId,
		starTime:      0,
		endTime:       0,
		AppTimsseries: make(map[int64]*AppAggregate),
		appBuReqs:     make(map[string]*AppAggregate),
	}
	//a.AppTimsseries[0] = NewAppAggregate(appId, lastTs)
	l_aggregatePool.appTimsseries[appId] = &a
	return &a
}

type AppAggregatePool struct {
	app           string
	starTime      int64
	endTime       int64
	AppTimsseries map[int64]*AppAggregate

	timsseriesLock sync.RWMutex

	appBuReqs map[string]*AppAggregate
}

func (pThis *AppAggregatePool) GetAppAggregateCount() int {
	return len(pThis.AppTimsseries)
}

/*
	func (pThis *AppAggregatePool) GetAppAggregateTimeList() []int64 {
		count := len(pThis.AppTimsseries)
		ret := make([]int64, count)
		for i:=0; i<count; i++{
			ret[i] = pThis.AppTimsseries[i].ts
		}
		return ret
	}
*/
func (pThis *AppAggregatePool) FindLastTimeAppAggregate() (*AppAggregate, bool) {
	//lock
	pThis.timsseriesLock.RLock()
	defer pThis.timsseriesLock.RUnlock()
	if len(pThis.AppTimsseries) == 0 {
		return nil, false
	}

	var lastTs int64 = 0
	for k, _ := range pThis.AppTimsseries {
		if lastTs == 0 || k > lastTs {
			lastTs = k
		}
	}
	ret, ok := pThis.AppTimsseries[lastTs]
	return ret, ok
}

func (pThis *AppAggregatePool) FindAppAggregateByTs(ts int64) (*AppAggregate, bool) {
	//lock
	pThis.timsseriesLock.RLock()
	defer pThis.timsseriesLock.RUnlock()
	key := GetCycleKey(ts)
	if a, ok := pThis.AppTimsseries[key]; ok {
		return a, true
	}
	return nil, false
}

func (pThis *AppAggregatePool) GetOrNewAppAggregateByTs(ts int64) (*AppAggregate, error) {

	if a, ok := pThis.FindAppAggregateByTs(ts); ok {
		return a, nil
	}
	lastTs := GetCycleKey(ts)
	//size := int(lastTs - pThis.endTime)
	////if size <= 0 || size > AGGREGATE_MAX_PERIED {
	//return nil, errors.New("out of the max size of pool")
	//}

	//lock
	pThis.timsseriesLock.Lock()
	defer pThis.timsseriesLock.Unlock()

	if a, ok := pThis.AppTimsseries[lastTs]; ok {
		return a, nil
	}

	ret := NewAppAggregate(pThis.app, lastTs)
	pThis.AppTimsseries[lastTs] = ret
	if pThis.starTime == 0 || lastTs < pThis.starTime {
		pThis.starTime = lastTs
	}
	endTs := lastTs + 60000 - 1
	if pThis.endTime == 0 || endTs > pThis.endTime {
		pThis.endTime = endTs
	}

	pThis.CheckAndClearTs(lastTs)
	return ret, nil
}

func (pThis *AppAggregatePool) CheckAndClearTs(ts int64) int {

	if len(pThis.AppTimsseries) >= CACHE_MAX_CYCLE {
		//clear CALC_CLEAR_TS_COUNT
		//start := GetNewStartTime()
		removeList := make([]int64, 0)
		start := GetNewStartTime(ts)
		for k, _ := range pThis.AppTimsseries {
			if k < start {
				removeList = append(removeList, k)
			}
		}
		for _, it := range removeList {
			delete(pThis.AppTimsseries, it)
		}
	}

	return 0
}

func (pThis *AppAggregatePool) FindAppAggregateByReq(reqId string) (*AppAggregate, bool) {
	if reqId == "" {
		return nil, false
	}
	v, ok := pThis.appBuReqs[reqId]
	return v, ok
}

func (pThis *AppAggregatePool) GetOrNewAppAggregateByReq(reqId string) (*AppAggregate, error) {
	if reqId == "" {
		return nil, errors.New("the reqid is null")
	}
	if appAggregate, ok := pThis.appBuReqs[reqId]; ok {
		return appAggregate, nil
	}

	appAggregate := NewAppAggregate(pThis.app, -1)
	appAggregate.reqId = reqId
	pThis.appBuReqs[reqId] = appAggregate
	return appAggregate, nil
}

type AppAggregate struct {
	app       string
	ts        int64
	reqId     string
	nodeTran  map[string]TransactionAggregate
	groupTran map[string]TransactionAggregate
	AppTran   TransactionAggregate

	nodeCounter  map[string]CounterAggregate
	groupCounter map[string]CounterAggregate
	AppCounter   CounterAggregate

	lock sync.RWMutex
	//second, last update time, when the calc, the time reset to 0
	LastUpdateTime int64
}

func NewAppAggregate(appId string, ts int64) *AppAggregate {
	p := AppAggregate{
		app:            appId,
		ts:             ts,
		nodeTran:       make(map[string]TransactionAggregate),
		groupTran:      make(map[string]TransactionAggregate),
		AppTran:        *NewTransactionAggregate(mtype.FILTER_TAG_APP, "", ""),
		nodeCounter:    make(map[string]CounterAggregate),
		groupCounter:   make(map[string]CounterAggregate),
		AppCounter:     *NewCounterAggregate(mtype.FILTER_TAG_APP, "", ""),
		LastUpdateTime: GetNowMs(),
	}

	return &p
}

func (pThis *AppAggregate) GetTs() int64 {
	return pThis.ts
}

func (pThis *AppAggregate) FindTranNode(namespace string, node string) *TransactionAggregate {
	if a, ok := pThis.nodeTran[node]; ok {
		return &a
	}
	return nil
}

func (pThis *AppAggregate) GetTranNodeOrNew(namespace string, node string) *TransactionAggregate {
	if a, ok := pThis.nodeTran[node]; ok {
		return &a
	}
	a := NewTransactionAggregate(mtype.FILTER_TAG_NODE, node, namespace)
	pThis.nodeTran[node] = *a
	return a
}

func (pThis *AppAggregate) FindTranGroup(namespace string, tagType string, tagValue string) *TransactionAggregate {
	if a, ok := pThis.groupTran[tagValue]; ok {
		return &a
	}
	return nil
}

func (pThis *AppAggregate) GetTranGroupOrNew(namespace string, tagType string, tagValue string) *TransactionAggregate {
	if a, ok := pThis.groupTran[tagValue]; ok {
		return &a
	}
	a := NewTransactionAggregate(tagType, tagValue, namespace)
	pThis.groupTran[TagGroupKey(tagType, tagValue)] = *a
	return a
}

func (pThis *AppAggregate) FindCounterNode(namespace string, node string) *CounterAggregate {
	if a, ok := pThis.nodeCounter[node]; ok {
		return &a
	}
	return nil
}

func (pThis *AppAggregate) GetCounterNodeOrNew(namespace string, node string) *CounterAggregate {
	if a, ok := pThis.nodeCounter[node]; ok {
		return &a
	}
	a := NewCounterAggregate(mtype.FILTER_TAG_NODE, node, namespace)
	pThis.nodeCounter[node] = *a
	return a
}

func (pThis *AppAggregate) FindCounterGroup(namespace string, tagType string, tagValue string) *CounterAggregate {
	if a, ok := pThis.groupCounter[tagValue]; ok {
		return &a
	}
	return nil
}

func (pThis *AppAggregate) GetCounterGroupOrNew(namespace string, tagType string, tagValue string) *CounterAggregate {
	if a, ok := pThis.groupCounter[tagValue]; ok {
		return &a
	}
	a := NewCounterAggregate(tagType, tagValue, namespace)
	pThis.groupCounter[TagGroupKey(tagType, tagValue)] = *a
	return a
}
