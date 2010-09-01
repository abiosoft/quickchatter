package util

import (
	"syscall"
)

const (
	AUTHOR      = "Abiosoft"
	VERSION     = "0.1"
	FOLDER_NAME = ".quickchat"
)

func UserEnvVar() string {
	if syscall.OS == "windows" {
		return "USERPROFILE"
	}
	return "HOME"
}
