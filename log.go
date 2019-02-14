package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

func chopPath(original string) string {
	i := strings.LastIndex(original, "/")
	if i == -1 {
		return original
	}
	return original[i+1:]
}

// Debug prints debug messages
func Debug(vars ...interface{}) {
	printWithColor("\x1b[44;1m", "DEBUG", vars)
}

// Error prints error messages
func Error(vars ...interface{}) {
	printWithColor("\x1b[41;1m", "ERROR", vars)
}

// Warning prints warning messages
func Warning(vars ...interface{}) {
	printWithColor("\x1b[43;1m", "WARNING", vars)
}

func printWithColor(color, ptype string, vars ...interface{}) {
	function, file, line, _ := runtime.Caller(2)

	msg := ""
	for _, v := range vars {
		txt := fmt.Sprintf("%v", v)
		txt = strings.Trim(txt, "[")
		txt = strings.Trim(txt, "]")
		msg = fmt.Sprintf("%v", txt)
	}
	fmt.Printf("%v[%v][%s:%d][%s][%v]\x1b[0;1m %v\n", color, time.Now().Format("15:04:05.000"), chopPath(file), line, runtime.FuncForPC(function).Name(), ptype, msg)
}
