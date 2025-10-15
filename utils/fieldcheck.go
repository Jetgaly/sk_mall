package utils

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

func CheckString(originStr string, min int, max int) (string, bool) {
	finalStr := strings.ReplaceAll(originStr, " ", "")
	if utf8.RuneCountInString(finalStr) > max || utf8.RuneCountInString(finalStr) < min {
		return "", false
	}
	return finalStr, true
}

func IsPositiveNumber(s string) (string, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", false
	}

	if s[0] == '-' {
		return "", false
	}

	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return "", false
	}

	if num <= 0 {
		return "", false
	}

	// 直接处理字符串
	return truncateToTwoDecimals(s), true
}

// 字符串层面截断到两位小数
func truncateToTwoDecimals(s string) string {
	if s[0] == '+' {
		s = s[1:]
	}
	parts := strings.Split(s, ".")

	if len(parts) == 1 {
		return parts[0]
	}
	integerPart := parts[0]
	decimalPart := parts[1]
	if len(decimalPart) > 2 {
		decimalPart = decimalPart[:2]
	}
	return integerPart + "." + decimalPart
}
