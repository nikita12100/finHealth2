package stats

import (
	"test2/internal/models"
	"time"
)


func GetReplenishmentPerMonth(moneyOperations []models.MoneyOperation) []TimeValue {
	var res []TimeValue
	curSum := 0.0
	var curMonth time.Time
	for i, op := range moneyOperations {
		if op.OperationType == models.Replenishment {
			if curMonth.IsZero() {
				curMonth = op.Time
			}
			if op.Time.Month() == curMonth.Month() {
				curSum += op.AmountIn
			} else {
				res = append(res, TimeValue{Time: curMonth, Value: curSum})
				curSum = 0
				curSum += op.AmountIn
				curMonth = op.Time
			}

		}
		if i == (len(moneyOperations) - 1) {
			res = append(res, TimeValue{Time: curMonth, Value: curSum})
		}
	}

	return res
}