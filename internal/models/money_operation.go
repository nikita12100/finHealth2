package models

import "time"

type OperationType int

const (
	Unknown       OperationType = iota // 0
	Replenishment                      // 1
	BuyOrSell                          // 2
	Commision                          // 3
	DFP                                // 4
	DVP                                // 5
	Coupon                             // 6
	BondExpire                         // 7
	Tax                                // 8
	TaxYear                            // 9
	TaxDiv                             // 10
	Dividends                          // 11
	Withdraw                           // 12
	Repo                               // 13
	Redemption                         // 14
)

type CommentMoneyOperation struct {
	Ticker string `json:"ticker"`
	Count  int    `json:"count"`
}

type MoneyOperation struct {
	Time          time.Time             `json:"time"`
	OperationType OperationType         `json:"operation_type"`
	AmountIn      float64               `json:"amount_in"`
	AmountOut     float64               `json:"amount_out"`
	CommentRaw    string                `json:"comment_raw"`
	Comment       CommentMoneyOperation `json:"comment"`
}

type StatsMoneyOperationSnapshoot struct {
	Time          time.Time
	Replenishment float64
	Coupon        float64
	Dividends     float64
}
