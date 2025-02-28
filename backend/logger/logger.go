package logger

import (
	// "SQL-OJ-BACKEND/conf"
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/RMS_V3/config"
	"github.com/gin-gonic/gin"
)

type logInfo struct {
	filename string
	msg      string
}

var logChan chan logInfo

// Init logger with channel size,
// and start a goroutine to receive and write log
func InitLogger(bufSize int) {
	logChan = make(chan logInfo, bufSize)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println("logger error", err)
			}
		}()
		var l logInfo
		for {
			l = <-logChan
			writeLog(l)
		}
	}()
}

func writeLog(l logInfo) {
	dstFile, err := os.OpenFile(l.filename, os.O_WRONLY|os.O_APPEND, 0666)
	defer dstFile.Close()

	if err != nil && os.IsNotExist(err) {
		dstFile, err = os.Create(l.filename)
	}
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	write := bufio.NewWriter(dstFile)
	write.WriteString(l.msg)
	write.Flush()
}

func DEBUG_LOG(s string, c *gin.Context) {
	runLog(s, "DEBUG", c)
}

func ERROR_LOG(s string, c *gin.Context) {
	runLog(s, "ERROR", c)
}

func runLog(s string, logType string, c *gin.Context) {
	ptrFunc, _, line, ok := runtime.Caller(2)
	if !ok {
		return
	}
	sFuncName := runtime.FuncForPC(ptrFunc).Name()
	tmp := strings.Split(sFuncName, "/")
	tmp = strings.Split(tmp[len(tmp)-1], ".")
	sPackageName := tmp[0]
	sDate := time.Now().Local().Format("2006-01-02")
	var tmpLogInfo logInfo
	tmpLogInfo.filename = config.GetGlobalConfig().LogConfig.LogPath + sDate + ".runlog." + sPackageName

	tmpLogInfo.msg = fmt.Sprintf("%s [%s] %s [%s:%d] %s\n",
		c.Writer.Header().Get("Request-uuid"),
		logType,
		time.Now().Local().Format("2006-01-02 15:04:05"),
		sFuncName,
		line,
		s)
	logChan <- tmpLogInfo
}

func Log2File(filename string, msg string) {
	logChan <- logInfo{filename: filename, msg: msg}
}
