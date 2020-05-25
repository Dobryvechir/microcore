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
	var err error
	request.Params = dvparser.CloneGlobalProperties()
	if options['p'] != "" || options['P'] != "" {
		dvproviders.PlaceProviderReferences(request)
	}
	buffer, err = dvparser.ConvertByteArrayBySpecificProperties(buffer, request.FileName, request.Params, 3, dvparser.CONFIG_PRESERVE_SPACE)
	if err == nil {
		if options['g'] != "" || options['G'] != "" {
			buffer, err = goTemplateHandler(buffer, request)
		}
	}
	request.Error = err
	request.Output = buffer
	request.HandleCommunication()
}
