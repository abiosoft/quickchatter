package database

import (
	"testing"
	"container/vector"
	"quickchat/crypt"
)

func TestProfile(t *testing.T) {
	var frnds = &vector.Vector{}
	vals := [][]string{[]string{"good", "boy"}, []string{"boy", "good"}}
	frnds.Push(&Friend{"good", "boy", "", ""})
	frnds.Push(&Friend{"boy", "good", "", ""})
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
	for i, v := range *s.Friends {
		f, ok := v.(*Friend)
		if !ok || f.Name != vals[i][0] || f.Hostname != vals[i][1] {
			t.Fail()
		}
	}
}
