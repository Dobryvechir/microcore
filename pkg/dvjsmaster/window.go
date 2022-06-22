/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjsmaster

import (
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
	"net/url"
	"strings"
)

var console_fns = []*dvevaluation.DvVariable{
	{
		Name: []byte("debug"),
		Kind: dvevaluation.FIELD_FUNCTION,
		Extra: &dvevaluation.DvFunction{
			Fn: Window_console_log,
		},
	},
	{
		Name: []byte("error"),
		Kind: dvevaluation.FIELD_FUNCTION,
		Extra: &dvevaluation.DvFunction{
			Fn: Window_console_log,
		},
	},
	{
		Name: []byte("info"),
		Kind: dvevaluation.FIELD_FUNCTION,
		Extra: &dvevaluation.DvFunction{
			Fn: Window_console_log,
		},
	},
	{
		Name: []byte("log"),
		Kind: dvevaluation.FIELD_FUNCTION,
		Extra: &dvevaluation.DvFunction{
			Fn: Window_console_log,
		},
	},
	{
		Name: []byte("trace"),
		Kind: dvevaluation.FIELD_FUNCTION,
		Extra: &dvevaluation.DvFunction{
			Fn: Window_console_log,
		},
	},
	{
		Name: []byte("warn"),
		Kind: dvevaluation.FIELD_FUNCTION,
		Extra: &dvevaluation.DvFunction{
			Fn: Window_console_log,
		},
	},
	{
		Name: []byte("timeLog"),
		Kind: dvevaluation.FIELD_FUNCTION,
		Extra: &dvevaluation.DvFunction{
			Fn: Window_console_log,
		},
	},
	{
		Name: []byte("table"),
		Kind: dvevaluation.FIELD_FUNCTION,
		Extra: &dvevaluation.DvFunction{
			Fn: Window_console_log,
		},
	},
	{
		Name: []byte("clear"),
		Kind: dvevaluation.FIELD_FUNCTION,
		Extra: &dvevaluation.DvFunction{
			Fn: Window_console_log,
		},
	},
}

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
			{
				Name:   []byte("console"),
				Kind:   dvevaluation.FIELD_OBJECT,
				Fields: console_fns,
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

func Window_console_log(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	r := ""
	for i := 0; i < n; i++ {
		s := dvevaluation.AnyToString(params[i])
		if i == 0 {
			r = s
		} else {
			r = r + " " + s
		}
	}
	dvlog.Println(r, r)
	return nil, nil
}
