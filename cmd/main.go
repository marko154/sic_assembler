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

	infilename := flag.Arg(0)
	infile, err := os.Open(infilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
	defer infile.Close()

	outfilename := strings.Replace(infilename, "asm", "obj", 1)
	outfile, err := os.Create(outfilename)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %s\n", err)
		return
	}
	defer outfile.Close()
	fmt.Printf("Output file: %s\n", outfilename)

	assembler := assembler.NewAssembler(outfile)
	assembler.Assemble(infile)
}
