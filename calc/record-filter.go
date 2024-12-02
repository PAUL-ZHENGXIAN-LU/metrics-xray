package calc

import (
	"context"
	"fmt"
	"metrics-xray/model"
	"metrics-xray/model/mtype"
	"metrics-xray/repository"
	"strconv"
	"time"

	"github.com/gohutool/log4go"
)

var logger = log4go.LoggerManager.GetLogger("filter")

var (
	l_cachePool RecordChainFilter = RecordChainFilter{
		historyWrite: repository.GetLocalHistoryWrite(),
		requestWrite: repository.GetLocalRequestStoreWrite(),
	}
	ctx = context.Background()
)

type RecordChainFilter struct {
	calcPeriodSecond int32
	historyWrite     repository.IHistoryStoreWrite
	requestWrite     repository.IRequestStoreWrite

	historyRead repository.IHistoryStoreRead
	requestRead repository.IRequestStoreRead

	model.SOEV1
	model.TransactionListener
	model.CounterListener
}

func InitCalc() bool {
	l_cachePool.Init()
	return true
}

func (pThis *RecordChainFilter) Init() {
	go pThis.calcRun()
}
func (pThis *RecordChainFilter) OnSOE(rec *model.SOEV1) int {

	appPool := GetOrNewAppSOEPool(rec.App, "")
	//bu request
	if rec.BuReqId != "" {
		appReq := appPool.GetOrNewReqSOEPool(rec.BuReqId)
		appReq.OnSOE(rec)
	}
	return 1
}
func (pThis *RecordChainFilter) OnReqStart(rec *model.SOEV1) int {

	appPool := GetOrNewAppSOEPool(rec.App, "")
	//bu request
	if rec.BuReqId != "" {
		appReq := appPool.GetOrNewReqSOEPool(rec.BuReqId)
		appReq.OnReqStart(rec)
	}
	return 1
}
func (pThis *RecordChainFilter) OnReqEnd(rec *model.SOEV1) int {

	appPool := GetOrNewAppSOEPool(rec.App, "")
	//bu request
	if rec.BuReqId != "" {
		appReq := appPool.GetOrNewReqSOEPool(rec.BuReqId)
		appReq.OnReqEnd(rec)
	}
	return 1
}

func (pThis *RecordChainFilter) OnTransaction(rec *model.TransactionV1) int {

	fmt.Println("onTransaction ts=" + strconv.FormatInt(rec.Ts, 10))
	//pThis.historyWrite.WriteTransaction(TestFilterNode(*rec.Tags), rec)
	appPool := GetOrNewAppPool(rec.App, "")
	//appAgregate, _ := appPool.FindAppAggregateByTs(rec.Ts)
	appAgregate, err := appPool.GetOrNewAppAggregateByTs(rec.Ts)
	if err == nil {
		a := appAgregate.GetTranNodeOrNew(rec.Namespace, TestFilterNode(*rec.Tags))
		a.AddTrans(NewTransactionStatistics(rec))
		appAgregate.LastUpdateTime = GetNowMs()
	}

	//bu request
	if rec.BuReqId != "" {
		appReq, _ := appPool.GetOrNewAppAggregateByReq(rec.BuReqId)
		reqs := appReq.GetTranNodeOrNew(rec.Namespace, TestFilterNode(*rec.Tags))
		reqs.AddTrans(NewTransactionStatistics(rec))
		appAgregate.LastUpdateTime = GetNowMs()
	}

	return 1
}

