package sql_tool

import (
	"fmt"
	"reflect"
	"strings"
)

const dbTag = "db"

// RawFieldNames converts golang struct field into slice string.
func FieldNameMap(in any) map[string]string {
	out := make(map[string]string)
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		panic(fmt.Errorf("ToMap only accepts structs; got %T", v))
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		tagv := fi.Tag.Get(dbTag)
		switch tagv {
		case "-":
			continue
		case "":
			out[fi.Name] = fi.Name
		default:
			// get tag name with the tag opton, e.g.:
			// `db:"id"`
			// `db:"id,type=char,length=16"`
			// `db:",type=char,length=16"`
			// `db:"-,type=char,length=16"`
			if strings.Contains(tagv, ",") {
				tagv = strings.TrimSpace(strings.Split(tagv, ",")[0])
			}
			if tagv == "-" {
				continue
			}
			if len(tagv) == 0 {
				tagv = fi.Name
			}
			out[fi.Name] = tagv
		}
	}
	return out
}

type UpdateParams struct {
	Table       string
	Data        any
	Condition   *WhereArgs
	IgnoreEmpty bool
	NameMap     map[string]string
}

func GetUpdate(params *UpdateParams) (string, []any) {
	condition := params.Condition
	data := params.Data
	table := params.Table
	nameMap := params.NameMap
	if condition == nil {
		condition = &WhereArgs{}
	}
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	typ := v.Type()
	args := make([]any, 0)
	setArr := make([]string, 0)
	for i := 0; i < v.NumField(); i++ {
		fi := typ.Field(i)
		if field, ok := nameMap[fi.Name]; ok {
			vi := v.Field(i).Interface()
			// 不用对主键ID进行更新, 主键放到where条件中
			if field == "id" {
				if len(condition.Where) > 0 {
					condition.Where = condition.Where + " AND id = ? "
				} else {
					condition.Where = " where id = ? "
				}
				condition.Args = append(condition.Args, vi)
				continue
			}
			setArr = append(setArr, fmt.Sprintf("`%s` = ?", field))
			args = append(args, vi)
		}
	}
	sql := fmt.Sprintf("update %s set %s %s", table, strings.Join(setArr, ","), condition.Where)
	return sql, append(args, condition.Args...)
}
