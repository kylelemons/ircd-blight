include $(GOROOT)/src/Make.inc

TARG=kevlar/ircd
GOFILES=\
	connproxy.go\
	datastore.go\
	datastore_link.go\
	parser.go\
	rfc2812.go\

# Also remove generated rfc files
CLEANFILES=rfc*.go

include $(GOROOT)/src/Make.pkg

rfc2812.go : ../../../doc/IRC-RFC2812.txt ../../../doc/IRC-CustomNumerics.txt
	../../../scripts/extract_numerics.py $^ > $@

