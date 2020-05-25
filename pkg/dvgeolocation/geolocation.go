/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvgeolocation

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvproviders"
)

var dictionaryPath string

func provideGeolocation(request *dvcontext.RequestContext) bool {
	PlaceGeoInfo(request.Reader, request.Params)
	return true
}

var geolocationRegistrationConfig *dvproviders.RegistrationConfig = &dvproviders.RegistrationConfig{
	Name:              "geolocation",
	EndPointHandler:   provideGeolocation,
	GlobalInitHandler: globalInitializer,
	ServerInitHandler: serverInitializer,
}

func globalInitializer(data map[string]string) error {
	for k, v := range data {
		switch k {
		case "dictionaryPath":
			dictionaryPath = v
		default:
			return errors.New("Global parameter " + k + " is not supposed for geolocation, only dictionaryPath")
		}
	}
	if dictionaryPath == "" {
		return errors.New("geolocation global dictionaryPath has not been initialized")
	}
	err := LookupDbInit(dictionaryPath+"/ip4.bin", dictionaryPath+"/ip6.bin", dictionaryPath+"/countries.properties", 7)
	if err != nil {
		return err
	}
	return nil
}

func serverInitializer(params []string) (map[string]string, error) {
	if dictionaryPath == "" {
		return nil, errors.New("specify providers:{geolocation:{ dictionaryPath: absolute or relative file path}} in the main config")
	}
	return nil, nil
}

func Init() bool {
	dvproviders.RegisterProvider(geolocationRegistrationConfig, false)
	return true
}

var inited = Init()
