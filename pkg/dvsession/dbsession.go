// package dvsecurity provides server security, including sessions, login, jwt token
// MicroCore Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvsession

import (
	"encoding/json"
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvcrypt"
	"github.com/Dobryvechir/microcore/pkg/dvdbdata"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"log"
	"strconv"
	"time"
)

type SessionMetaInfo struct {
	Db               string   `json:"db"`
	Table            string   `json:"table"`
	Id               string   `json:"id"`
	Data             string   `json:"data"`
	ModTime          string   `json:"modTime"`
	Prefix           string   `json:"prefix"`
	ReadSql          []string `json:"readSql"`
	CreateSql        []string `json:"createSql"`
	UpdateSql        []string `json:"updateSql"`
	DeleteExpireSql  string   `json:"deleteExpireSql"`
	SqlType          int
	lastCleaningTime int64
}

const (
	SessionId                         = "SESSION_ID"
	MODE_SESSION_MAY_BE               = 0
	MODE_SESSION_CREATE_IF_NOT_EXISTS = 1
	MODE_SESSION_MUST_EXIST           = 2
)

var sessionCleanInterval int64 = 3600 * 1000 * 24
var sessionMetaInfo *SessionMetaInfo = nil

func CreateSession(initialData map[string]string, env *dvevaluation.DvObject) (err error) {
	CleanOldSessions(-1)
	if sessionMetaInfo == nil {
		sessionMetaInfo, err = createSessionMetaInfo()
		if err != nil {
			return
		}
	}
	var sessionId string
	for i := 0; i < 1000000; i++ {
		sessionId = dvcrypt.GetRandomUuid()
		ok, err := IsSessionPresent(sessionId)
		if err != nil {
			log.Print(err)
		} else {
			if !ok {
				break
			}
		}
	}
	if sessionId == "" {
		return errors.New("Session cannot be created")
	}
	env.Set(SessionId, sessionId)
	orig, err := WriteToSession(sessionId, initialData)
	if err != nil {
		return err
	}
	updateEnvironmentBySession(env, orig)
	return nil
}

func ReadSession(sessionId string, force bool) (data map[string]string, err error) {
	if sessionMetaInfo == nil {
		sessionMetaInfo, err = createSessionMetaInfo()
		if err != nil {
			return
		}
	}
	sql := sessionMetaInfo.ReadSql[0] + sessionId + sessionMetaInfo.ReadSql[1]
	r, ok, err := dvdbdata.SqlSingleValueByConnectionName(sessionMetaInfo.Db, sql)
	if err != nil {
		return nil, err
	}
	if !ok {
		if force {
			return nil, nil
		}
		return nil, errors.New("Session does not exist or expired")
	}
	return dvcrypt.ConvertStringToMap(r)
}

func ReadSessionByEnvironment(env *dvevaluation.DvObject) (string, map[string]string, error) {
	sessionId := env.GetString(SessionId)
	if sessionId == "" {
		return "", nil, errors.New("No session id is provided")
	}
	data, err := ReadSession(sessionId, false)
	return sessionId, data, err
}

