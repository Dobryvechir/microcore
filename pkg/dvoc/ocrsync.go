/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvoc

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvmeta"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"github.com/Dobryvechir/microcore/pkg/dvtemp"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const (
	StrategyReadWriteSingleFileBest    = 0
	StrategyReadWriteSingleFileThruCat = 1
	StrategyReadWriteSingleFileThruDir = 2
)

var StrategyReadWriteSingleFile = StrategyReadWriteSingleFileThruCat
var StragegyReadWriteSingleFileThruDirExclude = make(map[string]bool)

func WriteDirectoryToPod(podName string, distributionFolder string, podFolder string, options string) bool {
	if _, err := os.Stat(distributionFolder); err != nil {
		dvlog.PrintfError("DISTRIBUTION_FOLDER must point to a real folder, but it points to %s", distributionFolder)
		return false
	}
	cmdLine := "rsync " + distributionFolder + " " + podName + ":" + podFolder
	if options != "" {
		cmdLine += " " + options
	}
	info, ok := RunOCCommand(cmdLine)
	if !ok {
		dvlog.PrintfError("Failed to execute %s", cmdLine)
		return false
	} else {
		if strings.Index(info, "error:") >= 0 {
			dvlog.PrintlnError(info)
			return false
		}
		if Log >= LogInfo {
			dvlog.PrintlnError(info)
		}
	}
	return true
}

func ReadDirectoryFromPod(podName string, distributionFolder string, podFolder string, options string) bool {
	if _, err := os.Stat(distributionFolder); err != nil {
		dvlog.PrintfError("DISTRIBUTION_FOLDER must point to a real folder, but it points to %s", distributionFolder)
		return false
	}
	cmdLine := "rsync " + podName + ":" + podFolder + " " + distributionFolder
	if options != "" {
		cmdLine += " " + options
	}
	info, ok := RunOCCommand(cmdLine)
	if !ok {
		dvlog.PrintfError("Failed to execute %s", cmdLine)
	} else {
		if strings.Index(info, "error:") >= 0 {
			dvlog.PrintlnError(info)
			return false
		}
		if Log >= LogInfo {
			dvlog.PrintlnError(info)
		}
	}
	return true
}

func OpenShiftGetPodNameAndPath(name string) (tmpFolder string, serviceName, podName string, pathName string, fileName string, err error) {
	pos := strings.Index(name, ":")
	if pos <= 0 {
		err = errors.New("Incorrect pod url " + name + ". It must be service:/path/to/dir")
		return
	}
	serviceName = strings.TrimSpace(name[:pos])
	pathName = strings.TrimSpace(name[pos+1:])
	var ok bool
	if len(pathName) == 0 {
		err = errors.New("Path must not be empty in " + name + ". It must be service:/path/to/dir")
	} else {
		pos := strings.LastIndex(pathName, "/")
		if pos <= 0 {
			err = errors.New("Cannot work with the file at the root folder in " + name)
		} else {
			fileName = pathName[pos+1:]
			pathName = pathName[:pos]
			if fileName == "" {
				err = errors.New("File name must not be empty in " + name)
			} else {
				podName, ok = GetPodName(serviceName, Log < LogInfo)
				if !ok {
					err = errors.New("Pod " + podName + " does not exist")
				} else {
					tmpFolder = dvtemp.GetUniqueTmpFolder()
					if tmpFolder == "" {
						err = errors.New("Cannot create temporary folder")
					}
				}
			}
		}
	}
	return
}

func GetRsyncOptionsForSingleFile(fileName string) string {
	return "--exclude=* --include=" + fileName
}

func OpenShiftReadTextFile(name string, strategy int) (data string, err error) {
	tmpFolder, serviceName, podName, podFolder, fileName, err := OpenShiftGetPodNameAndPath(name)
	if err != nil {
		return
	}
	switch strategy {
	case StrategyReadWriteSingleFileThruDir:
		if !ReadDirectoryFromPod(podName, tmpFolder, podFolder, GetRsyncOptionsForSingleFile(fileName)) {
			err = errors.New("Cannot read access " + name)
			dvtemp.CleanTmpFolder(tmpFolder)
			return
		}
		pos := strings.LastIndex(podFolder, "/")
		subFolder := ""
		if pos >= 0 {
			subFolder = podFolder[pos+1:]
		}
		tmpFileName := dvtemp.GetFirstExistingPath(tmpFolder+"/"+fileName, tmpFolder+"/"+subFolder+"/"+fileName)
		if tmpFileName == "" {
			return "", errors.New("File " + fileName + " does not exist in " + podFolder + " at " + podName)
		}
		var byteData []byte
		byteData, err = ioutil.ReadFile(tmpFileName)
		dvtemp.CleanTmpFolder(tmpFolder)
		if err != nil {
			return
		}
		return string(byteData), nil
	case StrategyReadWriteSingleFileThruCat:
		dvtemp.CleanTmpFolder(tmpFolder)
		info, err := ExecuteCommandOnPod(podName, "cat "+podFolder+"/"+fileName)
		if err != nil {
			return "", err
		}
		return info, nil
	case StrategyReadWriteSingleFileBest:
		if !StragegyReadWriteSingleFileThruDirExclude[serviceName] {
			data, err = OpenShiftReadTextFile(name, StrategyReadWriteSingleFileThruDir)
			if err == nil {
				return
			}
			StragegyReadWriteSingleFileThruDirExclude[serviceName] = true
		}
		return OpenShiftReadTextFile(name, StrategyReadWriteSingleFileThruCat)
	}
	return "", errors.New("No strategy to read/write single files")
}

