/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvaction

import (
	"bytes"
	"encoding/json"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvnet"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	Authorization = "Authorization"
)

const (
	M2MTokenPrefix          = "M2M_TOKEN_"
	M2MTokenPath            = "M2MTOKEN_PATH"
	M2MAuthorizationRequest = "AuthorizationM2MService"
)

func convertToHeader(list []string) (res map[string]string) {
	n := len(list)
	res = make(map[string]string)
	for i := 0; i < n; i++ {
		s := strings.TrimSpace(list[i])
		if s == "" {
			continue
		}
		pos := strings.Index(s, "=")
		if pos <= 0 {
			dvlog.PrintfError("Incorrect header %s (no =)", s)
			continue
		}
		k := strings.TrimSpace(s[:pos])
		v := strings.TrimSpace(s[pos+1:])
		res[k] = v
	}
	return
}

func processNetInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	params := dvtextutils.ConvertToNonEmptyList(command)
	url := params[0]
	headers := convertToHeader(params[1:])
	return []interface{}{url, headers}, true
}

func processNetRun(data []interface{}) bool {
	url := data[0].(string)
	headers := data[1].(map[string]string)
	_, err := dvnet.NewRequest("GET", url, "", headers, dvnet.AveragePersistentOptions)
	return err == nil
}

type SmartNetConfigTemplate struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

type SmartNetConfig struct {
	Url         string                 `json:"url"`
	Method      string                 `json:"method"`
	Headers     string                 `json:"headers"`
	Template    SmartNetConfigTemplate `json:"template"`
	ContentType string                 `json:"type"`
	Result      string                 `json:"result"`
	Body        string                 `json:"body"`
}

func DefaultInitWithObject(command string, result interface{}, env *dvevaluation.DvObject) bool {
	cmd := strings.TrimSpace(command[strings.Index(command, ":")+1:])
	if cmd == "" {
		log.Printf("Empty parameters in %s", command)
		return false
	}
	cmdDat := []byte(cmd)
	if cmd[0] != '{' || cmd[len(cmd)-1] != '}' {
		if cmd[0] == '@' && len(cmd) > 1 {
			dat, err := ioutil.ReadFile(cmd[1:])
			if err != nil {
				log.Printf("Wrong file name in %s %v", command, err)
				return false
			}
			dat = bytes.TrimSpace(dat)
			if len(dat) < 2 || dat[0] != '{' || dat[len(dat)-1] != '}' {
				log.Printf("Empty file in %s", command)
				return false
			}
			cmdDat = dat
		} else {
			log.Printf("Empty parameters in %s", command)
			return false
		}
	}
	err := json.Unmarshal(cmdDat, result)
	if err != nil {
		log.Printf("Error converting parameters: %v in %s", err, command)
		return false
	}
	return true
}

func SmartNetInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &SmartNetConfig{}
	if !DefaultInitWithObject(command, config, GetEnvironment(ctx)) {
		return nil, false
	}
	if config.Url == "" {
		log.Printf("Url must be specified in %s", command)
		return nil, false
	}
	switch config.Method {
	case "":
		config.Method = "GET"
	case "GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD", "CONNECT", "TRACE":
	default:
		log.Printf("Unknown method %s in %s", config.Method, command)
		return nil, false
	}
	return []interface{}{config, ctx}, true
}

