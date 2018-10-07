//=============================================================
// log.go
//-------------------------------------------------------------
// Various logging functions
//=============================================================
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
	} else {
		return original[i+1:]
	}
}
func Debug(msg string) {
	print_(msg, "\x1b[44;1m", "DEBUG")
}

func Error(msg string) {
	print_(msg, "\x1b[41;1m", "ERROR")
}

func Warning(msg string) {
	print_(msg, "\x1b[43;1m", "WARNING")
}

func print_(msg, color, type_ string) {
	function, file, line, _ := runtime.Caller(2)
	fmt.Printf("%v[%v][%s:%d][%s][%v]\x1b[0;1m %v\n", color, time.Now().Format("15:04:05.000"), chopPath(file), line, runtime.FuncForPC(function).Name(), type_, msg)
}
