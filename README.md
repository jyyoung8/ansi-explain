# ansi-explain

Terminal programs talk to each other with escape sequences: `ESC [ 31 m` for
red text, `ESC [ 2 J` to clear the screen, and dozens more. When that output
ends up somewhere other than a real terminal — a log file, a CI transcript,
a copy-pasted bug report — the sequences either render as garbage or
disappear entirely, and there's no easy way to tell what a program was
actually trying to do.

`ansi-explain` reads a byte stream, finds every escape sequence in it, and
prints a plain-English description of each one.

## Usage

Build it:

```
go build -o ansi-explain .
```

Pipe colored output through it:

```
$ printf '\033[31mhello\033[0m\n' | ansi-explain
"\x1b[31m"  select graphic rendition (SGR) (foreground red)
"\x1b[0m"  select graphic rendition (SGR) (reset)
```

Or point it at a captured log file:

```
$ ansi-explain build.log
"\x1b[2J"  erase in display (params=[2])
"\x1b[1;1H"  cursor position (params=[1 1])
```

Plain text between sequences is skipped in the output; only the escape
sequences themselves are reported.

Pass `-hex` to also print each sequence's raw bytes:

```
$ printf '\033[31m' | ansi-explain -hex
"\x1b[31m"  1b 5b 33 31 6d  select graphic rendition (SGR) (foreground red)
```

Pass `-strip` to skip the explanations and print the plain text with every
escape sequence removed instead:

```
$ printf '\033[31mhello\033[0m\n' | ansi-explain -strip
hello
```

## How it's built

The parsing and explanation logic (`token.go`, `explain.go`) is plain
functions over byte slices and structs — no I/O, no global state. Given the
same input they always return the same output, which makes them easy to
cover with table-driven tests without touching a real terminal. `main.go`
is the only file that reads stdin or files.

Currently understood:

- CSI sequences (`ESC [ ... final`), including SGR color/style codes
- OSC sequences (`ESC ] ... BEL` or `ESC ] ... ST`)
- Single-byte ESC sequences like save/restore cursor and terminal reset

Unrecognized sequences are still tokenized and reported, just without a
named description.

## License

MIT, see [LICENSE](LICENSE).
