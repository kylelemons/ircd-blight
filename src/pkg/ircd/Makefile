include $(GOROOT)/src/Make.inc

TARG=kevlar/ircd
GOFILES=\
	core.go\
	run.go\
	rfc2812.go\

# Also remove generated rfc files
CLEANFILES=rfc*.go

ovrd : sub all

include $(GOROOT)/src/Make.pkg

rfc2812.go : ../../../doc/IRC-RFC2812.txt ../../../doc/IRC-CustomNumerics.txt
	../../../scripts/extract_numerics.py $^ > $@

sub :
	$(MAKE) -f Makefile.sub install
