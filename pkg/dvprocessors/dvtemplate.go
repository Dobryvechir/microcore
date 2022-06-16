/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvprocessors

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"github.com/Dobryvechir/microcore/pkg/dvproviders"
)

func dvtemplateProcessing(request *dvcontext.RequestContext, buffer []byte, options map[byte]string) {
	paramsCloned := false
	if options['p'] != "" || options['P'] != "" {
		// TO DO: it was supposed to be done before reading the file
		request.Params = dvparser.CloneGlobalProperties()
		paramsCloned = true
		dvproviders.PlaceProviderReferences(request)
	}
	var err error
	if options['g'] != "" || options['G'] != "" {
		if !paramsCloned {
			request.Params = dvparser.CloneGlobalProperties()
		}
		buffer, err = goTemplateHandler(buffer, request)
	}
	request.Error = err
	request.Output = buffer
	request.HandleCommunication()
}
