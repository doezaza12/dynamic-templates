package util

import (
	"fmt"
	"strings"
)

func InterfaceToString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}

func IsRemoteTemplate(fp string) bool {
	if strings.HasSuffix(fp, ".git") && (strings.HasPrefix(fp, "https://") || strings.HasPrefix(fp, "git@")) {
		return true
	}
	return false
}
