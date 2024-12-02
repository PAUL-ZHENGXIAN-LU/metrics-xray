package mtype

//import "errors"

// Action represents the action which Wasm contexts expects hosts to take.

type IEnumType interface {
	Value() string
	Code() int32
}

type EnumDefine struct {
	items     map[int32]string
	itemTypes map[string]int32
	unknow    int32
	IEnumType
}

func (t EnumDefine) addType(value string, code int32) {
	if value == "" {
		t.unknow = code
	} else {
		t.items[code] = value
		t.itemTypes[value] = code
	}

}

func (t EnumDefine) getCode(value string) int32 {
	if value == "" {
		return t.unknow
	}
	if code, exists := t.itemTypes[value]; exists {
		return code
	} else {
		return t.unknow
	}
}

func (t EnumDefine) getValue(code int32) string {
	if code == t.unknow {
		return ""
	}
	if v, exists := t.items[code]; exists {
		return v
	} else {
		return ""
	}
}

type FType int32

func ParserFtType(name string) FType {
	return FType(l_FType.getCode(name))
}

func (p FType) Name() string {
	return l_FType.getValue(p.Code())
}
func (p FType) Code() int32 {
	return int32(p)
}

const (
	FT_SOE         FType = 1
	FT_TRANSACTION FType = 2
	FT_COUNTER     FType = 3
	FT_STATUS      FType = 4
	FT_UNKNOE      FType = 0
)

var (
	l_FType = &EnumDefine{
		items: map[int32]string{
			FT_SOE.Code():         "f-e-v1",
			FT_TRANSACTION.Code(): "f-ta-v1",
			FT_COUNTER.Code():     "f-ca-v1",
			FT_STATUS.Code():      "f-sa-v1",
			FT_UNKNOE.Code():      "",
		},
		itemTypes: map[string]int32{
			"f-e-v1":  1,
			"f-ta-v1": 2,
			"f-ca-v1": 3,
			"f-sa-v1": 4,
		},
	}
)

const (
	//ip/pod/ip+processid
	REC_TRANSACTION_STATISTICS string = "ta"
	REC_COUNTER_STATISTICS     string = "ca"
	REC_STATUS                 string = "sa"
	REC_SOE                    string = "soe"
	REC_SOE_GROUP              string = "eg"
)

// log format
const (
	HEADER_LENGTH uint32 = 6
	H_IDX_TS      uint32 = 0
	H_IDX_FORMAT  uint32 = 1
	H_IDX_PERIOD  uint32 = 2
	H_IDX_APP     uint32 = 3
	H_IDX_NS      uint32 = 4
	H_IDX_NAME    uint32 = 5

	H_IDX_TA     uint32 = 6
	H_IDX_TA_TAG uint32 = 7

	H_IDX_SOE_TRACEGROUP uint32 = 6
	H_IDX_SOE_TYPE       uint32 = 7
	H_IDX_SOE_STEP       uint32 = 8
	H_IDX_SOE_STATUS     uint32 = 9
	H_IDX_SOE_INFO       uint32 = 9
	H_IDX_SOE_TAG        uint32 = 10

	H_IDX_COUNTER_SUCC   uint32 = 6
	H_IDX_COUNTER_FAILED uint32 = 7
	H_IDX_COUNTER_TAG    uint32 = 8

	H_IDX_STATUS_VALUE uint32 = 6
	H_IDX_STATUS_TAG   uint32 = 7
)

func (p FType) GetTagsIdx() int32 {
	switch p {
	case FT_SOE:
		return int32(H_IDX_SOE_TAG)
	case FT_TRANSACTION:
		return int32(H_IDX_TA_TAG)
	case FT_COUNTER:
		return int32(H_IDX_COUNTER_TAG)
	case FT_STATUS:
		return int32(H_IDX_STATUS_TAG)
	default:
		return 7
	}
}

func (p FType) GetSectorSize() int32 {
	return p.GetTagsIdx() + 1
}

const (
	//ta:(count=4,totalTime=903,avgTime=225,minTime=172,maxTime=360,bucks=[0,0,0,0,0,3,1,0,0,0,0,0],)
	TA_TAG_COUNT      string = "count" //succ count, the time only calc the succ
	TA_TAG_TOTAL_TIME string = "totalTime"
	TA_TAG_AVG_TIME   string = "avgTime"
	TA_TAG_MIN_TIME   string = "minTime"
	TA_TAG_MAX_TIME   string = "maxTime"
	TA_TAG_FAILED     string = "failed"
	TA_TAG_BUCKS      string = "bucks"
)

//type FilterTagType string
const (
	//ip/pod/ip+processid
	FILTER_TAG_NODE string = "node"
	FILTER_TAG_IDC  string = "idc"
	//service version
	FILTER_TAG_VERSION        string = "ver"
	FILTER_TAG_ORG            string = "org"
	FILTER_TAG_CLIENT_VERSION string = "cver"

	FILTER_TAG_APP string = "app"
)

var (
	BUCKS_DEF_V1 = []int64{0, 1}
)
