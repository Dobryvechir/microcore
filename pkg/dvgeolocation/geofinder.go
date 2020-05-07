/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvgeolocation

import (
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"net/http"
	"strings"
)

func PlaceGeoInfo(r *http.Request, params map[string]string) {
	ip := GetIPAddress(r)
	ipBuf, ok := ReadIP(ip)
	if ok {
		params["USER_IP"] = ip
		countryCode, err := LookupCountryCode(ipBuf)
		if err == nil {
			params["USER_COUNTRY_CODE"] = countryCode
			countryInfo := CountryInfo[countryCode]
			if countryInfo != "" {
				info := strings.Split(countryInfo, ";")
				locale := info[1]
				pos := strings.Index(locale, ",")
				if pos > 0 {
					locale = locale[:pos]
				}
				params["USER_COUNTRY_NAME"] = info[2]
				params["USER_CONTINENT_NAME"] = info[0]
				params["USER_COUNTRY_LOCALE"] = locale
				params["USER_COUNTRY_LOCALES"] = info[1]
			}
		}
	}
}

func GetIPAddress(r *http.Request) string {
	ips := dvparser.ConvertToNonEmptyList(r.Header.Get("X-FORWARDED-FOR"))
	for _, v := range ips {
		if _, ok := ReadIP(v); ok {
			return v
		}
	}
	return GetAddressWithoutPort(r.RemoteAddr)
}

func GetAddressWithoutPort(ip string) string {
	pos := strings.Index(ip, "]")
	pos1 := strings.Index(ip, "[")
	if pos1 >= 0 && pos > pos1 {
		return ip[pos1+1 : pos]
	}
	pos = strings.LastIndex(ip, ":")
	if pos > 0 && ip[pos-1] != ':' {
		return ip[:pos]
	}
	return ip
}
