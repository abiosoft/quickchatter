package manager

import (
	"os"
	"exec"
	"io/ioutil"
	"strings"
	"syscall"
	"net"
)

func GetBroadcastAddr() (add string, err os.Error){
	hostname, err := os.Hostname()
	if err != nil {
		return
	}
	if syscall.OS == "windows" {
		//LookupHost(name string) (cname string, addrs []string, err os.Error)
		_, addrs, err := net.LookupHost(hostname)
		if err != nil{ return }
		for _, a := range addrs {
			if len(strings.Split(a, ".", -1)) == 4 {
				add = a
				return
			}
		}
		err = os.NewError("Could not retrieve connection details")
		return
	}
	out, _ := ioutil.TempFile("", "tmp")
	fd := []*os.File{out, out, out}
	path, err := exec.LookPath("arp")
	if err != nil {
		return
	}
	pid, err := os.ForkExec(path, []string{"", "-a"}, os.Envs, "", fd)
	if err != nil {
		return
	}
	os.Wait(pid, 0)
	out.Close()
	f, _ := ioutil.ReadFile(out.Name())
	info := string(f)
	start, end := strings.Index(info, "("), strings.Index(info, ")")
	if start == -1 || end == -1 {
		err = os.NewError("Not connected to any network")
		return
	}
	info = info[start+1 : end]
	ad := strings.Split(info, ".", -1)
	if len(ad) != 4 {
		err = os.NewError("Could not retrieve connection details")
		return
	}
	ad[3] = "255"
	add = strings.Join(ad, ".")
	os.Remove(out.Name())
	return
}
