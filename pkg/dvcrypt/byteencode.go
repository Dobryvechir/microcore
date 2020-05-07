/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvcrypt

import (
	"github.com/Dobryvechir/microcore/pkg/dvparser"
)

var SimpleEncodingTable = []byte{'D', 'o', 'b', 'r', 'y', 'v', 'e', 'c', 'h', 'i', 'R', 'V', 'O', 'l', 'd', 'I',
	'm', 'Y', 'B', 'a', 'n', 'E', 'L', 'z', 'H', 'N', 'A', 'S', 'p', 'f', 'g', 'G'}

func CalculateDecodingTable(encodingTable []byte) (res []byte) {
	n := len(encodingTable)
	if n != 32 {
		panic("encoding table must have 32 bytes")
	}
	p := 128
	res = make([]byte, p)
	for i := 0; i < p; i++ {
		res[i] = byte(127)
	}
	for i := 0; i < n; i++ {
		c := encodingTable[i]
		if res[c] != 127 {
			panic("Encoding table must not have repetitions")
		}
		res[c] = byte(i)
	}
	return
}

var SimpleDecodingTable = CalculateDecodingTable(SimpleEncodingTable)

func SetEncodingPattern(pattern string) {
	SimpleEncodingTable = []byte(pattern)
	SimpleDecodingTable = CalculateDecodingTable(SimpleEncodingTable)
}

func InitBySettings(parameters map[string]string, functionPool map[string]interface{}) {
	if parameters[dvparser.MicroCorePrexix+"DVCRYPT_ENCODING_PATTERN"] != "" {
		SetEncodingPattern(parameters[dvparser.MicroCorePrexix+"DVCRYPT_ENCODING_PATTERN"])
	}
	if parameters[dvparser.MicroCorePrexix+"DVCRYPT_KEY_FOLDER"] != "" {
		SetKeyFolder(parameters[dvparser.MicroCorePrexix+"DVCRYPT_KEY_FOLDER"])
	}
}

var initedBySettings bool = dvparser.RegisterInitBySettingFunc(InitBySettings)
