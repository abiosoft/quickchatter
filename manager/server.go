package manager

import (
	"net"
	"quickchat/database"	
	"quickchat/util"
	"log"
)

type Server struct{
	Server net.PacketConn
	Friends map[string]*database.Friend
}

var LocalServer *Server

func (this *Server) addFriend(friend *database.Friend){
	_, ok := this.Friends[friend.Name]
	if ok {
		return
	}
	this.Friends[friend.Name] = friend
}

func (this *Server) deleteFriend(friend *database.Friend){
	this.Friends[friend.Name] = nil;
}

func init(){
	l, err := net.ListenPacket("udp", ":"+util.LISTEN_PORT)
	if err != nil{
		log.Exit(err)
	}
	LocalServer = &Server{
		Server : l,
		Friends : make(map[string]*database.Friend) }
}
