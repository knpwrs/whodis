package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/fatih/color"
)

var (
	colorGreen   = color.New(color.FgGreen, color.Bold).SprintFunc()
	colorCyan    = color.New(color.FgCyan, color.Bold).SprintFunc()
	colorYellow  = color.New(color.FgYellow, color.Bold).SprintFunc()
	colorMagenta = color.New(color.FgMagenta, color.Bold).SprintFunc()
)

const helpTemplate = `{{call .Cyan "whodis"}} - pretty WHOIS lookup tool

{{call .Green "USAGE:"}}
  whodis [flags] domain [domain...]

{{call .Green "EXAMPLES:"}}
  {{call .Cyan "# Basic lookup"}}
  whodis example.com

  {{call .Cyan "# Multiple domains"}}
  whodis example.com github.com

  {{call .Cyan "# JSON output"}}
  whodis --json example.com

  {{call .Cyan "# Short table format"}}
  whodis --short example.com github.com

  {{call .Cyan "# Raw WHOIS output"}}
  whodis --raw example.com

  {{call .Cyan "# Custom timeout"}}
  whodis --timeout 30 example.com

{{call .Green "FLAGS:"}}
  -j, --json              Output results in JSON format
  -s, --short             Output minimal information in table format
  -r, --raw               Output raw WHOIS data
  -t, --timeout int       Query timeout in seconds (default: 30)
  -d, --debug             Show debug information
  -h, --help              Show this help message
  -v, --version           Show version information

{{call .Green "ABOUT:"}}
  whodis is to whois what doggo is to dig - a modern, colorful CLI tool
  that makes WHOIS queries more readable and enjoyable.

  Repository: {{call .Magenta "https://github.com/knpwrs/whodis"}}
  License: Public Domain (CC0/Unlicense)
`

type helpData struct {
	Cyan    func(...interface{}) string
	Green   func(...interface{}) string
	Yellow  func(...interface{}) string
	Magenta func(...interface{}) string
}

// printHelp renders and prints the help message
func printHelp() {
	tmpl := template.Must(template.New("help").Parse(helpTemplate))
	data := helpData{
		Cyan:    colorCyan,
		Green:   colorGreen,
		Yellow:  colorYellow,
		Magenta: colorMagenta,
	}

	if err := tmpl.Execute(os.Stdout, data); err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering help: %v\n", err)
	}
}

// printVersion prints version information
func printVersion() {
	fmt.Printf("%s %s\n", colorCyan("whodis"), colorGreen("v0.1.0"))
}
