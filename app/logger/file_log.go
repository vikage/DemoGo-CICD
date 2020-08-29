package logger

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

// FileLog engine that write log to file
type FileLog struct {
}

// Debug write log
func (filelog FileLog) Debug(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)

	filelog.WriteMessage(message, "log/debug.txt")
}

// Error Write log error type
func (filelog FileLog) Error(format string, v ...interface{}) {
	reformat := fmt.Sprintf("[Error] %s", format)
	message := fmt.Sprintf(reformat, v...)

	filelog.WriteMessage(message, "log/error.txt")
}

// ErrorReq write error contain req info
func (filelog FileLog) ErrorReq(req *http.Request, format string, v ...interface{}) {
	reformat := fmt.Sprintf("[Error] %s %s", req.URL.Path, format)
	message := fmt.Sprintf(reformat, v...)

	filelog.WriteMessage(message, "log/error.txt")
}

// ErrorWithFields write error with fields
func (filelog FileLog) ErrorWithFields(fields map[string]string, format string, v ...interface{}) {
	reformat := fmt.Sprintf("[Error] %s", format)
	message := fmt.Sprintf(reformat, v...)

	for key, value := range fields {
		message = fmt.Sprintf("%s, %s: %s", message, key, value)
	}

	filelog.WriteMessage(message, "log/error.txt")
}

// LogReq log request
func (filelog FileLog) LogReq(req *http.Request, duration int, status int) {
	message := fmt.Sprintf("Request %s %s: %d take %dms", req.Method, req.URL.Path, status, duration)
	WriteContent(message, "log/info.txt")
}

// WriteMessage write message
func (filelog FileLog) WriteMessage(message string, filePath string) {
	currentTime := time.Now()
	time := currentTime.Format("2006-01-02 15:04:05")
	function, _, line, _ := runtime.Caller(3)
	functionCaller := runtime.FuncForPC(function).Name()

	logMessage := fmt.Sprintf("%s[%s:%d] %s\n", time, functionCaller, line, strings.TrimSpace(message))
	WriteContent(logMessage, filePath)
}

// WriteContent append content to file
func WriteContent(content string, filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.Mkdir("log", os.ModePerm)
		os.Create(filePath)
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	if _, err := file.WriteString(content); err != nil {
		fmt.Println("Write log error ", err)
	}

	file.Close()
}
