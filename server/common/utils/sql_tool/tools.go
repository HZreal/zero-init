package sql_tool

import (
	"errors"
	"fmt"
	stringTool "overall/common/utils/string"
	"overall/modules/model"
	"reflect"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type WhereItem struct {
	Field     string      // 字段名称
	Opterater string      // 操作符
	Value     any         // 字段的值
	SubWhere  []WhereItem // or 类型的字段
}

const (
	SQLOptE        string = "="            // 等于
	SQLOptGT       string = ">"            // 大于
	SQLOptLT       string = "<"            // 小于
	SQLOptGET      string = ">="           // 大于等于
	SQLOptLET      string = "<="           // 小于等于
	SQLOptNE       string = "!="           // 不等
	SQLOptIn       string = " in "         // in
	SQLOptNIn      string = " not in "     // not in
	SQLOptNotNull  string = " is not null" // 不是空
	SQLOptIsNull   string = " is null"     // 是空
	SQLOptLike     string = " like "       // 模糊匹配
	SQLOptREGEXP   string = " REGEXP "     // 正则匹配
	SQLOptInvalid  string = ""             // 模糊匹配
	SQLOptCombined string = "Combined"
)

func isValidOpterater(opt string) bool {
	switch opt {
	case SQLOptE,
		SQLOptGT,
		SQLOptLT,
		SQLOptGET,
		SQLOptLET,
		SQLOptNE,
		SQLOptNotNull,
		SQLOptIsNull,
		SQLOptIn,
		SQLOptNIn,
		SQLOptLike,
		SQLOptREGEXP,
		SQLOptCombined:
		return true
	}

	return false
}

// GetUpdateSet 生成更新数据集合
//
//	@param dataset *map[string]any key是字段名字，val是字段值
//	@return sql
//	@return vl
func GetUpdateSet(dataset *map[string]any) (sql string, vl []any) {
	var argList []any = make([]any, 0)
	var setList []string = make([]string, 0)

	for k, v := range *dataset {
		// 扩展一下，v可以是[]interface{}
		vType := reflect.TypeOf(v).String()
		if vType == "[]interface {}" {
			val := v.([]interface{})
			if val == nil || len(val) < 2 || val[0] == nil || val[1] == nil {
				// 长度必须2或以上，即用两个数据描述
				continue
			}
			// 第一个值描述需要做什么
			if val[0] == "++" {
				// ++ 表示使用 col=col+n的方式更新字段
				setList = append(setList, k+"="+k+"+"+val[1].(string))
			}
		} else {
			setList = append(setList, k+"=?")
			argList = append(argList, v)
		}
	}

	return strings.Join(setList, ","), argList
}

// BuildUpdateSet
//
//	@Description: 基于结构体中指针字段是否为空判断是否待更新字段
//	@param data
//	@return string
//	@return []interface{}
//	@return error
func BuildUpdateSet(data interface{}) (string, []interface{}, error) {
	// 类型安全检查
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return "", nil, errors.New("data must be non-nil struct pointer")
	}

	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return "", nil, errors.New("data must be struct pointer")
	}

	typ := val.Type()
	var sets []string
	var args []interface{}

	// 遍历结构体字段
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// 跳过非指针或未设置(nil)的字段
		if field.Kind() != reflect.Ptr || field.IsNil() {
			continue
		}

		// 获取db标签，不存在则转为snake_case
		dbTag := fieldType.Tag.Get("db")
		if dbTag == "" {
			dbTag = stringTool.ToSnakeCase(fieldType.Name)
		}

		sets = append(sets, fmt.Sprintf("%s = ?", dbTag))
		args = append(args, field.Elem().Interface()) // 解引用指针获取实际值
	}

	if len(sets) == 0 {
		return "", nil, errors.New("no fields to update")
	}

	return strings.Join(sets, ", "), args, nil
}

// GetInsertSet
//
//	@param tableName
//	@param useReplace
//	@param dataset
//	@return sql
//	@return vl
func GetInsertSet(tableName string, useReplace bool, ignore bool, dataset any) (sql string, vl []any) {
	outField := make([]string, 0)
	outPalce := make([]string, 0)
	outValue := make([]any, 0)

	v := reflect.ValueOf(dataset)
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
			outField = append(outField, fmt.Sprintf("`%s`", fi.Name))
			outPalce = append(outPalce, "?")
			outValue = append(outValue, v.Field(i).Interface())
		default:
			// get tag name with the tag opton, e.g.:
			// `db:"id"`
			// `db:"id,type=char,length=16"`
			// `db:",type=char,length=16"`
			// `db:"-,type=char,length=16"`
			if strings.Contains(tagv, ",") {
				tagv = strings.TrimSpace(strings.Split(tagv, ",")[0])
			}
			if tagv == "-" || tagv == "id" {
				continue
			}
			if len(tagv) == 0 {
				tagv = fi.Name
			}

			outField = append(outField, fmt.Sprintf("`%s`", tagv))
			outPalce = append(outPalce, "?")
			outValue = append(outValue, v.Field(i).Interface())
		}
	}

	sql = strings.Join(outField, ",")
	place := strings.Join(outPalce, ",")

	cmd := "INSERT"
	if useReplace {
		cmd = "REPLACE"
	}

	opt := "IGNORE"
	if !ignore {
		opt = ""
	}

	sql = fmt.Sprintf("%s %s INTO %s (%s)VALUES(%s)", cmd, opt, tableName, sql, place)

	return sql, outValue
}

