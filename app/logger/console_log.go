package logger

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"
)

// ConsoleLog type that write log to console
type ConsoleLog struct {
}

// Debug write log
func (console ConsoleLog) Debug(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	console.PrintMessage(message)
}

// Error write error
func (console ConsoleLog) Error(format string, v ...interface{}) {
	reformat := fmt.Sprintf("Error %s", format)
	message := fmt.Sprintf(reformat, v...)
	console.PrintMessage(message)
}

// ErrorWithFields write error with fields
func (console ConsoleLog) ErrorWithFields(fields map[string]string, format string, v ...interface{}) {
	reformat := fmt.Sprintf("Error %s", format)
	message := fmt.Sprintf(reformat, v...)

	for key, value := range fields {
		message = fmt.Sprintf("%s, %s: %s", message, key, value)
	}

	console.PrintMessage(message)
}

// ErrorReq write error contain request info
func (console ConsoleLog) ErrorReq(req *http.Request, format string, v ...interface{}) {
	reformat := fmt.Sprintf("Error %s %s", req.URL.Path, format)
	message := fmt.Sprintf(reformat, v...)

	console.PrintMessage(message)
}

// LogReq log request
func (console ConsoleLog) LogReq(req *http.Request, duration int, status int) {
	message := fmt.Sprintf("Request %s %s: %d take %dms", req.Method, req.URL.Path, status, duration)
	console.PrintMessage(message)
}

// PrintMessage print message
func (console ConsoleLog) PrintMessage(message string) {
	currentTime := time.Now()
	time := currentTime.Format("2006-01-02 15:04:05")
	function, _, line, _ := runtime.Caller(3)
	functionCaller := runtime.FuncForPC(function).Name()

	logMessage := fmt.Sprintf("%s[%s:%d] %s\n", time, functionCaller, line, strings.TrimSpace(message))
	fmt.Print(logMessage)
}
