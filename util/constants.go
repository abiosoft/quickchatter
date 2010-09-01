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
	LISTEN_PORT     = "8489"
	SECRET_KEY_SIZE = 32
)

const (
	BCAST_INTERVAL = 5e8
)

func UserEnvVar() string {
	if syscall.OS == "windows" {
		return "USERPROFILE"
	}
	return "HOME"
}
