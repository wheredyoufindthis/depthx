package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/KyleBanks/depth"
)

const (
	outputPadding    = "  "
	outputPrefix     = "├ "
	outputPrefixLast = "└ "
)

var (
	outputJSON bool
)

func main() {
	var t depth.Tree

	flag.BoolVar(&t.ResolveInternal, "internal", false, "If set, resolves dependencies of internal (stdlib) packages.")
	flag.BoolVar(&t.ResolveTest, "test", false, "If set, resolves dependencies used for testing.")
	flag.IntVar(&t.MaxDepth, "max", 0, "Sets the maximum depth of dependencies to resolve.")
	flag.BoolVar(&outputJSON, "json", false, "If set, outputs the depencies in JSON format.")
	flag.Parse()

	for _, arg := range flag.Args() {
		err := t.Resolve(arg)
		if err != nil {
			log.Fatal(err)
		}

		if outputJSON {
			writePkgJSON(os.Stdout, *t.Root)
			continue
		}

		writePkg(os.Stdout, *t.Root, 0, false)
	}
}

func writePkgJSON(w io.Writer, p depth.Pkg) {
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(p)
}

// writePkg recursively prints a Pkg and its dependencies to the Writer provided.
func writePkg(w io.Writer, p depth.Pkg, indent int, isLast bool) {
	var prefix string
	if indent > 0 {
		prefix = outputPrefix

		if isLast {
			prefix = outputPrefixLast
		}
	}

	out := fmt.Sprintf("%v%v%v\n", strings.Repeat(outputPadding, indent), prefix, p.Name)
	w.Write([]byte(out))

	for idx, d := range p.Deps {
		writePkg(w, d, indent+1, idx == len(p.Deps)-1)
	}
}