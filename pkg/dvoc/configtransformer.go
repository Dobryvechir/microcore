/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvoc

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
)

type ConfigTransformer func(*dvevaluation.DvVariable) (*dvevaluation.DvVariable, error)

func ConfigTransformerDc(src *dvevaluation.DvVariable) (dst *dvevaluation.DvVariable, err error) {
	dst = dvjson.CreateDvFieldInfoObject()
	s := dvjson.ConvertDvFieldInfoArrayIntoMap(src.Fields)
	dst.AddStringField("kind", "DeploymentConfig")
	dst.AddStringField("apiVersion", "v1")
	if s["metadata"] == nil {
		err = errors.New("Required metadata not found")
		return
	}
	dst.AddField(s["metadata"])
	if s["spec"] == nil {
		err = errors.New("Required spec not found")
		return
	}
	dst.AddField(s["spec"])
	return
}

func ConfigTransformerSvc(src *dvevaluation.DvVariable) (dst *dvevaluation.DvVariable, err error) {
	dst = dvjson.CreateDvFieldInfoObject()
	s := dvjson.ConvertDvFieldInfoArrayIntoMap(src.Fields)
	dst.AddStringField("kind", "Service")
	dst.AddStringField("apiVersion", "v1")
	if s["metadata"] == nil {
		err = errors.New("Required metadata not found")
		return
	}
	dst.AddField(s["metadata"])
	if s["spec"] == nil {
		err = errors.New("Required spec not found")
		return
	}
	dst.AddField(s["spec"])
	return
}

func ConfigTransformerConfigMap(src *dvevaluation.DvVariable) (dst *dvevaluation.DvVariable, err error) {
	dst = dvjson.CreateDvFieldInfoObject()
	s := dvjson.ConvertDvFieldInfoArrayIntoMap(src.Fields)
	dst.AddStringField("kind", "ConfigMap")
	dst.AddStringField("apiVersion", "v1")
	if s["metadata"] == nil {
		err = errors.New("Required metadata not found")
		return
	}
	dst.AddField(s["metadata"])
	if s["data"] == nil {
		err = errors.New("Required data not found")
		return
	}
	dst.AddField(s["data"])
	return
}

func ConfigTransformerRoute(src *dvevaluation.DvVariable) (dst *dvevaluation.DvVariable, err error) {
	dst = dvjson.CreateDvFieldInfoObject()
	s := dvjson.ConvertDvFieldInfoArrayIntoMap(src.Fields)
	dst.AddStringField("kind", "Route")
	dst.AddStringField("apiVersion", "v1")
	if s["metadata"] == nil {
		err = errors.New("Required metadata not found")
		return
	}
	dst.AddField(s["metadata"])
	if s["spec"] == nil {
		err = errors.New("Required spec not found")
		return
	}
	dst.AddField(s["spec"])
	return
}
