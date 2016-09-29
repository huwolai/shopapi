package comm

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"time"
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"gitlab.qiyunxin.com/tangtao/utils/config"
)


type JWTCouponBackend struct {
	PublicKey  *rsa.PublicKey
}

const (
	expireOffset  = 600 //10分钟
)

var couponBackendInstance *JWTCouponBackend = nil

func InitJWTCouponBackend() *JWTCouponBackend {
	if couponBackendInstance == nil {
		couponBackendInstance = &JWTCouponBackend{
			PublicKey:  getCouponPublicKey(),
		}
	}

	return couponBackendInstance
}





func (backend *JWTCouponBackend) getCouponTokenRemainingValidity(timestamp interface{}) int {
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remainer := tm.Sub(time.Now())
		if remainer > 0 {
			return int(remainer.Seconds() + expireOffset)
		}
	}
	return expireOffset
}

func (backend *JWTCouponBackend)  FetchCouponToken(couponToken string) (token *jwt.Token,err error){
	token, err =jwt.Parse(couponToken, func(token *jwt.Token)(interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return backend.PublicKey, nil
	})
	return token,err;
}




func getCouponPublicKey() *rsa.PublicKey {
	publicKeyFile, err := os.Open(config.GetValue("coupon_publickey_path").ToString())
	if err != nil {
		panic(err)
	}
	pemfileinfo, _ := publicKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(pembytes)
	data, _ := pem.Decode([]byte(pembytes))
	publicKeyFile.Close()
	publicKeyImported, err := x509.ParsePKIXPublicKey(data.Bytes)
	if err != nil {
		panic(err)
	}

	rsaPub, ok := publicKeyImported.(*rsa.PublicKey)
	if !ok {
		panic(err)
	}
	return rsaPub
}
