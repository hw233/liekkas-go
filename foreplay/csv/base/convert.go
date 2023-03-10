package base

import (
	"strconv"
	"strings"
)

func String2Float(s string) (float32, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		s = "0.0"
	}
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return float32(0), false
	} else {
		return float32(n), true
	}
}

// 字符串转换成数字
func String2Uint64(s string) (uint64, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		s = "0"
	}
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, false
	} else {
		return n, true
	}
}

// adapter csv coder.
func String2Int64(s string) (int64, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		s = "0"
	}
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	} else {
		return int64(n), true
	}
}

func String2Uint32(s string) (uint32, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		s = "0"
	}
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, false
	} else {
		return uint32(n), true
	}
}
func String2Int32(s string) (int32, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		s = "0"
	}
	// 填浮点数的int有时候会包含.0
	s = strings.ReplaceAll(s, ".0", "")
	n, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, false
	} else {
		return int32(n), true
	}
}
func String2Bool(s string) (bool, bool) {
	sret := strings.TrimSpace(s)
	if s == "" {
		s = "false"
	}
	if sret == "true" || sret == "TRUE" || sret == "True" || sret == "1" {
		return true, true
	} else if sret == "false" || sret == "FALSE" || sret == "False" || sret == "0" || sret == "" || sret == " " {
		return false, true
	} else {
		n, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return false, false
		} else {
			if n != 0 {
				return true, true
			} else {
				return false, true
			}
		}
	}
}
