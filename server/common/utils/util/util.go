package util

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"overall/common/utils/encrypt"
	timeTools "overall/common/utils/time"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// 获取程序名称
func GetAppName() string {
	path, _ := os.Executable()
	_, exec := filepath.Split(path)
	// 兼容windows
	sysType := runtime.GOOS
	if sysType == "windows" {
		arr := strings.Split(exec, ".")
		exec = strings.Join(arr[0:len(arr)-1], ".")
	}
	return exec
}

// 获取程序当前的工作路径
func GetAppWd() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}

	if strings.Contains(wd, "/service") {
		index := strings.Index(wd, "/service")
		return wd[0 : index+len("/service")]
	}

	return wd
}

func GetClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if net.ParseIP(ip) != nil {
		return ip
	}

	if ip = r.Header.Get("X-Real-IP"); net.ParseIP(ip) != nil {
		return ip
	}

	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		if net.ParseIP(ip) != nil {
			return ip
		}
	}

	return ""
}

// TypeList2AnyList  将已经确定的数据类型的列表，转化为any类型的列表
//
//	@param typeList
//	@return []any
func TypeList2AnyList[T any](typeList any) []any {
	values := typeList.([]T)
	anyList := []any{}

	for _, v := range values {
		anyList = append(anyList, v)
	}

	return anyList
}

func GetRecordFileType(fileName string) string {
	if strings.Contains(fileName, ".wav") || strings.Contains(fileName, ".mp3") {
		return ""
	}

	// 判断是否带了其他后缀
	ext := filepath.Ext(fileName)
	fmt.Println("ext", ext)
	if ext != "" {
		return ""
	}

	// 切割录音文件路径，根据文件前缀获取对应的后缀是什么
	nameList := strings.Split(fileName, "-")
	fmt.Println("nameList", nameList)
	switch nameList[0] {
	case "M2":
		ext = ".mp3"
		break
	case "W1":
		ext = ".wav"
		break
	case "W2":
		ext = ".wav"
		break
	default:
		ext = ".mp3"
		break
	}

	return ext
}

func CdrFriendlyCodeConversion(code int64) (int64, int64) {
	// 转2进制
	codeBin := strconv.FormatInt(code, 2)
	// 不满32位 2进制，补足
	if len(codeBin) < 32 {
		codeBin = fmt.Sprintf("%0*s", 32, codeBin)
	}

	// 高位友好码
	codeHigh := codeBin[0:16]

	// 低位友好码
	codeLow := codeBin[16:32]

	realCodeHigh, _ := strconv.ParseInt(codeHigh, 2, 64)

	realCodeLow, _ := strconv.ParseInt(codeLow, 2, 64)

	return realCodeHigh, realCodeLow
}

func StructToJson(s interface{}) string {
	jsonData, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

func IsMap(v interface{}) bool {
	return reflect.ValueOf(v).Kind() == reflect.Map
}

func IsSlice(i interface{}) bool {
	kind := reflect.TypeOf(i).Kind()
	return kind == reflect.Slice
}

func GetSeq(n int) string {
	return timeTools.Format(timeTools.TIME_FORMAT_YYYYMMDDHHMMSS, time.Now()) + "_" + encrypt.Krand(n, encrypt.KC_RAND_KIND_LOWER)
}

func GetMilliSecond() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
