package random

import (
	"math/rand"
	"time"
)

func GetRandomAlias(size int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	const alph = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"0123456789"

	b := make([]byte, size)
	for i := range b {
		b[i] = alph[rnd.Intn(len(alph))]
	}
	return string(b)
}
