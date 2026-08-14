package security

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"math/big"
)

func GenerateOTP() (string, error) {
	n, err := rand.Int(
		rand.Reader,
		big.NewInt(1000000),
	)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", n.Int64()), nil
}

func HashOTP(otp string) string {
	hash := sha256.Sum256([]byte(otp))

	return hex.EncodeToString(hash[:])
}

func VerifyOTP(otp string, otpHash string) bool {
	inputHash := HashOTP(otp)

	if subtle.ConstantTimeCompare([]byte(inputHash), []byte(otpHash)) != 1 {
		return false
	}

	return true

}
