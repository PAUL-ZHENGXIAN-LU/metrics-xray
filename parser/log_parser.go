package parser

import (
	"bufio"
	"fmt"
	"metrics-xray/model"
	"metrics-xray/model/mtype"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gohutool/log4go"
)

const (
	TAG_NAME_BUREQ      = "bu-reqid"
	TAG_NAME_RPCREQ     = "rpc-reqid"
	TAG_SECTOR_TA       = "ta"
	TAG_SECTOR_TAGS     = "tags"
	TAG_TAGS_NAME_BUCKS = "bucks"
)

var logger = log4go.LoggerManager.GetLogger("parser")

type Parser interface {
	onLine(line string) int32
	onData(rec model.MetricsRecordBase) int32
}

type LogParser struct {
	queueClient model.IRecordQueueClient
	Parser
}

func (pThis *LogParser) onLine(line string) int32 {
	defer func() {
		//logger.Info("parser line is err:", line)
	}()

	fmt.Println("line=" + line)
	fields := strings.Split(line, "\t")
	if len(fields) <= int(mtype.H_IDX_FORMAT) {
		return 0
	}
	fType := mtype.ParserFtType(fields[mtype.H_IDX_FORMAT])
	iTags := fType.GetTagsIdx()
	count := int(iTags) + 1
	if len(fields) < count {
		return 0
	}

	fmt.Println(fields[0])

	switch fType {
	case mtype.FT_SOE:
		r, _ := pThis.parserSOE(fields)
		pThis.queueClient.PostSOE(r)
	case mtype.FT_TRANSACTION:
		r, _ := pThis.parserTransaction(fields)
		pThis.queueClient.PostTransaction(r)
	case mtype.FT_COUNTER:
		r, _ := pThis.parserCounter(fields)
		pThis.queueClient.PostCounter(r)
	case mtype.FT_STATUS:
		break
	default:

	}

	return 1
}

func (pThis *LogParser) onData(rec model.MetricsRecordBase) int32 {

	return 1
}

func (pThis *LogParser) parserRecord(fields []string, rec *model.MetricsRecordBase) bool {
	rec.Ts, _ = ParserTs(fields[mtype.H_IDX_TS])
	rec.TypeFormat = mtype.ParserFtType(fields[mtype.H_IDX_FORMAT])
	rec.App = strings.TrimSpace(fields[mtype.H_IDX_APP])
	rec.Namespace = strings.TrimSpace(fields[mtype.H_IDX_NS])
	rec.Name = strings.TrimSpace(fields[mtype.H_IDX_NAME])
	rec.TypePeriod = strings.TrimSpace(fields[mtype.H_IDX_PERIOD])
	iTags := rec.TypeFormat.GetTagsIdx()

	strTags, modified := strings.CutPrefix(fields[iTags], TAG_SECTOR_TAGS+":")
	if modified {
		rec.Tags, _ = pThis.parserMap(strTags)
	}

	rec.BuReqId = (*rec.Tags)[TAG_NAME_BUREQ]

	return true
}
func (pThis *LogParser) parserMap(str string) (*map[string]string, bool) {
	m := make(map[string]string)

	temp1, modified := strings.CutPrefix(str, "(")
	if !modified {
		return &m, false
	}
	strMap := strings.TrimRight(temp1, ") ")
	sectors := strings.Split(strMap, "=")
	var k string = sectors[0]
	count := len(sectors) - 1
	for i := 0; i < count; i++ {
		if i < count-1 {
			kv := strings.Split(sectors[i+1], ",")
			m[k] = kv[0]
			if len(kv) > 1 {
				k = kv[1]
			} else {
				continue
			}
		} else {
			m[k] = strings.TrimRight(sectors[i+1], ", ")
		}

	}
	return &m, true
}

