package file

/**
 * @Author nico
 * @Date 2024-12-27
 * @File: file.go
 * @Description:
 */

import (
	"os"
	"overall/common/utils/strtool"
)

func CheckFileDir(dir string) error {
	_, err := os.Stat(dir)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		return nil
	}
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	return nil
}

func GetRandName(prefix string, timestamp string) string {
	return prefix + "_" + timestamp + "_" + strtool.Krand(6, strtool.KC_RAND_KIND_LOWER)
}
