/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjsmaster

import (
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
	"net/url"
	"strings"
)

func window_init() {
	dvevaluation.WindowMaster.Prototype = &dvevaluation.DvVariable{
		Fields: []*dvevaluation.DvVariable{
			{
				Name: []byte("encodeURIComponent"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Window_encodeURIComponent,
				},
			},
			{
				Name: []byte("decodeURIComponent"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Window_decodeURIComponent,
				},
			},
			{
				Name: []byte("encodeURI"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Window_encodeURI,
				},
			},
			{
				Name: []byte("decodeURI"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Window_decodeURI,
				},
			},
		},
		Kind: dvevaluation.FIELD_OBJECT,
	}
}

func Window_encodeURIComponent(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return "", nil
	}
	s := dvevaluation.AnyToString(params[0])
	s = EncodeURIComponent(s)
	return s, nil
}

func EncodeURIComponent(s string) string {
	s = url.QueryEscape(s)
	if strings.Contains(s, "+") {
		s = strings.Replace(s, "+", "%20", -1)
	}
	return s
}

func Window_decodeURIComponent(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return "", nil
	}
	s := dvevaluation.AnyToString(params[0])
	s, err := url.QueryUnescape(s)
	return s, err
}

func Window_encodeURI(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return "", nil
	}
	s := dvevaluation.AnyToString(params[0])
	s = dvtextutils.EncodeURI(s)
	return s, nil
}

func Window_decodeURI(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return "", nil
	}
	s := dvevaluation.AnyToString(params[0])
	s, err := url.PathUnescape(s)
	return s, err
}
