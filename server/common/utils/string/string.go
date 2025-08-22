package stringTool

/**
 * @Author nico
 * @Date 2025-03-27
 * @File: string.go
 * @Description:
 */

import (
	"bytes"
	"unicode"
)

func ToSnakeCase(s string) string {
	var buf bytes.Buffer
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			buf.WriteByte('_')
		}
		buf.WriteRune(unicode.ToLower(r))
	}
	return buf.String()
}

func Get7Suffix(number string) string {
	if len(number) > 7 {
		return number[len(number)-7:]
	}
	return number
}
