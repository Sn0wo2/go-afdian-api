package afdian

import (
	"crypto"
	"crypto/md5"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"github.com/Sn0wo2/go-afdian-api/internal/helper"
	"github.com/Sn0wo2/go-afdian-api/pkg/payload"
)

// publicKeyPEM afdian public key
const publicKeyPEM = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwwdaCg1Bt+UKZKs0R54y
lYnuANma49IpgoOwNmk3a0rhg/PQuhUJ0EOZSowIC44l0K3+fqGns3Ygi4AfmEfS
4EKbdk1ahSxu7Zkp2rHMt+R9GarQFQkwSS/5x1dYiHNVMiR8oIXDgjmvxuNes2Cr
8fw9dEF0xNBKdkKgG2qAawcN1nZrdyaKWtPVT9m2Hl0ddOO9thZmVLFOb9NVzgYf
jEgI+KWX6aY19Ka/ghv/L4t1IXmz9pctablN5S0CRWpJW3Cn0k6zSXgjVdKm4uN7
jRlgSRaf/Ind46vMCm3N2sgwxu/g3bnooW+db0iLo13zzuvyn727Q3UDQ0MmZcEW
MQIDAQAB
-----END PUBLIC KEY-----`

// WebHookSignVerify
// if the signature is valid, return nil
func WebHookSignVerify(p *payload.WebHook) error {
	if p.Data.Sign == "" {
		return fmt.Errorf("sign is empty")
	}

	sigBytes, err := base64.StdEncoding.DecodeString(p.Data.Sign)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(helper.StringToBytes(publicKeyPEM))
	if block == nil || block.Type != "PUBLIC KEY" {
		return fmt.Errorf("invalid public key")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	pubKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("invalid public key")
	}

	order := p.Data.Order

	hashed := crypto.SHA256.New()
	hashed.Write(helper.StringToBytes(order.OutTradeNo +
		order.UserId +
		order.PlanId +
		order.TotalAmount))

	return rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed.Sum(nil), sigBytes)
}

// APISignParams performs MD5 signature on API parameters
//
// https://afdian.com/p/9c65d9cc617011ed81c352540025c377
func APISignParams(userID, apiToken string, params []byte, ts int64) (string, error) {
	return fmt.Sprintf("%x", md5.Sum(helper.StringToBytes(fmt.Sprintf("%sparams%sts%duser_id%s", apiToken, helper.BytesToString(params), ts, userID)))), nil
}
