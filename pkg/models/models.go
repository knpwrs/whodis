package models

import (
	"time"

	whoisparser "github.com/likexian/whois-parser"
)

// QueryFlags holds all the configuration for a WHOIS query
type QueryFlags struct {
	Domains []string `koanf:"domains"`
	JSON    bool     `koanf:"json"`
	Short   bool     `koanf:"short"`
	Raw     bool     `koanf:"raw"`
	Timeout int      `koanf:"timeout"`
	Debug   bool     `koanf:"debug"`
}

// WhoisResult represents the result of a WHOIS query
type WhoisResult struct {
	Domain     string
	RawText    string
	ParsedInfo *whoisparser.WhoisInfo
	QueryTime  time.Duration
	Error      error
}

// OutputFormat represents the different output formats
type OutputFormat string

const (
	OutputTerminal OutputFormat = "terminal"
	OutputJSON     OutputFormat = "json"
	OutputShort    OutputFormat = "short"
	OutputRaw      OutputFormat = "raw"
)
