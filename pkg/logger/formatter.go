package logger

import (
	"fmt"
	"net/http"
	"time"
)

// code based on gin/logger.go

type color string

const (
	green   = color("\033[97;42m")
	white   = color("\033[90;47m")
	yellow  = color("\033[90;43m")
	red     = color("\033[97;41m")
	blue    = color("\033[97;44m")
	magenta = color("\033[97;45m")
	cyan    = color("\033[97;46m")
	reset   = color("\033[0m")
)

type LogFormatter struct {
	AppName string

	Request *http.Request

	// TimeStamp shows the time after the server returns a response.
	TimeStamp time.Time
	// StatusCode is HTTP response code.
	StatusCode int
	// Latency is how much time the server cost to process a certain request.
	Latency time.Duration
	// ClientIP equals Context's ClientIP method.
	ClientIP string
	// Method is the HTTP method given to the request.
	Method string
	// Path is a path the client requests.
	Path string
	// ErrorMessage is set if error has occurred in processing the request.
	ErrorMessage string
	// isTerm shows whether gin's output descriptor refers to a terminal.
	isTerm bool
	// BodySize is the size of the Response Body
	BodySize int
	// Keys are the keys set on the request's context.
	Keys map[string]any
}

// statusCodeColor is the ANSI color for appropriately logging http status code to a terminal.
func (lf *LogFormatter) statusCodeColor() color {
	code := lf.StatusCode

	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return white
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return yellow
	default:
		return red
	}
}

// methodColor is the ANSI color for appropriately logging http method to a terminal.
func (lf *LogFormatter) methodColor() color {
	method := lf.Method

	switch method {
	case http.MethodGet:
		return blue
	case http.MethodPost:
		return cyan
	case http.MethodPut:
		return yellow
	case http.MethodDelete:
		return red
	case http.MethodPatch:
		return green
	case http.MethodHead:
		return magenta
	case http.MethodOptions:
		return white
	default:
		return reset
	}
}

// resetColor resets all escape attributes.
func (lf *LogFormatter) resetColor() color {
	return reset
}

// IsOutputColor indicates whether can colors be outputted to the log.
func (lf *LogFormatter) isOutputColor() bool {
	return true //consoleColorMode == forceColor || (consoleColorMode == autoColor && lf.isTerm)
}

// Format is the log format function for logger.
func (lf *LogFormatter) Format() string {
	var statusColor, methodColor, resetColor color
	if lf.isOutputColor() {
		statusColor = lf.statusCodeColor()
		methodColor = lf.methodColor()
		resetColor = lf.resetColor()
	}

	if lf.Latency > time.Minute {
		lf.Latency = lf.Latency.Truncate(time.Second)
	}
	return fmt.Sprintf("[%s] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
		lf.AppName,
		lf.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, lf.StatusCode, resetColor,
		lf.Latency,
		lf.ClientIP,
		methodColor, lf.Method, resetColor,
		lf.Path,
		lf.ErrorMessage,
	)
}
