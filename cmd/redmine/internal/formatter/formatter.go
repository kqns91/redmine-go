// Package formatter provides output formatting utilities for CLI commands.
package formatter

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/olekukonko/tablewriter"
)

// Format represents the output format type.
type Format string

const (
	// FormatJSON outputs data in JSON format.
	FormatJSON Format = "json"
	// FormatTable outputs data in table format.
	FormatTable Format = "table"
	// FormatText outputs data in structured text format.
	FormatText Format = "text"
)

// Styles defines the color scheme for output formatting.
var Styles = struct {
	Title      lipgloss.Style
	Header     lipgloss.Style
	Key        lipgloss.Style
	Value      lipgloss.Style
	StatusOpen lipgloss.Style
	StatusDone lipgloss.Style
	Priority   lipgloss.Style
	Subtle     lipgloss.Style
}{
	Title:      lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")),
	Header:     lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14")),
	Key:        lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10")),
	Value:      lipgloss.NewStyle().Foreground(lipgloss.Color("15")),
	StatusOpen: lipgloss.NewStyle().Foreground(lipgloss.Color("10")),
	StatusDone: lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
	Priority:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9")),
	Subtle:     lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
}

// IsColorEnabled checks if color output should be enabled.
func IsColorEnabled() bool {
	// Respect NO_COLOR environment variable
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	// Check if stdout is a terminal
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		return false
	}
	return true
}

// OutputJSON outputs data in JSON format.
func OutputJSON(data interface{}) error {
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
	}
	fmt.Println(string(output))
	return nil
}

// RenderTable renders data as a simple table with headers.
func RenderTable(headers []string, rows [][]string) {
	table := tablewriter.NewTable(os.Stdout)

	// Convert headers to []any
	headerAny := make([]any, len(headers))
	for i, h := range headers {
		headerAny[i] = h
	}
	table.Header(headerAny...)

	// Convert rows to [][]any
	_ = table.Bulk(convertToAny(rows))
	_ = table.Render()
}

// convertToAny converts [][]string to [][]any for tablewriter.
func convertToAny(rows [][]string) [][]any {
	result := make([][]any, len(rows))
	for i, row := range rows {
		result[i] = make([]any, len(row))
		for j, cell := range row {
			result[i][j] = cell
		}
	}
	return result
}

// FormatKeyValue formats a key-value pair with styling.
func FormatKeyValue(key, value string) string {
	if !IsColorEnabled() {
		return fmt.Sprintf("%s: %s", key, value)
	}
	return fmt.Sprintf("%s %s", Styles.Key.Render(key+":"), Styles.Value.Render(value))
}

// FormatSection formats a section header.
func FormatSection(title string) string {
	if !IsColorEnabled() {
		return "\n" + title + "\n" + strings.Repeat("=", len(title))
	}
	return "\n" + Styles.Header.Render(title)
}

// FormatTitle formats a main title.
func FormatTitle(title string) string {
	if !IsColorEnabled() {
		return title
	}
	return Styles.Title.Render(title)
}

// TruncateString truncates a string to the specified length and adds ellipsis.
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
