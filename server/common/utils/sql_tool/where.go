package sql_tool

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type (
	WhereArgs struct {
		Where string
		Args  []any
	}
)

/**
* @param  whereList
* @param  hasWhere  bool  true 会返回 where xxx 这样的字符串，false 返回 xxxx（没有where关键词)
 */
func GetAndWhere(whereList []any, hasWhere bool) (*WhereArgs, error) {
	return parseWhereList(whereList, "AND", hasWhere)
}

/**
  * 解析条件语句
   *
   * @param  $whereList 数组 [
   *          "id=".SP_ID
   *    或者  ["sql ? ? ",['a',1]],
   *    或者  ["sql :a :b",[":a"=>"xxx",":b"=>"bbbbb"]]
   *    或者  [" a.type = ?","xxxid"] //
   * ]
   *
*/

func parseWhereList(whereList []any, delimiter string, hasWhere bool) (*WhereArgs, error) {
	args := make([]any, 0)
	sqlArr := make([]string, 0)
	for _, v := range whereList {
		rtype := reflect.TypeOf(v).Kind().String()
		switch rtype {
		case "string":
			if strings.Count(v.(string), "?") != 0 {
				return nil, errors.New(fmt.Sprintf("%v", v) + " can't include ? when there is no params")
			}
			sqlArr = append(sqlArr, v.(string))
		case "array", "slice":
			iParamsLength := len(v.([]any))
			if iParamsLength == 2 {
				if strings.Count(v.([]any)[0].(string), "?") != len(v.([]any)[1].([]any)) {
					return nil, errors.New(fmt.Sprintf("%v", v) + "  ? num is not equal to params num")
				}
				sqlArr = append(sqlArr, v.([]any)[0].(string))
				args = append(args, (v.([]any)[1].([]any))...)
			} else if iParamsLength == 1 {
				if strings.Count(v.([]any)[0].(string), "?") != 0 {
					return nil, errors.New(fmt.Sprintf("%v", v) + " can't include ? when there is no params")
				}
				sqlArr = append(sqlArr, v.([]any)[0].(string))
			} else {
				return nil, errors.New("GetWhere params format Error: " + fmt.Sprintf("%v", v))
			}
		default:
			return nil, errors.New("GetWhere params type Error:" + rtype + " :" + fmt.Sprintf("%v", v))
		}

	}
	sql := ""
	if len(sqlArr) > 0 {
		sql = strings.Join(sqlArr, " "+delimiter+" ")
		if hasWhere {
			sql = " WHERE " + sql
		}
	}
	return &WhereArgs{
		Where: sql,
		Args:  args,
	}, nil
}

func GetOrWhere(whereList []any, hasWhere bool) (*WhereArgs, error) {
	return parseWhereList(whereList, "OR", hasWhere)
}
