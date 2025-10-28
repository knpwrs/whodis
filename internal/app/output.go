package app

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/fatih/color"
	"github.com/knpwrs/whodis/pkg/models"
)

// Color definitions
var (
	colorGreen  = color.New(color.FgGreen, color.Bold).SprintFunc()
	colorCyan   = color.New(color.FgCyan, color.Bold).SprintFunc()
	colorYellow = color.New(color.FgYellow, color.Bold).SprintFunc()
	colorRed    = color.New(color.FgRed, color.Bold).SprintFunc()
)

// Output writes results in the specified format
func Output(results []models.WhoisResult, format models.OutputFormat) error {
	switch format {
	case models.OutputTerminal:
		return outputTerminal(results)
	case models.OutputJSON:
		return outputJSON(results)
	case models.OutputShort:
		return outputShort(results)
	case models.OutputRaw:
		return outputRaw(results)
	default:
		return fmt.Errorf("unknown output format: %s", format)
	}
}

// outputTerminal prints colorized, formatted output
func outputTerminal(results []models.WhoisResult) error {
	for i, result := range results {
		if i > 0 {
			fmt.Println()
		}

		if result.Error != nil {
			fmt.Printf("%s %s\n", colorRed("ERROR:"), result.Error)
			continue
		}

		if result.ParsedInfo == nil {
			fmt.Printf("%s No parsed data available for %s\n", colorYellow("WARNING:"), result.Domain)
			continue
		}

		info := result.ParsedInfo

		// Domain Information Section
		fmt.Println(colorCyan("DOMAIN INFORMATION"))
		if info.Domain != nil {
			tree := &treeSection{}
			tree.add("Name", info.Domain.Domain)

			if len(info.Domain.Status) > 0 {
				tree.add("Status", strings.Join(info.Domain.Status, ", "))
			}

			if info.Domain.CreatedDate != "" {
				createdStr := info.Domain.CreatedDate
				if info.Domain.CreatedDateInTime != nil {
					age := time.Since(*info.Domain.CreatedDateInTime)
					years := int(age.Hours() / 24 / 365)
					if years > 0 {
						createdStr = fmt.Sprintf("%s (%d years ago)", createdStr, years)
					}
				}
				tree.add("Created", createdStr)
			}

			if info.Domain.UpdatedDate != "" {
				tree.add("Updated", info.Domain.UpdatedDate)
			}

			if info.Domain.ExpirationDate != "" {
				expiryStr := info.Domain.ExpirationDate
				if info.Domain.ExpirationDateInTime != nil {
					duration := time.Until(*info.Domain.ExpirationDateInTime)
					days := int(duration.Hours() / 24)
					if days > 0 {
						expiryStr = fmt.Sprintf("%s (in %d days)", expiryStr, days)
					} else if days < 0 {
						expiryStr = fmt.Sprintf("%s (%s)", expiryStr, colorRed("EXPIRED"))
					}
				}
				tree.add("Expires", expiryStr)
			}

			tree.render()
		}

		// Registrar Information Section
		if info.Registrar != nil && info.Registrar.Name != "" {
			fmt.Println()
			fmt.Println(colorCyan("REGISTRAR"))

			tree := &treeSection{}
			tree.add("Name", info.Registrar.Name)
			tree.add("Email", info.Registrar.Email)
			tree.add("WHOIS Server", info.Registrar.ReferralURL)
			tree.render()
		}

		// Nameservers Section
		if info.Domain != nil && len(info.Domain.NameServers) > 0 {
			fmt.Println()
			fmt.Println(colorCyan("NAMESERVERS"))

			tree := &treeSection{}
			for _, ns := range info.Domain.NameServers {
				tree.addPlain(ns)
			}
			tree.render()
		}

		// Registrant Contact Section
		if info.Registrant != nil && (info.Registrant.Name != "" || info.Registrant.Email != "" || info.Registrant.Organization != "") {
			fmt.Println()
			fmt.Println(colorCyan("REGISTRANT"))

			tree := &treeSection{}
			tree.add("Name", info.Registrant.Name)
			tree.add("Organization", info.Registrant.Organization)
			tree.add("Email", info.Registrant.Email)
			tree.render()
		}

		// Query timing
		fmt.Println()
		fmt.Printf("%s %v\n", colorCyan("Query time:"), result.QueryTime)
	}

	return nil
}

// treeSection helps build tree-style output sections
type treeSection struct {
	items []treeItem
}

type treeItem struct {
	label string
	value string
}

// add adds a field to the tree section
func (t *treeSection) add(label, value string) {
	if value != "" {
		t.items = append(t.items, treeItem{label: label, value: value})
	}
}

// addPlain adds a plain value (no label) to the tree section
func (t *treeSection) addPlain(value string) {
	if value != "" {
		t.items = append(t.items, treeItem{value: value})
	}
}

// render prints the tree section
func (t *treeSection) render() {
	for i, item := range t.items {
		isLast := i == len(t.items)-1
		prefix := "├─"
		if isLast {
			prefix = "└─"
		}

		if item.label != "" {
			fmt.Printf("%s %s: %s\n", prefix, colorGreen(item.label), item.value)
		} else {
			fmt.Printf("%s %s\n", prefix, item.value)
		}
	}
}

// outputJSON prints results as JSON
func outputJSON(results []models.WhoisResult) error {
	output := make([]map[string]interface{}, len(results))

	for i, result := range results {
		entry := map[string]interface{}{
			"domain":     result.Domain,
			"query_time": result.QueryTime.String(),
		}

		if result.Error != nil {
			entry["error"] = result.Error.Error()
		} else {
			entry["data"] = result.ParsedInfo
		}

		output[i] = entry
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}

// outputShort prints minimal information
func outputShort(results []models.WhoisResult) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Print header
	fmt.Fprintln(w, "DOMAIN\tREGISTRAR\tEXPIRES\tSTATUS")

	for _, result := range results {
		if result.Error != nil {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", result.Domain, "ERROR", "-", result.Error.Error())
			continue
		}

		if result.ParsedInfo == nil {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", result.Domain, "N/A", "-", "Parse failed")
			continue
		}

		info := result.ParsedInfo
		registrar := "N/A"
		expires := "N/A"
		status := "N/A"

		if info.Registrar != nil && info.Registrar.Name != "" {
			registrar = info.Registrar.Name
		}

		if info.Domain != nil {
			if info.Domain.ExpirationDate != "" {
				expires = info.Domain.ExpirationDate
			}
			if len(info.Domain.Status) > 0 {
				status = info.Domain.Status[0]
			}
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", result.Domain, registrar, expires, status)
	}

	return nil
}

// outputRaw prints raw WHOIS text
func outputRaw(results []models.WhoisResult) error {
	for i, result := range results {
		if i > 0 {
			fmt.Println("\n" + strings.Repeat("-", 80) + "\n")
		}

		fmt.Printf("%s %s\n", colorGreen("Domain:"), result.Domain)
		fmt.Println()

		if result.Error != nil {
			fmt.Printf("%s %s\n", colorRed("ERROR:"), result.Error)
			continue
		}

		fmt.Println(result.RawText)
	}

	return nil
}
