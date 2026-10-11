package main

import (
	"os"
	"slices"
	"testing"
)

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
