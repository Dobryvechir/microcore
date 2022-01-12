/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvsession

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvurl"
	"log"
	"strings"
)

type SessionRequestBlock struct {
	Prefix  string
	Request SessionRequest
	Urls    *dvurl.UrlPool
}

type SessionActionBlock struct {
	Prefix  string
	Id      string
	Storage SessionStorage
}

func GetServerSessionProvider(name string, option int, params map[string]string, urls map[string]string, prefix string) dvcontext.ServerSessionProvider {
	request := GetSessionRequest(name, params, option)
	if request == nil {
		return nil
	}
	actionBlock := &SessionRequestBlock{
		Request: request,
		Prefix:  prefix,
	}
	if urls != nil && len(urls) > 0 {
		urlPool := dvurl.GetUrlHandler()
		for k, v := range urls {
			urlPool.RegisterHandlerFunc(k, v)
		}
		actionBlock.Urls = urlPool
	}
	return actionBlock
}

func (action *SessionRequestBlock) GetSessionStorage(ctx *dvcontext.RequestContext, request *dvcontext.SessionActionRequest, sessionId string) (dvcontext.RequestSession, error) {
	if request == nil || ctx == nil {
		return nil, nil
	}
	option := strings.ToLower(request.Option)
	isErrorFatal := strings.Contains(option, "e")
	isCreateOnly := strings.Contains(option, "c")
	isUpdateOnly := strings.Contains(option, "u")
	isLoadAll := strings.Contains(option, "l")
	loadPrefix := "SESSION_"
	if strings.Contains(option, "n") {
		loadPrefix = ""
	}
	if action == nil || ctx == nil || request == nil {
		if !isErrorFatal {
			return nil, nil
		}
		return nil, errors.New("Session was not present on this server")
	}
	prefix := action.Prefix
	if action.Urls != nil {
		ok, res := dvurl.SingleSimplifiedUrlSearch(action.Urls, ctx.Url)
		if ok && res != nil && res.CustomObject != nil {
			pref, ok := res.CustomObject.(string)
			if ok && pref != "" {
				prefix = pref
			}
		}
	}
	if prefix == "" {
		if !isErrorFatal {
			return nil, nil
		}
		return nil, errors.New("Session is not available for this url")
	}
	if isCreateOnly {
		sessionId = PseudoUuid()
	} else if sessionId == "" {
		if !isErrorFatal {
			return nil, nil
		}
		return nil, errors.New("Session id is not provided")
	}
	session, err, isExisting := action.Request.Init(prefix+sessionId, isCreateOnly, isUpdateOnly)
	if err != nil {
		if isErrorFatal {
			return nil, err
		}
		log.Printf("Error session retrieving: %v", err)
		return nil, nil
	}
	if !isExisting && isUpdateOnly && !isCreateOnly {
		err = errors.New("Session is not present " + sessionId)
		if isErrorFatal {
			return nil, err
		}
		log.Printf("Error session retrieving: %v", err)
		return nil, nil
	}
	block := &SessionActionBlock{
		Prefix:  prefix,
		Id:      sessionId,
		Storage: session,
	}
	block.SetSessionVariables(ctx, isLoadAll, loadPrefix)
	return block, nil
}

func (block *SessionActionBlock) SetSessionVariables(ctx *dvcontext.RequestContext, isLoad bool, pref string) {
	if ctx.PrimaryContextEnvironment == nil {
		return
	}
	ctx.PrimaryContextEnvironment.Set("SESSION_ID", block.Id)
	if isLoad {
		data := block.Storage.Values()
		for k, v := range data {
			ctx.PrimaryContextEnvironment.Set(pref+k, v)
		}
	}
}

func PseudoUuid() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	uuid = fmt.Sprintf("%8X-%4X-%4X-%4X-%12X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return
}

func (block *SessionActionBlock) SetItem(key string, value interface{}) {
	if block != nil {
		block.Storage.SetItem(key, value)
	}
}

func (block *SessionActionBlock) GetItem(key string) interface{} {
	if block != nil {
		return block.Storage.GetItem(key)
	}
	return nil
}

func (block *SessionActionBlock) RemoveItem(key string) {
	if block != nil {
		block.Storage.RemoveItem(key)
	}
}
func (block *SessionActionBlock) Clear() {
	if block != nil {
		block.Storage.Clear()
	}
}
func (block *SessionActionBlock) Keys() []string {
	if block != nil {
		return block.Storage.Keys()
	}
	return nil
}

func (block *SessionActionBlock) Values() map[string]interface{} {
	if block != nil {
		return block.Storage.Values()
	}
	return nil
}

func (block *SessionActionBlock) GetId() string {
	if block != nil {
		return block.Id
	}
	return ""
}
