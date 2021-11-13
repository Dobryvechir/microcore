/***********************************************************************
MicroCore Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvcontext

type DvAction struct {
	Name        string             `json:"name"`
	Typ         string             `json:"type"`
	Url         string             `json:"url"`
	Method      string             `json:"method"`
	QueryParams map[string]string  `json:"query"`
	Body        map[string]string  `json:"body"`
	Result      string             `json:"result"`
	ResultMode  string             `json:"mode"`
	InnerParams string             `json:"params"`
	Conditions  map[string]string  `json:"conditions"`
	Validations []*ValidatePattern `json:"validations"`
}
