/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvoc

import (
	"bufio"
	"encoding/base64"
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"github.com/Dobryvechir/microcore/pkg/dvtemp"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	LogFatal   = 0
	LogError   = 1
	LogWarning = 2
	LogInfo    = 3
	LogDetail  = 4
	LogDebug   = 5
	LogTrace   = 6
)

var isOCLogined = false
var Log = LogError

const (
	openShiftLogin              = "login https://{{{OPENSHIFT_SERVER}}}.{{{OPENSHIFT_DOMAIN}}}:{{{OPENSHIFT_PORT}}} -u {{{OPENSHIFT_LOGIN}}} -p {{{OPENSHIFT_PASS}}} --insecure-skip-tls-verify=true -n {{{OPENSHIFT_NAMESPACE}}}"
	openShiftProject            = "\"{{{OPENSHIFT_NAMESPACE}}}\""
	openShiftSecrets            = "export --insecure-skip-tls-verify secret %1-client-credentials"
	openShiftExpose             = "expose svc/%1"
	openShiftDeleteSecret       = "delete secret "
	openShiftCacheProperty      = "OPENSHIFT_LOGIN_CACHE"
	openShiftStartCacheProperty = "OPENSHIFT_LOGIN_START_CACHE"
	ocGroupName                 = "dvoc"
	ocGroupStartName            = "dvoc_start"
	ocCachePrefix               = "dvoc."
)

func RunOCCommandFailureAllowed(params string, allowedFailureMessages []string) (string, int) {
	if !isOCLogined && !OcLogin() {
		return "", -1
	}
	return RunOCCommandLoginUnsafe(params, allowedFailureMessages, nil)
}

func RunOCCommandWithEditor(params string) (string, bool) {
	if !isOCLogined && !OcLogin() {
		return "", false
	}
	tmpFolder := dvtemp.GetUniqueTmpFolder()
	info := tmpFolder + "/i"
	content := tmpFolder + "/c"
	env := []string{
		"OC_EDITOR=dveditor",
		"DVEDITOR_INFO=" + info,
		"DVEDITOR_CONTENT=" + content,
	}
	exec, _ := RunOCCommandLoginUnsafe(params, nil, env)
	infoRes, err := ioutil.ReadFile(info)
	if err != nil {
		if Log >= LogError {
			dvlog.PrintfError("Cannot execute oc %s: %s", params, exec)
		}
		dvtemp.CleanTmpFolder(tmpFolder)
		return "", false
	}
	info = string(infoRes)
	if !strings.HasPrefix(info, "Ok") {
		if Log >= LogError {
			dvlog.PrintfError("Cannot execute oc %s: %s", params, info)
		}
		dvtemp.CleanTmpFolder(tmpFolder)
		return "", false
	}
	infoRes, err = ioutil.ReadFile(content)
	dvtemp.CleanTmpFolder(tmpFolder)
	if err != nil {
		if Log >= LogError {
			dvlog.PrintfError("Cannot read results of editing oc %s: %s %v", params, info, err)
		}
		return "", false
	}
	return string(infoRes), true
}

func RunOCCommandLoginUnsafe(params string, allowedFailureMessages []string, envParams []string) (string, int) {
	paramList := strings.Split(params, " ")
	if Log >= LogInfo {
		dvlog.Println(params, "Executing: oc "+params)
	}
	cmd := exec.Command("oc", paramList...)
	envLen := len(envParams)
	if envLen > 0 {
		cmd.Env = os.Environ()
		for i := 0; i < envLen; i++ {
			cmd.Env = append(cmd.Env, envParams[i])
		}
	}
	stdoutStderr, err := cmd.CombinedOutput()
	res := string(stdoutStderr)
	if Log >= LogInfo || Log >= LogError && err != nil {
		if Log < LogInfo {
			dvlog.Println(params, "Executed: oc "+params)
		}
		dvlog.Println("", "-------------------START EXECUTING OC RESULT --------------------\n"+res)
		dvlog.Println("", "____________________END EXECUTING OC RESULT______________________")
	}
	if err != nil {
		n := len(allowedFailureMessages)
		for i := 0; i < n; i++ {
			if strings.Index(res, allowedFailureMessages[i]) >= 0 {
				return res, i + 1
			}
		}
		errMessage := err.Error()
		if Log >= LogError {
			dvlog.PrintlnError("Error: " + errMessage)
		}
		if strings.Index(errMessage, "not found") >= 0 {
			dvlog.PrintlnError("You should have installed openshift client (oc) and put it path to PATH environment variable")
		}
		return "", -1
	}
	return res, 0
}

