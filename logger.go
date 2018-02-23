package webAppGo

import (
	"log"
	"os"
	"runtime"
	"strconv"
	"time"
)

// Logger is for custom leveled logging
type Logger struct {
	Level    int
	FilePath string
}

// V is for leveled verbose mode of logging
func (l *Logger) V(level int, msg string, args ...interface{}) {
	switch level {
	case 0:
		if level <= l.Level {
			l.Log("INFO(0)", msg, args...)
		}
	case 1:
		if level <= l.Level {
			l.Log("DEBUG(1)", msg, args...)
		}
	case 2:
		if level <= l.Level {
			l.Log("DEBUG(2)", msg, args...)
		}
	case 3:
		if level <= l.Level {
			l.Log("DEBUG(3)", msg, args...)
		}
	default:
	}
}

// Log function logs the msg with the set prefix
func (l *Logger) Log(prefix, msg string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(2)
	log.SetPrefix(prefix + ": ")
	log.Printf(file+":"+strconv.Itoa(line)+": "+msg, args...)
	return
}

// SetSource configures where to write (Default is stdout)
func (l *Logger) SetSource() {
	if l.FilePath != "" {
		f, err := os.OpenFile(l.FilePath+strconv.FormatInt(time.Now().Unix(), 10)+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		log.SetOutput(f)
	}
}
