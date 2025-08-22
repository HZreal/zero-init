package array

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
)

// 判断某一个值是否含在切片之中
func InArray(item interface{}, arr interface{}) bool {
	v := reflect.ValueOf(arr)

	for i := 0; i < v.Len(); i++ {
		if v.Index(i).Interface() == item {
			return true
		}
	}
	return false
}

// 和上述 InArray 功能一样，泛型，更明确类型
func IsInArray[T comparable](item T, arr []T) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}

// ReverseArray 数组逆序
func ReverseArray[T any](arr []T) []T {
	// 创建一个新的切片来存放反转后的数据
	reversed := make([]T, len(arr))

	left, right := 0, len(arr)-1
	for left <= right {
		reversed[left], reversed[right] = arr[right], arr[left]
		left++
		right--
	}

	return reversed
}

// 数组的交集
func Intersect[T comparable](a []T, b []T) []T {
	set := make([]T, 0)

	for _, v := range a {
		if containsGeneric(b, v) {
			set = append(set, v)
		}
	}

	return set
}

func containsGeneric[T comparable](b []T, e T) bool {
	for _, v := range b {
		if v == e {
			return true
		}
	}
	return false
}

// 数组reverse
func Reverse(s interface{}) {
	sort.SliceStable(s, func(i, j int) bool {
		return true
	})
}

// map函数
func MapSlice[T any, M any](a []T, f func(T) M) []M {
	n := make([]M, len(a))
	for i, e := range a {
		n[i] = f(e)
	}
	return n
}

func ArrayFill[T any](n int, value T) []T {
	s := make([]T, n)
	for i := range s {
		s[i] = value
	}
	return s
}

// GroupElements 按照指定的 groupSize 将一个切片分组，返回一个二维切片
func GroupElements[T any](elements []T, groupSize int) [][]T {
	var result [][]T

	// 计算分组的数量
	groupCount := (len(elements) + groupSize - 1) / groupSize

	// 遍历原始切片，按 groupSize 大小进行分组
	for i := 0; i < groupCount; i++ {
		// 计算每组的结束索引，确保不超过原始切片的长度
		start := i * groupSize
		end := start + groupSize
		if end > len(elements) {
			end = len(elements)
		}

		// 将每组的切片添加到结果中
		result = append(result, elements[start:end])
	}

	return result
}

func GetElemsFromIndex[T any](arr []T, index []int) []T {
	result := make([]T, 0, len(arr))
	for _, idx := range index {
		if idx >= 0 && idx < len(arr) {
			result = append(result, arr[idx])
		}
	}

	return result
}

func Filter[T any](items []T, condition func(T) bool) []T {
	var result []T
	for _, item := range items {
		if condition(item) {
			result = append(result, item)
		}
	}
	return result
}

func ExtractField[T any](objects []T, fieldName string) ([]interface{}, error) {
	var result []interface{}
	for _, obj := range objects {
		// 使用反射获取字段值
		v := reflect.ValueOf(obj)
		field := v.FieldByName(fieldName)
		if !field.IsValid() {
			return nil, fmt.Errorf("field %s not found", fieldName)
		}
		result = append(result, field.Interface()) // 将字段值添加到结果切片中
	}
	return result, nil
}

func FindIndex[T comparable](arr []T, target T) (index int, flag bool) {
	for i, item := range arr {
		if item == target {
			flag = true
			index = i
			return
		}
	}
	return
}

func ArrayToMap[T comparable](arr []T) map[T]int {
	result := make(map[T]int)
	for i, v := range arr {
		result[v] = i
	}
	return result
}

func ToAnySlice[T any](s []T) []any {
	result := make([]any, len(s))
	for i, v := range s {
		result[i] = v
	}
	return result
}

// PadToLength
//
//	@Description: 不足长度时，填补零值到指定长度
//	@param s
//	@param length
//	@return []string
func PadToLength[T any](s []T, length int) []T {
	if len(s) >= length {
		return s
	}
	return append(s, make([]T, length-len(s))...)
}

func StrArr2IntArr(arr []string) ([]int64, error) {
	intParts := make([]int64, len(arr))
	for i, s := range arr {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}

		intParts[i] = v
	}
	return intParts, nil

}

func IntArr2StrArr[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](arr []T) []string {
	strParts := make([]string, len(arr))
	for i, v := range arr {
		strParts[i] = fmt.Sprintf("%d", v)
	}
	return strParts
}

// EqualSlices [T comparable]
//
//	@Description: 两个数组等价
//	@param a
//	@param b
//	@return bool
func EqualSlices[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// IsSliceEqual
//
//	@Description: 长度相同、元素内容相同，顺序不一定相同
//	@param a
//	@param b
//	@return bool
func IsSliceEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	countMap := make(map[T]int)

	for _, v := range a {
		countMap[v]++
	}
	for _, v := range b {
		countMap[v]--
		if countMap[v] < 0 {
			return false
		}
	}
	return true
}
