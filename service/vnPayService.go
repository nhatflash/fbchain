package service

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"maps"
	"net/url"
	"os"
	"slices"
	"strings"
	"time"

	appErr "github.com/nhatflash/fbchain/error"
	"github.com/shopspring/decimal"
)

func GetVnPayUrl(price decimal.Decimal) (string, error) {
	vnpTmnCode := os.Getenv("VNPAY_TMNCODE")
	vnpHashSecret := os.Getenv("VNPAY_HASHSECRET")
	vnpUrl := os.Getenv("VNPAY_URL")
	vnpBaseUrl := os.Getenv("VNPAY_BASEURL")

	vnpParams := BuildVnpParams(price, vnpTmnCode, vnpBaseUrl)
	url, err := BuildVnPayUrl(vnpParams, vnpHashSecret, vnpUrl)
	if err != nil {
		return "", err
	}
	return url, nil
}


func BuildVnpParams(price decimal.Decimal, vnpTmnCode string, vnpBaseUrl string) map[string]string {

	vnpParams := make(map[string]string)
	vnpParams["vnp_Amount"] = (price.Mul(decimal.NewFromInt(100))).String()
	vnpParams["vnp_Command"] = "pay"
	vnpParams["vnp_CreateDate"] = GetVnPayDate(time.Now())
	vnpParams["vnp_CurrCode"] = "VND"
	vnpParams["vnp_ExpireDate"] = GetVnPayDate(time.Now().Add(15 * time.Minute))
	vnpParams["vnp_IpAddr"] = "192.168.1.1"
	vnpParams["vnp_Locale"] = "vn"
	vnpParams["vnp_OrderInfo"] = "Thanh toan don hang"
	vnpParams["vnp_OrderType"] = "Subscription"
	vnpParams["vnp_ReturnUrl"] = vnpBaseUrl + "/vnpay-return"
	vnpParams["vnp_TmnCode"] = vnpTmnCode
	vnpParams["vnp_TxnRef"] = fmt.Sprintf("%d", time.Now().UnixMilli())
	vnpParams["vnp_Version"] = "2.1.0"

	return vnpParams
}

func BuildVnPayUrl(vnpParams map[string]string, vnpHashSecret string, vnpUrl string) (string, error) {
	total := len(vnpParams)
	count := 0
	var hashedData strings.Builder
	var query strings.Builder
	for _, key := range slices.Sorted(maps.Keys(vnpParams)) {
		count++
		value := vnpParams[key]
		if value != "" {
			hashedData.WriteString(key)
			hashedData.WriteString("=")
			hashedData.WriteString(url.QueryEscape(value))

			query.WriteString(url.QueryEscape(key))
			query.WriteString("=")
			query.WriteString(url.QueryEscape(value))
		}
		if count < total {
			hashedData.WriteString("&")
			query.WriteString("&")
		}
	}
	queryUrl := query.String()
	vnpSecureHash, err := GetHmacSha512(vnpHashSecret, hashedData.String())
	if err != nil {
		return "", err
	}
	queryUrl += ("&vnp_SecureHash=" + vnpSecureHash)
	return (vnpUrl + "?" + queryUrl), nil
}

func GetHmacSha512(key string, data string) (string, error) {
	if key == "" || data == "" {
		return "", appErr.NotFoundError("Key or data is missing")
	}
	keyBytes := []byte(key)
	dataBytes := []byte(data)
	mac := hmac.New(sha512.New, keyBytes)
	mac.Write(dataBytes)
	signatureBytes := mac.Sum(nil)
	return hex.EncodeToString(signatureBytes), nil
} 


// for sha-256 the length will be 32 (256 bit), sha-512 will be 64 (512 bit)
func GenerateSecureKey(length int) ([]byte, error) {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func GetVnPayDate(dateTime time.Time) string {
	return dateTime.Format("20060102150405")
}


