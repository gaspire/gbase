package util

//获取随机数
import (
	"bytes"
	crand "crypto/rand"
	"crypto/sha1"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net/http"
	"path"
	"regexp"
	"sort"
	"strings"
	"time"
)

//GetFileExt 获取文件扩展名
func GetFileExt(disposition string) string {
	reg := regexp.MustCompile(`"(.*)"`)
	disposition = reg.FindString(disposition)
	reg = regexp.MustCompile(`"`)
	disposition = reg.ReplaceAllString(disposition, "")
	return path.Ext(disposition)
}

// ValidEnglishName 检测英文
func ValidEnglishName(name string) bool {
	name = strings.Replace(name, " ", "", -1)
	if m, _ := regexp.MatchString("^[a-zA-Z]*$", name); !m {
		return false
	}
	return true
}

//RandStringBytesMaskImprSrc 获取随机数
func RandStringBytesMaskImprSrc(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

//GetImageWH 获取图片宽高
func GetImageWH(uri string) (width, height int, err error) {
	resp, err := http.Get(uri)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fd := bytes.NewReader(data)
	im, _, err := image.DecodeConfig(fd)
	if err != nil {
		return
	}
	return im.Width, im.Height, nil
}

//GetMobileType 获取手机型号
func GetMobileType(http_user_agent string) map[string]string {
	if http_user_agent == "" {
		return nil
	}
	r := regexp.MustCompile(`\([^\(]+\)`)
	all := r.FindAllString(http_user_agent, -1)
	r = regexp.MustCompile(`(\(|\))`)
	str1 := r.ReplaceAllString(all[0], "")
	//str2 := r.ReplaceAllString(all[1], "")
	tmp := strings.Replace(str1, "Linux; ", "", -1)
	tmp = strings.Replace(tmp, "_", ".", -1)
	tmp = strings.Replace(tmp, " CPU iPhone OS ", "", -1)
	tmp = strings.Replace(tmp, " like Mac OS X", "", -1)
	tmp = strings.Replace(tmp, " Build", "", -1)
	osarr := strings.Split(tmp, ";")

	ios := []string{"iPhone", "iPad", "iPod", "iTouch"}
	//判断是否ios系统
	for _, v := range ios {
		if v == osarr[0] {
			mobile := strings.Split(http_user_agent, "Mobile/")
			mobile = strings.Split(mobile[1], " ")
			return map[string]string{
				"plateform": osarr[0],
				"version":   osarr[0],
				"system":    osarr[1],
			}
		}
	}
	//其他
	mobile := strings.Split(osarr[1], "/")
	if len(mobile) > 1 {
		return map[string]string{
			"plateform": mobile[0],
			"version":   mobile[1],
			"system":    osarr[0],
		}
	}

	if strings.Replace(osarr[0], " ", "", -1) == "U" {
		tmp := strings.Replace(str1, "Linux; U; ", "", -1)
		tmpsplit1 := strings.Split(tmp, " Build/")
		tmpsplit2 := strings.Split(tmpsplit1[0], "; zh-cn; ")
		return map[string]string{
			"plateform": tmpsplit2[1],
			"version":   "",
			"system":    tmpsplit2[0],
		}
	}
	return nil
}

//InArray 查询数据是否在数组中
func InArray(val int, array []int) (exist bool) {
	exist = false
	for _, v := range array {
		if val == v {
			exist = true
			return
		}
	}
	return
}

//Pow 次方
func Pow(x, n int) int {
	ret := 1 // 结果初始为0次方的值，整数0次方为1。如果是矩阵，则为单元矩阵。
	for n != 0 {
		if n%2 != 0 {
			ret = ret * x
		}
		n /= 2
		x = x * x
	}
	return ret
}

//UniqueSlice 数组去重
func UniqueSlice(slice *[]int) []int {
	found := make(map[int]bool)
	result := []int{}
	for _, val := range *slice {
		if _, ok := found[val]; !ok {
			found[val] = true
			result = append(result, val)
		}
	}

	return result
}

//UniqueSliceInt64 数组去重
func UniqueSliceInt64(slice []int64) []int64 {
	found := make(map[int64]bool)
	result := []int64{}
	for _, val := range slice {
		if _, ok := found[val]; !ok {
			found[val] = true
			result = append(result, val)
		}
	}

	return result
}

// Intersect 交集
func Intersect(array1, array2 []int) (intersect []int) {
	sort.Ints(array1)
	sort.Ints(array2)
	i, j := 0, 0
	for i < len(array1) && j < len(array2) {
		if array1[i] == array2[j] {
			intersect = append(intersect, array1[i])
			i++
			j++
		} else if array1[i] < array2[j] {
			i++
		} else {
			j++
		}
	}
	return
}

// LocalTimeToUTC 本地时间转UTC
func LocalTimeToUTC(localTime string, offset int) string {
	passTime, _ := time.Parse("2006-01-02 15:04:05", localTime)
	timeType := passTime.Add(time.Minute * time.Duration(offset))
	return timeType.Format("2006-01-02 15:04:05")
}

// UTCToLocalTime utc转本地时间
func UTCToLocalTime(utcTime time.Time, offset int) time.Time {
	return utcTime.Add(time.Minute * time.Duration(-offset))
}

// TimeOfWeek 当前日期所在一周开始和结束时间
func TimeOfWeek(localTime time.Time) (start string, end string) {
	week := int(localTime.Weekday())
	arr := []int{0, 1, 2, 3, 4, 5, 6}
	start = localTime.AddDate(0, 0, arr[0]-arr[week]).Format("2006-01-02") + " 00:00:00"
	end = localTime.AddDate(0, 0, 6-arr[week]).Format("2006-01-02") + " 23:59:59"
	return
}

//RandInt64 根据范围获取随机数
func RandInt64(min, max int64) int64 {
	maxBigInt := big.NewInt(max)
	i, _ := crand.Int(crand.Reader, maxBigInt)
	if i.Int64() < min {
		RandInt64(min, max)
	}
	return i.Int64()
}

// Krand 随机字符串 0-纯数字 1-小写字母 2-大写字母 3-数字、大小写字母
func Krand(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	isAll := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if isAll { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}

//Sha1 Sha1
func Sha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}
