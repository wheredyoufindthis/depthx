package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/wheredyoufindthis/depthx"
)

const (
	outputClosedPadding = "  "
	outputOpenPadding   = "│ "
	outputPrefix        = "├ "
	outputPrefixLast    = "└ "
)

var outputJSON bool
var explainPkg string

type summary struct {
	numInternal int
	numExternal int
	numTesting  int
}

var cpuprofile string

func main() {
	t, pkgs := parse(os.Args[1:])

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	if err := handlePkgs(t, pkgs, outputJSON, explainPkg); err != nil {
		os.Exit(1)
	}
}

// parse constructs a depth.Tree from command-line arguments, and returns the
// remaining user-supplied package names
func parse(args []string) (*depthx.Tree, []string) {
	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	var t depthx.Tree
	f.StringVar(&cpuprofile, "cpuprofile", "", "write cpu profile to `file`")
	f.BoolVar(&t.ResolveInternal, "internal", false, "If set, resolves dependencies of internal (stdlib) packages.")
	f.BoolVar(&t.ResolveTest, "test", false, "If set, resolves dependencies used for testing.")
	f.IntVar(&t.MaxDepth, "max", 0, "Sets the maximum depth of dependencies to resolve.")
	f.StringVar(&t.Domain, "domain", "", "Restrict to the certain code domain.")
	f.BoolVar(&outputJSON, "json", false, "If set, outputs the depencies in JSON format.")
	f.StringVar(&explainPkg, "explain", "", "If set, show which packages import the specified target")
	f.Parse(args)

	return &t, f.Args()
}

// handlePkgs takes a slice of package names, resolves a Tree on them,
// and outputs each Tree to Stdout.
func handlePkgs(t *depthx.Tree, pkgs []string, outputJSON bool, explainPkg string) error {
	for _, pkg := range pkgs {

		err := t.Resolve(pkg)
		if err != nil {
			fmt.Printf("'%v': FATAL: %v\n", pkg, err)
			return err
		}

		if outputJSON {
			writePkgJSON(os.Stdout, *t.Root)
			continue
		}

		if explainPkg != "" {
			writeExplain(os.Stdout, *t.Root, []string{}, explainPkg)
			continue
		}

		writePkg(os.Stdout, *t.Root)
		writePkgSummary(os.Stdout, *t.Root)
	}
	return nil
}

// writePkgSummary writes a summary of all packages in a tree
func writePkgSummary(w io.Writer, pkg depthx.Pkg) {
	var sum summary
	set := make(map[string]struct{})
	for _, p := range pkg.Deps {
		collectSummary(&sum, p, set)
	}
	fmt.Fprintf(w, "%d dependencies (%d internal, %d external, %d testing).\n",
		sum.numInternal+sum.numExternal,
		sum.numInternal,
		sum.numExternal,
		sum.numTesting)
}

func collectSummary(sum *summary, pkg depthx.Pkg, nameSet map[string]struct{}) {
	if _, ok := nameSet[pkg.Name]; !ok {
		nameSet[pkg.Name] = struct{}{}
		if pkg.Internal {
			sum.numInternal++
		} else {
			sum.numExternal++
		}
		if pkg.Test {
			sum.numTesting++
		}
		for _, p := range pkg.Deps {
			collectSummary(sum, p, nameSet)
		}
	}
}

// writePkgJSON writes the full Pkg as JSON to the provided Writer.
func writePkgJSON(w io.Writer, p depthx.Pkg) {
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(p)
}

func writePkg(w io.Writer, p depthx.Pkg) {
	fmt.Fprintf(w, "%s\n", p.String())

	for idx, d := range p.Deps {
		writePkgRec(w, d, []bool{true}, idx == len(p.Deps)-1)
	}
}

// writePkg recursively prints a Pkg and its dependencies to the Writer provided.
func writePkgRec(w io.Writer, p depthx.Pkg, closed []bool, isLast bool) {
	var prefix string

	for _, c := range closed {
		if c {
			prefix += outputClosedPadding
			continue
		}

		prefix += outputOpenPadding
	}

	closed = append(closed, false)
	if isLast {
		prefix += outputPrefixLast
		closed[len(closed)-1] = true
	} else {
		prefix += outputPrefix
	}

	fmt.Fprintf(w, "%v%v\n", prefix, p.String())

	for idx, d := range p.Deps {
		writePkgRec(w, d, closed, idx == len(p.Deps)-1)
	}
}

// writeExplain shows possible paths for a given package.
func writeExplain(w io.Writer, pkg depthx.Pkg, stack []string, explain string) {
	stack = append(stack, pkg.Name)
	if pkg.Name == explain {
		fmt.Fprintln(w, strings.Join(stack, " -> "))
	}
	for _, p := range pkg.Deps {
		writeExplain(w, p, stack, explain)
	}
}
