/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvgeolocation

import (
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"io/ioutil"
)

type dblib struct {
	src    string
	info   []byte
	amount int
	size   int
}

var ip4DbLib *dblib = &dblib{src: "ip4.bin", size: 4}
var ip6DbLib *dblib = &dblib{src: "ip6.bin", size: 16}
var CountryInfo map[string]string

func LookupDbInit(srcIP4, srcIP6 string, countryPropertiesSrc string, activate int) error {
	ip4DbLib.src = srcIP4
	ip6DbLib.src = srcIP6
	var err error
	if activate&1 != 0 {
		err = loadDbLib(ip4DbLib)
		if err != nil {
			return err
		}
	}
	if activate&2 != 0 {
		err = loadDbLib(ip6DbLib)
		if err != nil {
			return err
		}
	}
	if activate&4 != 0 {
		CountryInfo, err = dvparser.LoadSimpleMapFromFile(countryPropertiesSrc, '=')
		if err != nil {
			return err
		}
	}
	return nil
}

func loadDbLib(lib *dblib) error {
	buf, err := ioutil.ReadFile(lib.src)
	if err != nil {
		return err
	}
	lib.info = buf
	lib.amount = len(buf) / (lib.size + 2)
	return nil
}

func compareBuffers(buf1 []byte, buf2 []byte, size int) int {
	for i := size - 1; i >= 0; i-- {
		v := int(buf1[i]) - int(buf2[i])
		if v != 0 {
			return v
		}
	}
	return 0
}

func findInBuffer(pool []byte, ip []byte, size int, itemSize int, start int, end int) string {
	mid := (start + end) >> 1
	midPos := mid * itemSize
	if start == mid {
		endPos := end * itemSize
		comparison := compareBuffers(pool[endPos:endPos+size], ip, size)
		if comparison <= 0 {
			midPos = endPos
		}
		return string(pool[midPos+size : midPos+size+2])
	}
	comparison := compareBuffers(pool[midPos:midPos+size], ip, size)
	if comparison >= 0 {
		return findInBuffer(pool, ip, size, itemSize, start, mid)
	} else {
		return findInBuffer(pool, ip, size, itemSize, mid, end)
	}
}

func findInDbLib(lib *dblib, ip []byte) (string, error) {
	if lib.info == nil {
		err := loadDbLib(lib)
		if err != nil {
			return "", err
		}
	}
	return findInBuffer(lib.info, ip, lib.size, lib.size+2, 0, lib.amount-1), nil
}

func LookupCountryCode(ip []byte) (string, error) {
	if len(ip) < 16 || IsIP6UsedAsIP4Container(ip) {
		return findInDbLib(ip4DbLib, ip)
	}
	return findInDbLib(ip6DbLib, ip)
}
