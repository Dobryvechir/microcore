/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvprocessors

import (
	"github.com/Dobryvechir/microcore/pkg/dvcom"
	"github.com/Dobryvechir/microcore/pkg/dvmeta"
	"os"
	"strings"
)

const (
	comment_no = iota
	comment_singleLine
	comment_multiLine
	comment_HTML
)

func templateFileHandler(request *dvmeta.RequestContext) bool {
	f, err := os.Open(request.FileName)
	if err != nil {
		return false
	}
	buf := make([]byte, 4096)
	n, err1 := f.Read(buf)
	defer f.Close()
	if err1 != nil {
		return false
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
					return false
				}
				i += 2
			} else if buf[i] == '<' && i+10 < n && buf[i+1] == '!' && buf[i+2] == '-' && buf[i+3] == '-' {
				i += 4
				comment = comment_HTML
			} else {
				return false
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
		return false
	}
	if strings.ToLower(string(buf[i:i+8])) != "template" {
		return false
	}
	options := make(map[byte]string)
	for i += 8; i < n; i++ {
		if buf[i] > 'a' && buf[i] <= 'z' || buf[i] > 'A' && buf[i] <= 'Z' {
			options[buf[i]] = "_"
		} else if buf[i] != '_' {
			break
		}
	}
	stat, err3 := f.Stat()
	if err3 != nil {
		dvcom.HandleError(request, "500 "+err3.Error())
		return true
	}
	length := int(stat.Size())
	buffer := make([]byte, length)
	for j := i; j < n; j++ {
		buffer[j] = buf[j]
	}
	if n < length {
		k, err4 := f.Read(buffer[n:])
		if err4 != nil {
			dvcom.HandleError(request, "500 "+err4.Error())
			return true
		}
		n += k
	}
	switch comment {
	case comment_singleLine:
		for ; i < n; i++ {
			if buffer[i] == 10 || buffer[i] == 13 {
				i++
				if i < n && (buffer[i] == 10 || buffer[i] == 13) {
					i++
				}
				break
			}
		}
	case comment_multiLine:
		for ; i < n; i++ {
			if buffer[i] == '*' && i+1 < n && buffer[i+1] == '/' {
				i += 2
				break
			}
		}
	case comment_HTML:
		for ; i < n; i++ {
			if buffer[i] == '-' && i+2 < n && buffer[i+1] == '-' && buffer[i+2] == '>' {
				i += 3
				break
			}
		}
	}
	request.DataType = dvcom.GetContentTypeByFileName(request.FileName)
	dvtemplateProcessing(request, buffer[i:n], options)
	return true
}

var templateFileConfig *RegistrationConfig = &RegistrationConfig{
	Name:            "template",
	EndPointHandler: templateFileHandler,
}

var templateFileInited bool = RegisterProcessor(templateFileConfig, true)
