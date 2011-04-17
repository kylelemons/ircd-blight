include $(GOROOT)/src/Make.inc

TARG=kevlar/ircd
GOFILES=\
	core.go\
	run.go\

ovrd : sub all

include $(GOROOT)/src/Make.pkg

sub :
	$(MAKE) -f Makefile.sub install
