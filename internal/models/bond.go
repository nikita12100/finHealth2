package models

import "regexp"

type Bond struct {
	Ticker string
}

func IsBond(ticker string) bool {
	prefixPattern := `^[A-Z]{2}[0-9]{2}[A-Z0-9]+`
	re := regexp.MustCompile(prefixPattern)
	return re.MatchString(ticker) 
}
