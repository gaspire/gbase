package base

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Uploader struct {
}

func (me *Uploader) stringWithCharset(length int) string {
	var charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func (me *Uploader) getSavePath(mediaType string) string {
	now := time.Now()
	if mediaType == "" {
		mediaType = "images"
	}
	return fmt.Sprintf("/%s/%d%02d/%02d/", mediaType, now.Year(), now.Month(), now.Day())
}

func (me *Uploader) getSaveName() string {
	return fmt.Sprintf("%d%s", time.Now().Unix(), me.stringWithCharset(12))
}

func (me *Uploader) checkPath(path string) (err error) {
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, 0777)
	}
	return
}

func (me *Uploader) Upload(c *gin.Context, fileKey string) (err error, filePath, filename string) {
	mediaType := c.Param("media_type")
	file, _, err := c.Request.FormFile(fileKey)
	if err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		log.Error(file, err)
		return
	}

	// 文件上传相对路径
	filePath = me.getSavePath(mediaType)
	// 文件上传全路径
	fullPath := bytes.NewBufferString(os.Getenv("UPLOAD_PATH"))
	fullPath.WriteString(filePath)
	// 检查目录
	fullPathStr := fullPath.String()
	err = me.checkPath(fullPathStr)
	if err != nil {
		log.Error(fullPathStr, err)
		return
	}

	// 文件名
	filename = me.getSaveName()
	path := bytes.NewBufferString(fullPathStr)
	path.WriteString(filename)
	// 复制文件
	out, err := os.Create(path.String())
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		log.Error(err)
		return
	}
	return
}
