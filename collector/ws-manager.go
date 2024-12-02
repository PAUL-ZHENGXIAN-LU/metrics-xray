package collector

import (
	"errors"
	"sync"
)

type ReportApp struct {
	appId   string
	clients map[string]*ClientSession
	mu      sync.RWMutex
}

func NewReportApp(appId string) *ReportApp {
	p := &ReportApp{
		appId:   appId,
		clients: make(map[string]*ClientSession),
		//mu:      sync.RWMutex,
	}
	return p
}

func (pThis *ReportApp) getClient(instanceId string) *ClientSession {
	var client *ClientSession = nil
	pThis.mu.RLock()
	c, exists := pThis.clients[instanceId]
	if exists {
		client = c
	}
	pThis.mu.Unlock()
	return client
}

func (pThis *ReportApp) addClient(client *ClientSession) {
	logger.Info("new client connect, appId=" + client.appId + "instance=" + client.instanceId)
	pThis.mu.Lock()
	oldClient, exists := pThis.clients[client.instanceId]
	pThis.clients[client.instanceId] = client
	pThis.mu.Unlock()
	if exists {
		oldClient.Close(errors.New("have new connect, close the old connect"))
	}
}

func (pThis *ReportApp) removeClient(instanceId string) *ClientSession {
	var client *ClientSession = nil
	pThis.mu.Lock()
	c, exists := pThis.clients[instanceId]
	if exists {
		client = c
		delete(pThis.clients, instanceId)
	}
	pThis.mu.Unlock()
	return client
}

type ClientManager struct {
	apps sync.Map
}

var G_manager = ClientManager{
	//clients: make(map[string]*ClientSession),

}

func (pThis *ClientManager) Init() bool {
	pThis.apps.Store("mymapp", NewReportApp("mymapp"))
	return true
}

func (pThis *ClientManager) getApp(appId string) *ReportApp {
	value, exists := pThis.apps.Load(appId)
	if exists {
		return value.(*ReportApp)
	}
	return nil
}

func (pThis *ClientManager) addApp(app *ReportApp) {
	pThis.apps.Store(app.appId, app)
}

func (pThis *ClientManager) removeApp(appId string) *ReportApp {
	value, exists := pThis.apps.LoadAndDelete(appId)
	if exists {
		return value.(*ReportApp)
	}
	return nil
}

func (pThis *ClientManager) removeClient(appId string, instanceId string) bool {
	value, exists := pThis.apps.Load(appId)
	if exists {
		var app *ReportApp = value.(*ReportApp)
		c := app.removeClient(instanceId)
		if c != nil {
			return true
		}
	}
	return false
}
