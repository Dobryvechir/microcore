/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvoc

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvaction"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvdir"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	MicroServiceDeleteForced                = 0
	MicroServiceDeleteTrySaveAndForceDelete = 1
	MicroServiceDeleteSaveAndSafeDelete     = 2
	MicroServiceSaveTemporaryFolder         = "__dobryvechir_microcore"
	MicroServiceProperty                    = "MICROSERVICE"
)

func MicroServiceUpInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	pos := strings.Index(command, ":")
	command = strings.TrimSpace(command[pos+1:])
	pos = strings.Index(command, "{")
	name := command
	annotations := make(map[string]string)
	image := ""
	files := make(map[string]string)
	commands := make([]string, 0, 1)
	if pos >= 0 {
		name = strings.TrimSpace(name[:pos])
		options, err := dvjson.JsonFullParser([]byte(command[pos:]))
		if err != nil {
			dvlog.PrintfError("Error in json %s: %v", command, err)
			return nil, false
		}
		if options.Kind != dvjson.FIELD_OBJECT {
			dvlog.PrintfError("Expected object {} after microservice name in %s", command)
			return nil, false
		}
		for _, v := range options.Fields {
			switch string(v.Name) {
			case "annotations":
				err = v.ReadSimpleStringMap(annotations)
				if err != nil {
					dvlog.PrintfError("annotations must be a simple json string map in %s: %v", command, err)
					return nil, false
				}
			case "commands":
				commands, err = v.ReadSimpleStringList(commands)
				if err != nil {
					dvlog.PrintfError("files must be a simple json string map in %s: %v", command, err)
					return nil, false
				}
			case "files":
				err = v.ReadSimpleStringMap(files)
				if err != nil {
					dvlog.PrintfError("files must be a simple json string map in %s: %v", command, err)
					return nil, false
				}
			case "image":
				image, err = v.ReadSimpleString()
				if err != nil {
					dvlog.PrintfError("image must be a simple json string in %s: %v", command, err)
					return nil, false
				}
			default:
				dvlog.PrintfError("Unknown property %s in %s", string(v.Name), command)
				return nil, false
			}
		}
	}
	if name == "" {
		dvlog.PrintfError("Empty microservice name")
		return nil, false
	}
	params := map[string]string{
		MicroServiceProperty: name,
		"ANNOTATIONS":        dvjson.ConvertSimpleStringMapToJson(annotations, false),
	}
	if image != "" {
		params["TEMPLATE_IMAGE"] = image
	}
	return []interface{}{params, files, commands}, true
}

func DownMicroServiceCommands(deleteInfo []string) (ok bool) {
	n := len(deleteInfo)
	ok = true
	for i := 0; i < n; i++ {
		info, okSingle := RunOCCommand(deleteInfo[i])
		if strings.Index(info, "NotFound") >= 0 || strings.Index(info, "No resources") >= 0 {
			continue
		}
		if !okSingle {
			ok = false
			dvlog.PrintfError("Error deleting %s", info)
		} else if Log >= LogInfo {
			dvlog.PrintlnError(info)
		}
	}
	return
}

func DownWholeMicroServiceHard(microServiceName string, microServiceAliases string, includeList []string) bool {
	aliases := dvtextutils.ConvertToNonEmptyList(microServiceName + "," + microServiceAliases)
	n := len(aliases)
	deleteInfo := make([]string, 0, n*5)
	for i := 0; i < n; i++ {
		name := aliases[i]
		deleteInfo = append(deleteInfo, "delete all -l name="+name)
		deleteInfo = append(deleteInfo, "delete deploymentconfig "+name)
		deleteInfo = append(deleteInfo, "delete configmap "+name+".monitoring-config")
		svc := ensureOpenShiftServiceNameInParams(nil, name)
		deleteInfo = append(deleteInfo, "delete svc "+svc)
		route := ensureOpenShiftRouteNameInParams(nil, name)
		deleteInfo = append(deleteInfo, "delete route "+route)
	}
	deleteInfo = dvtextutils.MakeUniqueStringList(deleteInfo, includeList)
	return DownMicroServiceCommands(deleteInfo)
}

