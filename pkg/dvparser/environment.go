/***********************************************************************
MicroCore
Copyright 2020-2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvparser

import (
	"os"
	"strings"
)

var MicroCorePrexix = "MICROCORE_"
var MicroCorePathSuffix = "PATH"
var MicroCorePropertiesInCurrentFolderFileName = "MicroCore.properties"
var globalCommandLine []string
var globalPremap = makePremapOfEnvironment()

func makePremapOfEnvironment() map[string]string {
	mappa := make(map[string]string)
	for _, e := range os.Environ() {
		s := strings.TrimSpace(e)
		p := strings.IndexByte(s, '=')
		if p > 0 {
			k := s[:p]
			v := strings.TrimSpace(s[p+1:])
			mappa[k] = v
		}
	}
	commands := os.Args[1:]
	globalCommandLine = make([]string, 0, len(commands))
	for _, v := range commands {
		n := strings.Index(v, "=")
		if n > 2 && v[:2] == "-D" {
			k := v[2:n]
			vl := v[n+1:]
			if k == "MICROCORE" {
				SetMicroCorePrefix(vl)
			}
			mappa[k] = vl
		} else {
			globalCommandLine = append(globalCommandLine, v)
		}
	}
	return mappa
}

func ReadPropertiesFileWithEnvironmentVariablesInCurrentDirectory(propertiesName string) error {
	return ReadPropertiesFileWithEnvironmentVariables("", propertiesName, defaultSetFilePath)
}

func ReadPropertiesOrPanic(name string) map[string]string {
	err := ReadPropertiesFileWithEnvironmentVariablesInCurrentDirectory(name)
	if err != nil {
		panic(err.Error())
	}
	return GlobalProperties
}

func ReadPropertiesInEnvironment(name string) map[string]string {
	if name == "" || name == "." {
		name = MicroCorePropertiesInCurrentFolderFileName
		if _, err := os.Stat(name); err != nil {
			name = os.Getenv(MicroCorePrexix + MicroCorePathSuffix)
		}
	}
	if name == "" {
		panic("properties file must be either *.properties in current folder or specified in environment variable " + MicroCorePrexix + MicroCorePathSuffix)
	}
	if _, err := os.Stat(name); err != nil {
		panic("Properties file " + name + " not found")
	}
	return ReadPropertiesOrPanic(name)
}

func SetMicroCorePrefix(prefix string) {
	if prefix == "" {
		panic("MICROCORE prefix cannot be empty")
	}
	c := prefix[0]
	if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z') {
		panic("MICROCORE prefix must start with a Latin letter")
	}
	c = prefix[len(prefix)-1]
	if c != '_' {
		prefix += "_"
	}
	MicroCorePrexix = prefix
}

func GetCommandLine() []string {
	return globalCommandLine
}

func InitAndReadCommandLine() []string {
	_ = ReadPropertiesInEnvironment("")
	res, err := ConvertStringArrayByGlobalProperties(globalCommandLine, "")
	if err != nil {
		panic("Error: " + err.Error())
	}
	CallInitBySettingFuncDefault()
	return res
}
