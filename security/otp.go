package security


import (
	"crypto/rand"
	"math/big"
	"crypto/subtle"
)

const Seed = "0123456789"

func GenerateOTPCode(length int) (string, error) {
	otp := make([]byte, length)
	for i := range otp {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(Seed))))
		if err != nil {
			return "", err
		}
		otp[i] = Seed[num.Int64()]
	}
	return string(otp), nil
}


func VerifyOTPCode(inputOTP string, actualOTP string) bool  {
	inputBytes := []byte(inputOTP)
	actualBytes := []byte(actualOTP)

	if len(inputBytes) != len(actualBytes) {
		return false
	}

	match := subtle.ConstantTimeCompare(inputBytes, actualBytes)
	return match == 1
}