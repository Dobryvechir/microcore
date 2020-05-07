/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvoc

import (
	"encoding/json"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvmeta"
	"github.com/Dobryvechir/microcore/pkg/dvnet"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	Authorization = "Authorization"
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

func processNetInit(command string, ctx *dvmeta.RequestContext) ([]interface{}, bool) {
	params := dvparser.ConvertToNonEmptyList(command)
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
	Url      string                 `json:"url"`
	Method   string                 `json:"method"`
	Headers  string                 `json:"headers"`
	Template SmartNetConfigTemplate `json:"template"`
	Result   string                 `json:"result"`
	Body     string                 `json:"body"`
}

func SmartNetInit(command string, ctx *dvmeta.RequestContext) ([]interface{}, bool) {
	cmd := strings.TrimSpace(command[strings.Index(command, ":")+1:])
	if cmd == "" || cmd[0] != '{' || cmd[len(cmd)-1] != '}' {
		log.Printf("Empty net parameters", command)
		return nil, false
	}
	config := &SmartNetConfig{}
	err := json.Unmarshal([]byte(cmd), config)
	if err != nil {
		log.Printf("Error converting parameters: %v", err)
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
	var ctx *dvmeta.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvmeta.RequestContext)
	}
	return SmartNetRunByConfig(config, ctx)
}

func SmartNetRunByConfig(config *SmartNetConfig, ctx *dvmeta.RequestContext) bool {
	options := dvnet.GetAveragePersistentOptions()
	headers := make(map[string]string)
	if config.Headers != "" {
		dvparser.PutDescribedAttributesToMapFromCommaSeparatedList(dvparser.GlobalProperties, headers, config.Headers)
	}
	if strings.HasPrefix(headers[Authorization], "M2M_") {
		microServiceName := headers[Authorization][4:]
		options[M2MAuthorizationRequest] = microServiceName
	}
	res, err := NetRequest(config.Method, config.Url, config.Body, headers, options)
	if err != nil {
		log.Println(res)
		log.Printf("%s %s failed: %v", config.Method, config.Url, err)
		return false
	}
	if config.Result != "" {
		if ctx != nil {
			ctx.ExtraAsDvObject.SetProperty(config.Result, string(res))
		} else {
			dvparser.SetGlobalPropertiesValue(config.Result, string(res))
		}
	}
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

func processOsInit(command string, ctx *dvmeta.RequestContext) ([]interface{}, bool) {
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
	if Log >= LogInfo || Log >= LogError && err != nil {
		if Log < LogInfo {
			dvlog.Printf("Executed: %s ", command)
		}
		dvlog.Println("", "-------------------START EXECUTING "+command+"--------------------\n"+res)
		dvlog.Println("", "____________________END EXECUTING "+command+"______________________")
	}
	return err == nil
}

func portForwardInit(command string, ctx *dvmeta.RequestContext) ([]interface{}, bool) {
	p := strings.Index(command, ":")
	items := dvparser.ConvertToNonEmptyList(command[p+1:])
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