func RunOCCommand(params string) (string, bool) {
	res, status := RunOCCommandFailureAllowed(params, nil)
	return res, status >= 0
}

func ResetPodNameCache() {
	dvtemp.ResetCacheGroups(ocGroupName)
}

func OcLogin() bool {
	if isOCLogined {
		return true
	}
	dvtemp.SetCacheConfigByDescription(ocGroupName, dvparser.GlobalProperties[openShiftCacheProperty], nil)
	loginStart := dvparser.GlobalProperties[openShiftStartCacheProperty]
	if loginStart == "" {
		loginStart = "3600,3600"
	}
	dvtemp.SetCacheConfigByDescription(ocGroupStartName, loginStart, nil)
	projectName, err1 := dvparser.ConvertByteArrayByGlobalProperties([]byte(openShiftProject), "OPENSHIFT_NAME")
	if err1 != nil {
		dvlog.PrintfError("Make sure you specified OPENSHIFT_NAME (project name) in .properties")
		return false
	}
	cmdLine, err := dvparser.ConvertByteArrayByGlobalProperties([]byte(openShiftLogin), "oc login parameters")
	if err != nil {
		dvlog.PrintfError("Make sure you specified all constants related to oc login (%s)  (%v)", openShiftLogin, err)
		return false
	}
	ocCacheLogin := ocGroupStartName + ".project"
	info, err := dvtemp.GetCacheStringValue(ocCacheLogin)
	if err == nil && info == cmdLine {
		isOCLogined = true
		isOCLogined = OpenShiftEnsureExposeRoutes()
		return isOCLogined
	}
	res, status := RunOCCommandLoginUnsafe(cmdLine, nil, nil)
	if status < 0 {
		return false
	}
	pnt := strings.Index(res, projectName)
	if pnt < 0 {
		dvlog.PrintlnError(res)
		dvlog.PrintfError("Project %s is missing, specify it in OPENSHIFT_NAME property", projectName)
		return false
	}
	isOCLogined = true
	dvtemp.SetCacheValue(ocCacheLogin, cmdLine)
	isOCLogined = OpenShiftEnsureExposeRoutes()
	return isOCLogined
}

func GetOpenshiftSecrets(microserviceName string) (user string, ps string, okFinal bool) {
	cmdLine := strings.ReplaceAll(openShiftSecrets, "%1", microserviceName)
	info, ok := RunOCCommand(cmdLine)
	if !ok {
		return
	}
	scanner := bufio.NewScanner(strings.NewReader(info))
	for scanner.Scan() {
		s := strings.Split(strings.TrimSpace(scanner.Text()), " ")
		if len(s) == 2 {
			switch s[0] {
			case "password:":
				str, err := base64.StdEncoding.DecodeString(s[1])
				if err != nil {
					dvlog.PrintlnError(info)
					dvlog.PrintfError("Cannot get secret for microservice %s error: %v", microserviceName, err)
					return
				}
				ps = string(str)
				if user != "" {
					okFinal = true
					return
				}
			case "username:":
				str, err := base64.StdEncoding.DecodeString(s[1])
				if err != nil {
					dvlog.PrintlnError(info)
					dvlog.PrintfError("Cannot get secret for microservice %s error: %v", microserviceName, err)
					return
				}
				user = string(str)
				if ps != "" {
					okFinal = true
					return
				}
			}
		}
	}
	dvlog.PrintlnError(info)
	dvlog.PrintfError("Cannot get secret for microservice %s", microserviceName)
	return
}

func OpenShiftExposeSpecificRoutes(routes []string) bool {
	if len(routes) != 0 {
		for _, v := range routes {
			cmdLine := strings.ReplaceAll(openShiftExpose, "%1", v)
			res, err := RunOCCommandOrCacheFailureAllowed(cmdLine, []string{"AlreadyExist"})
			if err != nil {
				return false
			}
			if strings.Index(res, "AlreadyExist") < 0 && strings.Index(res, "exposed") < 0 {
				dvlog.PrintfError("Warning: Unrecognized response to %s : %s", cmdLine, res)
			}
		}
	}
	return true
}

func OpenShiftEnsureExposeRoutes() bool {
	routes := dvparser.ConvertToNonEmptyList(dvparser.GlobalProperties[openShiftEnsureRoutes])
	return OpenShiftExposeSpecificRoutes(routes)
}

