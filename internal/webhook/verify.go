package webhook

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"

	"github.com/Sn0wo2/go-afdian-api/pkg/payload"
	"github.com/parquet-go/bitpack/unsafecast"
)

var publicKey = func() *rsa.PublicKey {
	block, _ := pem.Decode(unsafecast.Bytes(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwwdaCg1Bt+UKZKs0R54y
lYnuANma49IpgoOwNmk3a0rhg/PQuhUJ0EOZSowIC44l0K3+fqGns3Ygi4AfmEfS
4EKbdk1ahSxu7Zkp2rHMt+R9GarQFQkwSS/5x1dYiHNVMiR8oIXDgjmvxuNes2Cr
8fw9dEF0xNBKdkKgG2qAawcN1nZrdyaKWtPVT9m2Hl0ddOO9thZmVLFOb9NVzgYf
jEgI+KWX6aY19Ka/ghv/L4t1IXmz9pctablN5S0CRWpJW3Cn0k6zSXgjVdKm4uN7
jRlgSRaf/Ind46vMCm3N2sgwxu/g3bnooW+db0iLo13zzuvyn727Q3UDQ0MmZcEW
MQIDAQAB
-----END PUBLIC KEY-----`))
	if block == nil {
		panic("afdian: invalid embedded webhook public key")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)

	pub, ok := key.(*rsa.PublicKey)
	if err != nil || !ok {
		panic("afdian: invalid embedded webhook public key")
	}

	return pub
}()

func Verify(p *payload.WebHook) error {
	if p.Data == nil || p.Data.Order == nil {
		return errors.New("invalid webhook payload")
	}

	sig, err := base64.StdEncoding.DecodeString(p.Data.Sign)
	if err != nil {
		return err
	}

	o := p.Data.Order
	sum := sha256.Sum256(unsafecast.Bytes(o.OutTradeNo + o.UserID + o.PlanID + o.TotalAmount))

	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, sum[:], sig)
}
