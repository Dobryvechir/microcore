/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvcrypt

import (
	"crypto/rand"
	"fmt"
	"log"
)

func GetUuidByByteArray(b []byte) string {
	uuid := fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

func IsValidUUID(uuid string) bool {
	n := len(uuid)
	if n != 36 {
		return false
	}
	for i := 0; i < n; i++ {
		c := uuid[i]
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if c != '-' {
				return false
			}
		} else if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F') {
			return false
		}
	}
	return true
}

func PutUInt64ToByteArray(b []byte, a uint64) {
	for i := 0; i < 8; i++ {
		b[i] = byte(a & 255)
		a <<= 8
	}
}

func GetUuidByLong(a1 uint64, a2 uint64) string {
	b := make([]byte, 16)
	PutUInt64ToByteArray(b, a1)
	PutUInt64ToByteArray(b[8:], a2)
	return GetUuidByByteArray(b)
}

func GetRandomUuid() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Print(err)
	}
	return GetUuidByByteArray(b)
}
