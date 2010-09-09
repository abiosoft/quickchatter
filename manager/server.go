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

var (
	LocalServer *Server
	Me          *database.Friend
)

type Server struct {
	Server        net.PacketConn
	Friends       map[string]*database.Friend
	Sender        *netchan.Exporter
	Messages      map[string]string
	AsyncMessages map[string]string
	ReceivedFiles map[string]string
}

func (this *Server) AddFriend(friend *database.Friend) bool {
	_, ok := this.Friends[friend.Name]
	if ok {
		return false
	}
	this.Friends[friend.Name] = friend
	go ListenToFriend(friend)
	return true
}

func (this *Server) HasFriend(name string) bool {
	_, ok := this.Friends[name]
	return ok
}

func (this *Server) DeleteFriend(name string) {
	this.Friends[name] = nil
}

func (this *Server) Send(conn *Conn, friend *database.Friend) {
	c := make(chan *Conn)
	LocalServer.Sender.Export(friend.Hostname, c, netchan.Send)
	c <- conn
}

func ListenToFriend(frnd *database.Friend) {
	c := make(chan *Conn)
	err := frnd.Receiver.Import(frnd.Name, c, netchan.Recv)
	if err != nil {
		log.Stderr(err)
		LocalServer.DeleteFriend(frnd.Name)
	}
	for {
		if closed(c) {
			break
		}
		conn := <-c
		if conn == nil {
			break
		}
		go ReceiveAndProcess(conn)
	}
}

func ReceiveAndProcess(c *Conn) {
	from := c.Friend
	if !LocalServer.HasFriend(from.Name) {
		log.Stderr("connection refused from", from.Name, "at", from.Hostname)
	}
	frnd := LocalServer.Friends[from.Name]
	if frnd == nil ||  frnd.TempHashKey != c.HashKey{
		log.Stderr("connection refused from", from.Name, "at", from.Hostname)
	}
	switch c.Type {
	case MESSAGE:
		{
			LocalServer.Messages[from.Name] = c.Data
			break
		}
	case ASYNC_MESSAGE:
		{
			LocalServer.AsyncMessages[from.Name] = c.Data
			break
		}
	case GOING_OFFLINE:
		{
			LocalServer.DeleteFriend(from.Name)
			break
		}
	case FILE_TRANSFER:
		{
			LocalServer.ReceivedFiles[from.Name] = c.Data
			break
		}
	}
}

func StartBroadcast() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Exit(err)
	}
	sync := &database.Sync{
		Name:        proto.String(Me.Name),
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

func StartListener() {
	for {
		data := make([]byte, util.BUFFER_SIZE)
		_, _, err := LocalServer.Server.ReadFrom(data)
		if err != nil {
			log.Exit(err)
		}
		sync := &database.Sync{}
		err = proto.Unmarshal(data, sync)
		if err != nil {
			log.Exit(err)
		}
		frnd := &database.Friend{
			Name:        proto.GetString(sync.Name),
			Hostname:    proto.GetString(sync.Hostname),
			SenderPort:  int(proto.GetInt32(sync.Port)),
			TempHashKey: Me.TempHashKey}
		imp, err := netchan.NewImporter("tcp", frnd.Hostname+":"+fmt.Sprint(frnd.SenderPort))
		if err != nil {
			log.Stderr(err)
			continue
		}
		frnd.Receiver = imp
		if !LocalServer.AddFriend(frnd) {
			log.Stderr("Connection failed")
		}
	}
}
//TODO after finishing, change to init
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
	if strings.LastIndex(sAdd, ":") == -1 || err != nil {
		log.Exit("address binding failure")
	}
	Me = &database.Friend{
		Name:        Settings.ProfileName,
		TempHashKey: crypt.Md5(crypt.GenerateKey(util.SECRET_KEY_SIZE)),
		SenderPort:  sPort}
}
