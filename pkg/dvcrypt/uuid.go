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
