package main

import (
	"flag"
	"fmt"
	"os"
	assembler "sic_assembler/internal"
	"strings"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <filename>\n", os.Args[0])
		flag.PrintDefaults()
	}
	debug := flag.Bool("debug", false, "Output debug information")

	if *debug {
		fmt.Println("Debug mode enabled")
	}

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		return
	}

	filename := flag.Arg(0)
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}

	outfilename := strings.Replace(filename, "asm", "obj", 1)
	outfile, err := os.Create(outfilename)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %s\n", err)
		return
	}

	ass := assembler.NewAssembler(outfile)
	ass.Assemble(file)
}
