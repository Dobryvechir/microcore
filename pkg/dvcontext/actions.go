/***********************************************************************
MicroCore Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvcontext


type DvAction struct {
	Name        string                `json:"name"`
	Typ         string                `json:"type"`
	Url         string                `json:"url"`
	Method      string                `json:"method"`
	QueryParams map[string]string     `json:"query"`
	Body        map[string]string     `json:"body"`
	Result      string                `json:"result"`
	ResultMode  string                `json:"mode"`
	Definitions map[string]string     `json:"definitions"`
	InnerParams string                `json:"params"`
	Conditions  map[string]string     `json:"conditions"`
	Validations []*ValidatePattern    `json:"validations"`
	ErrorPolicy string                `json:"error-policy"`
	Session     *SessionActionRequest `json:"session"`
}

func (action *DvAction) CloneFrom(other *DvAction) {
	action.Name = other.Name
	action.Typ = other.Typ
	action.Url = other.Url
	action.Method = other.Method
	action.QueryParams = other.QueryParams
	action.Body = other.Body
	action.Result = other.Result
	action.ResultMode = other.ResultMode
	action.Definitions = other.Definitions
	action.InnerParams = other.InnerParams
	action.Conditions = other.Conditions
	action.Validations = other.Validations
	action.ErrorPolicy = other.ErrorPolicy
	action.Session = other.Session
}
