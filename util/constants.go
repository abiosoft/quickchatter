package util

import (
	"syscall"
)

const (
	AUTHOR      = "Abiosoft"
	VERSION     = "0.1"
	FOLDER_NAME = ".quickchat"
)
//PORTS
const (
	LISTEN_PORT = "8489"
)

func UserEnvVar() string {
	if syscall.OS == "windows" {
		return "USERPROFILE"
	}
	return "HOME"
}
