package random

import (
	"crypto/rand"
	"math"
	"math/big"
)

type RandomGeneratorEngine struct{}

func NewRandomGeneratorEngine() *RandomGeneratorEngine {
	return &RandomGeneratorEngine{}
}

func (rge *RandomGeneratorEngine) GenerateRandomNumber(digits int) int {
	if digits <= 0 {
		return 0
	}

	min := int64(math.Pow10(digits - 1))
	max := int64(math.Pow10(digits)) - 1
	diff := max - min + 1

	n, err := rand.Int(rand.Reader, big.NewInt(diff))
	if err != nil {
		return 0
	}

	return int(n.Int64() + min)
}
