package common

import (
	"fmt"
	"log/slog"
)

func PrintTable2(headers []string, rows [][]string) string {
	if len(headers) != 2 {
		slog.Error("Wrong table size != 2")
		return ""
	}
	if len(headers) != len(rows[0]) {
		slog.Error("Wrong size headers!=rows[i]")
		return ""
	}

	table := fmt.Sprintf("```\n%-6s | %-8s\n", headers[0], headers[1])
	table += "-------------------\n"
	for _, row := range rows {
		table += fmt.Sprintf("%-6s | %-8s\n", row[0], row[1])
	}
	table += "```"
	return table
}
