#!/usr/bin/python

import sys
import re

def quot(str):
    return re.sub(r'"', r'\"', str)

def process_rfc(files):
    contents = ""
    slurp = re.compile(r'\s+', re.MULTILINE)
    numeric = re.compile(r'(\d{3})\s+([A-Z]+_[-_A-Z]+)\s+":?(.*?)"$',
        re.MULTILINE | re.DOTALL)

    for file in files:
        with open(file) as f:
            contents += f.read()

    numerics = numeric.findall(contents)
    print "package parser"
    print ""
    print "// Automatically generated from %s" % (file)
    numerics.sort(key=lambda x: x[0])
    print "const ("
    for num in numerics:
        print '\t%s = "%s"' % (num[1], quot(num[0]))
    print ")"
    print ""
    numerics.sort(key=lambda x: x[1])
    print "var NumericName = map[string]string{"
    for num in numerics:
        print '\t%s: "%s",' % (num[1], quot(num[1]))
    print "}"
    print ""
    print "var NumericText = map[string]string{"
    for num in numerics:
        slurped = slurp.sub(' ', num[2])
        print '\t%s: "%s",' % (num[1], quot(slurped))
    print "}"

process_rfc(sys.argv[1:])
