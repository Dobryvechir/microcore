/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvprocessors

import (
	"github.com/Dobryvechir/microcore/pkg/dvcom"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"os"
	"strings"
)

const (
	comment_no = iota
	comment_singleLine
	comment_multiLine
	comment_HTML
)

func readOptionsInFileHeader(fileName string) (options map[byte]string, pos int, err error) {
	f, err := os.Open(fileName)
	if err != nil {
		return
	}
	buf := make([]byte, 4096)
	n, err := f.Read(buf)
	defer f.Close()
	if err != nil {
		return
	}
	comment := comment_no
	i := 0
	for ; i < n; i++ {
		if buf[i] < 128 && buf[i] >= 33 {
			if buf[i] == '/' && i+8 < n {
				if buf[i+1] == '/' {
					comment = comment_singleLine
				} else if buf[i+1] == '*' {
					comment = comment_multiLine
				} else {
					return
				}
				i += 2
			} else if buf[i] == '<' && i+10 < n && buf[i+1] == '!' && buf[i+2] == '-' && buf[i+3] == '-' {
				i += 4
				comment = comment_HTML
			} else {
				return
			}
			break
		}
	}
	for ; i < n; i++ {
		if !(buf[i] <= 32 || buf[i] == '-' || buf[i] == '*') {
			break
		}
	}
	if i+8 >= n {
		return
	}
	if strings.ToLower(string(buf[i:i+8])) != "template" {
		return
	}
	options = make(map[byte]string)
	for i += 8; i < n; i++ {
		if buf[i] > 'a' && buf[i] <= 'z' || buf[i] > 'A' && buf[i] <= 'Z' {
			options[buf[i]] = "_"
		} else if buf[i] != '_' {
			break
		}
	}
	switch comment {
	case comment_singleLine:
		for ; i < n; i++ {
			if buf[i] == 10 || buf[i] == 13 {
				i++
				if i < n && (buf[i] == 10 || buf[i] == 13) {
					i++
				}
				pos = i
				break
			}
		}
	case comment_multiLine:
		for ; i < n; i++ {
			if buf[i] == '*' && i+1 < n && buf[i+1] == '/' {
				i += 2
				pos = i
				break
			}
		}
	case comment_HTML:
		for ; i < n; i++ {
			if buf[i] == '-' && i+2 < n && buf[i+1] == '-' && buf[i+2] == '>' {
				i += 3
				pos = i
				break
			}
		}
	}
	if pos == 0 {
		options = nil
	}
	return
}

func templateFileHandler(request *dvcontext.RequestContext) bool {
	options, pos, err := readOptionsInFileHeader(request.FileName)
	if err != nil || options == nil {
		return false
	}
	env := request.GetEnvironment()
    dat, err := dvparser.SmartReadLikeJsonTemplate(request.FileName, 3, env)
	if err != nil {
		request.Error = err
		request.HandleInternalServerError()
		return true
	}
	request.DataType = dvcom.GetContentTypeByFileName(request.FileName)
	dvtemplateProcessing(request, dat[pos:], options)
	return true
}

var templateFileConfig *RegistrationConfig = &RegistrationConfig{
	Name:            "template",
	EndPointHandler: templateFileHandler,
}

var templateFileInited bool = RegisterProcessor(templateFileConfig, true)
