package trades

import (
	"strings"
	"time"
)

type TransactionType string

const (
	TransactionTypeUnknown     TransactionType = "unknown"
	TransactionTypePurchase    TransactionType = "purchase"
	TransactionTypePartialSale TransactionType = "partial sale"
	TransactionTypeFullSale    TransactionType = "full sales"
)

func StringToTransactionType(text string) TransactionType {
	text = strings.ToLower(text)
	text = strings.Trim(text, " ")
	text = strings.ReplaceAll(text, "_", " ")

	return TransactionType(text)
}

func ParseDate(text string) *time.Time {

	if date, err := time.Parse("2006-01-02", text); err == nil {
		return &date
	}

	if date, err := time.Parse("02/2006", text); err == nil {
		return &date
	}

	if date, err := time.Parse("01/02/2006", text); err == nil {
		return &date
	}

	return nil

}
