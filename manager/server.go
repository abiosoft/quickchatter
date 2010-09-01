package manager

import (
	"net"
	"quickchat/database"
	"quickchat/util"
	"log"
	"quickchat/crypt"
	"strings"
	"netchan"
	"strconv"
)

type Server struct {
	Server  net.PacketConn
	Friends map[string]*database.Friend
	Sender  *netchan.Exporter
}

var (
	LocalServer *Server
	Me          *database.Friend
)

func (this *Server) AddFriend(friend *database.Friend) bool {
	_, ok := this.Friends[friend.Name]
	if ok {
		return false
	}
	this.Friends[friend.Name] = friend
	return true
}

func (this *Server) DeleteFriend(friend *database.Friend) {
	this.Friends[friend.Name] = nil
}

func (this *Server) Send(conn Conn, friend *database.Friend) {
	c := make(chan Conn)
	LocalServer.Sender.Export(friend.Hostname, c, netchan.Send)
	c <- conn
}

func Init() {
	l, err := net.ListenPacket("udp", ":"+util.LISTEN_PORT)
	if err != nil {
		log.Exit(err)
	}
	LocalServer = &Server{
		Server:  l,
		Friends: Settings.Friends}
	LocalServer.Sender, err = netchan.NewExporter("tcp", ":0")
	if err != nil {
		log.Exit(err)
	}
	sAdd := LocalServer.Sender.Addr().String()
	sPort, err := strconv.Atoi(sAdd[strings.LastIndex(sAdd, ":"):])
	if strings.LastIndex(sAdd, ":") == -1 {
		log.Exit("address binding failure")
	}
	Me = &database.Friend{
		Name:        Settings.ProfileName,
		TempHashKey: crypt.Md5(crypt.GenerateKey(util.SECRET_KEY_SIZE)),
		SenderPort:  sPort}
}
