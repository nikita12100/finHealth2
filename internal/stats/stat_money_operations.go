package stats

import (
	"log/slog"
	"test2/internal/models"
)

func GetStatMoneyOperations(moneyOperations []models.MoneyOperation) []models.StatsMoneyOperationSnapshoot {
	if len(moneyOperations) < 1 {
		slog.Error("GetReplenishmentPerMonth len < 1")
		return nil
	}

	var stats []models.StatsMoneyOperationSnapshoot
	var curStat models.StatsMoneyOperationSnapshoot
	curStat.Time = moneyOperations[0].Time

	for i, op := range moneyOperations {
		if op.Time.Month() != curStat.Time.Month() {
			stats = append(stats, curStat)
			curStat = models.StatsMoneyOperationSnapshoot{}
			curStat.Time = op.Time
		}

		if op.OperationType == models.Replenishment {
			curStat.Replenishment += op.AmountIn
		}
		if op.OperationType == models.Coupon {
			curStat.Coupon += op.AmountIn
		}
		if op.OperationType == models.Dividends {
			curStat.Dividends += op.AmountIn
		}

		if i == (len(moneyOperations) - 1) {
			stats = append(stats, curStat)
		}
	}

	return stats
}

func GetStatMoneyOperationsSumDivPerTicker(moneyOperations []models.MoneyOperation) map[string]float64 {
	stats := make(map[string]float64)
	for _, operation := range moneyOperations {
		if operation.OperationType == models.Dividends {
			if _, exists := stats[operation.Comment.Ticker]; !exists {
				stats[operation.Comment.Ticker] = 0.0
			}
			stats[operation.Comment.Ticker] += operation.AmountIn
		}
	}

	return stats
}
