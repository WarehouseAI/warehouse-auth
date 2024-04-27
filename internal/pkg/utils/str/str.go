package str

import "math/rand"

func RandomString(length int) string {
	res := make([]byte, length)
	for i := 0; i < length; i++ {
		res[i] = Alphabet[rand.Intn(len(Alphabet))]
	}
	return string(res)
}
