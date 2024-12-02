package monitor

import (
	"metrics-xray/monitor/vo"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gohutool/log4go"
)

var logger = log4go.LoggerManager.GetLogger("collecter")

func getHelp(c *gin.Context) {
	c.String(http.StatusOK, "req unready interface:"+c.FullPath())
}

func queryAppList(c *gin.Context) {

	//env := c.Query("tenant") //  get app
	//calc.GetOrNewAppPool(appId, env, lastTs)

	c.String(http.StatusOK, "queryAppList")
}

func mapRequestForm(c *gin.Context) *vo.QueryAggreagetForm {

	query := vo.QueryAggreagetForm{
		Tenant:    c.Query("tenant"),
		AppId:     c.Query("appId"),
		Namespace: c.Query("namespace"),
		//contision
		FilterTagType:  c.Query("filterTagType"),
		FilterTagVaule: c.Query("filterTagVaule"),
		Freq:           c.Query("freq"),

		Ts:      ParserLong(c.Query("ts")),
		BuReqId: c.Query("buReqId"),
	}
	return &query
}

func queryLastRequestList(c *gin.Context) {
	c.String(http.StatusOK, "queryLastRequestList")
}

func getTransaction(c *gin.Context) {
	if c.Query("appId") == "" || c.Query("namespace") == "" || c.Query("ts") == "" {
		c.String(http.StatusBadRequest, "the para is wrong")
		logger.Warning("instanceId is empty:")
		return
	}
	query := mapRequestForm(c) //
	ret := GetTransactionAggregation(query)
	c.JSON(http.StatusOK, ret)
}

func getCounter(c *gin.Context) {
	if c.Query("appId") == "" || c.Query("namespace") == "" || c.Query("ts") == "" {
		c.String(http.StatusBadRequest, "the para is wrong")
		logger.Warning("instanceId is empty:")
		return
	}
	query := mapRequestForm(c) //
	ret := GetCounterAggregation(query)
	c.JSON(http.StatusOK, ret)
}

func getTransactionByReq(c *gin.Context) {
	if c.Query("appId") == "" || c.Query("namespace") == "" || c.Query("buReqId") == "" {
		c.String(http.StatusBadRequest, "the para is wrong")
		logger.Warning("instanceId is empty:")
		return
	}
	query := mapRequestForm(c) //
	ret := GetTransactionByBuRequest(query)
	c.JSON(http.StatusOK, ret)
}
func getCounterByReq(c *gin.Context) {
	if c.Query("appId") == "" || c.Query("namespace") == "" || c.Query("buReqId") == "" {
		c.String(http.StatusBadRequest, "the para is wrong")
		logger.Warning("instanceId is empty:")
		return
	}
	query := mapRequestForm(c) //
	ret := GetCounterByBuRequest(query)
	c.JSON(http.StatusOK, ret)
}

func getSOEGoupByReq(c *gin.Context) {
	c.String(http.StatusOK, "req getSOEGoupByReq")
}
func getIncompleteSOEByReq(c *gin.Context) {
	c.String(http.StatusOK, "req getIncompleteSOEByReq")
}

func ExportMornitorApi(r *gin.Engine) {
	//global
	r.GET("/monitor/queryAppList", queryAppList)
	//app
	r.GET("/monitor/getAppInfo", getHelp)
	r.GET("/monitor/queryNamespaceList", getHelp)
	r.GET("/monitor/queryLastRequestList", queryLastRequestList)

	//namespace + ts
	r.GET("/monitor/getTransaction", getTransaction)
	r.GET("/monitor/getCounter", getCounter)
	r.GET("/monitor/getStatus", getHelp)

	//request
	r.GET("/monitor/getSOEGoupByReq", getSOEGoupByReq)
	r.GET("/monitor/getIncompleteSOEByReq", getIncompleteSOEByReq)

	r.GET("/monitor/getTransactionByReq", getTransactionByReq)
	r.GET("/monitor/getCounterByReq", getCounterByReq)
	r.GET("/monitor/getStatusByReq", getHelp)
}

func ParserInt(str string) int {
	v, _ := strconv.Atoi(str)
	return v
}

func ParserLong(str string) int64 {
	v, _ := strconv.ParseInt(str, 10, 64)
	return v
}
