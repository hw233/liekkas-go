package key

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	EtcdSplitByte  = '/'
	RedisSplitByte = ':'
)

func MakeEtcdKey(vals ...interface{}) string {
	return MakeKey(EtcdSplitByte, vals...)
}

func SplitEtcdKey(key string) []string {
	return strings.Split(key, string(EtcdSplitByte))
}

func SubEtcdKey(key string) string {
	keys := SplitEtcdKey(key)
	if len(keys) >= 2 {
		return keys[1]
	}

	return ""
}

func MakeRedisKey(vals ...interface{}) string {
	return MakeKey(RedisSplitByte, vals...)
}

func SplitRedisKey(key string) []string {
	return strings.Split(key, string(RedisSplitByte))
}

func MakeKey(split byte, vals ...interface{}) string {
	key := bytes.Buffer{}
	for i, val := range vals {
		if i >= 1 {
			_ = key.WriteByte(split)
		}
		_, _ = key.WriteString(fmt.Sprint(val))
	}

	return key.String()
}
