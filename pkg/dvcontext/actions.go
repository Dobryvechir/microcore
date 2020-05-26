/***********************************************************************
MicroCore Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvcontext

type DvAction struct {
	Name        string            `json:"name"`
	Typ         string            `json:"type"`
	Url         string            `json:"url"`
	Method      string            `json:"method"`
	QueryParams map[string]string `json:"query"`
	Body        string            `json:"body"`
	Result      string            `json:"result"`
	Conditions  map[string]string `json:"conditions"`
}