func DownWholeMicroService(microServiceName string, mode int, debugNotSaved bool) bool {
	saveMode := 2
	if debugNotSaved {
		saveMode = 1
	}
	switch mode {
	case MicroServiceDeleteForced:
		deleteInfo, _, _ := SaveMicroServiceOpenShiftInfo(microServiceName, -1)
		aliases := PropertyForSpecificMicroService(microServiceName, "ALIASES")
		DownWholeMicroServiceHard(microServiceName, aliases, deleteInfo)
		return true
	case MicroServiceDeleteTrySaveAndForceDelete:
		deleteInfo, _, ok := SaveMicroServiceOpenShiftInfo(microServiceName, saveMode)
		if ok {
			DownMicroServiceCommands(deleteInfo)
		} else {
			aliases := PropertyForSpecificMicroService(microServiceName, "ALIASES")
			DownWholeMicroServiceHard(microServiceName, aliases, deleteInfo)
		}
		return true
	case MicroServiceDeleteSaveAndSafeDelete:
		deleteInfo, _, ok := SaveMicroServiceOpenShiftInfo(microServiceName, saveMode)
		if !ok {
			return false
		}
		return DownMicroServiceCommands(deleteInfo)
	}
	dvlog.PrintfError("Unknown microservice delete mode %d", mode)
	return false
}

func GetParamByMapOrGlobal(params map[string]string, key string, baseName string) (res string) {
	res = ""
	if params != nil {
		res = params[key]
	}
	if res == "" {
		if baseName != "" {
			res = PropertyForSpecificMicroService(baseName, key)
		} else {
			res = dvparser.GlobalProperties[key]
		}
	}
	return
}

func PropertyForSpecificMicroService(microServiceName string, prefix string) string {
	name := strings.ReplaceAll(strings.ToUpper(microServiceName), "-", "_") + "_" + prefix
	return dvparser.GlobalProperties[name]
}

func ensureOpenShiftServiceNameInParams(params map[string]string, name string) string {
	service := GetParamByMapOrGlobal(params, "OPENSHIFT_SERVICE_NAME", name)
	n := len(service)
	if n == 0 || n > 32 {
		service = name
		if len(service) > 32 {
			service = service[:32]
			if service[31] == '-' {
				service = service[:31] + "v"
			}
		}
	}
	if params != nil {
		params["OPENSHIFT_SERVICE_NAME"] = service
	}
	return service
}

func ensureOpenShiftRouteNameInParams(params map[string]string, name string) string {
	route := GetParamByMapOrGlobal(params, "OPENSHIFT_ROUTE_NAME", name)
	if route == "" {
		route = name
	}
	if params != nil {
		params["OPENSHIFT_ROUTE_NAME"] = route
	}
	return route
}

func CreateMicroService(params map[string]string, files map[string]string, commands []string) bool {
	microServiceName := params[MicroServiceProperty]
	ensureOpenShiftServiceNameInParams(params, microServiceName)
	ensureOpenShiftRouteNameInParams(params, microServiceName)
	onBehalfMicroServiceName := GetParamByMapOrGlobal(params, "MANAGING_MICROSERVICE", "")
	if microServiceName == "" {
		dvlog.PrintfError("Microservice name to be created was not specified")
	}
	tmpSuffix := "__dvoc_create_microservice__" + microServiceName
	_, _, ok := GetOpenshiftSecrets(microServiceName)
	if !ok {
		user, pw, ok := dvaction.GetIdentityProviderClientCredentials(microServiceName)
		if !ok {
			return false
		}
		if !CreateClientCredentials(user, pw, microServiceName) {
			return false
		}
		if !RegisterUserCredentialsWithIdentityProvider(user, pw, microServiceName, onBehalfMicroServiceName, tmpSuffix) {
			return false
		}
		provideKeycloakAutoUpdate()
	}
	var json []byte
	var err error
	specTemplate := strings.TrimSpace(GetParamByMapOrGlobal(params, "OPENSHIFT_TEMPLATE", microServiceName))
	if specTemplate != "" {
		if specTemplate[0] != '{' {
			json, err = ioutil.ReadFile(specTemplate)
			if err != nil {
				dvlog.PrintlnError(specTemplate)
				dvlog.PrintfError("Your template for %s is neither file nor template itself ", microServiceName)
				return false
			}
			specTemplate = strings.TrimSpace(string(json))
			if specTemplate[0] != '{' {
				dvlog.PrintfError("Your template file for %s does not contain legal template", microServiceName)
				return false
			}
			requiredParams, ok := dvtextutils.ComposeParametersInTemplate(specTemplate, 2, templateRequired)
			if !ok {
				dvlog.PrintfError("Your template is corrupt: %s", specTemplate)
				return false
			}
			json, err = ComposeOpenShiftJsonTemplateBySample(specTemplate, requiredParams, params)
			if err != nil {
				dvlog.PrintfError("Cannot compose template for %s: %v", microServiceName, err)
				return false
			}
		}

	}

	if len(json) == 0 {
		json, err = ComposeOpenShiftJsonTemplate(params)
		if err != nil {
			dvlog.PrintfError("Cannot compose template for %s: %v", microServiceName, err)
			return false
		}
	}
	path := dvdir.GetTempPathSlashed() + microServiceName + "__debug_fragments_template.json"
	err = ioutil.WriteFile(path, json, 0664)
	if err != nil {
		dvlog.PrintfError("Cannot temporarily save a template file %s", path)
		return false
	}
	cmdLine := "new-app -f " + path
	info, ok := RunOCCommand(cmdLine)
	if !ok {
		dvlog.PrintfError("Failed to execute %s", cmdLine)
		return false
	} else {
		if Log >= LogInfo {
			dvlog.PrintlnError(info)
		}
	}
	time.Sleep(10 * time.Second)
	return UploadFiles(microServiceName, files, commands)
}

