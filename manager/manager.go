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
	"runtime"
	"path"
	"fmt"
)

var (
	Settings        *database.Settings
	LocalFileServer *FileServer
)

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

type File struct {
	Path string
	Size int64
}

type FileServer struct {
	Port  int
	Files map[string]*File
}

func (this *FileServer) SendFile(fPath string) (url string, err os.Error) {
	f, err := os.Open(fPath, os.O_RDONLY, 0666)
	if err != nil {
		log.Stderr(err)
		return
	}
	fStat, err := f.Stat()
	length := int64(0)
	if err == nil {
		length = fStat.Size
	}
	f.Close()
	fServer := func(c *http.Conn, req *http.Request) {
		f, ok := LocalFileServer.Files[req.URL.String()]
		if !ok {
			http.Error(c, "File not found", http.StatusNotFound)
			return
		}
		if f.Size > 0 {
			c.SetHeader("Content-Length", fmt.Sprint(f.Size))
		}
		http.ServeFile(c, req, f.Path)
	}
	url = "/" + crypt.GenerateNums(5) + "/" + path.Base(fPath)
	http.HandleFunc(url, fServer)
	LocalFileServer.Files[url] = &File{fPath, length}
	return
}

func StartFileServer() *FileServer {
	n, err := net.ListenPacket("udp", ":0")
	if err != nil {
		log.Exit(err)
	}
	i := strings.LastIndex(n.LocalAddr().String(), ":")
	port := n.LocalAddr().String()[i:]
	p, err := strconv.Atoi(port[1:])
	if err != nil {
		log.Exit(err)
	}
	go start(port)
	return &FileServer{p, make(map[string]*File)}
}

func start(port string) {
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Stderr(err)
	}
}

func init() {
	runtime.MemProfileRate = 0
	runtime.GOMAXPROCS(20)
	LocalFileServer = StartFileServer()
}
