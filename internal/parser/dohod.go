package parser

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func GetDivYield(ticker string) (float64, error) {
	resp, err := http.Get(fmt.Sprintf("https://www.dohod.ru/ik/analytics/dividend/%v", strings.ToLower(ticker)))
	if err != nil {
		slog.Error("Error fetching GetDivYield URL.", ticker, err)
		return 0.0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response body.", ticker, err)
		return 0.0, err
	}

	html := string(body)
	pattern := `Совокупные дивиденды в следующие 12m:.*?<span[^>]*>([\d\.]+)</span>`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(html)
	if len(matches) < 2 {
		slog.Error("Could not find dividend forecast in the page.", ticker, err)
		return 0.0, err
	}
	rawNumber := matches[1]
	cleanNumber := regexp.MustCompile(`\s+`).ReplaceAllString(rawNumber, "")

	dividend, err := strconv.ParseFloat(cleanNumber, 64)
	if err != nil {
		slog.Error("Error converting dividend to float.", ticker, err)
		return 0.0, err
	}

	return float64(dividend), nil
}