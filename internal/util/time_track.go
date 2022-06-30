package util

import (
	"log"
	"math/big"
	"time"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %dms", name, elapsed.Nanoseconds()/1000)
}

func Factorial(n *big.Int) (result *big.Int) {
	defer timeTrack(time.Now(), "factorial")
	result = big.NewInt(1)
	var one big.Int
	one.SetInt64(1)
	for n.Cmp(&big.Int{}) == 1 {
		result.Mul(result, n)
		n.Sub(n, &one)
	}
	return n
}
