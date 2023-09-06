package base

import "spyrosoft-recruitment-task/marshal"

type ExchangeRate struct {
	No            string             `json:"no"`
	EffectiveDate marshal.CustomTime `json:"effectiveDate"`
	Mid           float64            `json:"mid"`
}

type ExchangeRatesSummary struct {
	Table    string          `json:"table"`
	Currency string          `json:"currency"`
	Code     string          `json:"code"`
	Rates    []*ExchangeRate `json:"rates"`
}