func SmartNetRun(data []interface{}) bool {
	config := data[0].(*SmartNetConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return SmartNetRunByConfig(config, ctx)
}

func SaveActionResult(result string, data interface{}, ctx *dvcontext.RequestContext) {
	if result != "" {
		if ctx != nil {
			ctx.PrimaryContextEnvironment.SetProperty(result, data)
		} else {
			switch data.(type) {
			case string:
				dvparser.SetGlobalPropertiesValue(result, data.(string))
			}
		}
	}
}

func ProcessSavingActionResult(result string, data interface{}, ctx *dvcontext.RequestContext, err error, message1 string, message2 string) bool {
	if err != nil {
		log.Printf("%s %s %v", message1, message2, err)
		return false
	}
	SaveActionResult(result, data, ctx)
	return true
}

func SmartNetRunByConfig(config *SmartNetConfig, ctx *dvcontext.RequestContext) bool {
	options := dvnet.GetAveragePersistentOptions()
	headers := make(map[string]string)
	if config.Headers != "" {
		dvtextutils.PutDescribedAttributesToMapFromCommaSeparatedList(dvparser.GlobalProperties, headers, config.Headers)
	}
	if strings.HasPrefix(headers[Authorization], "M2M_") {
		microServiceName := headers[Authorization][4:]
		options[M2MAuthorizationRequest] = microServiceName
	}
	body := config.Body
	if len(body) > 2 && body[0] == '@' && body[1] == '@' {
		k := body[2:]
		n := len(k)
		strict := true
		filled := false
		if n > 1 && k[n-1] == '?' {
			strict = false
			if k[n-2] == '?' {
				filled = true
				k = k[:n-2]
			} else {
				k = k[:n-1]
			}
		}
		env := GetEnvironment(ctx)
		v, ok := env.Get(k)
		if !ok {
			if strict {
				log.Printf("Empty #{k} in net request")
				if ctx != nil {
					ctx.SetHttpErrorCode(400, "Unset "+k)
				}
				return true
			}
			v = nil
		}
		if filled {
			if dvjson.IsEmptyAny(v) {
				return true
			}
		}
	}
	res, err := NetRequest(config.Method, config.Url, config.Body, headers, options)
	if err != nil {
		log.Println(res)
		log.Printf("%s %s failed: %v", config.Method, config.Url, err)
		return false
	}
	var result interface{}
	switch config.ContentType {
	case "json":
		result, err = dvjson.JsonFullParser(res)
		if err != nil {
			result = string(res)
		}
	default:
		result = string(res)
	}
	SaveActionResult(config.Result, result, ctx)
	if config.Template.Dst != "" {
		if config.Template.Src != "" {
			s, err := dvparser.ConvertByteArrayByGlobalPropertiesRuntime([]byte(config.Template.Src), "net step")
			if err != nil {
				log.Printf("Cannot convert %s: %v", config.Template.Src, err)
				return false
			}
			isFile := len(s) < 300 && strings.Index(s, "<") < 0 && strings.Index(s, ">") < 0 && strings.Index(s, "{") < 0
			if isFile {
				if _, err = os.Stat(s); err != nil {
					isFile = false
				}
			}
			if isFile {
				res, err = ioutil.ReadFile(s)
				if err != nil {
					log.Printf("Cannot read file %s: %v", s, err)
					return false
				}
			} else {
				res = []byte(s)
			}
		}
		result, err := dvparser.ConvertByteArrayByGlobalPropertiesRuntime(res, config.Url)
		if err != nil {
			log.Printf("Cannot convert for file %s: %v", config.Template.Dst, err)
			ioutil.WriteFile(config.Template.Dst, res, 0664)
			return false
		}
		err = ioutil.WriteFile(config.Template.Dst, []byte(result), 0664)
		if err != nil {
			log.Printf("Cannot write file %s: %v", config.Template.Dst, err)
			return false
		}
	}
	return true
}

func processOsInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	cmd := strings.TrimSpace(command[strings.Index(command, ":")+1:])
	if cmd == "" {
		log.Printf("Empty net parameters", command)
		return nil, false
	}
	s, err := dvparser.ConvertByteArrayByGlobalPropertiesRuntime([]byte(cmd), "os step")
	if err != nil || s == "" {
		log.Printf("Error in %s: %v", cmd, err)
		return nil, false
	}
	return []interface{}{s}, true
}

func processOsRun(data []interface{}) bool {
	command := data[0].(string)
	cmd := exec.Command("cmd.exe", "/c", command)
	stdoutStderr, err := cmd.CombinedOutput()
	res := string(stdoutStderr)
	if Log >= dvlog.LogInfo || Log >= dvlog.LogError && err != nil {
		if Log < dvlog.LogInfo {
			dvlog.Printf("Executed: %s ", command)
		}
		dvlog.Println("", "-------------------START EXECUTING "+command+"--------------------\n"+res)
		dvlog.Println("", "____________________END EXECUTING "+command+"______________________")
	}
	return err == nil
}

func portForwardInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	p := strings.Index(command, ":")
	items := dvtextutils.ConvertToNonEmptyList(command[p+1:])
	if len(items) != 2 {
		dvlog.Printf("forward: <host>, <target> expected, but you specified %s", command)
		return nil, false
	}
	host := items[0]
	target := items[1]
	if !dvnet.ValidateHostTargetForPortForwarding(host, target) {
		return nil, false
	}
	return []interface{}{host, target}, true
}

func portForwardRun(data []interface{}) bool {
	host := data[0].(string)
	target := data[1].(string)
	dvnet.Forward(host, target)
	return true
}