func OpenShiftAddRoutesTOBeExposed(routeList string) {
	origRoutes := dvparser.ConvertToNonEmptyList(dvparser.GlobalProperties[openShiftEnsureRoutes])
	if len(origRoutes) == 0 {
		dvparser.GlobalProperties[openShiftEnsureRoutes] = routeList
		return
	}
	newRoutes := dvparser.ConvertToNonEmptyList(routeList)
	if len(newRoutes) == 0 {
		return
	}
	origRoutes = dvparser.AddStringListWithoutRepeats(origRoutes, newRoutes)
	dvparser.GlobalProperties[openShiftEnsureRoutes] = strings.Join(origRoutes, ",")
}

func GetPodName(microserviceName string, silent bool) (name string, ok bool) {
	info, err := RunOCCommandOrCache("get pods")
	if err != nil {
		dvlog.PrintfError("Failed to get pods %v", err)
		return
	}
	candidates := make([]string, 0, 2)
	pos := strings.Index(info, microserviceName)
	for pos >= 0 {
		c := uint8(0)
		if pos > 0 {
			c = info[pos-1]
		}
		if c <= ' ' {
			c = info[pos+len(microserviceName)]
			if c <= ' ' {
				return microserviceName, true
			}
			if c == '-' {
				endPos := pos + len(microserviceName) + 1
				for ; endPos < len(info); endPos++ {
					c := info[endPos]
					if !(c == '-' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9') {
						break
					}
				}
				candidate := info[pos:endPos]
				if !strings.HasSuffix(candidate, "-deploy") {
					candidates = append(candidates, candidate)
				}
				pos = endPos
			}
		}
		pos++
		npos := strings.Index(info[pos:], microserviceName)
		if npos < 0 {
			break
		}
		pos += npos
	}
	n := len(candidates)
	if n == 0 {
		ResetPodNameCache()
		if !silent {
			dvlog.PrintfError("Pod for microservice %s does not exist in the cloud's project", microserviceName)
		}
		return "", false
	}
	candidate := candidates[0]
	for j := 1; j < n; j++ {
		s := candidates[j]
		if len(s) < len(candidate) {
			candidate = s
		}
	}
	return candidate, true
}

func DeletePod(name string) bool {
	cmdLine := "delete pod " + name
	info, ok := RunOCCommand(cmdLine)
	if !ok {
		dvlog.PrintfError("Failed to delete pod %s", name)
		return false
	}
	if Log >= LogInfo {
		dvlog.PrintlnError(info)
	}
	return true
}

func OpenShiftScaleToReplicas(microService string, replicas int) bool {
	cmdLine := "scale dc " + microService + " --replicas=" + strconv.Itoa(replicas)
	info, ok := RunOCCommand(cmdLine)
	if !ok {
		dvlog.PrintfError("Failed to execute %s", cmdLine)
	} else {
		if strings.Index(info, "scaled") > 0 {
			return true
		}
		dvlog.PrintlnError(cmdLine)
		dvlog.PrintlnError(info)
	}
	return false
}

func OpenShiftSetEnvironment(microService string, env []string) bool {
	cmdLine := "env dc " + microService + " " + strings.Join(env, " ")
	info, ok := RunOCCommand(cmdLine)
	if !ok {
		dvlog.PrintfError("Failed to execute %s", cmdLine)
	} else {
		if strings.Index(info, "updated") > 0 {
			return true
		}
		dvlog.PrintlnError(cmdLine)
		dvlog.PrintlnError(info)
	}
	return false
}

func RunOCCommandOrCache(cmdLine string) (string, error) {
	return RunOCCommandOrCacheFailureAllowed(cmdLine, nil)
}

func RunOCCommandOrCacheFailureAllowed(cmdLine string, allowedFailureMessages []string) (string, error) {
	var status int
	ocCacheName := ocCachePrefix + dvtemp.GetSafeFileName(cmdLine)
	if !isOCLogined {
		OcLogin()
	}
	info, err := dvtemp.GetCacheStringValue(ocCacheName)
	if err != nil || info == "" {
		info, status = RunOCCommandFailureAllowed(cmdLine, allowedFailureMessages)
		if status < 0 {
			dvlog.PrintfError("Failed to execute %s", cmdLine)
			return "", errors.New("Failed to execute " + cmdLine)
		}
		dvtemp.SetCacheValue(ocCacheName, info)
	}
	return info, nil
}

func InitBySettings(parameters map[string]string, functionPool map[string]interface{}) {
	Log = dvlog.LogLevelFromProperties(parameters, dvparser.MicroCorePrexix+"DVOC_LOG", Log)
}

var initedBySettings = dvparser.RegisterInitBySettingFunc(InitBySettings)
