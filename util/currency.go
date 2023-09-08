package util

const (
	CAD = "CAD"
	EUR = "EUR"
	GBP = "GBP"
	USD = "USD"
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case CAD, EUR, GBP, USD:
		return true
	}
	return false
}
