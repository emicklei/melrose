package main

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestWaitForPlayInput(t *testing.T) {
	for _, test := range []struct {
		name       string
		keys       string
		missing    bool
		wantPrints int
		wantError  bool
	}{
		{name: "enter exits", keys: "\n"},
		{name: "carriage return exits", keys: "\r"},
		{name: "print repeatedly", keys: "pp\np", wantPrints: 2},
		{name: "ignore unknown keys", keys: "xp\n", wantPrints: 1},
		{name: "end of input exits"},
		{name: "control c exits", keys: "\x03p"},
		{name: "control d exits", keys: "\x04p"},
		{name: "missing file reports error", keys: "p\n", missing: true, wantError: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			filename := filepath.Join(dir, "source.mel")
			const content = "sequence(\"c d e\")"
			if !test.missing {
				if err := os.WriteFile(filename, []byte(content), 0600); err != nil {
					t.Fatal(err)
				}
			}
			keyfile := filepath.Join(dir, "keys")
			if err := os.WriteFile(keyfile, []byte(test.keys), 0600); err != nil {
				t.Fatal(err)
			}
			input, err := os.Open(keyfile)
			if err != nil {
				t.Fatal(err)
			}
			defer input.Close()
			var output bytes.Buffer
			if err := waitForPlayInput(filename, input, &output); err != nil {
				t.Fatal(err)
			}
			if got := strings.Count(output.String(), content); got != test.wantPrints {
				t.Fatalf("printed file %d times, want %d; output: %q", got, test.wantPrints, output.String())
			}
			if got := strings.Contains(output.String(), "Error reading file:"); got != test.wantError {
				t.Fatalf("file error reported = %t, want %t; output: %q", got, test.wantError, output.String())
			}
		})
	}
}

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name         string
		arguments    []string
		wantCommand  string
		wantFilename string
		wantFlagArgs []string
		wantError    bool
	}{
		{
			name:         "play with debug flag after filename",
			arguments:    []string{"play", "somefile.mel", "-d"},
			wantCommand:  "play",
			wantFilename: "somefile.mel",
			wantFlagArgs: []string{"-d"},
		},
		{
			name:         "load with debug flag before subcommand",
			arguments:    []string{"-d", "load", "somefile.mel"},
			wantCommand:  "load",
			wantFilename: "somefile.mel",
			wantFlagArgs: []string{"-d"},
		},
		{
			name:         "play with value flag after filename",
			arguments:    []string{"play", "somefile.mel", "-log", "melrose.log"},
			wantCommand:  "play",
			wantFilename: "somefile.mel",
			wantFlagArgs: []string{"-log", "melrose.log"},
		},
		{
			name:         "midi",
			arguments:    []string{"midi"},
			wantCommand:  "midi",
			wantFlagArgs: []string{},
		},
		{
			name:         "midi with debug flag",
			arguments:    []string{"midi", "-d"},
			wantCommand:  "midi",
			wantFlagArgs: []string{"-d"},
		},
		{
			name:      "midi with unexpected argument",
			arguments: []string{"midi", "extra"},
			wantError: true,
		},
		{
			name:      "missing filename",
			arguments: []string{"load"},
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			originalArgs := os.Args
			t.Cleanup(func() { os.Args = originalArgs })
			os.Args = []string{"melrose"}

			command, filename, err := parseCommand(test.arguments)
			if (err != nil) != test.wantError {
				t.Fatalf("parseCommand() error = %v, wantError %t", err, test.wantError)
			}
			if command != test.wantCommand || filename != test.wantFilename {
				t.Fatalf("parseCommand() = (%q, %q), want (%q, %q)", command, filename, test.wantCommand, test.wantFilename)
			}
			if !slices.Equal(os.Args[1:], test.wantFlagArgs) {
				t.Fatalf("flag arguments = %q, want %q", os.Args[1:], test.wantFlagArgs)
			}
		})
	}
}
