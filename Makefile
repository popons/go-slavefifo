include $(GOROOT)/src/Make.$(GOARCH)

TARG=slavefifo
GOFILES=$(TARG).go	\
		progdata.go

# CGO_LDFLAGS=-lusb

include $(GOROOT)/src/Make.pkg