func (pThis *LogParser) parserTransaction(fields []string) (*model.TransactionV1, bool) {
	rec := model.NewTransactionV1()
	pThis.parserRecord(fields, &rec.MetricsRecordBase)

	strMap, modified := strings.CutPrefix(fields[mtype.H_IDX_TA], TAG_SECTOR_TA+":")
	if modified {
		mta, _ := pThis.parserMap(strMap)
		for k, v := range *mta {
			if k == TAG_TAGS_NAME_BUCKS {
				rec.Bucks = pThis.parserBucks(v)
			} else {
				rec.Tas[k], _ = strconv.ParseInt(v, 10, 64)
			}
		}
	}
	return rec, false
}
func (pThis *LogParser) parserBucks(str string) []int64 {

	s1 := strings.TrimLeft(str, " [")
	s2 := strings.TrimRight(s1, "] ")
	sectors := strings.Split(s2, ",")
	ret := make([]int64, len(sectors))
	for i, it := range sectors {
		ret[i], _ = strconv.ParseInt(it, 10, 64)
	}
	return ret
}

func (pThis *LogParser) parserSOE(fields []string) (*model.SOEV1, bool) {

	rec := &model.SOEV1{}
	pThis.parserRecord(fields, &rec.MetricsRecordBase)

	rec.TraceGroup = strings.TrimSpace(fields[mtype.H_IDX_SOE_TRACEGROUP])
	rec.TypeEvent = strings.TrimSpace(fields[mtype.H_IDX_SOE_TYPE])
	rec.RpcStep = int32(ParserInt(fields[mtype.H_IDX_SOE_STEP]))
	rec.Status = int32(ParserInt(fields[mtype.H_IDX_SOE_STATUS]))
	rec.Info = ParserEncodeString(fields[mtype.H_IDX_SOE_INFO])
	rec.RpcReqId = (*rec.Tags)[TAG_NAME_RPCREQ]

	return rec, false
}

func (pThis *LogParser) parserCounter(fields []string) (*model.CounterV1, bool) {

	rec := &model.CounterV1{}
	pThis.parserRecord(fields, &rec.MetricsRecordBase)
	rec.Count = ParserLong(fields[mtype.H_IDX_COUNTER_SUCC])
	rec.Failed = ParserLong(fields[mtype.H_IDX_COUNTER_FAILED])
	return rec, false
}

var LOC, _ = time.LoadLocation("Local")

// 2024-10-29T02:42:52.644
func ParserTs(str string) (int64, bool) {
	sectors := strings.Split(str, ".")

	layout := "2006-01-02T15:04:05"
	t, err := time.ParseInLocation(layout, sectors[0], LOC)
	if err != nil {
		fmt.Println(err)
		//return
	}
	var v int64 = t.Unix() * 1000

	if len(sectors) >= 2 {
		ms := ParserInt(sectors[1])
		if ms > 0 {
			v = v + int64(ms)
		}
	}
	return v, true
}
func ParserInt(str string) int {
	v, _ := strconv.Atoi(str)
	return v
}

func ParserLong(str string) int64 {
	v, _ := strconv.ParseInt(str, 10, 64)
	return v
}

func ParserEncodeString(str string) string {
	temp1 := strings.TrimLeft(str, "\"")
	temp2 := strings.TrimRight(temp1, "\"")
	ret := strings.ReplaceAll(temp2, "\\\"", "\"")
	return ret
}

type LogSax struct {
	parser   *LogParser
	useFile  bool
	fileName string
}

func NewLogSax(queueClient model.IRecordQueueClient) *LogSax {
	parser := &LogParser{
		queueClient: queueClient,
	}
	sax := &LogSax{
		parser:   parser,
		useFile:  false,
		fileName: "",
	}
	return sax
}

func (pThis *LogSax) UseFile() bool {
	return pThis.useFile
}

func (pThis *LogSax) LoadFile(fileName string) (result int32) {
	result = 0

	// 打开文件
	file, err := os.Open(fileName)
	if err != nil {
		logger.Info("open file fail, err=", err.Error())
		return 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		pThis.parser.onLine(line)

		result++
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return result
}

func (pThis *LogSax) LoadData(data string) (result int32) {
	result = 0
	defer func() {
		logger.Warning("sax loadData exception")
	}()

	r := strings.NewReader(data)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		pThis.parser.onLine(line)

		result++
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return result
}
