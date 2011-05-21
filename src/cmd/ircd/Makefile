include $(GOROOT)/src/Make.inc

PKGDIR=$(GOROOT)/pkg/$(GOOS)_$(GOARCH)

TARG=ircd
GOFILES=\
	main.go\

# Rebuild when library changes
PREREQ=$(PKGDIR)/kevlar/ircd/*.a

include $(GOROOT)/src/Make.cmd
