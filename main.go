package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

const usage = "usage: ansi-explain [-hex] [-strip] [file]"

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "ansi-explain:", err)
		os.Exit(1)
	}
}

// run holds all of the program's I/O so main stays a thin wrapper. With no
// file argument it reads standard input; with one it reads that file.
func run(args []string, stdin io.Reader, stdout io.Writer) error {
	fs := flag.NewFlagSet("ansi-explain", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	showHex := fs.Bool("hex", false, "show each sequence's raw bytes as hex alongside the quoted text")
	strip := fs.Bool("strip", false, "remove escape sequences and print the remaining plain text, instead of explaining them")
	if err := fs.Parse(args); err != nil {
		return errors.New(usage)
	}

	var data []byte
	var err error
	switch fs.NArg() {
	case 0:
		data, err = io.ReadAll(stdin)
	case 1:
		data, err = os.ReadFile(fs.Arg(0))
	default:
		return errors.New(usage)
	}
	if err != nil {
		return err
	}

	if *strip {
		for _, tok := range Tokenize(data) {
			if tok.Kind == TokenText {
				stdout.Write(tok.Text)
			}
		}
		return nil
	}

	for _, tok := range Tokenize(data) {
		if tok.Kind == TokenText {
			continue
		}
		exp := Explain(tok)
		text := fmt.Sprintf("%q", tok.Text)
		if *showHex {
			text += "  " + fmt.Sprintf("% x", tok.Text)
		}
		if exp.Detail != "" {
			fmt.Fprintf(stdout, "%s  %s (%s)\n", text, exp.Summary, exp.Detail)
		} else {
			fmt.Fprintf(stdout, "%s  %s\n", text, exp.Summary)
		}
	}
	return nil
}
