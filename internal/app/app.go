package app

import (
	"fmt"
	"sync"
	"time"

	"github.com/knpwrs/whodis/pkg/models"
	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
)

// App is the main application struct
type App struct {
	QueryFlags models.QueryFlags
}

// NewApp creates a new App instance
func NewApp(qf models.QueryFlags) *App {
	return &App{
		QueryFlags: qf,
	}
}

// LookupDomains performs WHOIS lookups for all configured domains
func (app *App) LookupDomains() ([]models.WhoisResult, error) {
	if len(app.QueryFlags.Domains) == 0 {
		return nil, fmt.Errorf("no domains specified")
	}

	results := make([]models.WhoisResult, len(app.QueryFlags.Domains))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, domain := range app.QueryFlags.Domains {
		wg.Add(1)
		go func(idx int, dom string) {
			defer wg.Done()
			result := app.lookupDomain(dom)
			mu.Lock()
			results[idx] = result
			mu.Unlock()
		}(i, domain)
	}

	wg.Wait()
	return results, nil
}

// lookupDomain performs a WHOIS lookup for a single domain
func (app *App) lookupDomain(domain string) models.WhoisResult {
	result := models.WhoisResult{
		Domain: domain,
	}

	// Create WHOIS client with timeout
	client := whois.NewClient()
	timeout := time.Duration(app.QueryFlags.Timeout) * time.Second
	client.SetTimeout(timeout)

	// Perform the lookup
	startTime := time.Now()
	rawText, err := client.Whois(domain)
	result.QueryTime = time.Since(startTime)

	if err != nil {
		result.Error = err
		return result
	}

	result.RawText = rawText

	// Parse the WHOIS data unless we're in raw mode
	if !app.QueryFlags.Raw {
		parsedInfo, err := whoisparser.Parse(rawText)
		if err != nil {
			// Don't fail the entire result if parsing fails
			if app.QueryFlags.Debug {
				result.Error = fmt.Errorf("parse error: %w", err)
			}
		} else {
			result.ParsedInfo = &parsedInfo
		}
	}

	return result
}

// Run executes the application
func (app *App) Run() error {
	results, err := app.LookupDomains()
	if err != nil {
		return err
	}

	// Determine output format
	var format models.OutputFormat
	switch {
	case app.QueryFlags.Raw:
		format = models.OutputRaw
	case app.QueryFlags.JSON:
		format = models.OutputJSON
	case app.QueryFlags.Short:
		format = models.OutputShort
	default:
		format = models.OutputTerminal
	}

	return Output(results, format)
}
