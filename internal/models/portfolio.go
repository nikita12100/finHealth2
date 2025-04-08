package models

type Portfolio struct {
	ChatId     int64       `json:"chat_id"`
	Name       string      `json:"name"`
	Operations []Operation `json:"operations"`
}

func (p *Portfolio) GetCountPerTicker() map[string]int {
	countMap := make(map[string]int)
	for _, operation := range p.Operations {
		if operation.IsBuy {
			countMap[operation.Ticker] += operation.Count
		} else {
			countMap[operation.Ticker] -= operation.Count
		}
	}
	return countMap
}

func (p *Portfolio) GetAvgBuyPricePerTicker() map[string]float64 {
	sumPrice := make(map[string]float64)
	avgCount := make(map[string]int)
	for _, operation := range p.Operations {
		if operation.IsBuy {
			sumPrice[operation.Ticker] += operation.Price * float64(operation.Count)
			avgCount[operation.Ticker] += operation.Count
		}
	}
	avgPrice := make(map[string]float64)
	for ticker, sumPrice := range sumPrice {
		avgPrice[ticker] = sumPrice / float64(avgCount[ticker])
	}
	return avgPrice
}
