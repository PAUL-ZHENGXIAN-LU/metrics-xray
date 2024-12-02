package repository

import (
	"metrics-xray/model"
	"metrics-xray/repository/rtype"
	"strings"
	"time"
)

const (
	CACHE_EXPIRE_TS  int32 = 3600 * 24 * 7
	CACHE_EXPIRE_REQ int32 = 3600 * 24 * 2
)

var (
	localAppStore = LocalAppStore{
		hashTable: make(map[string]*StoreHashItem),
		kvTable:   make(map[string]*StoreValueItem),
	}
	localHistoryStore = LocalHistoryStore{
		hashTable: make(map[string]*StoreHashItem),
		kvTable:   make(map[string]*StoreValueItem),
	}
	localRequestStore = LocalRequestStore{
		hashTable: make(map[string]*StoreHashItem),
		kvTable:   make(map[string]*StoreValueItem),
	}
)

type StoreItemStatus struct {
	expire  int64
	deleted bool
}

type StoreValueItem struct {
	StoreItemStatus
	data string
}

func NewStoreValueItem(v string, expireSecond int32) *StoreValueItem {
	return &StoreValueItem{
		StoreItemStatus: StoreItemStatus{
			expire:  int64(expireSecond) + time.Now().Unix(),
			deleted: false,
		},
		data: v,
	}
}

type StoreHashItem struct {
	StoreItemStatus
	table *map[string]string
}

func NewStoreHashItem(expireSecond int32) *StoreHashItem {
	tt := make(map[string]string, 0)
	return &StoreHashItem{
		table: &tt,
		StoreItemStatus: StoreItemStatus{
			expire:  int64(expireSecond) + time.Now().Unix(),
			deleted: false,
		},
	}
}

func GetLocalHistoryWrite() IHistoryStoreWrite {
	return &localHistoryStore
}

func GetLocalHistoryRead() IHistoryStoreRead {
	return &localHistoryStore
}
func GetLocalTimeSeriesRead() ITimeSeriesStoreRead {
	return &localHistoryStore
}
func GetLocalRequestStoreWrite() IRequestStoreWrite {
	return &localRequestStore
}
func GetLocalRequestStoreRead() IRequestStoreRead {
	return &localRequestStore
}
func LoadAppList() *[]model.AppStoreInfo {
	var ret []model.AppStoreInfo = make([]model.AppStoreInfo, 0)

	for _, it := range localAppStore.hashTable {
		if app, err := Map2App(it.table); err == nil {
			ret = append(ret, *app)
		}
	}
	return &ret
}

type LocalAppStore struct {
	hashTable map[string]*StoreHashItem
	kvTable   map[string]*StoreValueItem
}

func InitLocalStore() error {

	addApp("mymapp", "mysys", []string{"ta:service", "ta:task", "ta:client", "ca:service", "ca:task", "ca:rpc-call"}, []string{"idc:chicago", "idc:aws", "node"})

	return nil
}

func addApp(appId string, sys string, namespaces []string, tags []string) {

	m := NewStoreHashItem(30 * 24 * 360)
	(*m.table)["app"] = appId
	(*m.table)["sys"] = sys
	(*m.table)["namespaces"] = strings.Join(namespaces, ",")
	(*m.table)["tags"] = strings.Join(tags, ",")
	localAppStore.hashTable[appId] = m
}

type LocalHistoryStore struct {
	hashTable map[string]*StoreHashItem
	kvTable   map[string]*StoreValueItem
	IHistoryStoreWrite
	IHistoryStoreRead
	ITimeSeriesStoreRead
}

func (pThis *LocalHistoryStore) WriteTransaction(unit string, rec *model.TransactionV1) error {
	mk, str := Trans2Json(rec)
	key := rtype.KeyTransactionHistory(&rec.MetricsRecordBase)
	setHashMap(&pThis.hashTable, key, mk, str, CACHE_EXPIRE_TS)
	return nil
}

func (pThis *LocalHistoryStore) WriteCounter(unit string, rec *model.CounterV1) error {
	mk, str := Counter2Json(rec)
	key := rtype.KeyTransactionHistory(&rec.MetricsRecordBase)
	setHashMap(&pThis.hashTable, key, mk, str, CACHE_EXPIRE_TS)
	return nil
}

func (pThis *LocalHistoryStore) GetTransactions(app string, namespace string, ts int64, unit string) ([]*model.TransactionV1, error) {
	var rets []*model.TransactionV1 = make([]*model.TransactionV1, 0)
	key := rtype.KeyTransactionHistory2(app, ts, namespace)
	if h, exist := getHashMap(&pThis.hashTable, key); exist {
		for it, mk := range h {
			if _, ok := Mkey2TansactionName(mk); ok {
				rec, err := TransFromJson(it)
				if err == nil {
					rets = append(rets, rec)
				} else {

				}

			}

		}
	}
	return rets, nil
}

func (pThis *LocalHistoryStore) GetdCounters(app string, namespace string, ts int64, unit string) ([]*model.CounterV1, error) {
	rets := make([]*model.CounterV1, 0)
	key := rtype.KeyTransactionHistory2(app, ts, namespace)
	if h, exist := getHashMap(&pThis.hashTable, key); exist {
		for it, mk := range h {
			if _, ok := Mkey2CounterName(mk); ok {
				rec, err := CounterFromJson(it)
				if err == nil {
					rets = append(rets, rec)
				} else {

				}

			}

		}
	}
	return rets, nil
}

func (pThis *LocalHistoryStore) QueryTransactions(app string, namespace string, name string, start int64, size int, unit string) ([]*model.TransactionV1, error) {
	s := make([]*model.TransactionV1, 0)
	return s, nil
}

func (pThis *LocalHistoryStore) QueryCounters(app string, namespace string, name string, start int64, size int, unit string) ([]*model.CounterV1, error) {
	s := make([]*model.CounterV1, 0)
	return s, nil
}

type LocalRequestStore struct {
	hashTable map[string]*StoreHashItem
	kvTable   map[string]*StoreValueItem
	IHistoryStoreWrite
}

func NewLocalRequestStore() *LocalRequestStore {
	return nil
}

func (pThis *LocalRequestStore) WriteTransactions(rec []*model.TransactionV1) error {
	return nil
}
func (pThis *LocalRequestStore) WriteCounters(rec []*model.CounterV1) error {
	return nil
}
func (pThis *LocalRequestStore) WriteTransaction(rec *model.TransactionV1) error {
	mk, str := Trans2Json(rec)
	key := rtype.KeyTransactionHistory(&rec.MetricsRecordBase)
	setHashMap(&pThis.hashTable, key, mk, str, CACHE_EXPIRE_REQ)
	return nil
}
func (pThis *LocalRequestStore) WriteCounter(rec *model.CounterV1) error {
	return nil
}

func (pThis *LocalRequestStore) WriteGroupSOEs(rec *model.GroupSOE) error {
	return nil
}

func (pThis *LocalRequestStore) GetTransactions(app string, namespace string, reqId string) ([]*model.TransactionV1, error) {
	s := make([]*model.TransactionV1, 0)
	return s, nil
}

func (pThis *LocalRequestStore) GetCounters(app string, namespace string, reqId string) ([]*model.CounterV1, error) {
	s := make([]*model.CounterV1, 0)
	return s, nil

}

func (pThis *LocalRequestStore) GetGroupSOEs(app string, namespace string, reqId string) ([]*model.GroupSOE, error) {
	s := make([]*model.GroupSOE, 0)
	return s, nil
}
