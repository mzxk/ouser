package ouser

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

//主要的加密手段，这里使用sha512模仿sha1的结果
func sha(s string) string {
	if s == "" {
		return s
	}
	k := sha512.Sum512([]byte(s))
	return hex.EncodeToString(k[:])[58:98]
}

//这里名字写的不大对，暂时不改动
func getType(s ...string) string {
	return strings.Join(s, ".")
}

//生成limit key之类的用的
func joinSmsType(s string, i int) string {
	return s + "." + fmt.Sprint(i)
}

//字符串转换成int64，如果为空或者任何意外，返回0
func s2i(s string) int64 {
	if s == "" {
		return 0
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		i = 0
	}
	return i
}
