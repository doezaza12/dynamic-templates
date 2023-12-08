package util

import "fmt"

func InterfaceToString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}
