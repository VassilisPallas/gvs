package flags

import (
	"flag"
	"fmt"

	terminalColors "github.com/fatih/color"
)

// Flag contains the information for the flag.
type Flag struct {
	// The name of the flag
	name string

	// A boolean that represents if the flag expects a value to be passed.
	//
	// e.g. --install-version 21.3
	acceptsVale bool
}

// getHelpName returns the name for the flag and an indicator
// if the flag expects a value of not.
//
// For example, when a flag does not expect a value the result will be like: --some-flag.
// If the flag expects a value, the result will be like: --some-flag=value
func (f Flag) getHelpName() string {
	flagName := fmt.Sprintf("--%s", f.name)
	if f.acceptsVale {
		flagName += "=value"
	}

	return flagName
}

// FlagSet is the struct that will be used to set up the flags for the CLI.
type FlagSet struct {
	// An array of the flags that are available.
	// flags is used as a store to easy iterate on the flags.
	flags []Flag
}

// FlagBool defines a bool flag with specified name, default value, and usage string.
// The argument p points to a bool variable in which to store the value of the flag.
// FlagBool also appends the flag to the FlagSet array.
func (s *FlagSet) FlagBool(p *bool, name string, value bool, usage string) {
	s.flags = append(s.flags, Flag{name: name, acceptsVale: false})
	flag.BoolVar(p, name, value, usage)
}

// StringVar defines a string flag with specified name, default value, and usage string.
// The argument p points to a string variable in which to store the value of the flag.
// FlagBool also appends the flag to the FlagSet array.
func (s *FlagSet) FlagStr(p *string, name string, value string, usage string) {
	s.flags = append(s.flags, Flag{name: name, acceptsVale: true})
	flag.StringVar(p, name, value, usage)
}

// printSynopsis returns back all the available flags without any description.
// All the flags are iterated from FlagSet array that contains the flags.
func (s *FlagSet) printSynopsis() {
	msg := "  gvs\n"
	for _, flag := range s.flags {
		msg += fmt.Sprintf("   [%s]\n", flag.getHelpName())
	}

	fmt.Printf("%s\n", msg)
}

// printFlags returns back all the available flags wiath a description.
// All the flags are iterated from FlagSet array that contains the flags.
func (s *FlagSet) printFlags() {
	flagSet := flag.CommandLine

	for _, flag := range s.flags {
		flagInfo := flagSet.Lookup(flag.name)
		fmt.Printf("  %s\n\t%s\n", flag.getHelpName(), flagInfo.Usage)
	}

}

// Parse is preparing the help command and parses the flags.
func (s *FlagSet) Parse() {
	flag.Usage = func() {
		bold := terminalColors.New().Add(terminalColors.Bold)

		gvsMessage := bold.Sprint("gvs")

		fmt.Println()
		bold.Println("NAME")
		fmt.Printf("  gvs - go version manager\n\n")

		bold.Println("DESCRIPTION")
		fmt.Printf("  the %s CLI is a command line tool to manage multiple active Go versions.\n\n", gvsMessage)

		bold.Println("SYNOPSIS")
		s.printSynopsis()

		bold.Println("FLAGS")
		s.printFlags()

		fmt.Println()
		fmt.Printf("Before start using the %s CLI, make sure to delete all the existing go versions\n", gvsMessage)
		fmt.Printf("and append to your profile file the export: %q.\n", "export PATH=$PATH:$HOME/bin")
		fmt.Printf("The profile file could be one of: (%s)\n", "~/.bash_profile, ~/.zshrc, ~/.profile, or ~/.bashrc")
	}

	flag.Parse()
}
