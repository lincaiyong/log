package log

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
)

var logPath string

func SetLogPath(logPath string) error {
	logPath = logPath
	logDir := path.Dir(logPath)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err = os.MkdirAll(logDir, 0755)
		if err != nil {
			return err
		}
	}
	file, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	log.SetOutput(file)
	return nil
}

func writeLogM(level, msg string, args ...any) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	msg = strings.TrimRight(msg, "\n")
	lines := strings.Split(msg, "\n")
	for _, line := range lines {
		if level != "" {
			line = fmt.Sprintf("%s %s", level, line)
		}
		fmt.Println(line)
		if logPath != "" {
			log.Println(line)
		}
	}
}

func InfoLog(msg string, args ...any) {
	writeLogM("[INFO ]", msg, args...)
}

func WarnLog(msg string, args ...any) {
	writeLogM("[WARN ]", msg, args...)
}

func ErrorLog(msg string, args ...any) {
	writeLogM("[ERROR]", msg, args...)
}

var fatalError = errors.New("fatal error")

func IsFatalError(err any) bool {
	if e, ok := err.(error); ok {
		return errors.Is(e, fatalError)
	}
	return false
}

func FatalLog(msg string, args ...any) {
	writeLogM("[FATAL]", msg, args...)
	pc := make([]uintptr, 100)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	for {
		frame, more := frames.Next()
		writeLogM("[FATAL]", "%s\n    %s:%d", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
	panic(fatalError)
}
