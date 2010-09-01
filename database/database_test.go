package database

import (
	"testing"
	"quickchat/crypt"
)

func TestProfile(t *testing.T) {
	var frnds = make(map[string]*Friend)
	vals := make(map[string][]string)
	vals["good"] = []string{"good", "boy"}
	vals["boy"] = []string{"boy", "good"}
	frnds["good"] = &Friend{"good", "boy", "", "", nil, 0}
	frnds["boy"] = &Friend{"boy", "good", "", "", nil, 0}
	set := &Settings{ProfileName: "testprofile", Friends: frnds, Password: crypt.Md5([]byte("store"))}
	err := SaveSettings(set)
	if err != nil {
		println(err.String())
		t.Fail()
	}
	s, err := LoadSettings("testprofile")
	if err != nil {
		println(err.String())
		t.Fail()
	}
	if s.ProfileName != set.ProfileName || s.Password != set.Password {
		t.Fail()
	}
	for i, f := range s.Friends {
		if f.Name != vals[i][0] || f.Hostname != vals[i][1] {
			println(f.Name, f.Hostname, vals[i][0], vals[i][1])
			t.Fail()
		}
	}
}