func (pThis *RecordChainFilter) OnCounter(rec *model.CounterV1) int {

	fmt.Println("onTransaction ts=" + strconv.FormatInt(rec.Ts, 10))
	//pThis.historyWrite.WriteCounter(TestFilterNode(*rec.Tags), rec)
	appPool := GetOrNewAppPool(rec.App, "")
	appAgregate, err := appPool.GetOrNewAppAggregateByTs(rec.Ts)
	if err == nil {
		a := appAgregate.GetCounterNodeOrNew(rec.Namespace, TestFilterNode(*rec.Tags))
		a.AddRec(NewCounterStatistics(rec))
		appAgregate.LastUpdateTime = GetNowMs()
	}

	//bu request
	if rec.BuReqId != "" {
		appReq, _ := appPool.GetOrNewAppAggregateByReq(rec.BuReqId)
		reqs := appReq.GetCounterNodeOrNew(rec.Namespace, TestFilterNode(*rec.Tags))
		reqs.AddRec(NewCounterStatistics(rec))
		appAgregate.LastUpdateTime = GetNowMs()
	}
	return 1
}

///////////////////////////////////////
//
//calc

func (pThis *RecordChainFilter) calcCycle(appAgregate *AppAggregate, app *model.AppStoreInfo) {
	//lock
	appAgregate.lock.Lock()
	defer appAgregate.lock.Unlock()

	{
		namespaceList := app.GetMapGroupValues(&app.Namespaces, mtype.REC_TRANSACTION_STATISTICS)
		for _, ns := range *namespaceList {
			//agregate  the transaction group
			pThis.calcTransactionGroup(appAgregate, ns, mtype.FILTER_TAG_APP, nil)
			groups := app.GetMapGroupValues(&app.FilterTags, mtype.FILTER_TAG_IDC)
			pThis.calcTransactionGroup(appAgregate, ns, mtype.FILTER_TAG_IDC, groups)
		}
	}

	//////////////////////counter
	{
		namespaceList := app.GetMapGroupValues(&app.Namespaces, mtype.REC_COUNTER_STATISTICS)
		for _, ns := range *namespaceList {
			//agregate  the transaction group
			pThis.calcCounterGroup(appAgregate, ns, mtype.FILTER_TAG_APP, nil)
			groups := app.GetMapGroupValues(&app.FilterTags, mtype.FILTER_TAG_IDC)
			pThis.calcCounterGroup(appAgregate, ns, mtype.FILTER_TAG_IDC, groups)
		}
	}

	appAgregate.LastUpdateTime = 0

}

func (pThis *RecordChainFilter) calcHistoryCycle(appAgregate *AppAggregate, app *model.AppStoreInfo) {
	pThis.calcCycle(appAgregate, app)
}

func (pThis *RecordChainFilter) calcTransactionGroup(appAgregate *AppAggregate, namespace string, tagType string, tagVaules *[]string) {

	if tagType == mtype.FILTER_TAG_APP {
		tran := &appAgregate.AppTran
		(*tran).Clear()
		for _, nodeTran := range appAgregate.nodeTran {
			tran.AddNodeTrans(&nodeTran)
		}
		return
	}

	if tagVaules == nil || len(*tagVaules) == 0 {
		return
	}
	for _, tagValue := range *tagVaules {
		tran := appAgregate.GetTranGroupOrNew(namespace, tagType, tagValue)
		tran.Clear()
		for _, nodeTran := range appAgregate.nodeTran {
			tran.AddNodeTrans(&nodeTran)
		}
	}

}

func (pThis *RecordChainFilter) calcCounterGroup(appAgregate *AppAggregate, namespace string, tagType string, tagVaules *[]string) {
	if tagType == mtype.FILTER_TAG_APP {
		ss := &appAgregate.AppCounter
		appAgregate.AppCounter.Clear()
		for _, nodeSS := range appAgregate.nodeCounter {
			ss.AddNode(&nodeSS)
		}
		return
	}

	if tagVaules == nil || len(*tagVaules) == 0 {
		return
	}
	for _, tagValue := range *tagVaules {
		tran := appAgregate.GetTranGroupOrNew(namespace, tagType, tagValue)
		tran.Clear()
		for _, nodeTran := range appAgregate.nodeTran {
			tran.AddNodeTrans(&nodeTran)
		}
	}
}

////////////////////////////////////////////////////////////////////
//
// job

