package secure

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"time"
)

func CalcSignature(secret string, message string) string {
	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

func CalcInternalId(requestId int64) string {
	hasher := sha512.New()
	hasher.Write([]byte(fmt.Sprintf("%d_%s", requestId, time.Now().Format(time.RFC3339))))
	return hex.EncodeToString(hasher.Sum(nil))
}
