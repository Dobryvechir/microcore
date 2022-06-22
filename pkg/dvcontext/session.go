/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvcontext

/****
SessionActionRequest options:
c - create session
u - update session only
e - error if no session
l - load all session variables with SESSION_ prefix (default) or without
n - specifies option "l" must be without prefix

id is a place where to find the session id
*/
const (
	ServerSessionStoringKey = "___SERVER__SESSION__STORING__KEY___"
)

type SessionActionRequest struct {
	Prefix string   `json:"prefix"`
	Option string   `json:"option"`
	Id     []string `json:"id"`
}

type RequestSession interface {
	SetItem(key string, value interface{})
	GetItem(key string) interface{}
	RemoveItem(key string)
	Clear()
	Keys() []string
	Values() map[string]interface{}
	GetId() string
}

type ServerSessionProvider interface {
	GetSessionStorage(ctx *RequestContext, request *SessionActionRequest, sessionId string) (RequestSession, error)
}
