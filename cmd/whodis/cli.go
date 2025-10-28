package main

import (
	"fmt"
	"os"

	"github.com/knpwrs/whodis/internal/app"
	"github.com/knpwrs/whodis/pkg/models"
	flag "github.com/spf13/pflag"
)

// initFlags sets up all CLI flags
func initFlags() *flag.FlagSet {
	f := flag.NewFlagSet("whodis", flag.ContinueOnError)
	f.Usage = printHelp

	// Output format flags
	f.BoolP("json", "j", false, "Output results in JSON format")
	f.BoolP("short", "s", false, "Output minimal information in table format")
	f.BoolP("raw", "r", false, "Output raw WHOIS data")

	// Query configuration
	f.IntP("timeout", "t", 30, "Query timeout in seconds")

	// Utility flags
	f.BoolP("debug", "d", false, "Show debug information")
	f.BoolP("help", "h", false, "Show help message")
	f.BoolP("version", "v", false, "Show version information")

	return f
}

// Run is the main CLI entry point
func Run() error {
	// Initialize flags
	f := initFlags()

	// Parse command line arguments
	if err := f.Parse(os.Args[1:]); err != nil {
		if err == flag.ErrHelp {
			printHelp()
			return nil
		}
		return err
	}

	// Handle help flag
	if help, _ := f.GetBool("help"); help {
		printHelp()
		return nil
	}

	// Handle version flag
	if version, _ := f.GetBool("version"); version {
		printVersion()
		return nil
	}

	// Get remaining arguments (domains)
	domains := f.Args()
	if len(domains) == 0 {
		fmt.Fprintln(os.Stderr, colorYellow("Error:")+" no domains specified")
		fmt.Fprintln(os.Stderr, "\nRun 'whodis --help' for usage information")
		return fmt.Errorf("no domains specified")
	}

	// Get flag values directly from pflag
	jsonFlag, _ := f.GetBool("json")
	shortFlag, _ := f.GetBool("short")
	rawFlag, _ := f.GetBool("raw")
	timeout, _ := f.GetInt("timeout")
	debug, _ := f.GetBool("debug")

	// Build QueryFlags
	qf := models.QueryFlags{
		Domains: domains,
		JSON:    jsonFlag,
		Short:   shortFlag,
		Raw:     rawFlag,
		Timeout: timeout,
		Debug:   debug,
	}

	// Validate flags
	outputFlags := 0
	if qf.JSON {
		outputFlags++
	}
	if qf.Short {
		outputFlags++
	}
	if qf.Raw {
		outputFlags++
	}
	if outputFlags > 1 {
		return fmt.Errorf("only one output format flag can be specified (--json, --short, or --raw)")
	}

	// Create and run the app
	application := app.NewApp(qf)
	return application.Run()
}

func main() {
	if err := Run(); err != nil {
		// Error already printed in most cases
		os.Exit(1)
	}
}
