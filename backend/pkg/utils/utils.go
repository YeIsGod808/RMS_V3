package utils

import (
	"bytes"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/RMS_V3/config"
	"github.com/RMS_V3/log"
)

// 随机生成字符
func RandomString() string {
	var letters = []byte("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM")
	result := make([]byte, 16)

	rand.Seed(int64(time.Now().Unix()))
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
func GetImageFile(videoPath string) (string, error) {
	// 分割视频路径，获取视频文件名。
	videoName := filepath.Base(videoPath)

	// 使用 filepath 包安全地移除文件扩展名，并添加 ".jpg" 后缀。
	baseName := strings.TrimSuffix(videoName, filepath.Ext(videoName))
	picName := baseName + ".jpg"

	// 获取全局配置中的图片存储路径。
	picPath := config.GetGlobalConfig().MinioConfig.PicPath

	// 确保图片存储路径存在。
	if err := os.MkdirAll(picPath, os.ModePerm); err != nil {
		log.Errorf("Failed to create directory %s: %v\n", picPath, err)
		return "", err
	}

	// 将图片文件名和路径拼接，生成完整的图片文件路径。
	picFullPath := filepath.Join(picPath, picName)

	// 使用 ffmpeg 命令从视频中提取一帧作为图片。
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-ss", "1", "-f", "image2", "-t", "0.01", "-y", picFullPath)

	// 执行命令并处理错误。
	// Capture both stdout and stderr.
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	// Execute the command.
	err := cmd.Run()
	if err != nil {
		log.Errorf("cmd.Run() failed with %v\n", err)
		log.Errorf("ffmpeg output: %s\n", outBuf.String())
		log.Errorf("ffmpeg error: %s\n", errBuf.String())
		return "", err
	}
	// 返回提取的图片文件路径。
	return picFullPath, nil
}
