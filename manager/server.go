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
	proto "goprotobuf.googlecode.com/hg/proto"
	"os"
	"time"
	"fmt"
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

func StartBroadcast() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Exit(err)
	}
	sync := &database.Sync{
		Name:	proto.String(Me.Name)
		Hostname:    proto.String(hostname),
		Port:        proto.Int(Me.SenderPort),
		TempHashKey: proto.String(Me.TempHashKey)}
	data, err := proto.Marshal(sync)
	if err != nil {
		log.Exit(err)
	}
	bAddr, err := GetBroadcastAddr()
	if err != nil {
		log.Exit(err)
	}
	add, err := net.ResolveUDPAddr(bAddr + ":" + util.LISTEN_PORT)
	if err != nil {
		log.Exit(err)
	}
	for {
		_, err = LocalServer.Server.WriteTo(data, add)
		if err != nil {
			log.Stderr(err)
		}
		time.Sleep(util.BCAST_INTERVAL)
	}
}

func StartListener(){
	for {
		data := make([]byte, util.BUFFER_SIZE)
		n, _, err := LocalServer.Server.ReadFrom(data)
		if err != nil {
			log.Exit(err)
		}
		sync := &database.Sync{}
		err = proto.Unmarshal(data, sync)
		if err != nil {
			log.Exit(err)
		}
		frnd := &database.Friend{
			Name : proto.GetString(sync.Name)
			Hostname : proto.GetString(sync.Hostname)
			SenderPort : proto.GetInt(sync.Port)
			TempHashKey : proto.GetString(Me.TempHashKey)
		}
		imp, err := netchan.NewImporter("tcp", frnd.Hostname+":"+fmt.Sprint(frnd.SenderPort))
		if err != nil { log.Stderr(err) }
		frnd.Receiver = imp
		if !AddFriend(frnd) {
			log.Stderr("Connection failed")
		}
		time.Sleep(util.BCAST_INTERVAL)
	}
}

func init() {
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