func IsSessionPresent(sessionId string) (bool, error) {
	var err error
	var ok bool
	if sessionMetaInfo == nil {
		sessionMetaInfo, err = createSessionMetaInfo()
		if err != nil {
			return false, err
		}
	}
	sql := sessionMetaInfo.ReadSql[0] + sessionId + sessionMetaInfo.ReadSql[1]
	_, ok, err = dvdbdata.SqlSingleValueByConnectionName(sessionMetaInfo.Db, sql)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func insertToSession(sessionId string, data map[string]string) error {
	s := dvcrypt.ConvertMapToString(data)
	sql := sessionMetaInfo.CreateSql[0] + sessionId + sessionMetaInfo.CreateSql[1] + s + sessionMetaInfo.CreateSql[2]
	err := dvdbdata.SqlUpdateByConnectionName(sessionMetaInfo.Db, sql)
	if err != nil {
		return err
	}
	return nil
}

func WriteToSession(sessionId string, data map[string]string) (orig map[string]string, err error) {
	if sessionMetaInfo == nil {
		sessionMetaInfo, err = createSessionMetaInfo()
		if err != nil {
			return
		}
	}
	orig, err = ReadSession(sessionId, true)
	if err != nil {
		return nil, err
	}
	if orig == nil {
		return data, insertToSession(sessionId, data)
	}
	if data != nil {
		for k, v := range data {
			orig[k] = v
		}
	}
	s := dvcrypt.ConvertMapToString(orig)
	sql := sessionMetaInfo.UpdateSql[0] + s + sessionMetaInfo.UpdateSql[1] + sessionId + sessionMetaInfo.UpdateSql[2]
	err = dvdbdata.SqlUpdateByConnectionName(sessionMetaInfo.Db, sql)
	if err != nil {
		return nil, err
	}
	return orig, nil
}

func updateEnvironmentBySession(env *dvevaluation.DvObject, data map[string]string) {
	if data != nil {
		pref := sessionMetaInfo.Prefix
		for k, v := range data {
			env.Set(pref+k, v)
		}
	}
}

func UpdateBySession(env *dvevaluation.DvObject, mode int) (err error) {
	if sessionMetaInfo == nil {
		sessionMetaInfo, err = createSessionMetaInfo()
		if err != nil {
			return
		}
	}
	sessionId := env.GetString(SessionId)
	if sessionId == "" {
		switch mode {
		case MODE_SESSION_MAY_BE:
			return nil
		case MODE_SESSION_MUST_EXIST:
			return errors.New("Session has not been created before")
		case MODE_SESSION_CREATE_IF_NOT_EXISTS:
			return CreateSession(nil, env)
		}
		return errors.New("Incorrect session mode:" + strconv.Itoa(mode))
	}
	data, err := ReadSession(sessionId, mode != MODE_SESSION_MUST_EXIST)
	if err != nil {
		return err
	}
	updateEnvironmentBySession(env, data)
	return nil
}

func UpdateSessionWithMap(env *dvevaluation.DvObject, info map[string]string) (err error) {
	if sessionMetaInfo == nil {
		sessionMetaInfo, err = createSessionMetaInfo()
		if err != nil {
			return
		}
	}
	sessionId, data, err := ReadSessionByEnvironment(env)
	if err != nil {
		return err
	}
	if data == nil {
		data = make(map[string]string)
	}
	if info != nil {
		for k, v := range info {
			data[k] = v
			env.Set(sessionMetaInfo.Prefix+k, v)
		}
	}
	_, err = WriteToSession(sessionId, data)
	return err
}

func UpdateSessionWithKeyValue(env *dvevaluation.DvObject, key string, value string) (err error) {
	if sessionMetaInfo == nil {
		sessionMetaInfo, err = createSessionMetaInfo()
		if err != nil {
			return
		}
	}
	sessionId, data, err := ReadSessionByEnvironment(env)
	if err != nil {
		return err
	}
	if data == nil {
		data = make(map[string]string)
	}
	data[key] = value
	env.Set(sessionMetaInfo.Prefix+key, value)
	_, err = WriteToSession(sessionId, data)
	return err
}

func CleanOldSessions(updateTime int64) (err error) {
	if sessionMetaInfo == nil {
		sessionMetaInfo, err = createSessionMetaInfo()
		if err != nil {
			return
		}
	}
	check := updateTime < 0
	if updateTime <= 0 {
		updateTime = time.Now().Unix()
	}
	if check && updateTime < sessionMetaInfo.lastCleaningTime+sessionCleanInterval {
		return nil
	}
	sessionMetaInfo.lastCleaningTime = updateTime
	sql := sessionMetaInfo.DeleteExpireSql
	err = dvdbdata.SqlUpdateByConnectionName(sessionMetaInfo.Db, sql)
	if err != nil {
		return err
	}
	return nil
}

func createSessionMetaInfo() (*SessionMetaInfo, error) {
	sessionDescription := dvparser.GlobalProperties["DV_SESSION_INFO"]
	meta := &SessionMetaInfo{}
	var err error
	if sessionDescription != "" {
		err = json.Unmarshal([]byte(sessionDescription), meta)
	}
	if meta.Id == "" {
		meta.Id = "id"
	}
	if meta.Data == "" {
		meta.Data = "data"
	}
	if meta.ModTime == "" {
		meta.ModTime = "modtime"
	}
	if meta.Table == "" {
		meta.Table = "DV_SESSIONS"
	}
	if meta.Db == "" {
		meta.Db = dvdbdata.GetDefaultDbConnection()
	}
	if meta.Prefix == "" {
		meta.Prefix = "SESS_"
	}
	if meta.ReadSql == nil || len(meta.ReadSql) != 2 {
		meta.ReadSql = []string{"SELECT " + meta.Data + " FROM " + meta.Table + " WHERE " + meta.Id + "='", "'"}
	}
	meta.SqlType = dvdbdata.GetConnectionType(meta.Db)
	dateNow := dvdbdata.GetDateNowFunction(meta.SqlType)
	timestampLessDay := dvdbdata.GetTimestampLessDay(meta.SqlType)
	if meta.UpdateSql == nil || len(meta.UpdateSql) != 3 {
		meta.UpdateSql = []string{"UPDATE " + meta.Table + " SET " + meta.ModTime + "=" + dateNow + "," + meta.Data + "='", "' WHERE " + meta.Id + "='", "'"}
	}
	if meta.CreateSql == nil || len(meta.CreateSql) != 3 {
		meta.CreateSql = []string{"INSERT INTO " + meta.Table + "(" + meta.Id + "," + meta.Data + "," + meta.ModTime + ") VALUES('", "','", "'," + dateNow + ")"}
	}
	if len(meta.DeleteExpireSql) != 2 {
		meta.DeleteExpireSql = "DELETE FROM " + meta.Table + " WHERE " + meta.ModTime + "<=" + timestampLessDay
	}
	if err != nil {
		return nil, err
	}
	return meta, nil
}