func MicroServiceUpRun(data []interface{}) bool {
	params := data[0].(map[string]string)
	files := data[1].(map[string]string)
	commands := data[2].([]string)
	DownWholeMicroService(params[MicroServiceProperty], MicroServiceDeleteForced, true)
	return CreateMicroService(params, files, commands)
}

func MicroServiceUpOnlyRun(data []interface{}) bool {
	params := data[0].(map[string]string)
	files := data[1].(map[string]string)
	commands := data[2].([]string)
	return CreateMicroService(params, files, commands)
}

func MicroServiceDownInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	save := MicroServiceDeleteTrySaveAndForceDelete
	debugNotSaved := true
	pos := strings.Index(command, "{")
	if pos >= 0 {
		saveOption := strings.ToLower(dvtextutils.QuickLookJsonOption(command[pos:], "save"))
		switch saveOption {
		case "false", "\"false\"":
			save = MicroServiceDeleteForced
		case "true", "\"true\"", "":
			save = MicroServiceDeleteTrySaveAndForceDelete
			debugNotSaved = false
		case "nondebug", "\"nondebug\"":
			save = MicroServiceDeleteTrySaveAndForceDelete
		case "safenondebug", "\"safenondebug\"":
			save = MicroServiceDeleteSaveAndSafeDelete
		case "safe", "\"safe\"":
			save = MicroServiceDeleteSaveAndSafeDelete
			debugNotSaved = false
		default:
			dvlog.PrintfError("Unknown option of save (only NonDebug, SafeNonDebug, True, False and Safe are allowed) %s", saveOption)
			return nil, false
		}
		command = command[:pos]
	}
	command = strings.TrimSpace(command)
	if command == "" {
		dvlog.PrintfError("microservice name was not specified")
		return nil, false
	}
	return []interface{}{command, save, debugNotSaved}, true
}

func MicroServiceDownRun(data []interface{}) bool {
	name := data[0].(string)
	save := data[1].(int)
	debugNotSaved := data[2].(bool)
	return DownWholeMicroService(name, save, debugNotSaved)
}

func SaveMicroServiceOpenShiftInfo(microServiceName string, save int) (deleteInfo []string, deployment string, ok bool) {
	notCritical := save < 0
	deployment, deleteInfo, _, _, _, err := GetExistingFullOpenShiftTemplate(microServiceName, notCritical)
	if err != nil {
		if notCritical {
			err = nil
			return
		}
		dvlog.PrintfError("Error: %v", err)
		return
	}
	if save == 1 && strings.Index(deployment, templateDebugSign) < 0 || save == 2 {
		fileName := dvdir.GetNextSaveFileName(MicroServiceSaveTemporaryFolder, microServiceName)
		err = ioutil.WriteFile(fileName, []byte(deployment), 0664)
		if err != nil {
			return
		}
	}
	ok = true
	return
}

