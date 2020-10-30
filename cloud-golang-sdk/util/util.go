package util

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func (f LogWriter) Write(p []byte) (n int, err error) {
	pc, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "?"
		line = 0
	}

	fn := runtime.FuncForPC(pc)
	var fnName string
	if fn == nil {
		fnName = "?()"
	} else {
		dotName := filepath.Ext(fn.Name())
		fnName = strings.TrimLeft(dotName, ".") + "()"
	}

	//logPath := "centrifysdk.log"
	logPath := ""

	if logPath != "" {
		logf, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer logf.Close()
		log.SetOutput(logf)
	}

	log.Printf("%s:%d %s: %s", filepath.Base(file), line, fnName, p)
	return len(p), nil
}

// LogD logs debug message
var LogD = log.New(LogWriter{}, "[DEBUG] ", 0)

// LogWriter struct
type LogWriter struct{}
