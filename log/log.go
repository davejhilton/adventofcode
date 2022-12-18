package log

import (
	"encoding/json"
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

func DebugEnabled() bool {
	return enableDebugMessages
}

func EnableColors(enable bool) {
	enableColors = enable
}

func Print(msg string) {
	fmt.Print(msg)
}

func Printf(msg string, a ...any) {
	fmt.Printf(msg, a...)
}

func Println(a ...any) {
	a = append(a, "") // stops go vet from complaining about trailing \n
	fmt.Println(a...)
}

func PrintColor(v any, c Color) {
	Println(Colorize(v, c, 0))
}

func JSON(v any, pretty bool) {
	var j []byte
	if pretty {
		j, _ = json.MarshalIndent(v, "", "  ")
	} else {
		j, _ = json.Marshal(v)
	}
	Println(j)
}

func Debug(msg string) {
	if enableDebugMessages {
		Print(msg)
	}
}

func Debugf(msg string, a ...any) {
	if enableDebugMessages {
		Printf(msg, a...)
	}
}

func Debugln(a ...any) {
	if enableDebugMessages {
		Println(a...)
	}
}

func DebugJSON(v any, pretty bool) {
	if enableDebugMessages {
		JSON(v, pretty)
	}
}

func DebugColor(v any, c Color) {
	if enableDebugMessages {
		PrintColor(v, c)
	}
}

func Colorize(v any, color Color, width int) string {
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
