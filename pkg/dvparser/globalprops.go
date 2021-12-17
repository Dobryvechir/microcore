/***********************************************************************
MicroCore
Copyright 2020-2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvparser

import (
	"log"
	"os"
	"runtime"
	"strings"
)

var MicroCorePrexix = "MICROCORE_"
var MicroCorePathSuffix = "PATH"
var IsDevelopmentVar = "IS_DEVELOPMENT"

// MicroCorePropertiesInCurrentFolderFileName is a main property file name, which can be customized
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
	osname := strings.ToLower(runtime.GOOS)
	switch osname {
	case "windows":
		mappa["IS_WINDOWS"] = "1"
	default:
		mappa["IS_LINUX"] = "1"
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
		log.Printf("No properties file: neither %s nor %s\n", MicroCorePropertiesInCurrentFolderFileName, MicroCorePrexix+MicroCorePathSuffix)
	} else if _, err := os.Stat(name); err != nil {
		panic("Properties file " + name + " not found")
	}
	return ReadPropertiesOrPanic(name)
}

func SetMicroCorePrefix(prefix string) {
	if prefix == "" {
		panic("MICROCORE prefix cannot be empty")
	}
	c := prefix[0]
	if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_') {
		panic("MICROCORE prefix must start with a Latin letter or underscore")
	}
	n := len(prefix)
	for i := 0; i < n; i++ {
		c := prefix[i]
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_' || c >= '0' && c <= '9') {
			panic("MICROCORE prefix contain only Latin letters, underscore and digits")
		}
	}
	c = prefix[n-1]
	if c != '_' {
		prefix += "_"
	}
	MicroCorePrexix = prefix
}

// SetPrefixesByApplicationName must be called at start of the application to customize prefixes and property file names
func SetPrefixesByApplicationName(name string) {
	MicroCorePropertiesInCurrentFolderFileName = name + ".properties"
	SetMicroCorePrefix(strings.ToUpper(name))
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

func GetByGlobalPropertiesOrDefault(propName, defValue string) string {
	value := defValue
	if propName != "" && GlobalProperties[propName] != "" {
		value = GlobalProperties[propName]
	}
	return value
}

func IsDevelopment() bool {
	if _, ok := GlobalProperties[IsDevelopmentVar]; ok {
		return true
	}
	return false
}
