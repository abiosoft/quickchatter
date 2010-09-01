package manager

import (
	"quickchat/database"
)

const (
	MESSAGE = iota
	STILL_ONLINE
	GOING_OFFLINE
	ASYNC_CHAT
	FILE_TRANSFER
)

var Settings *database.Settings

type FromFriend struct {
	Name     string
	Hostname string
}

type Conn struct {
	Type    int
	HashKey string
	Data    interface{}
	Friend  *FromFriend
}
