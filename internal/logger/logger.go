package logger

import (
	"fmt"
	"github.com/jkaninda/goma-gateway/util"
	"log"
	"os"
)

type Logger struct {
	msg  string
	args interface{}
}

// Info returns info log
func Info(msg string, args ...interface{}) {
	log.SetOutput(getStd(util.GetStringEnv("GOMA_PROXY_ACCESS_LOG", "/dev/stdout")))
	formattedMessage := fmt.Sprintf(msg, args...)
	if len(args) == 0 {
		log.Printf("INFO: %s\n", msg)
	} else {
		log.Printf("INFO: %s\n", formattedMessage)
	}
}

// Warn returns warning log
func Warn(msg string, args ...interface{}) {
	log.SetOutput(getStd(util.GetStringEnv("GOMA_PROXY_ACCESS_LOG", "/dev/stdout")))
	formattedMessage := fmt.Sprintf(msg, args...)
	if len(args) == 0 {
		log.Printf("WARN: %s\n", msg)
	} else {
		log.Printf("WARN: %s\n", formattedMessage)
	}
}

// Error error message
func Error(msg string, args ...interface{}) {
	log.SetOutput(getStd(util.GetStringEnv("GOMA_PROXY_ERROR_LOG", "/dev/stderr")))
	formattedMessage := fmt.Sprintf(msg, args...)
	if len(args) == 0 {
		log.Printf("ERROR: %s\n", msg)
	} else {
		log.Printf("ERROR: %s\n", formattedMessage)

	}
}
func Fatal(msg string, args ...interface{}) {
	log.SetOutput(getStd(util.GetStringEnv("GOMA_PROXY_ERROR_LOG", "/dev/stderr")))
	formattedMessage := fmt.Sprintf(msg, args...)
	if len(args) == 0 {
		log.Printf("ERROR: %s\n", msg)
	} else {
		log.Printf("ERROR: %s\n", formattedMessage)
	}

	os.Exit(1)
}

func Debug(msg string, args ...interface{}) {
	log.SetOutput(getStd(util.GetStringEnv("GOMA_PROXY_ACCESS_LOG", "/dev/stdout")))
	formattedMessage := fmt.Sprintf(msg, args...)
	if len(args) == 0 {
		log.Printf("INFO: %s\n", msg)
	} else {
		log.Printf("INFO: %s\n", formattedMessage)
	}
}
func getStd(out string) *os.File {
	switch out {
	case "/dev/stdout":
		return os.Stdout
	case "/dev/stderr":
		return os.Stderr
	case "/dev/stdin":
		return os.Stdin
	default:
		return os.Stdout

	}
}
