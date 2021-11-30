/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

/************************************
The Version function assumes the original version to be stored as either an object
with multiple versions, for example:
{
   "module1": "0.0.0.1",
   "module2": "1.2.3",
   "module3": "20",
}
or a single string (for example: "20.0")
Versions in the same format must be stored in 2 places: source and destination.
It executes actions only if any of those versions in the source are higher
than in the destination.
In case of a single string:
I. A specified subroutine ("mainUpdate")is executed, which receives as a parameter
the highest version
||. A specified subroutine ("versionUpdate") is executed to store
the version in the destination (it is stored in VERSION_VALUE variable")
In case of multiple versions, the execution depends on the mode as follows:
"mode": = "simple" or "combined"
In case of "simple", "mainUpdate"  is executed for each key where version is higher
Parameters are "VERSION_VALUE", "VERSION_KEY"
The final version of the object is stored by special action ("versionUpdate")
In case of combined, a filter is created for multiple versions to contain
only those keys where the version is higher
It is executed only once with this object of key:version as a parameter ("mainUpdate")
For both modes, the saving of the versions in the destinations is saved
as combined
*/
package dvaction

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvjsmaster"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"log"
	"strings"
)

type VersionConfig struct {
	Source         string `json:"src"`
	Destination    string `json:"dst"`
	SourcePath     string `json:"srcPath"`
	DstPath        string `json:"dstPath"`
	Mode           string `json:"mode"`
	MainUpdate     string `json:"mainUpdate"`
	VersionUpdate  string `json:"versionUpdate"`
	VersionVarName string `json:"versionVarName"`
	KeyVarName     string `json:"keyVarName"`
}

func versionInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &VersionConfig{}
	if !DefaultInitWithObject(command, config, GetEnvironment(ctx)) {
		return nil, false
	}
	if config.Source == "" {
		log.Printf("source must be specified in %s", command)
		return nil, false
	}
	if config.Destination == "" {
		log.Printf("destination must be present in %s", command)
		return nil, false
	}
	if config.VersionVarName == "" {
		config.VersionVarName = "VERSION_VALUE"
	}
	if config.KeyVarName == "" {
		config.KeyVarName = "VERSION_KEY"
	}
	if config.Mode == "" {
		config.Mode = "simple"
	} else {
		config.Mode = strings.ToLower(config.Mode)
	}
	return []interface{}{config, ctx}, true
}

func versionRun(data []interface{}) bool {
	config := data[0].(*VersionConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return VersionRunByConfig(config, ctx)
}

func VersionRunByConfig(config *VersionConfig, ctx *dvcontext.RequestContext) bool {
	src, ok := ctx.LocalContextEnvironment.Get(config.Source)
	if !ok || src == nil {
		if ActionLog {
			dvlog.PrintfError("Empty version source %s", config.Source)
		}
		return true
	}
	srcVersion := src
	var err error
	if config.SourcePath != "" {
		srcVersion, err = dvjson.ReadPathOfAny(src, config.SourcePath, false, ctx.LocalContextEnvironment)
		if err != nil {
			dvlog.PrintfError("Failed to read %s %v", config.SourcePath, err)
			return true
		}
	}
	dst, ok := ctx.LocalContextEnvironment.Get(config.Destination)
	if !ok {
		dst = nil
	}
	dstVersion := dst
	if dst != nil && config.DstPath != "" {
		dstVersion, err = dvjson.ReadPathOfAny(dst, config.DstPath, false, ctx.LocalContextEnvironment)
		if err != nil {
			dstVersion = nil
			dst = nil
			if ActionLog {
				dvlog.PrintfError("Warning reading %s %v", config.SourcePath, err)
			}
		}
	}
	srcMap, ok := dvjson.ConvertInterfaceIntoStringMap(srcVersion)
	if !ok {
		srcMap = nil
	}
	ctx.LocalContextEnvironment.Set(config.KeyVarName, "")
	ctx.LocalContextEnvironment.Set("VERSION_SRC_OBJ", src)
	ctx.LocalContextEnvironment.Set("VERSION_DST_OBJ", dst)
	if srcMap != nil {
		dstMap, ok := dvjson.ConvertInterfaceIntoStringMap(dstVersion)
		if !ok {
			dstMap = nil
		}
		VersionKeyValueExecution(srcMap, dstMap, config, ctx)
	} else {
		s := dvevaluation.AnyToString(srcVersion)
		d := dvevaluation.AnyToString(dstVersion)
		VersionStringExecution(s, d, config, ctx)
	}
	return true
}

func VersionKeyValueExecution(src map[string]string, dst map[string]string, config *VersionConfig, ctx *dvcontext.RequestContext) {
	res := make(map[string]string)
	found := false
	if dst == nil {
		dst = make(map[string]string)
	}
	var vdst string
	var ok bool = false
	for k, vsrc := range src {
		vdst, ok = dst[k]
		if ok && vdst != "" && dvjsmaster.MathCompareVersions(vsrc, vdst, "") <= 0 {
			continue
		}
		found = true
		res[k] = vsrc
		dst[k] = vsrc
	}
	if found {
		if config.Mode == "simple" {
			for key, val := range res {
				ctx.LocalContextEnvironment.Set(config.KeyVarName, key)
				ctx.LocalContextEnvironment.Set(config.VersionVarName, val)
				versionMainUpdate(config, ctx)
			}
		} else {
			ctx.LocalContextEnvironment.Set(config.VersionVarName, res)
			versionMainUpdate(config, ctx)
		}
		ctx.LocalContextEnvironment.Set(config.VersionVarName, dst)
		versionVersionUpdate(config, ctx)
	}
}

func VersionStringExecution(src string, dst string, config *VersionConfig, ctx *dvcontext.RequestContext) {
	if dvjsmaster.MathCompareVersions(src, dst, "") > 0 {
		ctx.LocalContextEnvironment.Set(config.VersionVarName, src)
		versionMainUpdate(config, ctx)
		versionVersionUpdate(config, ctx)
	}
}

func versionMainUpdate(config *VersionConfig, ctx *dvcontext.RequestContext) {
	ExecuteAddSubsequenceShort(ctx, config.MainUpdate)
}

func versionVersionUpdate(config *VersionConfig, ctx *dvcontext.RequestContext) {
	ExecuteAddSubsequenceShort(ctx, config.VersionUpdate)
}
