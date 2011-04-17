include $(GOROOT)/src/Make.inc

TARG=kevlar/ircd/parser
GOFILES=\
	parser.go\
	rfc2812.go\

# Also remove generated rfc files
CLEANFILES=rfc*.go

rfc2812.go : ../../../../doc/IRC-RFC2812.txt ../../../../doc/IRC-CustomNumerics.txt
	../../../../scripts/extract_numerics.py $^ > $@
	gofmt -w $@

include $(GOROOT)/src/Make.pkg

