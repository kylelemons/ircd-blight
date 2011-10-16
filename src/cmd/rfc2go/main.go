package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

var (
	pkg = flag.String("pkg", "parser", "The package for the generated source file")
	out = flag.String("out", "-", "File to write the generated source file or - for stdout")
)

func main() {
	flag.Parse()

	numeric2name := make(map[string]string)
	name2text := make(map[string]string)

	numerics := make([]string, 0)
	names := make([]string, 0)

	extract := regexp.MustCompile(`([0-9][0-9][0-9])[ \t]+([A-Z]+_[\-_A-Z]+)` +
		`[ \t\r\n]+":?(([^"]|"[^"\r\n]*")+)"\n`)
	joiner := regexp.MustCompile(`\n[ \t\n]+`)

	var fout io.Writer = os.Stdout
	if *out != "-" {
		fout = bytes.NewBuffer(nil)
	}

	for _, file := range flag.Args() {
		bytes, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatalf("ReadFile(%q): %s", file, err)
		}

		matches := extract.FindAllSubmatch(bytes, -1)
		for _, match := range matches {
			numeric, name, text := string(match[1]), string(match[2]), string(match[3])
			text = joiner.ReplaceAllString(text, " ")

			if _, overwrite := numeric2name[numeric]; overwrite {
				o, n := numeric2name[numeric], name
				log.Printf("Overwriting numeric %s (%s) with %s", o, numeric, n)

				// Remove the old text mapping
				name2text[o] = "", false
			}

			numerics = append(numerics, numeric)
			numeric2name[numeric] = name
			names = append(names, name)
			name2text[name] = text
		}
	}

	if len(numeric2name) == 0 {
		log.Fatalf("No numerics loaded.")
	}

	sort.Strings(numerics)
	sort.Strings(names)

	fmt.Fprintf(fout, "package %s\n\n", *pkg)
	fmt.Fprintf(fout, "// Automatically generated from %s\n", strings.Join(flag.Args(), " "))
	fmt.Fprintf(fout, "const (\n")
	for _, numeric := range numerics {
		fmt.Fprintf(fout, "\t%s = %q\n", numeric2name[numeric], numeric)
	}
	fmt.Fprintf(fout, ")\n\n")
	fmt.Fprintf(fout, "// Automatically generated from %s\n", strings.Join(flag.Args(), " "))
	fmt.Fprintf(fout, "var NumericName = map[string]string{\n")
	for _, name := range names {
		if _, ok := name2text[name]; !ok {
			continue
		}
		fmt.Fprintf(fout, "\t%s: %q,\n", name, name)
	}
	fmt.Fprintf(fout, "}\n\n")
	fmt.Fprintf(fout, "// Automatically generated from %s\n", strings.Join(flag.Args(), " "))
	fmt.Fprintf(fout, "var NumericText = map[string]string{\n")
	for _, name := range names {
		if _, ok := name2text[name]; !ok {
			continue
		}
		fmt.Fprintf(fout, "\t%s: %#q,\n", name, name2text[name])
	}
	fmt.Fprintf(fout, "}\n")

	if buf, ok := fout.(*bytes.Buffer); ok {
		err := ioutil.WriteFile(*out, buf.Bytes(), 0666)
		if err != nil {
			log.Fatalf("WriteFile(%q): %s", *out, err)
		}
	}
}
