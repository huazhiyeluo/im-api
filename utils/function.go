package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"
)

func StringToUint32(str string) uint32 {
	tempStr, _ := strconv.Atoi(str)
	return uint32(tempStr)
}

func GenMd5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func ToNumber(v interface{}) float64 {
	var i float64
	switch v.(type) {
	case string:
		s := v.(string)
		i, _ = strconv.ParseFloat(s, 64)
	case float64:
		i = v.(float64)
	case int:
		i = float64(v.(int))
	case json.Number:
		i, _ = v.(json.Number).Float64()
	default:
		panic(fmt.Sprintf("ToNumber not support type: %T", v))
	}
	return i
}

func ReverseStringArray(arr []string) []string {
	// 获取数组长度
	length := len(arr)

	// 创建一个新的数组，用于存储倒序的元素
	reversedArr := make([]string, length)

	// 倒序遍历原数组，并将元素存储到新数组中
	for i := 0; i < length; i++ {
		reversedArr[length-i-1] = arr[i]
	}

	return reversedArr
}

// 生成Guid
func GenGUID() string {
	u := uuid.NewV4()
	return strings.Replace(u.String(), "-", "", 4)
}
