include $(GOROOT)/src/Make.inc

DEPS=\
		 sms

TARG=testsms
GOFILES=\
	testsms.go

include $(GOROOT)/src/Make.cmd