func OpenShiftWriteTextFile(name string, data string, strategy int) error {
	tmpFolder, _, podName, podFolder, fileName, err := OpenShiftGetPodNameAndPath(name)
	if err != nil {
		return err
	}
	switch strategy {
	case StrategyReadWriteSingleFileThruDir:
		pos := strings.LastIndex(podFolder, "/")
		filePlace := tmpFolder + "/" + fileName
		writeFolder := tmpFolder
		if pos >= 0 {
			lastPath := podFolder[pos+1:]
			if lastPath != "" {
				writeFolder = tmpFolder + "/" + lastPath
				os.Mkdir(writeFolder, 0777)
				filePlace = writeFolder + "/" + fileName
				podFolder = podFolder[:pos+1]
			}
		}
		err = ioutil.WriteFile(filePlace, []byte(data), 0664)
		if err != nil {
			dvtemp.CleanTmpFolder(tmpFolder)
			return err
		}
		if !WriteDirectoryToPod(podName, writeFolder, podFolder, "") {
			err = errors.New("Cannot write access " + name)
			dvtemp.CleanTmpFolder(tmpFolder)
			return err
		}
		return nil
	case StrategyReadWriteSingleFileThruCat:
		dvtemp.CleanTmpFolder(tmpFolder)
		cmd := "echo VOL >/tmp/vol.txt"
		info, err := ExecuteCommandOnSpecificPod(podName, cmd)
		dvlog.PrintfError("info: %s err: %v", info, err)
		cmd = "echo \"VOLKov\" >>/tmp/vol.txt"
		info, err = ExecuteCommandOnSpecificPod(podName, cmd)
		dvlog.PrintfError("info: %s err: %v", info, err)

	case StrategyReadWriteSingleFileBest:
		err = OpenShiftWriteTextFile(name, data, StrategyReadWriteSingleFileThruDir)
		return err
		//if !StragegyReadWriteSingleFileThruDirExclude[serviceName] {
		//	err = OpenShiftWriteTextFile(name, data, StrategyReadWriteSingleFileThruDir)
		//	return err
		//if err == nil {
		//	return nil
		//}
		//StragegyReadWriteSingleFileThruDirExclude[serviceName] = true
		//}
		//return OpenShiftWriteTextFile(name, data, StrategyReadWriteSingleFileThruCat)

	}
	return errors.New("No strategy to write a single file")
}

func ExecuteCommandOnSpecificPod(podName string, command string) (info string, err error) {
	cmdLine := "exec " + podName + " -- " + command
	info, ok := RunOCCommand(cmdLine)
	if !ok {
		dvlog.PrintfError("Failed to execute %s", cmdLine)
		return info, errors.New("Failed to execute " + cmdLine)
	}
	if Log >= LogInfo {
		dvlog.PrintlnError(info)
	}
	return info, nil
}

func ExecuteCommandOnPod(name string, command string) (info string, err error) {
	pos := strings.Index(name, ":")
	if pos >= 0 {
		name = name[:pos]
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return "", errors.New("Pod name is empty")
	}
	podName, ok := GetPodName(name, Log < LogInfo)
	if !ok {
		return "", errors.New("There is no pod for " + podName)
	}
	return ExecuteCommandOnPodByPodName(podName, command)
}

func ExecuteCommandOnPodByPodName(podName string, command string) (info string, err error) {
	cmdLine := "exec " + podName + " -- " + command
	info, ok := RunOCCommand(cmdLine)
	if !ok {
		dvlog.PrintfError("Failed to execute %s", cmdLine)
		return info, errors.New("Failed to execute " + cmdLine)
	}
	if Log >= LogInfo {
		dvlog.Println("", info)
	}
	return info, nil
}

