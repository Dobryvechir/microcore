/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvcrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
)

func EncodeByteLine(src []byte) []byte {
	n := len(src)
	origLen := n * 8 / 5
	if n%5 != 0 {
		origLen++
	}
	buf := make([]byte, origLen)
	already := uint(0)
	pos := 0
	last := 0
	for i := 0; i < n; i++ {
		b := int(src[i])
		last |= int(b) << already
		already += 3
		buf[pos] = SimpleEncodingTable[last&0x1f]
		pos++
		last >>= 5
		if already >= 5 {
			already -= 5
			buf[pos] = SimpleEncodingTable[last&0x1f]
			pos++
			last >>= 5
		}
	}
	if already > 0 {
		buf[pos] = SimpleEncodingTable[last&0x1f]
	}
	return buf
}

func EncodeLine(s string) string {
	return string(EncodeByteLine([]byte(s)))
}

func ConvertMapToString(src map[string]string) string {
	buf := make([]byte, 0, 1024)
	if src != nil {
		for k, v := range src {
			buf = append(buf, []byte(k)...)
			buf = append(buf, '=')
			buf = append(buf, EncodeByteLine([]byte(v))...)
			buf = append(buf, ';')
		}
	}
	return string(buf)
}

func EncryptString(inKey string, s string, final bool) (string, error) {
	pubKey, err := LoadPublicKey(inKey)
	if err != nil {
		return "", errors.New("No public key " + inKey + ": " + err.Error())
	}
	buf, errEncrypt := rsa.EncryptPKCS1v15(rand.Reader, pubKey, []byte(s))
	if errEncrypt != nil {
		return "", errEncrypt
	}
	if !final {
		return string(buf), nil
	}
	return string(EncodeByteLine(buf)), nil
}

func EncryptMap(inKey string, src map[string]string, final bool) (string, error) {
	data := ConvertMapToString(src)
	return EncryptString(inKey, data, final)
}