// GetWhereSet  获取where条件
//
//	@param wheres
//	@param inSubWhere
//	@return sql
//	@return vl
func GetWhereSet(wheres *[]WhereItem, inSubWhere bool) (sql string, vl []any) {
	var argList []any = make([]any, 0)
	var whereList []string = make([]string, 0)

	for _, v := range *wheres {
		// 既没有字段，也没有子查询，不要了
		if v.Field == "" && len(v.SubWhere) == 0 {
			continue
		}

		if v.Opterater != SQLOptInvalid && !isValidOpterater(v.Opterater) {
			continue
		}

		if v.Field != "" {
			rtype := reflect.TypeOf(v.Value).Kind().String()
			if rtype == "string" || strings.Contains(rtype, "int") {
				if v.Opterater == SQLOptNotNull || v.Opterater == SQLOptIsNull {
					whereList = append(whereList, v.Field+v.Opterater)
				} else if v.Opterater == SQLOptCombined {
					// 已经构造好地组合条件
					// 这里不做预处理，预处理由外面组合的时候去做，这里只吧条件拼接，然后值加到预处理参数里面
					whereList = append(whereList, v.Field)
					argList = append(argList, v.Value)
				} else if v.Opterater == SQLOptLike {

					if rtype == "string" {
						whereList = append(whereList, v.Field+v.Opterater+"?")
						argList = append(argList, "%"+v.Value.(string)+"%")
					}

				} else {
					whereList = append(whereList, v.Field+v.Opterater+"?")
					argList = append(argList, v.Value)
				}

			} else if rtype == "array" || rtype == "slice" {
				if v.Opterater != SQLOptIn && v.Opterater != SQLOptNIn {
					continue
				}

				wl := []string{}
				for _, lv := range v.Value.([]any) {
					wl = append(wl, "?")
					argList = append(argList, lv)
				}

				whereList = append(whereList, v.Field+v.Opterater+"("+strings.Join(wl, ",")+")")
			}

			continue
		}

		if len(v.SubWhere) > 0 {
			subSQL, subArgList := GetWhereSet(&v.SubWhere, true)
			if subSQL == "" || len(subArgList) == 0 {
				continue
			}

			whereList = append(whereList, "("+subSQL+")")
			argList = append(argList, subArgList...)
		}
	}

	if inSubWhere {
		return strings.Join(whereList, " OR "), argList
	} else {
		return strings.Join(whereList, " AND "), argList
	}
}

// GetPageInfo  获取分页信息
//
//	@param pagination 分页数据
//	@return offset
//	@return pageSize
func GetPageInfo(pagination model.Pagination) (offset int64, pageSize int64) {
	offset = (pagination.Page - 1) * pagination.PageSize
	if offset < 0 {
		offset = 0
	}

	if pagination.PageSize == 0 {
		pagination.PageSize = 10
	}

	return offset, pagination.PageSize
}

// GetBatchInsertSet
//
//	@param tableName
//	@param useReplace
//	@param dataset
//	@return sql
//	@return vl
func GetBatchInsertSet(tableName string, useReplace bool, ignore bool, dataset []any, update bool) (sql string, vl []any) {
	if len(dataset) == 0 {
		logx.Error("GetBatchInsertSet dataset is empty")
		return "", nil
	}

	outField := make([]string, 0)
	outValue := make([]any, 0)
	updateField := make([]string, 0)

	v := reflect.ValueOf(dataset[0])
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		// panic(fmt.Errorf("ToMap only accepts Slice; got %T", v))
		logx.Errorf("GetBatchInsertSet only accepts Slice; got %T", v)
		return "", nil
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
			outField = append(outField, fmt.Sprintf("`%s`", fi.Name))
		default:
			// get tag name with the tag opton, e.g.:
			// `db:"id"`
			// `db:"id,type=char,length=16"`
			// `db:",type=char,length=16"`
			// `db:"-,type=char,length=16"`
			if strings.Contains(tagv, ",") {
				tagv = strings.TrimSpace(strings.Split(tagv, ",")[0])
			}

			if tagv == "-" || (!update && tagv == "id") {
				continue
			}

			if len(tagv) == 0 {
				tagv = fi.Name
			}

			if update && tagv != "id" {
				updateField = append(updateField, fmt.Sprintf("%s=values(%s)", tagv, tagv))
			}

			outField = append(outField, fmt.Sprintf("`%s`", tagv))
		}
	}

	sql = strings.Join(outField, ",")

	var valueList []string
	for _, value := range dataset {
		objV := reflect.ValueOf(value)
		if objV.Kind() == reflect.Ptr {
			objV = objV.Elem()
		}

		valueString := "("
		for i := 0; i < objV.NumField(); i++ {
			tagV := objV.Type().Field(i).Tag.Get(dbTag)
			// 如果有 id 字段，过掉，会多一个 ?
			if tagV == "-" || (!update && tagV == "id") {
				continue
			}
			if i == objV.NumField()-1 {
				valueString += "?"
			} else {
				valueString += "?,"
			}
			outValue = append(outValue, objV.Field(i).Interface())
		}
		valueString += ")"
		valueList = append(valueList, valueString)
	}
	place := strings.Join(valueList, ",")

	cmd := "INSERT"
	if useReplace {
		cmd = "REPLACE"
	}

	opt := "IGNORE"
	if !ignore {
		opt = ""
	}

	sql = fmt.Sprintf("%s %s INTO %s (%s)VALUES %s", cmd, opt, tableName, sql, place)
	if update {
		updateStr := " on duplicate key update " + strings.Join(updateField, ",")
		sql = sql + updateStr
	}

	return sql, outValue
}

func GetFormatField(objV reflect.Value, index int, t string, sep string) string {
	v := ""
	if t == "string" {
		v += fmt.Sprintf("'%s'%s", objV.Field(index).String(), sep)
	} else if t == "int64" {
		v += fmt.Sprintf("%d%s", objV.Field(index).Int(), sep)
	}

	return v
}
