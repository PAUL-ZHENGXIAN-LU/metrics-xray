package repository

import "metrics-xray/model"

/**
filter unit  is the tags for opreater unit
for sample: ip, pod node, instance(ip+process)
**/

type IHistoryStoreWrite interface {
	//WriteTransactions(unit string, rec []*model.TransactionV1) error
	//WriteCounters(unit string, rec []*model.CounterV1) error

	WriteTransaction(unit string, rec *model.TransactionV1) error
	WriteCounter(unit string, rec *model.CounterV1) error
}

type IHistoryStoreRead interface {
	GetTransactions(app string, namespace string, ts int64, unit string) ([]*model.TransactionV1, error)
	GetCounters(app string, namespace string, ts int64, unit string) ([]*model.CounterV1, error)

	GetTransactionByAllNodes(app string, namespace string, ts int64, nodes []string) (*model.TransactionGroup, error)
	GetCounterByAllNodes(app string, namespace string, ts int64, nodes []string) (*model.TransactionGroup, error)
}

type ITimeSeriesStoreRead interface {
	QueryTransactions(app string, namespace string, name string, start int64, size int, unit string) ([]*model.TransactionV1, error)
	QueryCounters(app string, namespace string, name string, start int64, size int, unit string) ([]*model.CounterV1, error)
}

type IRequestStoreWrite interface {
	WriteTransactions(rec []*model.TransactionV1) error
	WriteCounters(rec []*model.CounterV1) error
	//WriteGroupSOEs(rec []*model.GroupSOE) error

	WriteTransaction(rec *model.TransactionV1) error
	WriteCounter(rec *model.CounterV1) error
	WriteGroupSOEs(rec *model.GroupSOE) error
}

type IRequestStoreRead interface {
	GetTransactions(app string, namespace string, reqId string) ([]*model.TransactionV1, error)
	GetCounters(app string, namespace string, reqId string) ([]*model.CounterV1, error)
	GetGroupSOEs(app string, namespace string, reqId string) ([]*model.GroupSOE, error)
	//GetIncompleteSOEs(app string, namespace string, reqId string) ([]*model.SOEV1, error)
	//GetSOEs(app string, namespace string, ts int64) ([]*model.SOEV1, error)
}
