package util

import (
	"fmt"
	"os"
)

const VERSION = "0.0.1"

func BUILD(def string) string {
	build := os.Getenv("VERSION")
	if build == "" {
		return def
	}
	return build
}
func FullVersion() string {
	ver := VERSION
	if b := BUILD(""); b != "" {
		ver = fmt.Sprintf("%s.%s", ver, b)
	}
	return ver
}
