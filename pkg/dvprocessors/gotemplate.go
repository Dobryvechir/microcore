/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvprocessors

import (
	"bytes"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"text/template"
)

func goTemplateHandler(data []byte, request *dvcontext.RequestContext) ([]byte, error) {
	tpl, err1 := template.New(request.Url).Parse(string(data))
	if err1 != nil {
		return data, err1
	}
	b := bytes.Buffer{}
	err := tpl.Execute(&b, request.Params)
	return b.Bytes(), err
}
