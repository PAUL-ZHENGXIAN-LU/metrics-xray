package rtype

import (
	"metrics-xray/model"
	"strconv"
)

const (
	MY_HEAD int32 = 6
)

/*
**********************************

system:
key, mt:[tenerid]:sys:format  ash
mk name
mk ...
mk, ta-fd:[version]:[property]:[type int/double/string/ints/doubles/strings/link]

key, mt:[tenerid]:sys:app-queue   zqueue  the app queue  the err--->, trancation(service+uri)-->
mk system, ta-fd:[version]:[property]:[type int/double/string/ints/doubles/strings/link]

app:
key, mt:[app]:app  def


transaction:
key, mt:[app]:[freq]:[ts]:ta:[ns]   hash,
--->property

key, mt:[app]:[freq]:[ts]:ta:[ns]:qurue   zqueue,

counter:
key, mt:[app]:[freq]:[ts]:ca:[ns]   hash,  count+failed
mk,  count, failed

counter:
key, mt:[app]:[freq]:[ts]:sa:[ns]   string
key, mt:[app]:[freq]:[ts]:sa:[ns]:detail   list-->string[6]



*
*/

func KeyTransactionHistory(rec *model.MetricsRecordBase) string {
	key := "mt:" + rec.App + ":1m:" + Ts2Tag(rec.Ts) + ":ta:" + rec.Namespace
	return key
}
func KeyTransactionHistory2(app string, ts int64, namespace string) string {
	key := "mt:" + app + ":1m:" + Ts2Tag(ts) + ":ta:" + namespace
	return key
}

func KeyEventGroup(rec *model.MetricsRecordBase) string {
	key := "mt:" + rec.App + ":-1s:" + rec.BuReqId + ":e:" + rec.Namespace
	return key
}

func KeyEventGroup2(app string, buReqId string, namespace string) string {
	key := "mt:" + app + ":-1s:" + buReqId + ":e:" + namespace
	return key
}

func Ts2Tag(ts int64) string {
	return strconv.FormatInt(ts, 10)
}

type TransactionEntity struct {
	Tas   map[string]uint64
	Bucks []int64
}
