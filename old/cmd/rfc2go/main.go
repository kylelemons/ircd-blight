package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	pkg    = flag.String("pkg", "parser", "The package for the generated source file")
	prefix = flag.String("prefix", "rfc", "Filename prefix (.go will be appended)")
)

func main() {
	flag.Parse()

	filename := *prefix + ".go"

	buf := new(bytes.Buffer)

	rfc := parse()
	generate(buf, rfc)
	write(filename, buf)
}

func generate(f io.Writer, rfc parsedRFC) {
	fmt.Fprintf(f, "// Automatically generated from %s\n", strings.Join(flag.Args(), " "))
	fmt.Fprintf(f, "// DO NOT EDIT\n\n")
	fmt.Fprintf(f, "package %s\n\n", *pkg)
	fmt.Fprintf(f, "// IRC Numerics\n")
	fmt.Fprintf(f, "const (\n")
	for _, numeric := range rfc.numerics {
		fmt.Fprintf(f, "\t%s = %q\n", rfc.numeric2name[numeric], numeric)
	}
	fmt.Fprintf(f, ")\n\n")
	fmt.Fprintf(f, "// NumericName maps IRC numerics to their human-readable names.\n")
	fmt.Fprintf(f, "var NumericName = map[string]string{\n")
	for _, name := range rfc.names {
		if _, ok := rfc.name2text[name]; !ok {
			continue
		}
		fmt.Fprintf(f, "\t%s: %q,\n", name, name)
	}
	fmt.Fprintf(f, "}\n\n")
	fmt.Fprintf(f, "// NumericText maps IRC numerics to their text descriptions.\n")
	fmt.Fprintf(f, "var NumericText = map[string]string{\n")
	for _, name := range rfc.names {
		if _, ok := rfc.name2text[name]; !ok {
			continue
		}
		fmt.Fprintf(f, "\t%s: %#q,\n", name, rfc.name2text[name])
	}
	fmt.Fprintf(f, "}\n")
}

var (
	extract = regexp.MustCompile(`([0-9][0-9][0-9])[ \t]+([A-Z]+_[\-_A-Z]+)[ \t\r\n]+":?(([^"]|"[^"\r\n]*")+)"\n`)
	joiner  = regexp.MustCompile(`\n[ \t\n]+`)
)

type parsedRFC struct {
	names     []string
	name2text map[string]string

	numerics     []string
	numeric2name map[string]string
}

func parse() parsedRFC {
	rfc := parsedRFC{
		name2text:    map[string]string{},
		numeric2name: map[string]string{},
	}

	for _, file := range flag.Args() {
		log.Printf("Parsing RFC: %s", filepath.Base(file))
		bytes, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatalf("ReadFile(%q): %s", file, err)
		}

		matches := extract.FindAllSubmatch(bytes, -1)
		for _, match := range matches {
			numeric, name, text := string(match[1]), string(match[2]), string(match[3])
			text = joiner.ReplaceAllString(text, " ")

			if _, overwrite := rfc.numeric2name[numeric]; overwrite {
				o, n := rfc.numeric2name[numeric], name
				log.Printf("Overwriting numeric %s (%s) with %s", o, numeric, n)

				// Remove the old text mapping
				delete(rfc.name2text, o)
			}

			rfc.numerics = append(rfc.numerics, numeric)
			rfc.numeric2name[numeric] = name
			rfc.names = append(rfc.names, name)
			rfc.name2text[name] = text
		}
	}

	if len(rfc.numeric2name) == 0 {
		log.Fatalf("No numerics loaded.")
	}
	sort.Strings(rfc.numerics)
	sort.Strings(rfc.names)

	return rfc
}

func write(filename string, src *bytes.Buffer) {
	source, err := format.Source(src.Bytes())
	if err != nil {
		log.Fatalf("Failed to format code: %s\n\n%s", err, src)
	}

	log.Printf("Generating %q", filename)
	if err := os.WriteFile(filename, source, 0644); err != nil {
		log.Fatalf("Writing output: %s", err)
	}
}
