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
		panic(fmt.Sprintf("Error creating output file: %s", err))
	}
	defer outfile.Close()

	lstoutfilename := strings.Replace(infilename, "asm", "lst", 1)
	lstoutfile, err := os.Create(lstoutfilename)

	if err != nil {
		panic(fmt.Sprintf("Error creating output file: %s", err))
	}
	defer lstoutfile.Close()

	assembler := assembler.NewAssembler(outfile, lstoutfile)
	assembler.Assemble(infile)
}
