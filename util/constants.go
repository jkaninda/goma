package util

import (
	"os"
)

var Version string

func VERSION(def string) string {
	build := os.Getenv("VERSION")
	if build == "" {
		return def
	}
	return build
}
func FullVersion() string {
	ver := Version
	if b := VERSION(""); b != "" {
		return b
	}
	return ver
}
