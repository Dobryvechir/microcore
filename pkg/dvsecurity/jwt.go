// package dvsecurity provides server security, including sessions, login, jwt token
// MicroCore Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvsecurity

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"strconv"
	"strings"
	"time"
)

func EncodeJwtBase64(s string) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString([]byte(s)), "=")
}

func getJwtHeader(valid int64) string {
	expired := valid + GetCurrentSeconds()
	data := "{\"typ\": \"JWT\",\"alg\":\"HS256\",\"exp\":" + strconv.FormatInt(expired, 10) + "}"
	return EncodeJwtBase64(data)
}

func getJwtContent(data map[string]string) string {
	jsonStr := dvjson.ConvertSimpleStringMapToJson(data, true)
	return EncodeJwtBase64(jsonStr)
}

func GenerateHs256Jwt(claims map[string]string, key string, valid int64) string {
	main := getJwtHeader(valid) + "." + getJwtContent(claims)
	return main + "." + getJwtSign(main, key)
}

func GetCurrentSeconds() int64 {
	return time.Now().UnixNano() / int64(time.Second)
}

func getJwtSign(s string, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(s))
	return EncodeJwtBase64(string(h.Sum(nil)))
}

/************************ Decoding Part **************************************/

func DecodeJwtBase64(s string) (string, error) {
	n := len(s) % 4
	if n > 0 {
		s += strings.Repeat("=", 4-n)
	}
	res, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
