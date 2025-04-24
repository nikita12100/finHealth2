package stats

import (
	"math"
	"test2/internal/fetcher"
	"test2/internal/models"
)

const (
	USD_TO_RUB = 85
	CNY_TO_RUB = 11
)

func GetLastStatShare(operations []models.Operation) map[string]models.StatsShare {
	stats := make(map[string]models.StatsShare)
	for _, operation := range operations {
		if models.IsShare(operation.Ticker) {
			if _, exists := stats[operation.Ticker]; !exists {
				stats[operation.Ticker] = models.StatsShare{}
			}

			currStats := stats[operation.Ticker]
			if operation.IsBuy {
				currStats.Count += operation.Count
				currStats.CountBuy += operation.Count
				currStats.SumPriceBuy += operation.Price * float64(operation.Count)
				currStats.AvgPriceBuy = math.Round((currStats.SumPriceBuy/float64(currStats.CountBuy))*100) / 100
			} else {
				currStats.Count -= operation.Count
			}

			stats[operation.Ticker] = currStats
		}
	}

	for ticker, stat := range stats {
		currStats := stat
		currStats.LastPrice, _ = fetcher.GetLastPriceShareCached(ticker)
		currStats.SumPriceTotal = math.Round(currStats.LastPrice * float64(currStats.Count))

		currStats.Div, _ = fetcher.GetDivYieldCached(ticker)
		currStats.Div = math.Round(currStats.Div*100) / 100

		currStats.SumDiv = currStats.Div * float64(currStats.Count)
		currStats.SumDiv = math.Round(currStats.SumDiv*100) / 100

		currStats.DivPerc = (currStats.Div / currStats.AvgPriceBuy) * 100
		currStats.DivPerc = math.Round(currStats.DivPerc*100) / 100
		stats[ticker] = currStats
	}

	return stats
}

func GetLastStatBond(operations []models.Operation) map[string]models.StatsBond {
	stats := make(map[string]models.StatsBond)
	for _, operation := range operations {
		if models.IsBond(operation.Ticker) {
			if _, exists := stats[operation.Ticker]; !exists {
				stats[operation.Ticker] = models.StatsBond{}
			}

			currStats := stats[operation.Ticker]
			if operation.IsBuy {
				currStats.Count += operation.Count
				currStats.CountBuy += operation.Count
				currStats.SumPriceBuy += operation.Price * float64(operation.Count)
				currStats.AvgPriceBuy = (currStats.SumPriceBuy / float64(currStats.CountBuy)) / 100
			} else {
				currStats.Count -= operation.Count
			}

			stats[operation.Ticker] = currStats
		}
	}

	for ticker, stat := range stats {
		// todo тут нужно не получать инфу по бумаге если она заэкспирировалась. сейчас цена 0 и поэтому они не учитываются в стате
		currStats := stat
		bondInfo, _ := fetcher.GetLastPriceBondCached(ticker)
		currStats.CouponValue = bondInfo.CouponValue
		if bondInfo.FaceUnit == "USD" {
			bondInfo.LastPrice = bondInfo.LastPrice * USD_TO_RUB
			currStats.AvgPriceBuy = currStats.AvgPriceBuy * USD_TO_RUB
			currStats.CouponValue = currStats.CouponValue * USD_TO_RUB
		}
		if bondInfo.FaceUnit == "CNY" {
			bondInfo.LastPrice = bondInfo.LastPrice * CNY_TO_RUB
			currStats.AvgPriceBuy = currStats.AvgPriceBuy * CNY_TO_RUB
			currStats.CouponValue = currStats.CouponValue * CNY_TO_RUB
		}
		var couponPeriodPerYear int
		if bondInfo.CouponPeriod != 0.0 {
			couponPeriodPerYear = 12 / int(math.Round(bondInfo.CouponPeriod/30.4))
		} else {
			couponPeriodPerYear = 0
		}
		bondInfo.LastPrice = bondInfo.LastPrice / 100

		currStats.AvgPriceBuy *= bondInfo.FaceValue // % -> RUB
		currStats.LastPrice = bondInfo.LastPrice * bondInfo.FaceValue
		currStats.CouponPeriodPerYear = couponPeriodPerYear
		currStats.SumPriceTotal = float64(currStats.Count) * bondInfo.LastPrice * bondInfo.FaceValue
		currStats.Coup2025 = float64(currStats.Count) * currStats.CouponValue * float64(couponPeriodPerYear)

		stats[ticker] = currStats
	}

	return stats
}

func GetLastStatTOM(operations []models.Operation) map[string]models.StatsTOM {
	stats := make(map[string]models.StatsTOM)
	for _, operation := range operations {
		if models.IsCurrency(operation.Ticker) {
			if _, exists := stats[operation.Ticker]; !exists {
				stats[operation.Ticker] = models.StatsTOM{}
			}

			currStats := stats[operation.Ticker]
			if operation.IsBuy {
				currStats.Count += operation.Count
				currStats.CountBuy += operation.Count
				currStats.SumPriceBuy += operation.Price * float64(operation.Count)
				currStats.AvgPriceBuy = math.Round((currStats.SumPriceBuy/float64(currStats.CountBuy))*100) / 100
			} else {
				currStats.Count -= operation.Count
			}

			stats[operation.Ticker] = currStats
		}
	}

	for ticker, stat := range stats {
		currStats := stat
		currStats.LastPrice, _ = fetcher.GetLastPriceTOM(ticker)
		currStats.SumPriceTotal = math.Round(currStats.LastPrice * float64(currStats.Count))

		stats[ticker] = currStats
	}

	return stats
}
