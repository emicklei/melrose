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
	defer system.TearDown(ctx)
	if command == "play" {
		if err := cli.ExecuteFile(ctx, filename); err != nil {
			log.Fatalln(err)
		}
		// wait for the user to press enter before exiting
		fmt.Println("Press Enter to exit...")
		fmt.Scanln()
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
	if len(positional) != 2 {
		return "", "", errors.New("usage: melrose <load|play> <file>")
	}
	if positional[0] != "load" && positional[0] != "play" {
		return "", "", fmt.Errorf("unknown subcommand %q; expected load, play, or check", positional[0])
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
