package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/emicklei/melrose/dsl"
	"github.com/emicklei/melrose/server"
	"github.com/emicklei/melrose/system"
	"github.com/emicklei/melrose/ui/cli"
	"golang.org/x/term"
)

var (
	BuildTag = "dev"
)

func main() {
	command, filename, err := parseCommand(os.Args[1:])
	if err != nil {
		log.Fatalln(err)
	}
	if command == "check" {
		if err := validateSource(filename); err != nil {
			log.Fatalln(err)
		}
		return
	}
	ctx, err := system.Setup(BuildTag)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := system.TearDown(ctx); err != nil {
			log.Printf("error during shutdown: %v", err)
		}
	}()
	if command == "midi" {
		if device, ok := ctx.Device().(interface{ PrintInfoWithoutInstructions() }); ok {
			device.PrintInfoWithoutInstructions()
		} else {
			ctx.Device().Command(nil)
		}
		return
	}
	if command == "play" {
		if err := cli.ExecuteFile(ctx, filename); err != nil {
			log.Fatalln(err)
		}
		if err := waitForPlayInput(filename, os.Stdin, os.Stdout); err != nil {
			log.Println(err)
		}
		return
	}
	if command == "load" {
		if err := cli.LoadFile(ctx, filename); err != nil {
			log.Fatalln(err)
		}
	}
	// start the server and the REPL
	server.Start(ctx)
	cli.StartREPL(ctx)
}

func waitForPlayInput(filename string, input *os.File, output io.Writer) error {
	for {
		if _, err := fmt.Fprintln(output, "Press 'Enter' to exit or 'p' to show the file..."); err != nil {
			return err
		}
		key, err := readPlayKey(input)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		switch key {
		case '\r', '\n', 3, 4: // Enter, Ctrl-C, or Ctrl-D.
			return nil
		case 'p':
			content, err := os.ReadFile(filename)
			if err != nil {
				if _, err := fmt.Fprintf(output, "Error reading file: %v\n", err); err != nil {
					return err
				}
				continue
			}
			fmt.Fprintln(output, ">>>")
			if _, err := output.Write(content); err != nil {
				return err
			}
			fmt.Fprintln(output, "<<<")
			if len(content) > 0 && content[len(content)-1] != '\n' {
				if _, err := fmt.Fprintln(output); err != nil {
					return err
				}
			}
		}
	}
}

func readPlayKey(input *os.File) (key byte, err error) {
	fd := int(input.Fd())
	if term.IsTerminal(fd) {
		state, err := term.MakeRaw(fd)
		if err != nil {
			return 0, fmt.Errorf("setting terminal input mode: %w", err)
		}
		defer func() {
			if restoreErr := term.Restore(fd, state); restoreErr != nil {
				err = errors.Join(err, fmt.Errorf("restoring terminal input mode: %w", restoreErr))
			}
		}()
	}
	var buffer [1]byte
	_, err = io.ReadFull(input, buffer[:])
	return buffer[0], err
}

func validateSource(filename string) error {
	var source []byte
	var err error
	if filename == "" || filename == "-" {
		source, err = io.ReadAll(os.Stdin)
	} else {
		source, err = os.ReadFile(filename)
	}
	if err != nil {
		return err
	}
	return dsl.Validate(string(source))
}

func parseCommand(arguments []string) (string, string, error) {
	flags, positional := splitFlags(arguments)
	os.Args = append([]string{os.Args[0]}, flags...)

	if len(positional) == 0 {
		return "", "", nil
	}
	if positional[0] == "check" {
		if len(positional) == 1 {
			return "check", "", nil
		}
		if len(positional) == 2 {
			return "check", positional[1], nil
		}
		return "", "", errors.New("usage: melrose check [file]")
	}
	if positional[0] == "midi" {
		if len(positional) == 1 {
			return "midi", "", nil
		}
		return "", "", errors.New("usage: melrose midi")
	}
	if len(positional) != 2 {
		return "", "", errors.New("usage: melrose <load|play> <file>")
	}
	if positional[0] != "load" && positional[0] != "play" {
		return "", "", fmt.Errorf("unknown subcommand %q; expected load, play, check, or midi", positional[0])
	}
	return positional[0], positional[1], nil
}

func splitFlags(arguments []string) ([]string, []string) {
	var flags, positional []string
	for index := 0; index < len(arguments); index++ {
		argument := arguments[index]
		if argument == "--" {
			positional = append(positional, arguments[index+1:]...)
			break
		}
		name := strings.TrimLeft(argument, "-")
		if name == argument || name == "" {
			positional = append(positional, argument)
			continue
		}
		name, _, hasValue := strings.Cut(name, "=")
		registered := flag.CommandLine.Lookup(name)
		if registered == nil {
			positional = append(positional, argument)
			continue
		}
		flags = append(flags, argument)
		if hasValue || isBoolFlag(registered) {
			continue
		}
		if index+1 < len(arguments) {
			index++
			flags = append(flags, arguments[index])
		}
	}
	return flags, positional
}

func isBoolFlag(value *flag.Flag) bool {
	boolean, ok := value.Value.(interface{ IsBoolFlag() bool })
	return ok && boolean.IsBoolFlag()
}
