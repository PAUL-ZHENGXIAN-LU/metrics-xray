package collector

import (
	"metrics-xray/calc"
	"metrics-xray/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gohutool/log4go"
)

var logger = log4go.LoggerManager.GetLogger("collecter")

func getIpByApp(c *gin.Context) {
	c.String(http.StatusOK, "getIpByApp")
}

func getAppConnectStatus(c *gin.Context) {
	c.String(http.StatusOK, "getAppConnectStatus")
}
func getConnectSummary(c *gin.Context) {
	c.String(http.StatusOK, "getConnectSummary")
}
func getConnectList(c *gin.Context) {
	c.String(http.StatusOK, "getConnectList")
}

func ExportApi(r *gin.Engine) {
	r.GET("/collector/getIpByApp", getIpByApp)
	r.GET("/collector/getAppConnectStatus", getAppConnectStatus)
	r.GET("/collector/getConnectSummary", getConnectSummary)
	r.GET("/collector/getConnectList", getConnectList)
	r.GET("/collector/ws", wsHandle)
	repository.InitLocalStore()
	calc.InitCalc()
	udpInit()
}
