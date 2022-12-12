/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const replaceServer = "{{{SRV}}}";

type PostMnConfig struct {
	InputFolder string `json:"input_folder"`
	RecordFrom int `json:"record_from"`
	RecordTo int `json:"record_to"`
	Output string `json:"output"`
	Server string `json:"server"`
	HostReplacement map[string]string `json:"host_replacement"`
	HeaderReplacement map[string]string `json:"header_replacement"`
}

type RecordInfoBlock struct {

}

type PostCollectionBlock struct {

}

func readConfig() *PostMnConfig {
   if len(os.Args)<1 {
	   fmt.Errorf("postmn.exe <configName>")
	   os.Exit(1);
   }
   s:=os.Args[0]
   dat, err:=ioutil.ReadFile(s)
   if err!=nil || len(dat)==0 {
	   fmt.Printf("Cannot read file %s (%v)", s, err)
	   os.Exit(1);
   }
   conf:=&PostMnConfig{}
   err = json.Unmarshal(dat, conf)
   if err!=nil {
	   fmt.Printf("Cannot convert file %s (%v)", s, err)
	   os.Exit(1);
   }
   if conf.Server!="" {
	   m := len(conf.Server)
	   for k,v:=range conf.HostReplacement {
			pos:=strings.Index(v, replaceServer)
			if pos>0 {
				v = v[:pos] + conf.Server + v[pos+m:]
				conf.HostReplacement[k] = v
			}
	   }
   }
   return conf;
}

func readInputs() []*RecordInfoBlock {
	return nil
}

func outputCollection() {

}

func generateCollectionInfo(conf *PostMnConfig) {

}

func main() {
	conf:=readConfig()
	readInputs()
	generateCollectionInfo(conf)
	outputCollection()
}

