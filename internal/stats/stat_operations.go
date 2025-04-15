package stats

import (
	"math"
	"test2/internal/fetcher"
	"test2/internal/models"
	"time"
)

const (
	USD_TO_RUB = 90
	CNY_TO_RUB = 11
)

type TimeValue struct {
	Time  time.Time
	Value float64
}

func GetCountPerTicker(operations []models.Operation) map[string]int {
	countMap := make(map[string]int)
	for _, operation := range operations {
		if operation.IsBuy {
			countMap[operation.Ticker] += operation.Count
		} else {
			countMap[operation.Ticker] -= operation.Count
		}
	}
	return countMap
}

func GetAvgBuyPricePerTicker(operations []models.Operation) map[string]float64 {
	sumPrice := make(map[string]float64)
	avgCount := make(map[string]int)
	for _, operation := range operations {
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
		currStats.LastPrice, _ = fetcher.GetLastPriceShare(ticker)
		currStats.SumPriceTotal = math.Round(currStats.LastPrice * float64(currStats.Count))
		currStats.Weight = math.Round((currStats.SumPriceTotal/float64(models.WEIGHT_NORM))*100) / 100

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
