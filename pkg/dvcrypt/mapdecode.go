/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvcrypt

import (
	"crypto/rand"
	"errors"
	"strings"
)

func DecodeByteLine(src []byte) ([]byte, error) {
	n := len(src)
	origLen := n * 5 >> 3
	buf := make([]byte, origLen)
	already := uint(0)
	last := 0
	pos := 0
	for i := 0; i < n; i++ {
		b := SimpleDecodingTable[src[i]&127]
		if b >= 32 {
			if b != 32 {
				return nil, errors.New("Unexpected character " + string(src[i:i+1]))
			}
		}
		last |= int(b) << already
		already += 5
		if already >= 8 {
			already -= 8
			buf[pos] = byte(last & 255)
			pos++
			last >>= 8
		}
	}
	return buf, nil
}

func DecodeLine(s string) (string, error) {
	buf, err := DecodeByteLine([]byte(s))
	return string(buf), err
}

func ConvertStringToMap(s string) (map[string]string, error) {
	lines := strings.Split(s, ";")
	res := make(map[string]string)
	for _, v := range lines {
		pos := strings.Index(v, "=")
		if pos > 0 {
			key := v[:pos]
			value, err := DecodeLine(v[pos+1:])
			if err != nil {
				return res, err
			}
			res[key] = value
		}
	}
	return res, nil
}

func DecryptMap(inKey string, s string, predecode bool) (map[string]string, error) {
	s = strings.Replace(s, "_", "/", -1)
	str, err := DecryptString(inKey, s, predecode)
	if err != nil {
		return nil, err
	}
	return ConvertStringToMap(str)
}

func DecryptString(inKey string, s string, predecode bool) (string, error) {
	data := []byte(s)
	if predecode {
		var er error
		data, er = DecodeByteLine(data)
		if er != nil {
			return "", er
		}
	}
	privKey, err := LoadPrivateKey(inKey)
	if err != nil {
		return "", err
	}
	buf, errDecryptPrivate := privKey.Decrypt(rand.Reader, data, nil)
	return string(buf), errDecryptPrivate
}
