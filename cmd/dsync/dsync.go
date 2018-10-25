package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/siddontang/go/ioutil2"
)

type Arg struct {
	set   bool
	value string
}

func (a *Arg) String() string     { return a.value }
func (a *Arg) Set(v string) error { a.value = v; a.set = true; return nil }

var (
	dest Arg
)

func main() {
	defer func() {
		err := recover()
		if err == nil {
			return
		}

		fmt.Fprintln(flag.CommandLine.Output())
		fmt.Fprintln(flag.CommandLine.Output(), err)
		fmt.Fprintln(flag.CommandLine.Output())
	}()

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [-dest value] srcdir1 [srcdir2 [...]]", os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output())
		fmt.Fprintln(flag.CommandLine.Output(), "  srcdir1, ..., srcdirN: sources directories")
		fmt.Fprintln(flag.CommandLine.Output())

		flag.PrintDefaults()
	}

	flag.Var(&dest, "dest", "Set the destination directory")
	flag.Parse()

	if !dest.set {
		fmt.Fprintln(flag.CommandLine.Output(), "No destination specified")
		fmt.Fprintln(flag.CommandLine.Output())

		flag.Usage()

		return
	}

	checkdir("destination", dest.value)

	if flag.NArg() == 0 {
		fmt.Fprintln(flag.CommandLine.Output(), "No source specified")
		fmt.Fprintln(flag.CommandLine.Output())

		flag.Usage()

		return
	}

	sources := flag.Args()
	dest, err := filepath.Abs(dest.value)
	if err != nil {
		panic(err)
	}

	fmt.Println("Destination:", dest)
	for i, src := range sources {
		src, err := filepath.Abs(src)
		if err != nil {
			panic(err)
		}

		if src == dest {
			panic(fmt.Errorf("Destination cannot be in source"))
		}

		if strings.HasPrefix(src, dest+"/") {
			panic(fmt.Errorf("`%s` is a sub-directory of destination", src))
		}

		checkdir("source", src)

		sources[i] = src
	}

	fmt.Println("Sources:")
	for _, src := range sources {
		fmt.Println("  -", src)
	}

	doBackup(dest, sources)
}
