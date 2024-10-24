// Package util /
/*****
@author    Jonas Kaninda
@license   MIT License <https://opensource.org/licenses/MIT>
@Copyright © 2024 Jonas Kaninda
**/
package util

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

// Info message
func Info(msg string, args ...any) {
	log.SetOutput(GetStd(GetStringEnv("GOMA_PROXY_ACCESS_LOG", "/dev/stdout")))
	log.SetOutput(GetStd(GetStringEnv("", "")))
	formattedMessage := fmt.Sprintf(msg, args...)
	if len(args) == 0 {
		log.Printf("INFO: %s\n", msg)
	} else {
		log.Printf("INFO: %s\n", formattedMessage)
	}
}

// Warn Warning message
func Warn(msg string, args ...any) {
	log.SetOutput(GetStd(GetStringEnv("GOMA_PROXY_ACCESS_LOG", "/dev/stdout")))
	formattedMessage := fmt.Sprintf(msg, args...)
	if len(args) == 0 {
		log.Printf("WARN: %s\n", msg)
	} else {
		log.Printf("WARN: %s\n", formattedMessage)
	}
}

// Error error message
func Error(msg string, args ...any) {
	log.SetOutput(GetStd(GetStringEnv("GOMA_PROXY_ERROR_LOG", "/dev/stderr")))
	formattedMessage := fmt.Sprintf(msg, args...)
	if len(args) == 0 {
		log.Printf("ERROR: %s\n", msg)
	} else {
		log.Printf("ERROR: %s\n", formattedMessage)

	}
}
func Fatal(msg string, args ...any) {
	log.SetOutput(GetStd(GetStringEnv("GOMA_PROXY_ERROR_LOG", "/dev/stderr")))
	formattedMessage := fmt.Sprintf(msg, args...)
	if len(args) == 0 {
		log.Printf("ERROR: %s\n", msg)
	} else {
		log.Printf("ERROR: %s\n", formattedMessage)
	}

	os.Exit(1)
}

func Debug(msg string, args ...any) {
	log.SetOutput(GetStd(GetStringEnv("GOMA_PROXY_ACCESS_LOG", "/dev/stdout")))
	formattedMessage := fmt.Sprintf(msg, args...)
	if len(args) == 0 {
		log.Printf("INFO: %s\n", msg)
	} else {
		log.Printf("INFO: %s\n", formattedMessage)
	}
}
func GetStd(out string) *os.File {
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
func GetStringEnv(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}

func GetIntEnv(key string, defaultValue int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue

	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue

	}
	return i

}
