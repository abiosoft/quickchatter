package database

import (
	"os"
	"quickchat/util"
	"path"
	proto "goprotobuf.googlecode.com/hg/proto"
	"io/ioutil"
	"strings"
	"netchan"
)

var (
	USER_ENV_VAR  = util.UserEnvVar()
	ENV_VAR_ERROR = os.NewError("Environment Variable " + USER_ENV_VAR + " not set")
	PROFILE_ERROR = os.NewError("User profile doesn't exist")
)

type Friend struct {
	Name        string
	Hostname    string
	TempHashKey string
	AsyncMsg    string
	Receiver    *netchan.Importer
	SenderPort  int
}

type Settings struct {
	ProfileName string
	Friends     map[string]*Friend
	Password    string
}

func serializeFriends(friends map[string]*Friend) []byte {
	serial := ""
	fSep := "<<<<>>>>"
	vSep := "{{{}}}"
	if friends == nil {
		return make([]byte, 0)
	}
	for _, f := range friends {
		serial += f.Name + vSep + f.Hostname + fSep
	}
	return []byte(serial)
}

func unserializeFriends(data []byte) (friends map[string]*Friend) {
	str := string(data)
	fSep := "<<<<>>>>"
	vSep := "{{{}}}"
	frnds := strings.Split(str, fSep, -1)
	friends = make(map[string]*Friend)
	for _, f := range frnds {
		frnd := strings.Split(f, vSep, -1)
		if len(frnd) != 2 {
			continue
		}
		if frnd[0] == "" || frnd[1] == "" {
			continue
		}
		friends[frnd[0]] = &Friend{Name: frnd[0], Hostname: frnd[1]}
	}
	return
}

func LoadSettings(profileName string) (settings *Settings, err os.Error) {
	root := os.Getenv(USER_ENV_VAR)
	if root == "" {
		err = ENV_VAR_ERROR
		return
	}
	file, err := os.Open(path.Join(root, util.FOLDER_NAME, "settings", profileName+".dat"), os.O_RDONLY, 0666)
	if err != nil {
		err = PROFILE_ERROR
		return
	}
	file.Close()
	data, err := ioutil.ReadFile(file.Name())
	if err != nil {
		return
	}
	var friends = &SettingsWire{}
	err = proto.Unmarshal(data, friends)
	if err != nil {
		return
	}
	return &Settings{
		proto.GetString(friends.ProfileName),
		unserializeFriends(friends.Friends),
		proto.GetString(friends.Password)},
		nil
}

func SaveSettings(settings *Settings) (err os.Error) {
	wire := &SettingsWire{
		ProfileName: proto.String(settings.ProfileName),
		Password:    proto.String(settings.Password),
		Friends:     serializeFriends(settings.Friends)}

	root := os.Getenv(USER_ENV_VAR)
	if root == "" {
		err = ENV_VAR_ERROR
		return
	}
	root = path.Join(root, util.FOLDER_NAME, "settings")
	os.MkdirAll(root, 0755)
	file, err := os.Open(path.Join(root, settings.ProfileName+".dat"), os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		println(err.String())
		err = PROFILE_ERROR
		return
	}
	file.Close()
	data, err := proto.Marshal(wire)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(file.Name(), data, 0666)
	return
}
