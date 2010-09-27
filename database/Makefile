include $(GOROOT)/src/Make.inc

TARG=quickchatter/database
GOFILES=\
	database.go\
	settings.pb.go\

include $(GOROOT)/src/Make.pkg
include $(GOROOT)/src/pkg/goprotobuf.googlecode.com/hg/Make.protobuf

%.pb.go:	%.proto
	protoc --go_out=. $<
