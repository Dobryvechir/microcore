/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvsession

type SessionActionBlock struct {
	Prefix    string
	SessionId string
	Request   SessionRequest
	Storage   SessionStorage
}

