package security

import (
	"url-shortener/internal/lib/random"
)

const tokenLength = 16

func GenerateToken() string {
	return random.String(tokenLength)
}
