package utils

import (
	"APIGateWay/pkg/secure"
	"github.com/google/uuid"
	"math/rand"
	"strings"
)

func GenerateTrackerID(clientIdString string) string {
	trackerId := secure.CalcSignature(clientIdString, uuid.New().String())
	return trackerId
}

const letterBytes = "abcdefghijklmnopqrstuvwxyz1234567890"
const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

var src = rand.NewSource(int64(uuid.New().ID()))

func GenerateKey(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}
