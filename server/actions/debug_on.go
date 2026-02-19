//go:build debug

package actions

import "log"

func debugLog(format string, args ...any) {
	log.Printf(format, args...)
}