func readPodStreamingStrategy(strategy string) (int, bool) {
	switch strings.TrimSpace(strategy) {
	case "0":
		return 0, true
	case "1":
		return 1, true
	case "2":
		return 2, true
	}
	return 0, false
}

func parseSyncCommand(command string) (files []string, strategies []int, ok bool) {
	pos := strings.Index(command, ":")
	if pos <= 6 {
		return
	}
	operation := command[4:6]
	sample := "command " + command + " must start with copyToPod:{fileName1>podPath1;fileName2>podPath2;fileName3>podPath3[optional strategy])"
	if operation != "To" {
		operation = "From"
		sample = "command " + command + " must start with copyFromPod:{podPath1>fileName1;podPath2>fileName2;podPath3>fileName3[optional strategy])"
	}
	command = strings.TrimSpace(command[pos+1:])
	n := len(command) - 1
	if n < 0 || command[0] != '{' || command[n] != '}' {
		dvlog.PrintfError(sample)
		return
	}
	fileEntries := dvparser.ConvertToNonEmptySemicolonList(command[1:n])
	n = len(fileEntries)
	if n == 0 {
		dvlog.PrintfError(sample)
		return
	}
	files = make([]string, 2*n)
	strategies = make([]int, n)
	for i := 0; i < n; i++ {
		t := strings.TrimSpace(fileEntries[i])
		pos = strings.Index(t, "[")
		if pos >= 0 {
			endPos := strings.Index(t, "]")
			if endPos < pos {
				dvlog.PrintfError("No closing ] bracket in %s", t)
				ok = false
				return
			}
			strategies[i], ok = readPodStreamingStrategy(t[pos+1 : endPos])
			if !ok {
				dvlog.PrintfError("Not correct strategy in %s", t)
				ok = false
				return
			}
			t = t[:pos] + t[endPos+1:]
		}
		pos = strings.Index(t, ">")
		if pos <= 0 {
			dvlog.PrintfError("No > in %s", t)
			dvlog.PrintlnError(sample)
			return
		}
		f1 := strings.TrimSpace(t[:pos])
		f2 := strings.TrimSpace(t[pos+1:])
		if f1 == "" || f2 == "" {
			dvlog.PrintfError("Empty files in %s", t)
			dvlog.PrintlnError(sample)
			return
		}
		files[i<<1] = f1
		files[i<<1|1] = f2
	}
	ok = true
	return
}

func OpenShiftCopyFileFromPod(fileName string, podName string, strategy int) error {
	stat, err := os.Stat(fileName)
	if err == nil && stat.IsDir() || podName != "" && podName[len(podName)-1] == '/' {
		return OpenShiftCopyFolderFromPod(fileName, podName)
	}
	data, err := OpenShiftReadTextFile(podName, strategy)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fileName, []byte(data), 0664)
	return err
}

func OpenShiftCopyFileToPod(fileName string, podName string, strategy int) error {
	stat, err := os.Stat(fileName)
	if err != nil {
		return err
	}
	if stat.IsDir() {
		return OpenShiftCopyFolderToPod(fileName, podName)
	}
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	err = OpenShiftWriteTextFile(podName, string(data), strategy)
	return err
}

func executeSyncReadCommand(files []string, strategies []int) bool {
	n := len(strategies)
	for i := 0; i < n; i++ {
		err := OpenShiftCopyFileFromPod(files[i<<1|1], files[i<<1], strategies[i])
		if err != nil {
			dvlog.PrintfError("Failed to copy from pod (%s > %s) %s", files[i<<1], files[i<<1|1], err.Error())
			return false
		}
	}
	return true
}

func executeSyncWriteCommand(files []string, strategies []int) bool {
	n := len(strategies)
	for i := 0; i < n; i++ {
		err := OpenShiftCopyFileToPod(files[i<<1], files[i<<1|1], strategies[i])
		if err != nil {
			dvlog.PrintfError("Failed to copy to pod (%s > %s) %s", files[i<<1], files[i<<1|1], err.Error())
			return false
		}
	}
	return true
}

func copyToPodInit(command string, ctx *dvmeta.RequestContext) ([]interface{}, bool) {
	files, strategies, ok := parseSyncCommand(command)
	return []interface{}{files, strategies}, ok
}

func copyToPodRun(data []interface{}) bool {
	files := data[0].([]string)
	strategies := data[1].([]int)
	return executeSyncWriteCommand(files, strategies)
}

func copyFromPodInit(command string, ctx *dvmeta.RequestContext) ([]interface{}, bool) {
	files, strategies, ok := parseSyncCommand(command)
	return []interface{}{files, strategies}, ok
}

