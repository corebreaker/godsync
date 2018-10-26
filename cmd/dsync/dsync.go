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

func main() {
	defer func() {
		err := recover()
		if err == nil {
			return
		}

		out := flag.CommandLine.Output()

		fmt.Fprintln(out)
		fmt.Fprintln(out, err)
		fmt.Fprintln(out)
	}()

	flag.Usage = func() {
		out := flag.CommandLine.Output()

		fmt.Fprintf(out, "Usage: %s [-dest dirpath] [-root rootpath] srcdir1 [srcdir2 [...]]", os.Args[0])
		fmt.Fprintln(out)
		fmt.Fprintln(out, "  srcdir1, ..., srcdirN: source directories")
		fmt.Fprintln(out)

		flag.PrintDefaults()
	}

	var dest, base Arg

	flag.Var(&dest, "dest", "Set the destination directory")
	flag.Var(&base, "base", "Set the base directory for all source directories")
	flag.Parse()

	checkdir("destination", dest.value)
	checkdir("source base", base.value)

	if flag.NArg() == 0 {
		out := flag.CommandLine.Output()

		fmt.Fprintln(out, "No source specified")
		fmt.Fprintln(out)

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
