package TimeTools

import (
	"fmt"
	"time"
)

const (
	// 预设时间格式：2024225180219
	TIME_FORMAT_YYYYMMDDHHMMSS      int = 1
	TIME_FORMAT_YYYYMMDD            int = 2
	TIME_FORMAT_YYYY_MM_DD          int = 3
	TIME_FORMAT_YYYY_MM_DD_HH_MM_SS int = 4
	TIME_FORMAT_YYYYMM                  = 5
)

//time.Now().Format("20060102150405")

// @param lFormat 格式
// @param stTime
// @return string
func Format(lFormat int, stTime time.Time) string {
	switch lFormat {
	case TIME_FORMAT_YYYYMMDDHHMMSS:
		return stTime.Format("20060102150405")
	case TIME_FORMAT_YYYYMMDD:
		return stTime.Format("20060102")
	case TIME_FORMAT_YYYY_MM_DD:
		return stTime.Format("2006_01_02")
	case TIME_FORMAT_YYYY_MM_DD_HH_MM_SS:
		return stTime.Format("2006-01-02 15:04:05")
	case TIME_FORMAT_YYYYMM:
		return stTime.Format("200601")
	default:
		return ""
	}
}

// Sleep 延时lMilliSecond毫秒
//
//	@param lMilliSecond
func Sleep(lMilliSecond int) {
	time.Sleep(time.Duration(lMilliSecond * 1000 * 1000))
}

// StartOfDay  获取当天开始时间
//
//	@param t
//	@return time.Time
func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func SecToHMS(seconds int64) string {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
}