func (pThis *RecordChainFilter) onJobLastCycle(startTs int64) {

	apps := repository.LoadAppList()
	for _, app := range *apps {
		if appPool, ok := FindAppPool(app.AppId, ""); ok {
			if agg, aggOk := appPool.FindLastTimeAppAggregate(); aggOk {
				if agg.LastUpdateTime >= startTs {
					pThis.calcCycle(agg, &app)
				}
			}
		}
	}

}

func (pThis *RecordChainFilter) onJobHistoryWithChanged(startTs int64) {

	apps := repository.LoadAppList()
	for _, app := range *apps {
		if appPool, ok := FindAppPool(app.AppId, ""); ok {
			for _, agg := range appPool.AppTimsseries {
				if agg.LastUpdateTime >= startTs {
					pThis.calcHistoryCycle(agg, &app)
				}
			}
		}
	}
}

///////////////////////////////////////////////////////////////////

/*
func (pThis *RecordChainFilter) onJobCounter1M(tsMunit int64) {
	ts := tsMunit * 60

	apps := repository.LoadAppList()
	for _, app := range *apps {
		if appPool, ok := FindAppPool(app.AppId, ""); ok {
			pThis.calcCounter(ts, &app)
			for k, it := range appPool.AppTimsseries {
				if k != ts {
					pThis.calcCounter(it.ts, &app)
				}
			}
		}
	}

}

func (pThis *RecordChainFilter) onJobTransaction1M(tsMunit int64) {
	ts := tsMunit * 60

	apps := repository.LoadAppList()
	for _, app := range *apps {
		if appPool, ok := FindAppPool(app.AppId, ""); ok {
			pThis.calcTransaction(ts, &app)
			for k, it := range appPool.AppTimsseries {
				if k != ts {
					pThis.calcTransaction(it.ts, &app)
				}
			}
		}
	}
}

func (pThis *RecordChainFilter) onJobSOE1M(tsMunit int64) {
	ts := tsMunit * 60

	apps := repository.LoadAppList()
	for _, app := range *apps {
		if appPool, ok := FindAppPool(app.AppId, ""); ok {
			pThis.calcTransaction(tsMunit, &app)
			for k, it := range appPool.AppTimsseries {
				if k != ts {
					pThis.calcTransaction(it.ts, &app)
				}
			}
		}
	}
}

func (pThis *RecordChainFilter) onJobCounter(ts int64) {
	tsMunit := (ts / 60)

	apps := repository.LoadAppList()
	for _, app := range *apps {
		pThis.calcCounter(tsMunit, &app)
	}
}

func (pThis *RecordChainFilter) onJobTransaction(ts int64) {
	tsMunit := (ts / 60)

	apps := repository.LoadAppList()
	for _, app := range *apps {
		pThis.calcTransaction(tsMunit, &app)
	}
}

*/

//定时器， 处理transaction 和counter/10s

// ts --second of unix
func (pThis *RecordChainFilter) OnTick10s(ts int64, idx int) {
	logger.Info("onJob init")
	defer func() {
		logger.Info("onJob end")
	}()

	tsMunit := (ts / 60)

	if idx == 2 {
		//pThis.onJobCounter1M(tsMunit)
		//pThis.onJobTransaction1M(tsMunit)
		//pThis.onJobSOE1M(tsMunit)
		pThis.onJobHistoryWithChanged((tsMunit - 1) * 6000)
	} else {
		//pThis.onJobCounter(ts)
		//pThis.onJobTransaction(ts)
		pThis.onJobLastCycle(tsMunit * 6000)
	}

}

func (pThis *RecordChainFilter) calcRun() {

	ticker := time.NewTicker(1 * time.Second)
	var waiting int64 = (time.Now().Unix()/10 + 1) * 10
	var i int = 0
	for {
		select {
		case <-ticker.C:
			t := time.Now().Unix()
			i++
			if t >= waiting {
				pThis.OnTick10s(waiting, i/10)
				waiting += 10
				i = 0
			}

		default:
			//fmt.Println("default")
		}
	}
}
