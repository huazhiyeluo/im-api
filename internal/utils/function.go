package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"

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

func ToString(obj interface{}) string {
	switch obj.(type) {
	case float32, float64:
		return fmt.Sprintf("%s", obj)
	case int, int32, int64:
		return fmt.Sprintf("%d", obj)
	default:
		return fmt.Sprintf("%v", obj)
	}
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

func GetRandNum(min uint32, max uint32) uint32 {
	if min == max {
		return min
	}
	seedNum := time.Now().UnixNano()
	rand.Seed(seedNum)
	num := uint32(rand.Intn(int(max-min+1))) + min
	return num
}

func IsContainUint32(str uint32, arr []uint32) bool {
	for _, element := range arr {
		if str == element {
			return true
		}
	}
	return false
}

// 发送GET请求
// url:请求地址
// response:请求返回的内容
func HttpGet(targetUrl string) (response string, err error) {
	// 解析代理URL
	proxyURL := "http://127.0.0.1:8118"
	proxy, err := url.Parse(proxyURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse proxy URL %s: %v", proxyURL, err)
	}
	// 创建带有代理设置的Transport
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}
	client := &http.Client{Timeout: 60 * time.Second, Transport: transport}
	// client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(targetUrl)
	if err != nil {
		return "", fmt.Errorf("failed to get URL %s: %v", targetUrl, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}
	response = string(body)
	return response, nil
}

// CamelToSnakeCase 将驼峰命名转换为下划线命名
func CamelToSnakeCase(input string) string {
	var buf bytes.Buffer
	buf.Grow(len(input) * 2) // 预分配足够的空间，避免动态分配

	for i, r := range input {
		if unicode.IsUpper(r) {
			if i > 0 {
				buf.WriteByte('_')
			}
			buf.WriteRune(unicode.ToLower(r))
		} else {
			buf.WriteRune(r)
		}
	}

	return buf.String()
}

func GetDate(params ...interface{}) string {
	var timestamp int64
	var format string
	if len(params) > 0 {
		timestamp = params[0].(int64)
		if timestamp == 0 {
			timestamp = time.Now().Unix()
		}
	} else {
		timestamp = time.Now().Unix()
	}
	if len(params) > 1 {
		format = params[1].(string)
	} else {
		format = "20060102"
	}
	timeStr := time.Unix(timestamp, 0).Format(format)
	return timeStr
}