func copyFromPodRun(data []interface{}) bool {
	files := data[0].([]string)
	strategies := data[1].([]int)
	return executeSyncReadCommand(files, strategies)
}

func UploadFiles(microServiceName string, fileMap map[string]string, commands []string) bool {
	n := len(fileMap)
	if n == 0 {
		return true
	}
	strategies := make([]int, n)
	files := make([]string, n<<1)
	i := 0
	for k, v := range fileMap {
		strategies[i] = StrategyReadWriteSingleFileThruDir
		files[i<<1] = v
		files[i<<1|1] = microServiceName + ":" + k
		i++
	}
	if !WaitPodIsReady(microServiceName) {
		return false
	}
	ok := executeSyncWriteCommand(files, strategies)
	if !ok {
		return false
	}
	return ExecuteCommandsOnPod(microServiceName, commands)
}

func ExecuteCommandsOnPod(microServiceName string, commands []string) bool {
	m := len(commands)
	if m > 0 {
		podName, ok := GetPodName(microServiceName, Log < LogInfo)
		if !ok {
			dvlog.PrintlnError("There is no pod for " + podName)
			return false
		}
		for i := 0; i < m; i++ {
			_, err := ExecuteCommandOnPodByPodName(podName, commands[i])
			if err != nil {
				dvlog.PrintfError("Failed to execute %s: %v", commands[i], err)
				return false
			}
		}
	}
	return true
}

func WaitPodIsReady(microServiceName string) bool {
	time.Sleep(5 * time.Second)
	var podName string
	for i := 0; i < 100; i++ {
		time.Sleep(2 * time.Second)
		ResetPodNameCache()
		podName, _ = GetPodName(microServiceName, true)
		if podName != "" {
			break
		}
	}
	if podName == "" {
		dvlog.PrintfError("Waiting for pod %s getting up is timed out", microServiceName)
		return false
	}
	if Log >= LogInfo {
		dvlog.PrintfError("Waiting for 10 seconds until the pod %s is ready", podName)
	}
	time.Sleep(10 * time.Second)
	return true
}

func openShiftPodNameForFolderCopy(fullPodName string, action string, folderName string) (tmpFolder string, podName string, localFolderName string, podFolder string, extraFolder string, err error) {
	if fullPodName == "" {
		err = errors.New("Empty pod name")
		return
	}
	if fullPodName[len(fullPodName)-1] != '/' {
		fullPodName += "/"
	}
	tmpFolder, _, podName, podFolder, _, err = OpenShiftGetPodNameAndPath(fullPodName + action)
	folderName = dvtemp.GetFolderNameWithoutLastSlash(folderName)
	lastLocalPart := dvtemp.GetFolderLastPartName(folderName)
	lastPodPart := dvtemp.GetFolderLastPartName(podFolder)
	if lastLocalPart != "" && lastLocalPart == lastPodPart {
		localFolderName = dvtemp.GetFolderNameWithoutLastPart(folderName)
		dvtemp.CleanTmpFolder(tmpFolder)
		tmpFolder = ""
		return
	}
	localFolderName = tmpFolder
	extraFolder = tmpFolder + "/" + lastPodPart
	os.Mkdir(extraFolder, 0777)
	return
}

func OpenShiftCopyFolderToPod(folderName string, fullPodName string) error {
	tmpFolder, podName, localFolderName, podFolder, extraFolder, err := openShiftPodNameForFolderCopy(fullPodName, "to", folderName)
	if err != nil {
		return err
	}
	if extraFolder != "" {
		err = dvtemp.CopyWholeFolder(folderName, extraFolder)
		if err != nil {
			return err
		}
	}
	if !WriteDirectoryToPod(podName, localFolderName, podFolder, "") {
		err = errors.New("Cannot write access " + fullPodName)
		dvtemp.CleanTmpFolder(tmpFolder)
		return err
	}
	dvtemp.CleanTmpFolder(tmpFolder)
	return nil
}

func OpenShiftCopyFolderFromPod(folderName string, fullPodName string) error {
	tmpFolder, podName, localFolderName, podFolder, extraFolder, err := openShiftPodNameForFolderCopy(fullPodName, "to", folderName)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	if !ReadDirectoryFromPod(podName, localFolderName, podFolder, "") {
		err = errors.New("Cannot read access " + fullPodName)
		dvtemp.CleanTmpFolder(tmpFolder)
		return err
	}
	if extraFolder != "" {
		err = dvtemp.CopyWholeFolder(extraFolder, folderName)
		if err != nil {
			dvtemp.CleanTmpFolder(tmpFolder)
			return err
		}
	}
	dvtemp.CleanTmpFolder(tmpFolder)
	return nil
}
