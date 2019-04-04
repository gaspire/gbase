package util

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func StrRand(length int) string {
	return StringWithCharset(length, charset)
}

func GetUploadPath(mediaType string) string {
	now := time.Now()
	if mediaType == "" {
		mediaType = "images"
	}
	return fmt.Sprintf("/upload/%s/%d%02d/%02d/",
		mediaType,
		now.Year(),
		now.Month(),
		now.Day(),
	)
}

func GetUploadFileName() string {
	return fmt.Sprintf("%d%s",
		time.Now().Unix(),
		StrRand(12),
	)
}

func FileIsExists(file string) bool {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func CreateDir(mediaType, uploadPath string) (err error, dirStr, filePath string) {
	filePath = GetUploadPath(mediaType)
	dir := bytes.NewBufferString(uploadPath)
	dir.WriteString(filePath)
	dirStr = dir.String()
	isExist := FileIsExists(dirStr)
	if !isExist {
		err = os.MkdirAll(dirStr, 0777)
		if err != nil {
			return
		}
	}
	return
}
