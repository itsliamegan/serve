package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/itsliamegan/serve"
)

func main() {
	var port uint
	var publish bool

	flag.UintVar(&port, "port", 4000, "port to bind the server to")
	flag.BoolVar(&publish, "publish", false, "bind the server to the local IP address")

	flag.Usage = func() {
		fmt.Println("serve - make a directory and its descendants available via HTTP")
		fmt.Println()
		fmt.Println("usage: serve [options] [<directory>]")
		fmt.Println()
		fmt.Println("options:")
		flag.PrintDefaults()
	}
	flag.Parse()

	var rootDir string
	if len(flag.Args()) != 0 {
		rootDir = flag.Arg(0)
	} else {
		workingDir, err := os.Getwd()
		exitIf(err)
		rootDir = workingDir
	}

	var addr string
	if publish {
		addr = fmt.Sprintf("0.0.0.0:%d", port)
	} else {
		addr = fmt.Sprintf(":%d", port)
	}

	err := serve.Start(rootDir, addr)
	exitIf(err)
}

func exitIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
