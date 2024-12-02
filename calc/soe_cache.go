package calc

var (
	l_soeCache LocalSoeCache = LocalSoeCache{
		AppSOEPool: make(map[string]*AppSOEPool),
	}
)

type LocalSoeCache struct {
	AppSOEPool map[string]*AppSOEPool
	//groupListener model.SoeGroupListener
}

func FindAppSOEPool(appId string, env string) (*AppSOEPool, bool) {
	if a, ok := l_soeCache.AppSOEPool[appId]; ok {
		return a, true
	}
	return nil, false
}

func GetOrNewAppSOEPool(appId string, env string) *AppSOEPool {
	if a, ok := FindAppSOEPool(appId, env); ok {
		return a
	}
	a := NewAppSOEPool(appId)
	l_soeCache.AppSOEPool[appId] = a
	return a
}

type AppSOEPool struct {
	app        string
	reqSOEPool map[string]*ReqSOE
}

func NewAppSOEPool(appId string) *AppSOEPool {
	a := AppSOEPool{
		app:        appId,
		reqSOEPool: make(map[string]*ReqSOE),
	}
	return &a
}

func (pThis *AppSOEPool) FindReqSOEPool(buReqId string) (*ReqSOE, bool) {
	if a, ok := pThis.reqSOEPool[buReqId]; ok {
		return a, true
	}
	return nil, false
}

func (pThis *AppSOEPool) GetOrNewReqSOEPool(buReqId string) *ReqSOE {
	if a, ok := pThis.FindReqSOEPool(buReqId); ok {
		return a
	}
	a := NewReqSOE(buReqId)
	pThis.reqSOEPool[buReqId] = a
	return a
}
