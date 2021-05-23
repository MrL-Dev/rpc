package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var infoLogger *log.Logger

func init() {
	infoLogger = log.New(os.Stdout, "[INFO] ", 0)
}

func Info(v ...interface{}) {
	infoLogger.Println(buildPrefix(v)...)
}

func Infof(format string, v ...interface{}) {
	infoLogger.Printf(format, buildPrefix(v)...)
}

func Fatal(v ...interface{}) {
	log.Fatalln(buildPrefix(v)...)
}

func buildPrefix(v ...interface{}) []interface{} {
	vs := make([]interface{}, 0, len(v)+1)
	_, file, line, ok := runtime.Caller(2)
	if ok {
		vs = append(vs, fmt.Sprintf("%s %s:%d", time.Now().Format("2006/01/02 15:04:05"), filepath.Base(file), line))
	}
	vs = append(vs, v...)
	return vs
}
