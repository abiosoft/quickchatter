package manager

import (
	"quickchat/database"
	"net"
	"http"
	"os"
	"quickchat/crypt"
	"log"
	"strconv"
	"strings"
)

var Settings *database.Settings

type FromFriend struct {
	Name     string
	Hostname string
}

func FromMe() *FromFriend {
	return &FromFriend{Me.Name, Me.Hostname}
}

type Conn struct {
	Type    int
	HashKey string
	Data    string
	Friend  *FromFriend
}

type FileServer struct {
	Port  int
	Files map[string]string
}

func (this *FileServer) SendFile(path string) {
	f, err := os.Open(path, os.O_RDONLY, 0666)
	if err != nil {
		log.Stderr(err)
		return
	}
	f.Close()
	fServer := func(c *http.Conn, req *http.Request) {
		http.ServeFile(c, req, req.URL.String()[1:])
	}
	url := "/" + crypt.GenerateNums(5)
	http.HandleFunc(url, fServer)
	this.Files[url] = path
}


func StartFileServer() *FileServer {
	n, err := net.ListenPacket("udp", ":0")
	if err != nil {
		log.Exit(err)
	}
	i := strings.LastIndex(n.LocalAddr().String(), ":")
	port := n.LocalAddr().String()[i:]
	p, err := strconv.Atoi(port)
	if err != nil {
		log.Exit(err)
	}
	go start(port)
	return &FileServer{p, make(map[string]string)}
}

func start(port string) {
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Exit(err)
	}
}
