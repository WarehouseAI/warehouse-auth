package encode

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashedPassword(password string) (string, error) {
	pass, err := hex.DecodeString(password)
	if err != nil {
		return "", err
	}
	toCheck := sha256.Sum256(pass)
	return hex.EncodeToString(toCheck[:]), nil
}
