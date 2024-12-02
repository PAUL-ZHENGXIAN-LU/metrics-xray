package calc

import (
	"metrics-xray/model"
)

var (
	G_queueClient *LogQueueClient = &LogQueueClient{}

	l_logQueue *LogQueue = NewQueue()
)

type IFMetricsRecord interface {
}

type LogQueueClient struct {
	model.IRecordQueueClient
}

func (pThis *LogQueueClient) PostTransaction(rec *model.TransactionV1) int {
	var i IFMetricsRecord = rec
	l_logQueue.queue <- i
	return 1
}
func (pThis *LogQueueClient) PostSOE(rec *model.SOEV1) int {
	var i IFMetricsRecord = rec
	l_logQueue.queue <- i
	return 1
}
func (pThis *LogQueueClient) PostCounter(rec *model.CounterV1) int {
	var i IFMetricsRecord = rec
	l_logQueue.queue <- i
	return 1
}

type LogQueue struct {
	queue             chan IFMetricsRecord
	soeHandle         model.SoeListener
	transactionHandle model.TransactionListener
	counterHandle     model.CounterListener
}

func NewQueue() *LogQueue {
	r := &LogQueue{
		queue:             make(chan IFMetricsRecord, 10000),
		soeHandle:         &l_cachePool,
		transactionHandle: &l_cachePool,
		counterHandle:     &l_cachePool,
	}

	go r.queueRun()

	return r
}

func (pThis *LogQueue) queueRun() {

	for {
		select {
		case iface := <-pThis.queue:
			if ta, ok := iface.(*model.TransactionV1); ok {
				pThis.transactionHandle.OnTransaction(ta)
			} else if soe, ok := iface.(*model.SOEV1); ok {
				pThis.soeHandle.OnSOE(soe)
			} else if ca, ok := iface.(*model.CounterV1); ok {
				pThis.counterHandle.OnCounter(ca)
			}
			//case num2 := <-ch2:
			//  fmt.Println(num2)
		default:
			//fmt.Println("default")
		}
	}
}