func ReduceMicroServiceSaveInfo(microServiceName string) {
	dvdir.DeleteAllSavedFilesExceptZero(MicroServiceSaveTemporaryFolder, microServiceName)
}

func MicroServiceSaveInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	pos := strings.Index(command, ":")
	params := dvtextutils.ConvertToNonEmptyList(command[pos+1:])
	n := len(params)
	if n == 0 {
		dvlog.PrintlnError("microserviceSave:<microservice name>,<folder where information will be saved, optional, defaults to the current folder>")
		dvlog.PrintlnError(command + " requires the microservice name as the first parameter (optionally dir name as the second")
		return nil, false
	}
	dir := "."
	if n > 1 && params[1] != "." {
		dir = params[1]
		if dir[len(dir)-1] == '/' {
			dir = dir[:len(dir)-1]
		}
		os.MkdirAll(dir, 0755)
	}
	return []interface{}{params[0], dir}, true
}

func MicroServiceSaveSingle(microServiceName string, folder string) (createFile string, deleteFile string, ok bool) {
	createFile = folder + "/" + microServiceName + "_up.cmd"
	deleteFile = folder + "/" + microServiceName + "_down.cmd"
	templateFile := folder + "/" + microServiceName + "_template.json"
	err := ioutil.WriteFile(createFile, []byte("oc new-app -f "+dvdir.MakeLastPathIfNotAbsolute(templateFile)), 0664)
	if err != nil {
		dvlog.PrintfError("Cannot write file %s: %v", createFile, err)
		return
	}
	deleteInfo, deployment, ok := SaveMicroServiceOpenShiftInfo(microServiceName, 0)
	if !ok {
		return
	}
	err = ioutil.WriteFile(templateFile, []byte(deployment), 0664)
	if err != nil {
		dvlog.PrintfError("Cannot write file %s: %v", templateFile, err)
		return
	}
	deleteBlock := ""
	n := len(deleteInfo)
	for i := 0; i < n; i++ {
		deleteBlock += "oc " + deleteInfo[i] + "\n"
	}
	err = ioutil.WriteFile(deleteFile, []byte(deleteBlock), 0664)
	if err != nil {
		dvlog.PrintfError("Cannot write file %s: %v", deleteFile, err)
		return
	}
	ok = true
	return
}

func MicroServiceSaveRun(data []interface{}) bool {
	microServiceName := data[0].(string)
	folderName := data[1].(string)
	if microServiceName == "*" {
		list, err := GetMicroServiceFullList()
		if err != nil {
			dvlog.PrintfError("Failed to get list of microservices %v", err)
			return false
		}
		n := len(list)
		ok := true
		crInfo := ""
		dlInfo := ""
		for i := 0; i < n; i++ {
			cr, dl, isOk := MicroServiceSaveSingle(list[i], folderName)
			if !isOk {
				ok = false
			} else {
				crInfo += "call " + dvdir.MakeLastPathIfNotAbsolute(cr) + "\n"
				dlInfo += "call " + dvdir.MakeLastPathIfNotAbsolute(dl) + "\n"
			}
		}
		err = ioutil.WriteFile(folderName+"/__allUp.cmd", []byte(crInfo), 0664)
		if err != nil {
			dvlog.PrintfError("Error writing __allUp.cmd: %v", err)
			ok = false
		}
		err = ioutil.WriteFile(folderName+"/__allDown.cmd", []byte(dlInfo), 0664)
		if err != nil {
			dvlog.PrintfError("Error writing __allDown.cmd: %v", err)
			ok = false
		}
		return ok
	}
	_, _, ok := MicroServiceSaveSingle(microServiceName, folderName)
	return ok
}

func MicroServiceRestoreInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	pos := strings.Index(command, ":")
	microServiceName := strings.TrimSpace(command[pos+1:])
	pos = strings.Index(microServiceName, "{")
	if pos >= 0 {
		microServiceName = strings.TrimSpace(microServiceName[:pos])
	}
	if microServiceName == "" {
		dvlog.PrintlnError("microserviceRestore:<microservice name>")
		dvlog.PrintlnError(command + " requires the microservice name as the first parameter")
		return nil, false
	}
	fileName := dvdir.GetLastSavedFileName(MicroServiceSaveTemporaryFolder, microServiceName)
	if fileName == "" {
		dvlog.PrintfError("No configuration was saved for %s", microServiceName)
		return nil, false
	}
	return []interface{}{microServiceName, fileName}, true
}

