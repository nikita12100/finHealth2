package models

import "regexp"

type Share struct {
	Ticker    string
	HistPrice map[string]float32
}

func IsShare(ticker string) bool {
	prefixPattern := `[A-Z0-9]+`
	re := regexp.MustCompile(prefixPattern)
	return re.MatchString(ticker) 
}