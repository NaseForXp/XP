package debug

import (
	"fmt"
	"strings"
)
import "os"

func Println(a ...interface{}) (n int, err error) {
	if len(os.Args) >= 2 {
		if strings.ToLower(os.Args[1]) == "debug" {
			return fmt.Println(a...)
		}
	}
	return n, nil
}

func Printf(format string, a ...interface{}) (n int, err error) {
	if len(os.Args) >= 2 {
		if strings.ToLower(os.Args[1]) == "debug" {
			return fmt.Printf(format, a...)
		}
	}
	return n, nil
}
