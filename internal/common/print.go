package common

import (
	"fmt"
	"log/slog"
)

func PrintTable4(headers []string, rows [][]string) string {
	if len(headers) != 4 {
		slog.Error("Wrong table size != 4")
		return ""
	}
	if len(headers) != len(rows[0]) {
		slog.Error("Wrong size headers!=rows[i]")
		return ""
	}

	table := fmt.Sprintf("```\n%-3s | %-6s | %-6s | %-6s\n", headers[0], headers[1], headers[2], headers[3])
	table += "------------------------------\n"
	for _, row := range rows {
		table += fmt.Sprintf("%-3s | %-6s | %-6s | %-6s\n", row[0], row[1], row[2], row[3])
	}
	table += "```"
	return table
}
