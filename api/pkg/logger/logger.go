package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type LogLevel int

const (
	FATAL LogLevel = iota
	ERROR
	WARN
	INFO
	DEBUG
)

var (
	maxLogSize   int64 = 10000000 // 100 kbite
	ginLogsFile  *os.File
	apiLogsFile  *os.File
	ginLogsMutex sync.Mutex
	apiLogsMutex sync.Mutex
	logLevel     LogLevel = DEBUG
)

// Opens log files and set log output. Handle logger graceful shutdown and log files rotation
func init() {
	var ok bool
	var err error

	ginLogsFile, err = os.OpenFile("/logs/gin.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to open gin log file: %v", err)
	}
	apiLogsFile, err = os.OpenFile("/logs/api.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to open api log file: %v", err)
	}

	gin.DefaultWriter = io.MultiWriter(ginLogsFile)
	log.SetOutput(io.MultiWriter(apiLogsFile))

	logLevelStr, ok := os.LookupEnv("LOG_LEVEL")
	if ok {
		switch strings.ToLower(logLevelStr) {
		case "debug":
			logLevel = DEBUG
		case "info":
			logLevel = INFO
		case "warn":
			logLevel = WARN
		case "error":
			logLevel = ERROR
		case "fatal":
			logLevel = FATAL
		default:
			log.Fatalf("Invalid log level: %s", logLevelStr)
		}
	}

	go rotateLogsRoutine()
}

func rotateLogsRoutine() {
	for {
		if err := rotateLogIfNeeded(&ginLogsFile, &ginLogsMutex, "gin"); err != nil {
			log.Printf("Error rotating gin log: %v", err)
		}
		if err := rotateLogIfNeeded(&apiLogsFile, &apiLogsMutex, "api"); err != nil {
			log.Printf("Error rotating api log: %v", err)
		}
		time.Sleep(1 * time.Hour)
	}
}

// Should be called on caller module termination
func Shutdown() {
	ginLogsMutex.Lock()
	apiLogsMutex.Lock()
	defer ginLogsMutex.Unlock()
	defer apiLogsMutex.Unlock()

	if ginLogsFile != nil {
		ginLogsFile.Close()
	}
	if apiLogsFile != nil {
		apiLogsFile.Close()
	}
}

// If size of the logs file, rotates them
func rotateLogIfNeeded(filePtr **os.File, mutex *sync.Mutex, logName string) error {
	mutex.Lock()
	defer mutex.Unlock()

	file := *filePtr
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file info: %w", err)
	}
	if info.Size() <= maxLogSize {
		return nil
	}

	file.Close()
	newName := fmt.Sprintf("%s.%s.old", file.Name(), time.Now().Format("2006-01-02_15:04"))
	if err := os.Rename(file.Name(), newName); err != nil {
		return fmt.Errorf("error renaming log file: %w", err)
	}

	newFile, err := os.OpenFile(file.Name(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("error creating new log file: %w", err)
	}

	*filePtr = newFile
	switch logName {
	case "gin":
		gin.DefaultWriter = io.MultiWriter(newFile)
	case "api":
		log.SetOutput(io.MultiWriter(newFile))
	}

	log.Printf("Rotated %s log file", logName)
	return nil
}

func (l LogLevel) String() string {
	return [...]string{"FATAL", "ERROR", "WARN", "INFO", "DEBUG"}[l]
}

func logWithLevel(level LogLevel, format string, v ...any) {
	_, f, l, _ := runtime.Caller(2)
	fullFuncName := fmt.Sprintf("%s:%d", f, l)
	logMsg := fmt.Sprintf(format, v...)
	log.Printf("[%s]\n%s: %s", level, fullFuncName, logMsg)
}

// Wraps error with runtime.Caller. If err is nil, returns nil
func WrapError(err error) error {
	if err == nil {
		return nil
	}
	_, f, l, _ := runtime.Caller(1)
	return fmt.Errorf("\n%s:%d: %w", f, l, err)
}

func WrapErrMsg(err error, format string, v ...any) error {
	_, f, l, _ := runtime.Caller(1)
	if format != "" {
		logMsg := fmt.Sprintf(format, v...)
		return fmt.Errorf("\n%s:%d: %s: %v", f, l, logMsg, err)
	}
	return fmt.Errorf("\n%s:%d: %v", f, l, err)
}

func WrapMsg(format string, v ...any) error {
	_, f, l, _ := runtime.Caller(1)
	logMsg := fmt.Sprintf(format, v...)
	return fmt.Errorf("\n%s:%d: %s", f, l, logMsg)
}

func Debug(format string, v ...any) {
	if logLevel >= DEBUG {
		logWithLevel(DEBUG, format, v...)
	}
}

func Info(format string, v ...any) {
	if logLevel >= INFO {
		logWithLevel(INFO, format, v...)
	}
}

func Warn(format string, v ...any) {
	if logLevel >= WARN {
		logWithLevel(WARN, format, v...)
	}
}

func Error(format string, v ...any) {
	logWithLevel(ERROR, format, v...)
}

func Fatal(format string, v ...any) {
	logWithLevel(FATAL, format, v...)
	os.Exit(1)
}
