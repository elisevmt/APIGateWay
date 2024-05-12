package utils

import (
	"math"
)

func GetTokenRound(amountToRaw float64, decimal int64) float64 {
	multiplier := math.Pow(10, float64(decimal))
	amountTo := math.Ceil(amountToRaw*multiplier) / multiplier
	return amountTo
}
