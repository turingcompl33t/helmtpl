package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"helmtpl/internal/engine"
	"helmtpl/internal/logger"
)

func main() {
	flags := readFlags()
	if err := flags.validate(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	logger := logger.New(flags.verbose)
	logger.Debug(fmt.Sprintf("templating input from %s to %s", flags.input, flags.output))

	engine := engine.New("vars", logger)

	data, err := readInput(flags.input)
	if err != nil {
		logger.Fatal(err.Error())
	}

	r, err := engine.Run(data)
	if err != nil {
		logger.Fatal(err.Error())
	}

	if err := os.WriteFile(flags.output, r, 0644); err != nil {
		logger.Fatal(err.Error())
	}

	// Write the output file name on success
	fmt.Println(flags.output)
}

// ----------------------------------------------------------------------------
// Flag Parsing
// ----------------------------------------------------------------------------

type Flags struct {
	input   string
	output  string
	force   bool
	verbose bool
}

func readFlags() *Flags {
	// The path to the input file
	input := flag.String("input", "", "The path to the input (*.tpl) file")
	// The path to the output file
	output := flag.String("output", "", "The path to the output (*.yaml) file")
	// Force overwrite of the output file
	force := flag.Bool("force", false, "Force overwrite of the output file")
	// Enable verbose output
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	flag.Parse()

	return &Flags{
		input:   *input,
		output:  *output,
		force:   *force,
		verbose: *verbose,
	}
}

// Validate the input flags.
func (f *Flags) validate() error {
	if err := f.validateInput(); err != nil {
		return err
	}
	if err := f.validateOutput(); err != nil {
		return err
	}
	return nil
}

// Validate the input argument.
func (f *Flags) validateInput() error {
	if f.input == "" {
		return errors.New("valid input file is required")
	}
	if !fileExists(f.input) {
		return fmt.Errorf("input file %s does not exist", f.input)
	}
	return nil
}

// Validate the output argument.
func (f *Flags) validateOutput() error {
	if f.output == "" {
		name := strings.Split(filepath.Base(f.input), ".")[0]
		f.output = filepath.Join(filepath.Dir(f.input), fmt.Sprintf("%s.yaml", name))
	}

	if fileExists(f.output) {
		if f.force {
			os.Remove(f.output)
		} else {
			return fmt.Errorf("output file %s already exists", f.output)
		}
	}

	return nil
}

// Determine if a file exists.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// ----------------------------------------------------------------------------
// General Utilities
// ----------------------------------------------------------------------------

// Read the contents of the input file.
func readInput(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return data, nil
}
