package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunPlainOutput(t *testing.T) {
	var out bytes.Buffer
	in := strings.NewReader("\x1b[31mred\x1b[0m")
	if err := run(nil, in, &out); err != nil {
		t.Fatalf("run: %v", err)
	}
	got := out.String()
	if strings.Contains(got, "1b 5b") {
		t.Fatalf("output has hex bytes without -hex flag: %q", got)
	}
	if !strings.Contains(got, `"\x1b[31m"`) {
		t.Fatalf("output missing quoted sequence: %q", got)
	}
}

func TestRunHexFlag(t *testing.T) {
	var out bytes.Buffer
	in := strings.NewReader("\x1b[31m")
	if err := run([]string{"-hex"}, in, &out); err != nil {
		t.Fatalf("run: %v", err)
	}
	got := out.String()
	if !strings.Contains(got, "1b 5b 33 31 6d") {
		t.Fatalf("output missing hex bytes: %q", got)
	}
}

func TestRunStripFlag(t *testing.T) {
	var out bytes.Buffer
	in := strings.NewReader("\x1b[31mred\x1b[0m and \x1b]0;title\x07plain")
	if err := run([]string{"-strip"}, in, &out); err != nil {
		t.Fatalf("run: %v", err)
	}
	got := out.String()
	if want := "red and plain"; got != want {
		t.Fatalf("stripped output = %q, want %q", got, want)
	}
}

func TestRunTooManyArgs(t *testing.T) {
	var out bytes.Buffer
	err := run([]string{"a", "b"}, strings.NewReader(""), &out)
	if err == nil {
		t.Fatal("expected error for too many arguments")
	}
}

func TestRunUnknownFlag(t *testing.T) {
	var out bytes.Buffer
	err := run([]string{"-nope"}, strings.NewReader(""), &out)
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
}
