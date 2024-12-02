package parser

import (
	"fmt"
	"metrics-xray/model"
	"testing"

	"github.com/stretchr/testify/require"
)

type tt_QueueClient struct {
	taCount  int32
	caCount  int32
	saCount  int32
	soeCount int32
	count    int32
	model.IRecordQueueClient
}

func (pThis *tt_QueueClient) PostTransaction(rec *model.TransactionV1) int {
	pThis.taCount++
	fmt.Println("on data transaction")
	return 1
}
func (pThis *tt_QueueClient) PostSOE(rec *model.SOEV1) int {
	pThis.soeCount++
	fmt.Println("on data soe")
	return 1
}
func (pThis *tt_QueueClient) PostCounter(rec *model.CounterV1) int {
	pThis.caCount++
	fmt.Println("on data counter")
	return 1
}

var (
	tt_client = tt_QueueClient{
		taCount:  0,
		caCount:  0,
		saCount:  0,
		soeCount: 0,
		count:    0,
	}
)

func TestLoadFile(t *testing.T) {

	sax := NewLogSax(&tt_client)
	sax.LoadFile("./../m.log")

	//require.(t, 32, tt_client.count)
	require.LessOrEqual(t, int32(31), tt_client.soeCount)
	fmt.Println(" test end")

}
