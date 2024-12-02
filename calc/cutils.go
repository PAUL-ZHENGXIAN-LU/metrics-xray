package calc

import "time"

func TestFilterNode(tags map[string]string) string {
	return tags["node"]
}

func TestFilterGroup(tags map[string]string, groupName string) string {
	return tags["node"]
}

func TagGroupKey(tagType string, tagValue string) string {
	return tagType + ":" + tagValue
}

func GetNewStartTime(ts int64) int64 {
	minute := ts / 60000
	hour := minute / 60
	ret := (hour - 1) * 60 * 60000
	return ret
}

func GetCycleKey(ts int64) int64 {
	minute := ts / 60000
	ret := minute * 60000
	return ret
}
func GetNowMs() int64 {
	return time.Now().UnixMilli()
}

/**
func ScanNodes(group []*TransactionGroup) []string {

	var m map[string]int = make(map[string]int)
	for _, g := range group {
		for _, name := range (*g).List {
			if _, ok := m[name]; !ok {
				m[name] = 1
			}
		}
	}
	var rets []string = make([]string, 0)
	rets = append(rets, m...)
	return rets
}
**/
