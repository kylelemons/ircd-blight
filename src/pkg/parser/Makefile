include $(GOROOT)/src/Make.inc

ROOT=../../../

TARG=kevlar/ircd/parser
GOFILES=\
	parser.go\
	commands.go\
	numeric.go\
	util.go\
	rfc2812.go\

# Also remove generated rfc files
CLEANFILES=rfc*.go

rfc2812.go : $(ROOT)/doc/IRC-RFC2812.txt $(ROOT)/doc/IRC-CustomNumerics.txt
	$(ROOT)/scripts/extract_numerics.py $^ > $@
	gofmt -w $@

include $(GOROOT)/src/Make.pkg

