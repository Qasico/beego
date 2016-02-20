package helper

import (
	"strings"

	"github.com/qasico/beego/orm"
)

//
// Grouping query filter
//
func QueryDetail(oc *orm.Condition, k string, v string, f string) (cond *orm.Condition) {
	if strings.Contains(k, "__in") {
		vArr := strings.Split(v, ".")
		switch f {
		case "or":
			cond = oc.Or(k, vArr)
		case "ornot":
			cond = oc.OrNot(k, vArr)
		case "andnot":
			cond = oc.AndNot(k, vArr)
		default:
			cond = oc.And(k, vArr)
		}
	} else if strings.Contains(k, "__between") {
		vArr := strings.Split(v, ".")
		switch f {
		case "or":
			cond = oc.Or(k, vArr)
		case "ornot":
			cond = oc.OrNot(k, vArr)
		case "andnot":
			cond = oc.AndNot(k, vArr)
		default:
			cond = oc.And(k, vArr)
		}
	} else if strings.Contains(k, "__null") {
		k = strings.Replace(k, "__null", "__isnull", -1)
		switch f {
		case "or":
			cond = oc.Or(k, true)
		case "ornot":
			cond = oc.OrNot(k, true)
		case "andnot":
			cond = oc.AndNot(k, true)
		default:
			cond = oc.And(k, true)
		}
	} else if strings.Contains(k, "__notnull") {
		k = strings.Replace(k, "__notnull", "__isnull", -1)
		switch f {
		case "or":
			cond = oc.Or(k, false)
		case "ornot":
			cond = oc.OrNot(k, false)
		case "andnot":
			cond = oc.AndNot(k, false)
		default:
			cond = oc.And(k, false)
		}
	} else {
		switch f {
		case "or":
			cond = oc.Or(k, v);
		case "ornot":
			cond = oc.OrNot(k, v);
		case "andnot":
			cond = oc.AndNot(k, v);
		default:
			cond = oc.And(k, v)
		}
	}

	return cond
}

//
// Formating query string request for match with orm style
//
func QueryCondition(query map[int]map[string]string) (cond *orm.Condition) {
	cond = orm.NewCondition()
	for _, q := range query {
		condition := orm.NewCondition()
		for k, v := range q {
			if strings.Contains(k, "And.") {
				k = strings.Replace(k, "And.", "", -1)
				k = strings.Replace(k, ".", "__", -1)

				condition = QueryDetail(condition, k, v, "and")
			} else if strings.Contains(k, "Ex.") {
				k = strings.Replace(k, "Ex.", "", -1)
				k = strings.Replace(k, ".", "__", -1)

				condition = QueryDetail(condition, k, v, "andnot")
			} else if strings.Contains(k, "Or.") {
				k = strings.Replace(k, "Or.", "", -1)
				k = strings.Replace(k, ".", "__", -1)

				condition = QueryDetail(condition, k, v, "or")
			} else if strings.Contains(k, "OrNot.") {
				k = strings.Replace(k, "OrNot.", "", -1)
				k = strings.Replace(k, ".", "__", -1)

				condition = QueryDetail(condition, k, v, "ornot")
			} else {
				k = strings.Replace(k, ".", "__", -1)

				condition = QueryDetail(condition, k, v, "and")
			}
		}

		cond = cond.AndCond(condition)
	}

	return cond
}

//
// Get join parameters
//
func QueryJoin(joins []string) (field interface{}) {
	if len(joins) > 0 {
		return joins
	}

	return nil;
}

//
// Check if request need to join
//
func IsJoin(joins []string) bool {
	if (len(joins) > 0) && (joins[0] == "none") {
		return false
	}

	return true
}

//
// Set sorting for orm
// its combine between sortby field and order case
//
func SetSorting(sortby []string, order []string) (sortFields []string) {
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else {
					orderby = v
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) == 1 {
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else {
					orderby = v
				}

				sortFields = append(sortFields, orderby)
			}
		}
	}

	return
}