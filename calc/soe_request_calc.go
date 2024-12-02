package calc

import (
	"container/list"
	"metrics-xray/model"
	"time"
)

type ReqSOE struct {
	BuReqId string
	EndTime int64
	soeList list.List //*model.SOEV1
	soeMap  map[string]*model.SOEV1
	//namespace-tracegroup
	Group map[string]*model.GroupSOE

	GroupTran map[string]*model.TransactionItemV1
	model.SoeListener
}

func NewReqSOE(reqId string) *ReqSOE {
	ll := list.New()
	r := &ReqSOE{
		soeList: *ll,
		BuReqId: reqId,
		EndTime: 0,
	}
	return r
}

func (pThis *ReqSOE) addSoe(soe *model.SOEV1) {
	pThis.soeList.PushBack(soe)
}

func (pThis *ReqSOE) setEnd(endTime int64) {
	pThis.EndTime = endTime
}

func (pThis *AppSOEPool) OnSOE(soe *model.SOEV1) int {
	//pThis.setSoe(soe)
	return 1
}

func (pThis *AppSOEPool) setSoe(soe *model.SOEV1) *ReqSOE {

	if req, exist := pThis.reqSOEPool[soe.BuReqId]; exist {
		req.addSoe(soe)
		return req
	}
	r := NewReqSOE(soe.BuReqId)
	pThis.reqSOEPool[soe.BuReqId] = r
	r.addSoe(soe)

	return r
}

// ????
func (pThis *AppSOEPool) OnReqStart(soe *model.SOEV1) int {
	pThis.setSoe(soe)
	return 1
}

// waiting 1-5 minute
func (pThis *AppSOEPool) OnReqEnd(soe *model.SOEV1) int {
	req := pThis.setSoe(soe)
	t := time.Now()
	req.setEnd(t.UnixNano())
	return 1
}
