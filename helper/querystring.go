package helper

import (
	"net/url"
	"strings"
	"strconv"
	"encoding/json"
)

//
// Get body request keys
//
func GetInputKeys(input []byte) []string {
	var objmap map[string]*json.RawMessage

	json.Unmarshal(input, &objmap)
	keys := make([]string, 0, len(objmap))
	for k := range objmap {
		keys = append(keys, k)
	}

	return keys
}

//
// Formating query string requests
//
func QueryString(qs url.Values) (query map[int]map[string]string, fields []string, groupby []string, sortby []string, order []string,
offset int64, limit int64, join []string) {
	var cq map[string]string = make(map[string]string)
	query = make(map[int]map[string]string)
	limit = 10
	offset = 0

	if param, ok := qs["fields"]; ok {
		fields = strings.Split(param[0], ",")
	}

	if param, ok := qs["join"]; ok {
		k := strings.Replace(param[0], ".", "__", -1)
		join = strings.Split(k, ",")
	}

	if param, ok := qs["groupby"]; ok {
		k := strings.Replace(param[0], ".", "__", -1)
		groupby = strings.Split(k, ",")
	}

	if param, ok := qs["sortby"]; ok {
		k := strings.Replace(param[0], ".", "__", -1)
		sortby = strings.Split(k, ",")
	}

	if param, ok := qs["order"]; ok {
		order = strings.Split(param[0], ",")
	}

	if param, ok := qs["limit"]; ok {
		x, _ := strconv.Atoi(param[0])
		limit = int64(x)
	}

	if param, ok := qs["offset"]; ok {
		x, _ := strconv.Atoi(param[0])
		offset = int64(x)
	}

	if param, ok := qs["query"]; ok {
		var index int = 0
		for _, cond := range strings.Split(param[0], "|") {

			for _, partcond := range strings.Split(cond, ",") {
				kv := strings.Split(partcond, ":")
				if len(kv) > 1 {
					k, val := kv[0], kv[1]
					cq[k] = val
				} else {
					cq[partcond] = "true"
				}
			}

			index = index + 1
			query[index] = cq

			cq = make(map[string]string)
		}
	}

	return query, fields, groupby, sortby, order, offset, limit, join
}

//
// Making query from conditional map
//
func ModelCondition(c map[string]string) (query map[int]map[string]string, fields []string, groupby []string, sortby []string, order []string,
offset int64, limit int64, join []string) {
	query = make(map[int]map[string]string)
	query[1] = c
	return query, fields, groupby, sortby, order, offset, limit, join
}