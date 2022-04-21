package hash

import (
	"crypto/sha512"
	"fmt"
)

func Sha512(value string) string {
	hash := sha512.New()
	hash.Write([]byte(value))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
