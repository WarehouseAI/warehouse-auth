package random

import (
	"math/rand"
	"time"

	"github.com/warehouse/auth-service/internal/pkg/utils/str"
)

type (
	Adapter interface {
		RandomString(length int) string
		RandomStringWithTimeNanoSeed(length int) string
		RandomIntn(n int) int
	}

	adapter struct {
	}
)

func NewAdapter() Adapter {
	return &adapter{}
}

func (a *adapter) RandomString(length int) string {
	return str.RandomString(length)
}

func (a *adapter) RandomIntn(n int) int {
	return rand.Intn(n)
}

func (a *adapter) RandomStringWithTimeNanoSeed(length int) string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}
