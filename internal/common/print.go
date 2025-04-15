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

func PrintTable3(headers []string, rows [][]string) string {
	if len(headers) != 3 {
		slog.Error("Wrong table size != 3")
		return ""
	}
	if len(headers) != len(rows[0]) {
		slog.Error("Wrong size headers!=rows[i]")
		return ""
	}

	table := fmt.Sprintf("```\n%-6s | %-8s | %-6s\n", headers[0], headers[1], headers[2])
	table += "--------------------------\n"
	for _, row := range rows {
		table += fmt.Sprintf("%-6s | %-8s | %-6s\n", row[0], row[1], row[2])
	}
	table += "```"
	return table
}

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
