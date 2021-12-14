package dvaction

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvnet"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
	"log"
	"os/exec"
	"strings"
)

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

