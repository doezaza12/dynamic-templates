package util

import (
	"fmt"
	"strings"

	constants "github.com/doezaza12/dynamic-templates/constant"
)

func InterfaceToString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}

func IsRemoteTemplate(fp string) bool {
	if strings.Contains(fp, ".git") && (strings.HasPrefix(fp, "https://") || strings.HasPrefix(fp, "git@")) {
		return true
	}
	return false
}

func HasRevision(url string) bool {
	if strings.Contains(url, constants.REVISION_PATTERN) {
		return true
	}
	return false
}
