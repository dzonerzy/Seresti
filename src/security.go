package main

import (
	"strings"
)

func EscapeShell(variable string) string {
	forbidden := []string{";", "$", "`", "|", "\r", "\n", "[", "]",
		"(", ")", ".", ">", "<", "/", "\\", "'", "\"", "&", "-", "!"}
	for _, bad := range forbidden {
		variable = strings.Replace(variable, bad, "", -1)
	}
	return variable
}