func MicroServiceRestoreRun(data []interface{}) bool {
	microServiceName := data[0].(string)
	templateFileName := data[1].(string)
	err := MicroServiceRestoration(microServiceName, templateFileName)
	if err != nil {
		dvlog.PrintfError("Error: %v", err)
		return false
	}
	return true
}

func MicroServiceCacheCleanInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	pos := strings.Index(command, ":")
	params := dvtextutils.ConvertToNonEmptyList(command[pos+1:])
	n := len(params)
	if n == 0 {
		dvlog.PrintlnError("microserviceCacheClean:<microservice name>,<optionally, number of last cache versions to be removed, defaults to all>")
		dvlog.PrintlnError(command + " requires the microservice name as the first parameter")
		return nil, false
	}
	microServiceName := params[0]
	count := 0
	var err error
	if n > 1 {
		s := strings.ToLower(params[1])
		if s != "all" {
			count, err = strconv.Atoi(s)
			if err != nil {
				dvlog.PrintfError("The second parameter must be either all or a integer number, but it %s", s)
				return nil, false
			}
		}
	}
	return []interface{}{microServiceName, count}, true
}

func MicroServiceCacheCleanRun(data []interface{}) bool {
	microServiceName := data[0].(string)
	count := data[1].(int)
	return dvdir.DeleteLastSavedFile(MicroServiceSaveTemporaryFolder, microServiceName, count)
}

func MicroServiceRestoration(microServiceName string, templateFileName string) error {
	if templateFileName == "" {
		templateFileName = dvdir.GetLastSavedFileName(MicroServiceSaveTemporaryFolder, microServiceName)
		if templateFileName == "" {
			return errors.New("No configuration was stored for " + microServiceName)
		}
	}
	DownWholeMicroService(microServiceName, MicroServiceDeleteForced, true)
	cmdLine := "new-app -f " + templateFileName
	info, ok := RunOCCommand(cmdLine)
	if !ok {
		dvlog.PrintlnError(info)
		dvlog.PrintfError("Error restoring %s", microServiceName)
		return errors.New("Error restoring " + microServiceName)
	}
	return nil
}

func MicroServiceExecInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	pos := strings.Index(command, ":")
	command = command[pos+1:]
	pos = strings.Index(command, "[")
	if pos < 0 {
		dvlog.PrintlnError("No [: Expected format: microserviceExec:microserviceName[\"command1\",\"command2\"...]")
	}
	microServiceName := strings.TrimSpace(command[:pos])
	options, err := dvjson.JsonFullParser([]byte(command[pos:]))
	if err != nil {
		dvlog.PrintfError("Error in json %s: %v", command, err)
		return nil, false
	}
	if options.Kind != dvjson.FIELD_ARRAY {
		dvlog.PrintfError("Expected object [] after microservice name in %s", command)
		return nil, false
	}
	commands, err := options.ReadSimpleStringList(nil)
	if err != nil {
		dvlog.PrintfError("Error in json %s: %v", command, err)
		return nil, false
	}
	return []interface{}{microServiceName, commands}, true
}

func MicroServiceExecRun(data []interface{}) bool {
	microServiceName := data[0].(string)
	commands := data[1].([]string)
	return ExecuteCommandsOnPod(microServiceName, commands)
}

func ExposeMicroServiceInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	pos := strings.Index(command, ":")
	command = command[pos+1:]
	microServices := dvtextutils.ConvertToNonEmptyList(command)
	if len(microServices) == 0 {
		dvlog.PrintfError("Empty expose executor %s", command)
		return nil, false
	}
	return []interface{}{microServices}, true
}

func ExposeMicroServiceRun(data []interface{}) bool {
	microServices := data[0].([]string)
	return OpenShiftExposeSpecificRoutes(microServices)
}

func IsMicroServicePresent(name string) (bool, error) {
	list, err := GetMicroServiceFullList()
	if err != nil {
		return false, err
	}
	n := len(list)
	for i := 0; i < n; i++ {
		if list[i] == name {
			return true, nil
		}
	}
	return false, nil
}
