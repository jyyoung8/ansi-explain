package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "ansi-explain:", err)
		os.Exit(1)
	}
}

// run holds all of the program's I/O so main stays a thin wrapper. With no
// arguments it reads standard input; with one argument it reads that file.
func run(args []string, stdin io.Reader, stdout io.Writer) error {
	var data []byte
	var err error
	switch len(args) {
	case 0:
		data, err = io.ReadAll(stdin)
	case 1:
		data, err = os.ReadFile(args[0])
	default:
		return fmt.Errorf("usage: ansi-explain [file]")
	}
	if err != nil {
		return err
	}
	for _, tok := range Tokenize(data) {
		if tok.Kind == TokenText {
			continue
		}
		exp := Explain(tok)
		if exp.Detail != "" {
			fmt.Fprintf(stdout, "%q  %s (%s)\n", tok.Text, exp.Summary, exp.Detail)
		} else {
			fmt.Fprintf(stdout, "%q  %s\n", tok.Text, exp.Summary)
		}
	}
	return nil
}
