// Copyright 2020-2021 Tetrate
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
/*
  Wasm Virtual Machine
                      (.vm_config.code)
┌────────────────────────────────────────────────────────────────┐
│  Your program (.vm_config.code)                TcpContext      │
│          │                                  ╱ (Tcp stream)     │
│          │ 1: 1                            ╱                   │
│          │         1: N                   ╱ 1: N               │
│      VMContext  ──────────  PluginContext                      │
│                                (Plugin)   ╲ 1: N               │
│                                            ╲                   │
│                                             ╲  HttpContext     │
│                                               (Http stream)    │
└────────────────────────────────────────────────────────────────┘
*/
package main

import (
	"fmt"
	"metrics-xray/collector"
	"metrics-xray/monitor"
	"net/http"
	"strconv"

	"github.com/gohutool/log4go"

	"github.com/gin-gonic/gin"
)

var logger = log4go.LoggerManager.GetLogger("com.hello")

func helloHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hello, metrics-xray!")
}

func webInit() *gin.Engine {
	fmt.Println("Websocket Server!")

	r := gin.Default()

	r.GET("/hello", helloHandler)

	return r

}

func web_proc(r *gin.Engine, port int) {
	r.Run("localhost:" + strconv.Itoa(port))
}

func main() {
	collector.G_manager.Init()

	r := webInit()

	collector.ExportApi(r)
	monitor.ExportMornitorApi(r)

	web_proc(r, 6080)
}
