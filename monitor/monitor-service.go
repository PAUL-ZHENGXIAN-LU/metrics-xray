package monitor

import (
	"metrics-xray/calc"
	"metrics-xray/model/mtype"
	"metrics-xray/monitor/vo"
	"time"
)

func toMunite(t time.Time) int64 {
	ss := t.Unix()
	return ss / 60
}

func GetTransactionAggregation(query *vo.QueryAggreagetForm) *vo.TransactionAggreagtionVo {

	if appPool, ok := calc.FindAppPool(query.AppId, ""); ok {
		var appAgregate *calc.AppAggregate = nil
		var existAgg bool = false
		if query.Ts == 0 {
			appAgregate, existAgg = appPool.FindLastTimeAppAggregate()
		} else if query.Ts == -1 {
			appAgregate, existAgg = appPool.FindAppAggregateByTs(query.Ts)
			if existAgg {
				query.Ts = appAgregate.GetTs() - 60000
				appAgregate, existAgg = appPool.FindAppAggregateByTs(query.Ts)
			}
		} else {
			appAgregate, existAgg = appPool.FindAppAggregateByTs(query.Ts)
		}
		if existAgg {
			if query.FilterTagType == mtype.FILTER_TAG_APP {
				taVo := vo.Conver2TransactionAggreagteVo(query, &appAgregate.AppTran)
				return taVo
			} else if query.FilterTagType == mtype.FILTER_TAG_NODE {
				agg := appAgregate.FindTranNode(query.Namespace, query.FilterTagVaule)
				taVo := vo.Conver2TransactionAggreagteVo(query, agg)
				return taVo
			}
		}
	}
	return nil
}

func GetCounterAggregation(query *vo.QueryAggreagetForm) *vo.CounterAggreagtionVo {

	if appPool, ok := calc.FindAppPool(query.AppId, ""); ok {
		if appAgregate, okTs := appPool.FindAppAggregateByTs(query.Ts); okTs {
			if query.FilterTagType == mtype.FILTER_TAG_APP {
				countersVo := vo.Conver2CounterAggreagteVo(query, &appAgregate.AppCounter)
				return countersVo
			} else if query.FilterTagType == mtype.FILTER_TAG_NODE {
				agg := appAgregate.FindCounterNode(query.Namespace, query.FilterTagVaule)
				taVo := vo.Conver2CounterAggreagteVo(query, agg)
				return taVo
			}
		}
	}
	return nil
}

func GetTransactionByBuRequest(query *vo.QueryAggreagetForm) *vo.TransactionAggreagtionVo {

	if appPool, ok := calc.FindAppPool(query.AppId, ""); ok {
		if appAgregate, okData := appPool.FindAppAggregateByReq(query.BuReqId); okData {
			if query.FilterTagType == mtype.FILTER_TAG_APP {
				taVo := vo.Conver2TransactionAggreagteVo(query, &appAgregate.AppTran)
				return taVo
			}
		}
	}
	return nil
}

func GetCounterByBuRequest(query *vo.QueryAggreagetForm) *vo.CounterAggreagtionVo {

	if appPool, ok := calc.FindAppPool(query.AppId, ""); ok {
		if appAgregate, okData := appPool.FindAppAggregateByReq(query.BuReqId); okData {
			if query.FilterTagType == mtype.FILTER_TAG_APP {
				countersVo := vo.Conver2CounterAggreagteVo(query, &appAgregate.AppCounter)
				return countersVo
			}
		}
	}
	return nil
}

func GetBuRequestSoeGroup(query *vo.QueryAggreagetForm) *vo.SoeGroupRequestVo {

	return nil
}
