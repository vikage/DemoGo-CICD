package logger

import (
	"net/http"
)

// LogEngine behavior for log engine
type LogEngine interface {
	Debug(format string, v ...interface{})
	Error(format string, v ...interface{})
	ErrorWithFields(fields map[string]string, format string, v ...interface{})
	ErrorReq(req *http.Request, format string, v ...interface{})
	LogReq(req *http.Request, duration int, status int)
}

// DefaultFileLogger default file logger
var DefaultFileLogger = FileLog{}

// LogEngines store engines for logger
var LogEngines []LogEngine

// AddLogEngine store engine
func AddLogEngine(engine LogEngine) {
	LogEngines = append(LogEngines, engine)
}

// Debug handle write log
func Debug(format string, v ...interface{}) {
	for _, engine := range LogEngines {
		engine.Debug(format, v...)
	}
}

// Error write error log
func Error(format string, v ...interface{}) {
	for _, engine := range LogEngines {
		engine.Error(format, v...)
	}
}

// ErrorWithFields write error with fields
func ErrorWithFields(fields map[string]string, format string, v ...interface{}) {
	for _, engine := range LogEngines {
		engine.ErrorWithFields(fields, format, v...)
	}
}

// ErrorWithType write error with type in fields
func ErrorWithType(errorType string, format string, v ...interface{}) {
	for _, engine := range LogEngines {
		engine.ErrorWithFields(map[string]string{"type": errorType}, format, v...)
	}
}

// ErrorReq write error log when request
func ErrorReq(r *http.Request, format string, v ...interface{}) {
	for _, engine := range LogEngines {
		engine.ErrorReq(r, format, v...)
	}
}

// LogReq log request
func LogReq(req *http.Request, duration int, status int) {
	for _, engine := range LogEngines {
		engine.LogReq(req, duration, status)
	}
}
