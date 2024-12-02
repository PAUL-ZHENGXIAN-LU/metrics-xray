package repository

import (
	"encoding/json"
	"errors"
	"metrics-xray/model"
	"strings"
	"time"
)

func Trans2Json(rec *model.TransactionV1) (string, string) {
	// to JSON
	jsonData, err := json.Marshal(rec)
	if err != nil {
		//log.Fatalf("JSON marshaling failed: %s", err)
	}
	return rec.Name + ":ta", string(jsonData)
}

func TransFromJson(jsonData string) (*model.TransactionV1, error) {
	// from JSON
	var rec model.TransactionV1
	err := json.Unmarshal([]byte(jsonData), rec)
	if err != nil {
		//log.Fatalf("JSON marshaling failed: %s", err)
		return nil, err

	}
	return &rec, nil
}

func Mkey2TansactionName(mk string) (string, bool) {
	if name, modified := strings.CutSuffix(mk, ":ta"); modified {
		return name, true
	}
	return "", false
}

func Mkey2CounterName(mk string) (string, bool) {
	if name, modified := strings.CutSuffix(mk, ":ca"); modified {
		return name, true
	}
	return "", false
}

func Counter2Json(rec *model.CounterV1) (string, string) {
	jsonData, err := json.Marshal(rec)
	if err != nil {
		//logger.Warning("JSON marshaling failed: %s", err)
	}
	return rec.Name + ":ca", string(jsonData)
}

func CounterFromJson(jsonData string) (*model.CounterV1, error) {
	// from JSON
	var rec model.CounterV1
	err := json.Unmarshal([]byte(jsonData), rec)
	if err != nil {
		//log.Fatalf("JSON marshaling failed: %s", err)
		return nil, err

	}
	return &rec, nil
}

func setHashMap(hashTable *(map[string]*StoreHashItem), key string, mk string, value string, expire int32) {
	if h, ok := (*hashTable)[key]; ok {
		(*h.table)[mk] = value
	} else {
		h := NewStoreHashItem(expire)
		(*h.table)[mk] = value
		(*hashTable)[key] = h
	}
}

func getHashMap(hashTable *(map[string]*StoreHashItem), key string) (ret map[string]string, ok bool) {
	if h, ok := (*hashTable)[key]; ok {
		t := time.Now().Unix()
		if h.expire <= t {
			return ret, false
		}
		ret = *(h.table)
		return ret, true
	}
	return ret, false
}

func App2Map(appId string, sys string, namespaces []string, tags []string) *map[string]string {

	m := make(map[string]string, 0)
	m["app"] = appId
	m["sys"] = sys
	m["namespaces"] = strings.Join(namespaces, ",")
	m["tags"] = strings.Join(tags, ",")
	return &m
}

func Map2App(hash *map[string]string) (*model.AppStoreInfo, error) {

	if appId, ok := (*hash)["app"]; ok {
		app := model.AppStoreInfo{
			AppId:      appId,
			SysId:      "",
			Namespaces: make(map[string]*[]string),
			FilterTags: make(map[string]*[]string),
		}

		for k, v := range *hash {
			if k == "sys" {
				app.SysId = v
			} else if k == "tags" {
				app.FilterTags = *ParserMapGroup(v)
			} else if k == "namespaces" {
				app.Namespaces = *ParserMapGroup(v)
			}
		}
		return &app, nil

	}
	return nil, errors.New("didn't exist app data")
}

func ParserMapGroup(str string) *map[string](*[]string) {
	m := make(map[string](*[]string))
	ss := strings.Split(str, ",")
	for _, it := range ss {
		ite := strings.Split(it, ":")
		k := strings.TrimSpace(ite[0])
		if k == "" || len(ite) != 2 {
			continue
		}
		v := strings.TrimSpace(ite[1])
		if group, ok := m[k]; ok {
			*group = append(*group, v)
		} else {
			g := []string{v}
			m[k] = &g
		}
	}

	return &m

}
