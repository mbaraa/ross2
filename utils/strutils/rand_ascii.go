package strutils

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

func GetRandomChar() byte {
	rand.Seed(time.Now().UnixMicro())
	return alphabet[rand.Intn(len(alphabet))]
}

func GetRandomString(length int) string {
	sb := strings.Builder{}

	for i := 0; i < length; i++ {
		sb.WriteByte(GetRandomChar())
	}

	if IsBadWord(sb.String()) {
		return GetRandomString(3)
	}

	return sb.String()
}
