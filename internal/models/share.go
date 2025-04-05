package models

type Share struct {
	Ticker    string
	HistPrice map[string]float32
}
