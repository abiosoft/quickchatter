// Code generated by protoc-gen-go from "settings.proto"
// DO NOT EDIT!

package database

import proto "goprotobuf.googlecode.com/hg/proto"

// Reference proto import to suppress error if it's not otherwise used.
var _ = proto.GetString

type SettingsWire struct {
	ProfileName      *string "PB(bytes,1,req,name=profile_name)"
	Password         *string "PB(bytes,2,req,name=password)"
	Friends          []byte  "PB(bytes,3,req,name=friends)"
	XXX_unrecognized []byte
}

func (this *SettingsWire) Reset() {
	*this = SettingsWire{}
}

func init() {
}