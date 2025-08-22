package validate

import (
	"math/rand"
	"net"
	"regexp"
	"strings"
	"time"
)

func IsNumeric(str string) bool {
	for _, ch := range str {
		if !strings.ContainsRune("0123456789", ch) {
			return false
		}
	}
	return true
}

// 判断字符串是否为字母数字
func IsAlphaNumeric(str string) bool {
	// 遍历字符串，判断每个字符是否为字母数字
	for _, ch := range str {
		if !strings.ContainsRune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", ch) {
			return false
		}
	}
	return true
}

func RandomString(size int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < size; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func IsIp(ip string) bool {
	if ip == "" {
		return false
	}

	if address := net.ParseIP(ip); address == nil {
		return false
	}
	return true
}

func IsDomain(domain string) bool {
	if domain == "" {
		return false
	}

	reg := regexp.MustCompile("[\\w.-]+(?:\\.[\\w\\.-]+)+[\\w\\-\\._~:\\/?#[\\]@!\\$&'\\(\\)\\*\\+,;=.]+$")
	return reg.MatchString(domain)
}

func IsDomainOrIp(str string) bool {
	if !IsDomain(str) {
		return false
	}
	ipStr := strings.ReplaceAll(str, ".", "")
	if IsNumeric(ipStr) {
		if !IsIp(str) {
			return false
		}
	}
	return true
}

func CheckEmail(email string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return regex.MatchString(email)
}

func CheckPhone(phone string) bool {
	reg := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return reg.MatchString(phone)
}

func CheckLettersNumbers(s string) bool {
	reg := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return reg.MatchString(s)
}
