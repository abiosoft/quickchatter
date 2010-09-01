package manager

import (
	"net"
	"quickchat/database"	
)

type Server struct{
	Server *net.Listener
	Friends map[String]*database.Friend
}

var Server *Server

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
	Sever = &Server{
		Server : 
	}
}
