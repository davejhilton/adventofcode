package log

import (
	"fmt"
)

type Color string

const Black Color = "\033[1;30m%s\033[0m"
const Red Color = "\033[1;31m%s\033[0m"
const Green Color = "\033[1;32m%s\033[0m"
const Yellow Color = "\033[1;33m%s\033[0m"
const Purple Color = "\033[1;34m%s\033[0m"
const Magenta Color = "\033[1;35m%s\033[0m"
const Teal Color = "\033[1;36m%s\033[0m"
const White Color = "\033[1;37m%s\033[0m"
const Normal Color = "%s"

var enableDebugMessages = false
var enableColors = true

func EnableDebugLogs(enable bool) {
	enableDebugMessages = enable
}

func EnableColors(enable bool) {
	enableColors = enable
}

func Printf(msg string, a ...interface{}) {
	fmt.Printf(msg, a...)
}

func Println(a ...interface{}) {
	fmt.Println(a...)
}

func Debugf(msg string, a ...interface{}) {
	if enableDebugMessages {
		Printf(msg, a...)
	}
}

func Debug(a ...interface{}) {
	if enableDebugMessages {
		Println(a...)
	}
}

func Colorize(v interface{}, color Color, width int) string {
	result := ""
	switch v := v.(type) {
	case string:
		result = v
	case int:
		result = fmt.Sprintf("%d", v)
	default:
		result = fmt.Sprintf("%v", v)
	}

	if width != 0 {
		result = fmt.Sprintf("%*s", width, result)
	}

	if color != "" && enableColors {
		result = fmt.Sprintf(string(color), result)
	}

	return result
}
